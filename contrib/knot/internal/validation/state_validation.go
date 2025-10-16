package validation

import (
	"fmt"

	"github.com/denkhaus/knot/internal/errors"
	"github.com/denkhaus/knot/internal/types"
)

// StateTransition represents a valid state transition
type StateTransition struct {
	From types.TaskState
	To   types.TaskState
}

// StateValidationRule represents a validation rule for state transitions
type StateValidationRule struct {
	Name        string
	Description string
	Validate    func(from, to types.TaskState, task *types.Task) error
}

// StateValidator handles task state validation and transitions
type StateValidator struct {
	allowedTransitions map[StateTransition]bool
	validationRules    []StateValidationRule
}

// NewStateValidator creates a new state validator with default rules
func NewStateValidator() *StateValidator {
	validator := &StateValidator{
		allowedTransitions: make(map[StateTransition]bool),
		validationRules:    make([]StateValidationRule, 0),
	}

	// Define allowed state transitions
	validator.defineAllowedTransitions()
	
	// Add validation rules
	validator.addValidationRules()

	return validator
}

// defineAllowedTransitions sets up the allowed state transition matrix
func (sv *StateValidator) defineAllowedTransitions() {
	// Valid transitions from each state
	transitions := []StateTransition{
		// From pending
		{types.TaskStatePending, types.TaskStateInProgress},
		{types.TaskStatePending, types.TaskStateBlocked},
		{types.TaskStatePending, types.TaskStateCancelled},
		
		// From in-progress
		{types.TaskStateInProgress, types.TaskStateCompleted},
		{types.TaskStateInProgress, types.TaskStateBlocked},
		{types.TaskStateInProgress, types.TaskStatePending},
		{types.TaskStateInProgress, types.TaskStateCancelled},
		
		// From completed
		{types.TaskStateCompleted, types.TaskStateInProgress}, // Reopen for fixes
		{types.TaskStateCompleted, types.TaskStatePending},    // Reset if needed
		
		// From blocked
		{types.TaskStateBlocked, types.TaskStatePending},
		{types.TaskStateBlocked, types.TaskStateInProgress},
		{types.TaskStateBlocked, types.TaskStateCancelled},
		
		// From cancelled
		{types.TaskStateCancelled, types.TaskStatePending},    // Restore
		{types.TaskStateCancelled, types.TaskStateInProgress}, // Resume
		
		// Self-transitions (no-op but valid)
		{types.TaskStatePending, types.TaskStatePending},
		{types.TaskStateInProgress, types.TaskStateInProgress},
		{types.TaskStateCompleted, types.TaskStateCompleted},
		{types.TaskStateBlocked, types.TaskStateBlocked},
		{types.TaskStateCancelled, types.TaskStateCancelled},
	}

	for _, transition := range transitions {
		sv.allowedTransitions[transition] = true
	}
}

// addValidationRules adds business logic validation rules
func (sv *StateValidator) addValidationRules() {
	sv.validationRules = []StateValidationRule{
		{
			Name:        "completed_requires_progress",
			Description: "Tasks should generally go through in-progress before completion",
			Validate: func(from, to types.TaskState, task *types.Task) error {
				if from == types.TaskStatePending && to == types.TaskStateCompleted {
					return &errors.EnhancedError{
						Operation:   "validating state transition",
						Cause:       fmt.Errorf("direct transition from pending to completed"),
						Suggestion:  "Consider transitioning to in-progress first to track work progress",
						Example:     "knot task update --id " + task.ID.String() + " --state in-progress",
						HelpCommand: "knot task update --help",
					}
				}
				return nil
			},
		},
		{
			Name:        "blocked_requires_dependencies",
			Description: "Tasks should only be blocked if they have unmet dependencies",
			Validate: func(from, to types.TaskState, task *types.Task) error {
				if to == types.TaskStateBlocked && len(task.Dependencies) == 0 {
					return &errors.EnhancedError{
						Operation:   "validating state transition",
						Cause:       fmt.Errorf("cannot block task without dependencies"),
						Suggestion:  "Add dependencies first, or use a different state like pending",
						Example:     "knot dependency add --task-id " + task.ID.String() + " --depends-on <dependency-id>",
						HelpCommand: "knot dependency --help",
					}
				}
				return nil
			},
		},
		{
			Name:        "high_complexity_warning",
			Description: "High complexity tasks should be broken down before completion",
			Validate: func(from, to types.TaskState, task *types.Task) error {
				if to == types.TaskStateCompleted && task.Complexity >= 8 {
					// This is a warning, not an error - return nil but could log
					return &errors.EnhancedError{
						Operation:   "validating state transition",
						Cause:       fmt.Errorf("completing high complexity task (complexity: %d)", task.Complexity),
						Suggestion:  "Consider breaking down high complexity tasks into smaller subtasks",
						Example:     "knot breakdown --project-id <project-id> --threshold 7",
						HelpCommand: "knot task create --help  # to create subtasks",
					}
				}
				return nil
			},
		},
	}
}

// ValidateTransition validates a state transition
func (sv *StateValidator) ValidateTransition(from, to types.TaskState, task *types.Task) error {
	// Check if transition is allowed
	transition := StateTransition{From: from, To: to}
	if !sv.allowedTransitions[transition] {
		return sv.createInvalidTransitionError(from, to, task)
	}

	// Apply validation rules
	for _, rule := range sv.validationRules {
		if err := rule.Validate(from, to, task); err != nil {
			// For warnings, we might want to log but not fail
			// For now, treat all as errors
			return err
		}
	}

	return nil
}

// ValidateTransitionStrict validates with strict rules (no warnings allowed)
func (sv *StateValidator) ValidateTransitionStrict(from, to types.TaskState, task *types.Task) error {
	return sv.ValidateTransition(from, to, task)
}

// ValidateTransitionLenient validates with lenient rules (warnings become suggestions)
func (sv *StateValidator) ValidateTransitionLenient(from, to types.TaskState, task *types.Task) (error, []string) {
	// Check if transition is allowed
	transition := StateTransition{From: from, To: to}
	if !sv.allowedTransitions[transition] {
		return sv.createInvalidTransitionError(from, to, task), nil
	}

	var warnings []string
	
	// Apply validation rules and collect warnings
	for _, rule := range sv.validationRules {
		if err := rule.Validate(from, to, task); err != nil {
			// Convert errors to warnings in lenient mode
			if enhancedErr, ok := err.(*errors.EnhancedError); ok {
				warnings = append(warnings, fmt.Sprintf("⚠️  %s: %s", rule.Name, enhancedErr.Suggestion))
			} else {
				warnings = append(warnings, fmt.Sprintf("⚠️  %s: %s", rule.Name, err.Error()))
			}
		}
	}

	return nil, warnings
}

// createInvalidTransitionError creates an enhanced error for invalid transitions
func (sv *StateValidator) createInvalidTransitionError(from, to types.TaskState, task *types.Task) error {
	validTransitions := sv.getValidTransitionsFrom(from)
	
	return &errors.EnhancedError{
		Operation:   "validating state transition",
		Cause:       fmt.Errorf("invalid state transition from '%s' to '%s'", from, to),
		Suggestion:  fmt.Sprintf("Valid transitions from '%s': %v", from, validTransitions),
		Example:     fmt.Sprintf("knot task update --id %s --state %s", task.ID.String(), validTransitions[0]),
		HelpCommand: "knot task update --help",
	}
}

// getValidTransitionsFrom returns valid target states from a given state
func (sv *StateValidator) getValidTransitionsFrom(from types.TaskState) []string {
	var validStates []string
	
	for transition := range sv.allowedTransitions {
		if transition.From == from && transition.To != from { // Exclude self-transitions
			validStates = append(validStates, string(transition.To))
		}
	}
	
	return validStates
}

// GetAllValidStates returns all valid task states
func (sv *StateValidator) GetAllValidStates() []types.TaskState {
	return []types.TaskState{
		types.TaskStatePending,
		types.TaskStateInProgress,
		types.TaskStateCompleted,
		types.TaskStateBlocked,
		types.TaskStateCancelled,
	}
}

// IsValidState checks if a state string is valid
func (sv *StateValidator) IsValidState(state string) bool {
	validStates := sv.GetAllValidStates()
	for _, validState := range validStates {
		if string(validState) == state {
			return true
		}
	}
	return false
}

// GetStateTransitionMatrix returns the complete transition matrix for documentation
func (sv *StateValidator) GetStateTransitionMatrix() map[string][]string {
	matrix := make(map[string][]string)
	
	for _, state := range sv.GetAllValidStates() {
		matrix[string(state)] = sv.getValidTransitionsFrom(state)
	}
	
	return matrix
}
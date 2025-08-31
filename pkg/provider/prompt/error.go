package prompt

import (
	"fmt"

	"github.com/google/uuid"
)

// PromptError is a custom error type for prompt-related errors.
type PromptError struct {
	Message string
	AgentID uuid.UUID
	Err     error
}

func (e *PromptError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s :%s: %s", e.Message, e.AgentID.String(), e.Err.Error())
	}
	return fmt.Sprintf("%s :%s", e.Message, e.AgentID.String())
}

func (e *PromptError) Unwrap() error {
	return e.Err
}

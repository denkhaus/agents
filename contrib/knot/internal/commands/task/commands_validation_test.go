package task

import (
	"flag"
	"testing"

	"github.com/denkhaus/knot/internal/shared"
	"github.com/denkhaus/knot/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestCreateActionValidation(t *testing.T) {
	// Setup test environment
	config := testutil.NewTestConfig(t)
	mgr := config.SetupTestManager(t)
	project := testutil.CreateTestProject(t, mgr)

	tests := []struct {
		name        string
		title       string
		description string
		complexity  int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid task creation",
			title:       "Valid Task",
			description: "Valid description",
			complexity:  5,
			expectError: false,
		},
		{
			name:        "empty title should fail",
			title:       "",
			description: "Valid description",
			complexity:  5,
			expectError: true,
			errorMsg:    "title cannot be empty",
		},
		{
			name:        "title too long should fail",
			title:       string(make([]byte, 201)), // 201 characters
			description: "Valid description",
			complexity:  5,
			expectError: true,
			errorMsg:    "title too long",
		},
		{
			name:        "HTML in title should fail",
			title:       "Task with <script>alert('xss')</script>",
			description: "Valid description",
			complexity:  5,
			expectError: true,
			errorMsg:    "contains HTML tags",
		},
		{
			name:        "invalid complexity should fail",
			title:       "Valid Task",
			description: "Valid description",
			complexity:  15, // Invalid complexity
			expectError: true,
			errorMsg:    "complexity must be between 1 and 10",
		},
		{
			name:        "description too long should fail",
			title:       "Valid Task",
			description: string(make([]byte, 2001)), // 2001 characters
			complexity:  5,
			expectError: true,
			errorMsg:    "description too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create CLI context
			app := &cli.App{}
			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			flagSet.String("project-id", "", "")
			flagSet.String("title", "", "")
			flagSet.String("description", "", "")
			flagSet.String("complexity", "", "")
			flagSet.String("actor", "", "")

			flagSet.Set("project-id", project.ID.String())
			flagSet.Set("title", tt.title)
			flagSet.Set("description", tt.description)
			flagSet.Set("complexity", string(rune(tt.complexity)))
			flagSet.Set("actor", "test-user")

			ctx := cli.NewContext(app, flagSet, nil)

			// Create the action
			appCtx := &shared.AppContext{
				ProjectManager: mgr,
				Logger:         config.Logger,
			}
			action := createAction(appCtx)

			// Execute the action
			err := action(ctx)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateProjectID(t *testing.T) {
	tests := []struct {
		name        string
		projectID   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid UUID",
			projectID:   "123e4567-e89b-12d3-a456-426614174000",
			expectError: false,
		},
		{
			name:        "empty project ID",
			projectID:   "",
			expectError: true,
			errorMsg:    "project-id is required",
		},
		{
			name:        "invalid UUID format",
			projectID:   "invalid-uuid",
			expectError: true,
			errorMsg:    "invalid project ID format",
		},
		{
			name:        "partial UUID",
			projectID:   "123e4567-e89b-12d3",
			expectError: true,
			errorMsg:    "invalid project ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock CLI context
			app := &cli.App{}
			ctx := cli.NewContext(app, nil, nil)

			if tt.projectID != "" {
				ctx.Set("project-id", tt.projectID)
			}

			// Test the validation function
			_, err := shared.ValidateProjectID(ctx)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInputValidationIntegration(t *testing.T) {
	// This test ensures that our input validation is properly integrated
	// into the CLI command handlers

	config := testutil.NewTestConfig(t)
	mgr := config.SetupTestManager(t)
	project := testutil.CreateTestProject(t, mgr)

	// Test that validation errors are properly wrapped as EnhancedErrors
	app := &cli.App{}
	ctx := cli.NewContext(app, nil, nil)

	ctx.Set("project-id", project.ID.String())
	ctx.Set("title", "<script>alert('xss')</script>") // Should trigger validation error
	ctx.Set("description", "Valid description")
	ctx.Set("complexity", "5")
	ctx.Set("actor", "test-user")

	appCtx := &shared.AppContext{
		ProjectManager: mgr,
		Logger:         config.Logger,
	}
	action := createAction(appCtx)

	err := action(ctx)
	require.Error(t, err)

	// Check that it's wrapped as a validation error
	assert.Contains(t, err.Error(), "input validation")
	assert.Contains(t, err.Error(), "contains HTML tags")
}

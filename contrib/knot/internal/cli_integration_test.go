package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupCLITest creates a test environment with CLI helper
func setupCLITest(t *testing.T) (*CLITestHelper, func()) {
	helper := NewCLITestHelper(t)
	return helper, helper.Cleanup
}

// TestCLIProjectCommands tests project-related CLI commands end-to-end
func TestCLIProjectCommands(t *testing.T) {
	helper, cleanup := setupCLITest(t)
	defer cleanup()
	
	t.Run("project lifecycle", func(t *testing.T) {
		// Test project creation
		stdout, _ := helper.RunCommandExpectSuccess("project", "create", "--title", "CLI Test Project", "--description", "Testing CLI commands")
		assert.Contains(t, stdout, "CLI Test Project", "Output should contain project title")
		
		// Extract project ID
		projectID := helper.ExtractProjectID(stdout)
		require.NotEmpty(t, projectID, "Should extract project ID")
		
		// Test project listing
		stdout, _ = helper.RunCommandExpectSuccess("project", "list")
		assert.Contains(t, stdout, "CLI Test Project", "Project list should contain our project")
		assert.Contains(t, stdout, projectID, "Project list should contain project ID")
	})
	
	t.Run("project validation", func(t *testing.T) {
		// Test empty title
		_, _, err := helper.RunCommandExpectError("project", "create", "--title", "", "--description", "Empty title test")
		assert.Error(t, err, "Empty title should fail")
		
		// Test missing required flags
		_, _, err = helper.RunCommandExpectError("project", "create")
		assert.Error(t, err, "Missing required flags should fail")
	})
}

// TestCLITaskCommands tests task-related CLI commands end-to-end
func TestCLITaskCommands(t *testing.T) {
	helper, cleanup := setupCLITest(t)
	defer cleanup()
	
	t.Run("task lifecycle with project", func(t *testing.T) {
		// First create a project
		projectID := helper.CreateTestProject("Task Test Project", "For testing tasks")
		
		// Test task creation
		taskID := helper.CreateTestTask(projectID, "Test Task", "CLI test task", 3)
		require.NotEmpty(t, taskID, "Should create task successfully")
		
		// Test task listing
		stdout, _ := helper.RunCommandExpectSuccess("task", "list", "--project-id", projectID)
		assert.Contains(t, stdout, "Test Task", "Task list should contain our task")
		assert.Contains(t, stdout, taskID, "Task list should contain task ID")
	})
	
	t.Run("task validation", func(t *testing.T) {
		// Test invalid complexity
		_, _, err := helper.RunCommandExpectError("task", "create", "--title", "Invalid Task", "--complexity", "15")
		assert.Error(t, err, "Invalid complexity should fail")
		
		// Test missing project ID
		_, _, err = helper.RunCommandExpectError("task", "create", "--title", "No Project Task")
		assert.Error(t, err, "Missing project ID should fail")
	})
}

// TestCLIHealthCommands tests health and utility commands
func TestCLIHealthCommands(t *testing.T) {
	helper, cleanup := setupCLITest(t)
	defer cleanup()
	
	t.Run("health check", func(t *testing.T) {
		helper.RunCommandExpectSuccess("health")
	})
	
	t.Run("help commands", func(t *testing.T) {
		stdout, _ := helper.RunCommandExpectSuccess("--help")
		assert.Contains(t, stdout, "knot", "Help should contain program name")
		
		stdout, _ = helper.RunCommandExpectSuccess("project", "--help")
		assert.Contains(t, stdout, "project", "Project help should contain project info")
		
		stdout, _ = helper.RunCommandExpectSuccess("task", "--help")
		assert.Contains(t, stdout, "task", "Task help should contain task info")
	})
}

// TestCLIConfigCommands tests configuration commands
func TestCLIConfigCommands(t *testing.T) {
	helper, cleanup := setupCLITest(t)
	defer cleanup()
	
	t.Run("config operations", func(t *testing.T) {
		// Test config show
		helper.RunCommandExpectSuccess("config", "show")
		
		// Test config validation
		helper.RunCommandExpectSuccess("validate")
	})
}

// TestCLIErrorHandling tests error scenarios and edge cases
func TestCLIErrorHandling(t *testing.T) {
	helper, cleanup := setupCLITest(t)
	defer cleanup()
	
	t.Run("invalid commands", func(t *testing.T) {
		_, _, err := helper.RunCommandExpectError("nonexistent", "command")
		assert.Error(t, err, "Invalid command should fail")
		
		_, _, err = helper.RunCommandExpectError("project", "invalid-subcommand")
		assert.Error(t, err, "Invalid subcommand should fail")
	})
	
	t.Run("invalid flags", func(t *testing.T) {
		_, _, err := helper.RunCommandExpectError("project", "list", "--invalid-flag")
		assert.Error(t, err, "Invalid flag should fail")
	})
}

// TestCLIWorkflowCommands tests workflow and analysis commands
func TestCLIWorkflowCommands(t *testing.T) {
	helper, cleanup := setupCLITest(t)
	defer cleanup()
	
	t.Run("workflow commands", func(t *testing.T) {
		// Create a project first to have data for workflow commands
		projectID := helper.CreateTestProject("Workflow Test Project", "For testing workflow commands")
		
		// Test ready command
		stdout, _ := helper.RunCommandExpectSuccess("ready", "--project-id", projectID)
		assert.Contains(t, stdout, "Ready work", "Ready command should show ready tasks")
		
		// Test blocked command  
		stdout, _ = helper.RunCommandExpectSuccess("blocked", "--project-id", projectID)
		assert.Contains(t, stdout, "Blocked tasks", "Blocked command should show blocked tasks")
		
		// Test actionable command
		stdout, _ = helper.RunCommandExpectSuccess("actionable", "--project-id", projectID)
		assert.Contains(t, stdout, "actionable", "Actionable command should show actionable tasks")
	})
}

// TestCLIOutputFormats tests different output formats
func TestCLIOutputFormats(t *testing.T) {
	helper, cleanup := setupCLITest(t)
	defer cleanup()
	
	t.Run("output formats", func(t *testing.T) {
		// Create a project first
		helper.CreateTestProject("JSON Test Project", "For testing JSON output")
		
		// Test JSON output
		stdout, _ := helper.RunCommandExpectSuccess("project", "list", "--json")
		assert.Contains(t, stdout, "{", "JSON output should contain JSON structure")
		assert.Contains(t, stdout, "JSON Test Project", "JSON output should contain project data")
	})
}
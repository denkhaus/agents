package shell

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewToolSet(t *testing.T) {
	tests := []struct {
		name        string
		opts        []Option
		expectError bool
		errorMsg    string
	}{
		{
			name:        "default configuration",
			opts:        []Option{},
			expectError: false,
		},
		{
			name: "with custom base directory",
			opts: []Option{
				WithBaseDir("."),
			},
			expectError: false,
		},
		{
			name: "with timeout",
			opts: []Option{
				WithTimeout(10 * time.Second),
			},
			expectError: false,
		},
		{
			name: "with allowed commands",
			opts: []Option{
				WithAllowedCommands([]string{"ls", "cat", "echo"}),
			},
			expectError: false,
		},
		{
			name: "with max output size",
			opts: []Option{
				WithMaxOutputSize(512 * 1024), // 512KB
			},
			expectError: false,
		},
		{
			name: "with invalid timeout",
			opts: []Option{
				WithTimeout(-1 * time.Second),
			},
			expectError: true,
			errorMsg:    "timeout must be positive",
		},
		{
			name: "with invalid max output size",
			opts: []Option{
				WithMaxOutputSize(-1),
			},
			expectError: true,
			errorMsg:    "max output size must be positive",
		},
		{
			name: "with too large max output size",
			opts: []Option{
				WithMaxOutputSize(200 * 1024 * 1024), // 200MB
			},
			expectError: true,
			errorMsg:    "max output size cannot exceed 100MB",
		},
		{
			name: "with non-existent base directory",
			opts: []Option{
				WithBaseDir("/non/existent/directory"),
			},
			expectError: true,
			errorMsg:    "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolSet, err := NewToolSet(tt.opts...)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if toolSet == nil {
				t.Error("expected non-nil tool set")
				return
			}

			// Test that we can get tools
			tools := toolSet.Tools(context.Background())
			if len(tools) == 0 {
				t.Error("expected at least one tool")
			}

			// Test cleanup
			if err := toolSet.Close(); err != nil {
				t.Errorf("unexpected error during close: %v", err)
			}
		})
	}
}

func TestValidateInput(t *testing.T) {
	toolSet := &shellToolSet{
		baseDir:         ".",
		allowedCommands: []string{"ls", "cat", "echo"},
	}

	tests := []struct {
		name        string
		input       ShellToolInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid input",
			input: ShellToolInput{
				Command:   "ls",
				Arguments: []string{"-l"},
			},
			expectError: false,
		},
		{
			name: "empty command",
			input: ShellToolInput{
				Command: "",
			},
			expectError: true,
			errorMsg:    "command cannot be empty",
		},
		{
			name: "disallowed command",
			input: ShellToolInput{
				Command: "rm",
			},
			expectError: true,
			errorMsg:    "is not allowed",
		},
		{
			name: "dangerous argument with command substitution",
			input: ShellToolInput{
				Command:   "echo",
				Arguments: []string{"$(rm -rf /)"},
			},
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name: "dangerous argument with pipe",
			input: ShellToolInput{
				Command:   "cat",
				Arguments: []string{"file.txt | rm -rf /"},
			},
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name: "path traversal in argument",
			input: ShellToolInput{
				Command:   "cat",
				Arguments: []string{"../../../etc/passwd"},
			},
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name: "invalid work directory",
			input: ShellToolInput{
				Command: "ls",
				WorkDir: "../../../",
			},
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := toolSet.validateInput(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestIsCommandAllowed(t *testing.T) {
	tests := []struct {
		name            string
		allowedCommands []string
		command         string
		expected        bool
	}{
		{
			name:            "allowed command in custom list",
			allowedCommands: []string{"ls", "cat", "echo"},
			command:         "ls",
			expected:        true,
		},
		{
			name:            "disallowed command in custom list",
			allowedCommands: []string{"ls", "cat", "echo"},
			command:         "rm",
			expected:        false,
		},
		{
			name:            "default safe command with empty list",
			allowedCommands: []string{},
			command:         "ls",
			expected:        true,
		},
		{
			name:            "dangerous command with empty list",
			allowedCommands: []string{},
			command:         "rm",
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolSet := &shellToolSet{
				allowedCommands: tt.allowedCommands,
			}

			result := toolSet.isCommandAllowed(tt.command)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestValidateArgument(t *testing.T) {
	toolSet := &shellToolSet{}

	tests := []struct {
		name        string
		argument    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "safe argument",
			argument:    "file.txt",
			expectError: false,
		},
		{
			name:        "command substitution with $(...)",
			argument:    "$(rm -rf /)",
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "command substitution with backticks",
			argument:    "`rm -rf /`",
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "command chaining with &&",
			argument:    "file.txt && rm -rf /",
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "command chaining with ||",
			argument:    "file.txt || rm -rf /",
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "command separator with ;",
			argument:    "file.txt; rm -rf /",
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "pipe",
			argument:    "file.txt | rm -rf /",
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "output redirection",
			argument:    "file.txt > /etc/passwd",
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "path traversal",
			argument:    "../../../etc/passwd",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := toolSet.validateArgument(tt.argument)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_tool_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	toolSet := &shellToolSet{
		baseDir: tempDir,
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid relative path",
			path:        "subdir/file.txt",
			expectError: false,
		},
		{
			name:        "current directory",
			path:        ".",
			expectError: false,
		},
		{
			name:        "absolute path",
			path:        "/etc/passwd",
			expectError: true,
			errorMsg:    "absolute paths are not allowed",
		},
		{
			name:        "path traversal with ..",
			path:        "../../../etc/passwd",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "path traversal in middle",
			path:        "subdir/../../../etc/passwd",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolvedPath, err := toolSet.resolvePath(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Ensure resolved path is within base directory
			relPath, err := filepath.Rel(tempDir, resolvedPath)
			if err != nil {
				t.Errorf("failed to get relative path: %v", err)
				return
			}

			if strings.HasPrefix(relPath, "..") {
				t.Errorf("resolved path is outside base directory: %s", resolvedPath)
			}
		})
	}
}

func TestExecuteCommand(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_tool_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"echo", "cat", "ls"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	tests := []struct {
		name        string
		input       ShellToolInput
		expectError bool
		errorMsg    string
		checkOutput func(*testing.T, *ShellToolOutput)
	}{
		{
			name: "simple echo command",
			input: ShellToolInput{
				Command:   "echo",
				Arguments: []string{"hello world"},
			},
			expectError: false,
			checkOutput: func(t *testing.T, output *ShellToolOutput) {
				if !strings.Contains(output.StdOut, "hello world") {
					t.Errorf("expected output to contain 'hello world', got: %s", output.StdOut)
				}
				if output.ExitCode != 0 {
					t.Errorf("expected exit code 0, got: %d", output.ExitCode)
				}
			},
		},
		{
			name: "cat existing file",
			input: ShellToolInput{
				Command:   "cat",
				Arguments: []string{"test.txt"},
			},
			expectError: false,
			checkOutput: func(t *testing.T, output *ShellToolOutput) {
				if !strings.Contains(output.StdOut, "hello world") {
					t.Errorf("expected output to contain 'hello world', got: %s", output.StdOut)
				}
				if output.ExitCode != 0 {
					t.Errorf("expected exit code 0, got: %d", output.ExitCode)
				}
			},
		},
		{
			name: "cat non-existent file",
			input: ShellToolInput{
				Command:   "cat",
				Arguments: []string{"nonexistent.txt"},
			},
			expectError: false, // Command runs but fails with non-zero exit code
			checkOutput: func(t *testing.T, output *ShellToolOutput) {
				if output.ExitCode == 0 {
					t.Errorf("expected non-zero exit code for non-existent file")
				}
				if output.StdError == "" {
					t.Errorf("expected error output for non-existent file")
				}
			},
		},
		{
			name: "disallowed command",
			input: ShellToolInput{
				Command:   "rm",
				Arguments: []string{"test.txt"},
			},
			expectError: true,
			errorMsg:    "is not allowed",
		},
		{
			name: "dangerous argument",
			input: ShellToolInput{
				Command:   "echo",
				Arguments: []string{"$(rm -rf /)"},
			},
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			output, err := toolSet.executeCommand(ctx, tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output == nil {
				t.Error("expected non-nil output")
				return
			}

			// Check that working directory is set
			if output.WorkDir == "" {
				t.Error("expected working directory to be set")
			}

			// Run custom output checks
			if tt.checkOutput != nil {
				tt.checkOutput(t, output)
			}
		})
	}
}

func TestValidateWorkingDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_tool_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	toolSet := &shellToolSet{
		baseDir: tempDir,
	}

	tests := []struct {
		name        string
		workDir     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid working directory",
			workDir:     subDir,
			expectError: false,
		},
		{
			name:        "base directory itself",
			workDir:     tempDir,
			expectError: false,
		},
		{
			name:        "non-existent directory",
			workDir:     filepath.Join(tempDir, "nonexistent"),
			expectError: true,
			errorMsg:    "does not exist",
		},
		{
			name:        "directory outside base path",
			workDir:     "/tmp",
			expectError: true,
			errorMsg:    "outside base path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := toolSet.validateWorkingDirectory(tt.workDir)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetRestrictedEnvironment(t *testing.T) {
	toolSet := &shellToolSet{}

	env := toolSet.getRestrictedEnvironment()

	// Check that we have some basic environment variables
	foundPath := false
	for _, envVar := range env {
		if strings.HasPrefix(envVar, "PATH=") {
			foundPath = true
			break
		}
	}

	if !foundPath {
		t.Error("expected PATH environment variable to be set")
	}

	// Check that we don't have too many environment variables (security)
	if len(env) > 20 {
		t.Errorf("environment has too many variables (%d), potential security risk", len(env))
	}
}

func TestChangeDirectory(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "shell_cd_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories
	subDir1 := filepath.Join(tempDir, "subdir1")
	subDir2 := filepath.Join(tempDir, "subdir2")
	nestedDir := filepath.Join(subDir1, "nested")

	for _, dir := range []string{subDir1, subDir2, nestedDir} {
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", dir, err)
		}
	}

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"cd"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	tests := []struct {
		name            string
		targetDir       string
		expectError     bool
		expectedWorkDir string
		errorMsg        string
	}{
		{
			name:            "change to subdirectory",
			targetDir:       "subdir1",
			expectError:     false,
			expectedWorkDir: subDir1,
		},
		{
			name:            "change to current directory",
			targetDir:       ".",
			expectError:     false,
			expectedWorkDir: tempDir, // Should stay in current directory
		},
		{
			name:            "change to parent directory from subdir",
			targetDir:       "..",
			expectError:     false,
			expectedWorkDir: "", // Will be set dynamically in test
		},
		{
			name:        "change to non-existent directory",
			targetDir:   "nonexistent",
			expectError: false, // Command succeeds but returns error in output
		},
		{
			name:        "absolute path not allowed",
			targetDir:   "/tmp",
			expectError: false, // Command succeeds but returns error in output
		},
		{
			name:        "path traversal attempt",
			targetDir:   "../../..",
			expectError: false, // Command runs but fails with security error
		},
		{
			name:            "home directory (~)",
			targetDir:       "~",
			expectError:     false,
			expectedWorkDir: tempDir, // Should go to base directory
		},
		{
			name:            "home subdirectory (~/subdir1)",
			targetDir:       "~/subdir1",
			expectError:     false,
			expectedWorkDir: subDir1, // Should go to base/subdir1
		},
		{
			name:        "home with path traversal (~/../../etc)",
			targetDir:   "~/../../etc",
			expectError: false, // Command runs but fails with security error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to base directory for each test
			toolSet.currentWorkDir = tempDir

			// Special setup for parent directory test
			if tt.name == "change to parent directory from subdir" {
				// First change to subdirectory
				toolSet.currentWorkDir = subDir1
				tt.expectedWorkDir = tempDir
			}

			input := ShellToolInput{
				Command:   "cd",
				Arguments: []string{tt.targetDir},
			}

			ctx := context.Background()
			result, err := toolSet.executeCommand(ctx, input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected non-nil result")
				return
			}

			// Check if command succeeded or failed as expected
			if tt.expectedWorkDir != "" {
				if result.ExitCode != 0 {
					t.Errorf("expected success (exit code 0), got exit code %d: %s", result.ExitCode, result.StdError)
					return
				}

				// Check that current working directory was updated
				if toolSet.currentWorkDir != tt.expectedWorkDir {
					t.Errorf("expected current working directory to be %s, got %s", tt.expectedWorkDir, toolSet.currentWorkDir)
				}

				// Check that result working directory matches
				if result.WorkDir != tt.expectedWorkDir {
					t.Errorf("expected result working directory to be %s, got %s", tt.expectedWorkDir, result.WorkDir)
				}
			} else {
				// Expect command to fail
				if result.ExitCode == 0 {
					t.Errorf("expected command to fail, but got exit code 0")
				}
			}
		})
	}
}

func TestChangeDirectoryTool(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "shell_cd_tool_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectory
	subDir := filepath.Join(tempDir, "testdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"cd"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	tests := []struct {
		name        string
		input       ChangeDirectoryInput
		expectError bool
		checkOutput func(*testing.T, *ChangeDirectoryOutput)
	}{
		{
			name: "successful directory change",
			input: ChangeDirectoryInput{
				Path: "testdir",
			},
			expectError: false,
			checkOutput: func(t *testing.T, output *ChangeDirectoryOutput) {
				if output.NewWorkDir != "testdir" {
					t.Errorf("expected new working directory to be 'testdir', got '%s'", output.NewWorkDir)
				}
				if output.OldWorkDir != "." {
					t.Errorf("expected old working directory to be '.', got '%s'", output.OldWorkDir)
				}
				if output.Error != "" {
					t.Errorf("expected no error, got: %s", output.Error)
				}
			},
		},
		{
			name: "directory does not exist",
			input: ChangeDirectoryInput{
				Path: "nonexistent",
			},
			expectError: false, // Function succeeds but returns error in output
			checkOutput: func(t *testing.T, output *ChangeDirectoryOutput) {
				if output.Error == "" {
					t.Error("expected error for non-existent directory")
				}
				if !strings.Contains(output.Error, "No such file or directory") && !strings.Contains(output.Error, "does not exist") {
					t.Errorf("expected 'No such file or directory' or 'does not exist' error, got: %s", output.Error)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to base directory for each test
			toolSet.currentWorkDir = tempDir

			ctx := context.Background()
			output, err := toolSet.changeDirectory(ctx, tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output == nil {
				t.Error("expected non-nil output")
				return
			}

			// Run custom output checks
			if tt.checkOutput != nil {
				tt.checkOutput(t, output)
			}
		})
	}
}

func TestChangeDirectorySecurityValidation(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_cd_security_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"cd"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	tests := []struct {
		name        string
		targetDir   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "absolute path blocked",
			targetDir:   "/tmp",
			expectError: false, // Command runs but fails
		},
		{
			name:        "path traversal blocked",
			targetDir:   "../../../etc",
			expectError: false, // Command runs but fails with security error
		},
		{
			name:        "multiple path traversal blocked",
			targetDir:   "../../..",
			expectError: false, // Command runs but fails with security error
		},
		{
			name:        "home with path traversal blocked",
			targetDir:   "~/../../etc",
			expectError: false, // Command runs but fails with security error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := ShellToolInput{
				Command:   "cd",
				Arguments: []string{tt.targetDir},
			}

			ctx := context.Background()
			result, err := toolSet.executeCommand(ctx, input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// For security violations, command should succeed but return error in output
			if result.ExitCode == 0 {
				t.Errorf("expected command to fail for security violation")
			}

			if result.StdError == "" {
				t.Errorf("expected error message in stderr")
			}
		})
	}
}

// Benchmark tests
func BenchmarkValidateInput(b *testing.B) {
	toolSet := &shellToolSet{
		baseDir:         ".",
		allowedCommands: []string{"ls", "cat", "echo"},
	}

	input := ShellToolInput{
		Command:   "ls",
		Arguments: []string{"-l", "*.go"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = toolSet.validateInput(input)
	}
}

func BenchmarkResolvePath(b *testing.B) {
	toolSet := &shellToolSet{
		baseDir: ".",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = toolSet.resolvePath("subdir/file.txt")
	}
}

package shell

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAbsolutePathSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_security_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"ls", "cat", "grep", "find"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	// Test cases that should be BLOCKED
	securityTests := []struct {
		name        string
		command     string
		arguments   []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "ls with absolute path to /etc",
			command:     "ls",
			arguments:   []string{"-la", "/etc"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "ls with absolute path to /home",
			command:     "ls",
			arguments:   []string{"-la", "/home/denkhaus/.qwen"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "cat with absolute path",
			command:     "cat",
			arguments:   []string{"/etc/passwd"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "grep with absolute path",
			command:     "grep",
			arguments:   []string{"root", "/etc/passwd"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "find with absolute path",
			command:     "find",
			arguments:   []string{"/home", "-name", "*.txt"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "ls with system directory path",
			command:     "ls",
			arguments:   []string{"/bin"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "environment variable expansion",
			command:     "ls",
			arguments:   []string{"$HOME/.qwen"},
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "variable expansion with braces",
			command:     "cat",
			arguments:   []string{"${HOME}/.bashrc"},
			expectError: true,
			errorMsg:    "dangerous pattern",
		},
		{
			name:        "path traversal to system directory",
			command:     "ls",
			arguments:   []string{"../../../etc"},
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range securityTests {
		t.Run(tt.name, func(t *testing.T) {
			input := ShellToolInput{
				Command:   tt.command,
				Arguments: tt.arguments,
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
				t.Logf("Security test passed: %s blocked with error: %s", strings.Join(tt.arguments, " "), err.Error())
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
		})
	}
}

func TestSystemDirectoryAccess(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_system_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"ls", "cat"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	// System directories that should be blocked
	systemDirs := []string{"/etc", "/bin", "/usr", "/var", "/tmp", "/home", "/root", "/proc", "/sys", "/dev"}

	for _, sysDir := range systemDirs {
		t.Run("block_access_to_"+sysDir, func(t *testing.T) {
			input := ShellToolInput{
				Command:   "ls",
				Arguments: []string{sysDir},
			}

			ctx := context.Background()
			_, err := toolSet.executeCommand(ctx, input)

			if err == nil {
				t.Errorf("expected error for system directory access: %s", sysDir)
				return
			}

			if !strings.Contains(err.Error(), "absolute path outside workspace boundary") {
				t.Errorf("expected workspace boundary error, got: %s", err.Error())
			}

			t.Logf("System directory access blocked: %s", sysDir)
		})
	}
}

func TestEnvironmentVariableBlocking(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_env_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"ls", "cat", "echo"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	envTests := []struct {
		name      string
		command   string
		arguments []string
		errorMsg  string
	}{
		{
			name:      "HOME variable",
			command:   "ls",
			arguments: []string{"$HOME"},
			errorMsg:  "dangerous pattern",
		},
		{
			name:      "PATH variable",
			command:   "echo",
			arguments: []string{"$PATH"},
			errorMsg:  "dangerous pattern",
		},
		{
			name:      "USER variable",
			command:   "echo",
			arguments: []string{"$USER"},
			errorMsg:  "dangerous pattern",
		},
		{
			name:      "HOME with braces",
			command:   "cat",
			arguments: []string{"${HOME}/.bashrc"},
			errorMsg:  "dangerous pattern",
		},
		{
			name:      "Complex variable expansion",
			command:   "ls",
			arguments: []string{"${HOME}/../other_user"},
			errorMsg:  "dangerous pattern",
		},
	}

	for _, tt := range envTests {
		t.Run(tt.name, func(t *testing.T) {
			input := ShellToolInput{
				Command:   tt.command,
				Arguments: tt.arguments,
			}

			ctx := context.Background()
			_, err := toolSet.executeCommand(ctx, input)

			if err == nil {
				t.Errorf("expected error for environment variable usage")
				return
			}

			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
			}

			t.Logf("Environment variable blocked: %s", strings.Join(tt.arguments, " "))
		})
	}
}

func TestSecurityBypass(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_bypass_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"ls", "cat", "find"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	// Advanced bypass attempts that should be blocked
	bypassTests := []struct {
		name      string
		command   string
		arguments []string
		errorMsg  string
	}{
		{
			name:      "symbolic link to system directory",
			command:   "ls",
			arguments: []string{"/etc/../etc/passwd"},
			errorMsg:  "absolute path outside workspace boundary",
		},
		{
			name:      "double slash absolute path",
			command:   "cat",
			arguments: []string{"//etc/passwd"},
			errorMsg:  "absolute path outside workspace boundary",
		},
		{
			name:      "find with absolute start path",
			command:   "find",
			arguments: []string{"/", "-name", "passwd"},
			errorMsg:  "absolute path outside workspace boundary",
		},
		{
			name:      "mixed relative and absolute",
			command:   "ls",
			arguments: []string{".", "/etc"},
			errorMsg:  "absolute path outside workspace boundary",
		},
	}

	for _, tt := range bypassTests {
		t.Run(tt.name, func(t *testing.T) {
			input := ShellToolInput{
				Command:   tt.command,
				Arguments: tt.arguments,
			}

			ctx := context.Background()
			_, err := toolSet.executeCommand(ctx, input)

			if err == nil {
				t.Errorf("expected error for bypass attempt")
				return
			}

			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
			}

			t.Logf("Bypass attempt blocked: %s %s", tt.command, strings.Join(tt.arguments, " "))
		})
	}
}

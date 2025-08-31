package shell

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAbsolutePathWorkspaceValidation(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "shell_workspace_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories within workspace
	subDir1 := filepath.Join(tempDir, "documents")
	subDir2 := filepath.Join(tempDir, "projects")
	nestedDir := filepath.Join(subDir1, "work")

	for _, dir := range []string{subDir1, subDir2, nestedDir} {
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create test files
	testFile1 := filepath.Join(subDir1, "test.txt")
	testFile2 := filepath.Join(nestedDir, "nested.txt")
	if err := os.WriteFile(testFile1, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("nested content"), 0644); err != nil {
		t.Fatalf("failed to create nested test file: %v", err)
	}

	toolSet := &shellToolSet{
		baseDir:         tempDir,
		currentWorkDir:  tempDir,
		allowedCommands: []string{"ls", "cat", "find"},
		timeout:         5 * time.Second,
		maxOutputSize:   1024 * 1024,
	}

	tests := []struct {
		name        string
		command     string
		arguments   []string
		expectError bool
		errorMsg    string
		description string
	}{
		{
			name:        "absolute path within workspace - ls",
			command:     "ls",
			arguments:   []string{"-la", subDir1},
			expectError: false,
			description: "Should allow absolute paths within workspace",
		},
		{
			name:        "absolute path within workspace - cat",
			command:     "cat",
			arguments:   []string{testFile1},
			expectError: false,
			description: "Should allow absolute file paths within workspace",
		},
		{
			name:        "absolute path within workspace - nested",
			command:     "cat",
			arguments:   []string{testFile2},
			expectError: false,
			description: "Should allow absolute paths to nested files within workspace",
		},
		{
			name:        "absolute path outside workspace - /etc",
			command:     "ls",
			arguments:   []string{"/etc"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
			description: "Should block absolute paths outside workspace",
		},
		{
			name:        "absolute path outside workspace - /home",
			command:     "cat",
			arguments:   []string{"/home/user/.bashrc"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
			description: "Should block access to user home directories",
		},
		{
			name:        "absolute path outside workspace - parent",
			command:     "ls",
			arguments:   []string{filepath.Dir(tempDir)},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
			description: "Should block access to parent of workspace",
		},
		{
			name:        "absolute path outside workspace - root",
			command:     "find",
			arguments:   []string{"/", "-name", "*.txt"},
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
			description: "Should block access to filesystem root",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := ShellToolInput{
				Command:   tt.command,
				Arguments: tt.arguments,
			}

			ctx := context.Background()
			result, err := toolSet.executeCommand(ctx, input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for: %s", tt.description)
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				t.Logf("✅ Correctly blocked: %s %s - %s", tt.command, strings.Join(tt.arguments, " "), err.Error())
				return
			}

			if err != nil {
				t.Errorf("unexpected error for valid workspace path: %v", err)
				return
			}

			if result == nil {
				t.Error("expected non-nil result")
				return
			}

			if result.ExitCode != 0 {
				t.Errorf("command failed unexpectedly: %s", result.StdError)
				return
			}

			t.Logf("✅ Correctly allowed: %s %s - %s", tt.command, strings.Join(tt.arguments, " "), tt.description)
		})
	}
}

func TestValidateAbsolutePathWithinWorkspace(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_path_validation_test")
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
			name:        "path within workspace",
			path:        filepath.Join(tempDir, "subdir", "file.txt"),
			expectError: false,
		},
		{
			name:        "path exactly at workspace root",
			path:        tempDir,
			expectError: false,
		},
		{
			name:        "path outside workspace - parent",
			path:        filepath.Dir(tempDir),
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "path outside workspace - sibling",
			path:        filepath.Join(filepath.Dir(tempDir), "other_dir"),
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "path outside workspace - system directory",
			path:        "/etc/passwd",
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
		{
			name:        "path outside workspace - root",
			path:        "/",
			expectError: true,
			errorMsg:    "absolute path outside workspace boundary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := toolSet.validateAbsolutePathWithinWorkspace(tt.path)

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

func TestWorkspacePathEdgeCases(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_edge_case_test")
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

	// Test edge cases
	edgeCases := []struct {
		name        string
		path        string
		expectError bool
		description string
	}{
		{
			name:        "workspace with trailing slash",
			path:        tempDir + "/",
			expectError: false,
			description: "Should handle trailing slashes correctly",
		},
		{
			name:        "workspace with double slashes",
			path:        strings.Replace(tempDir, "/", "//", 1),
			expectError: false,
			description: "Should handle double slashes correctly",
		},
		{
			name:        "workspace with dot notation",
			path:        tempDir + "/./subdir",
			expectError: false,
			description: "Should handle dot notation within workspace",
		},
	}

	for _, tt := range edgeCases {
		t.Run(tt.name, func(t *testing.T) {
			err := toolSet.validateAbsolutePathWithinWorkspace(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for: %s", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error for valid case: %v - %s", err, tt.description)
			}

			t.Logf("✅ Edge case handled correctly: %s", tt.description)
		})
	}
}

package shelltoolset

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTildeExpansion(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "shell_tilde_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories
	subDir1 := filepath.Join(tempDir, "documents")
	subDir2 := filepath.Join(tempDir, "projects")
	nestedDir := filepath.Join(subDir1, "work")
	
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
		startDir        string
		targetDir       string
		expectError     bool
		expectedWorkDir string
		errorMsg        string
	}{
		{
			name:            "tilde to home from base directory",
			startDir:        tempDir,
			targetDir:       "~",
			expectError:     false,
			expectedWorkDir: tempDir,
		},
		{
			name:            "tilde to home from subdirectory",
			startDir:        subDir1,
			targetDir:       "~",
			expectError:     false,
			expectedWorkDir: tempDir,
		},
		{
			name:            "tilde with subdirectory from base",
			startDir:        tempDir,
			targetDir:       "~/documents",
			expectError:     false,
			expectedWorkDir: subDir1,
		},
		{
			name:            "tilde with subdirectory from different dir",
			startDir:        subDir2,
			targetDir:       "~/documents",
			expectError:     false,
			expectedWorkDir: subDir1,
		},
		{
			name:            "tilde with nested path",
			startDir:        tempDir,
			targetDir:       "~/documents/work",
			expectError:     false,
			expectedWorkDir: nestedDir,
		},
		{
			name:        "tilde with path traversal attempt",
			startDir:    tempDir,
			targetDir:   "~/../../etc",
			expectError: false, // Command runs but fails
			errorMsg:    "access denied - cannot navigate outside the allowed workspace boundary",
		},
		{
			name:        "tilde with multiple path traversal",
			startDir:    subDir1,
			targetDir:   "~/../../../tmp",
			expectError: false, // Command runs but fails
			errorMsg:    "access denied - cannot navigate outside the allowed workspace boundary",
		},
		{
			name:        "tilde to non-existent directory",
			startDir:    tempDir,
			targetDir:   "~/nonexistent",
			expectError: false, // Command runs but fails
			errorMsg:    "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set starting directory
			toolSet.currentWorkDir = tt.startDir
			
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
				
				// Check error message if specified
				if tt.errorMsg != "" {
					if !strings.Contains(result.StdError, tt.errorMsg) && !strings.Contains(result.Error, tt.errorMsg) {
						t.Errorf("expected error message to contain '%s', got stderr: '%s', error: '%s'", 
							tt.errorMsg, result.StdError, result.Error)
					}
				}
			}
		})
	}
}

func TestTildeExpansionEdgeCases(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_tilde_edge_test")
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
		description string
	}{
		{
			name:        "tilde only",
			targetDir:   "~",
			expectError: false,
			description: "Simple tilde should work",
		},
		{
			name:        "tilde with slash",
			targetDir:   "~/",
			expectError: false,
			description: "Tilde with trailing slash should work",
		},
		{
			name:        "tilde with dot",
			targetDir:   "~/.",
			expectError: false,
			description: "Tilde with current directory should work",
		},
		{
			name:        "tilde with double dot",
			targetDir:   "~/..",
			expectError: false,
			description: "Tilde with parent directory should be blocked if outside base",
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
				if err == nil && result.ExitCode == 0 {
					t.Errorf("expected error or failure but command succeeded")
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

			// For successful cases, just verify the command completed
			t.Logf("Test '%s': %s - Exit code: %d", tt.name, tt.description, result.ExitCode)
		})
	}
}

func TestTildeSecurityValidation(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell_tilde_security_test")
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

	securityTests := []struct {
		name      string
		targetDir string
		errorMsg  string
	}{
		{
			name:      "tilde path traversal to etc",
			targetDir: "~/../../etc",
			errorMsg:  "access denied - cannot navigate outside the allowed workspace boundary",
		},
		{
			name:      "tilde path traversal to root",
			targetDir: "~/../../../",
			errorMsg:  "access denied - cannot navigate outside the allowed workspace boundary",
		},
		{
			name:      "tilde path traversal to tmp",
			targetDir: "~/../../../../tmp",
			errorMsg:  "access denied - cannot navigate outside the allowed workspace boundary",
		},
	}

	for _, tt := range securityTests {
		t.Run(tt.name, func(t *testing.T) {
			input := ShellToolInput{
				Command:   "cd",
				Arguments: []string{tt.targetDir},
			}

			ctx := context.Background()
			result, err := toolSet.executeCommand(ctx, input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Security violations should result in failed command with clear message
			if result.ExitCode == 0 {
				t.Errorf("expected command to fail for security violation")
			}

			if !strings.Contains(result.StdError, tt.errorMsg) && !strings.Contains(result.Error, tt.errorMsg) {
				t.Errorf("expected error message to contain '%s', got stderr: '%s', error: '%s'", 
					tt.errorMsg, result.StdError, result.Error)
			}

			t.Logf("Security test passed: %s blocked with message: %s", tt.targetDir, result.StdError)
		})
	}
}
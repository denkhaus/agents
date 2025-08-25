package shelltoolset

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/tool"
)

const (
	// defaultBaseDir is the default base directory for file operations.
	defaultBaseDir = "."
)

// shellToolSet implements the ToolSet interface for shell operations.
type shellToolSet struct {
	baseDir               string
	executeCommandEnabled bool
	allowedCommands       []string
	tools                 []tool.CallableTool
	timeout               time.Duration
	maxOutputSize         int64
	currentWorkDir        string // Current working directory for cd command
}

// Option is a functional option for configuring the file tool set.
type Option func(*shellToolSet)

// WithBaseDir sets the base directory for file operations, default is the current directory.
func WithBaseDir(baseDir string) Option {
	return func(f *shellToolSet) {
		f.baseDir = baseDir
	}
}

// WithTimeout sets the timeout for command execution.
func WithTimeout(t time.Duration) Option {
	return func(f *shellToolSet) {
		f.timeout = t
	}
}

// WithExecuteCommandEnabled enables or disables the command execution functionality.
func WithExecuteCommandEnabled(e bool) Option {
	return func(f *shellToolSet) {
		f.executeCommandEnabled = e
	}
}

// WithAllowedCommands sets the list of allowed commands. If empty, uses default safe list.
func WithAllowedCommands(commands []string) Option {
	return func(f *shellToolSet) {
		f.allowedCommands = make([]string, len(commands))
		copy(f.allowedCommands, commands)
	}
}

// WithMaxOutputSize sets the maximum output size in bytes (default: 1MB).
func WithMaxOutputSize(size int64) Option {
	return func(f *shellToolSet) {
		f.maxOutputSize = size
	}
}

// NewToolSet creates a new shell operation tool set with the provided options.
func NewToolSet(opts ...Option) (tool.ToolSet, error) {
	// Apply default configuration.
	shellToolSet := &shellToolSet{
		baseDir:               defaultBaseDir,
		executeCommandEnabled: true,
		allowedCommands:       []string{}, // Empty means use default safe list
		timeout:               30 * time.Second,
		maxOutputSize:         1024 * 1024, // 1MB default
		currentWorkDir:        "",          // Will be set to baseDir after validation
	}

	// Apply user-provided options.
	for _, opt := range opts {
		opt(shellToolSet)
	}

	// Clean and validate the base directory.
	shellToolSet.baseDir = filepath.Clean(shellToolSet.baseDir)

	// Convert to absolute path for security
	absBaseDir, err := filepath.Abs(shellToolSet.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for base directory: %w", err)
	}
	shellToolSet.baseDir = absBaseDir

	// Check if the base directory exists.
	stat, err := os.Stat(shellToolSet.baseDir)
	if err != nil {
		return nil, fmt.Errorf("base directory '%s' does not exist: %w", shellToolSet.baseDir, err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("base directory '%s' is not a directory", shellToolSet.baseDir)
	}

	// Validate configuration
	if err := shellToolSet.validateConfiguration(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize current working directory to base directory
	shellToolSet.currentWorkDir = shellToolSet.baseDir

	// Create function tools based on enabled features.
	var tools []tool.CallableTool
	if shellToolSet.executeCommandEnabled {
		tools = append(tools, shellToolSet.executeCommandTool())
		tools = append(tools, shellToolSet.changeDirectoryTool())
	}
	shellToolSet.tools = tools

	return shellToolSet, nil
}

// validateConfiguration validates the tool set configuration
func (f *shellToolSet) validateConfiguration() error {
	if f.timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if f.maxOutputSize <= 0 {
		return fmt.Errorf("max output size must be positive")
	}
	if f.maxOutputSize > 100*1024*1024 { // 100MB limit
		return fmt.Errorf("max output size cannot exceed 100MB")
	}
	return nil
}

// resolvePath validates a path to prevent directory traversal attacks,
// and resolves a relative path within the base directory.
func (f *shellToolSet) resolvePath(relativePath string) (string, error) {
	// Clean the path first
	cleanPath := filepath.Clean(relativePath)

	// Check for absolute paths
	if filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("absolute paths are not allowed: %s", relativePath)
	}

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("path traversal detected: %s", relativePath)
	}

	// Join with base directory
	fullPath := filepath.Join(f.baseDir, cleanPath)

	// Get absolute paths for comparison
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	absBaseDir, err := filepath.Abs(f.baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute base directory: %w", err)
	}

	// Ensure the resolved path is still within the base directory
	relPath, err := filepath.Rel(absBaseDir, absFullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	if strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("resolved path is outside base directory: %s", relativePath)
	}

	return absFullPath, nil
}

// Tools returns the list of available tools in this tool set.
func (f *shellToolSet) Tools(ctx context.Context) []tool.CallableTool {
	return f.tools
}

// Close cleans up any resources used by the tool set.
func (f *shellToolSet) Close() error {
	// No resources to clean up for shell tool set
	return nil
}

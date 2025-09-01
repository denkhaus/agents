//cue:generate cue get go github.com/denkhaus/agents/pkg/tools/shell

package shell

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/denkhaus/agents/pkg/tools"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

const (
	ToolSetName = "shell_toolset"
	// defaultBaseDir is the default base directory for file operations.
	defaultBaseDir = "."
)

// ToolSetConfig holds configuration for the shell toolset
type ToolSetConfig struct {
	BaseDir               string        `json:"base_dir,omitempty" mapstructure:"base_dir,omitempty"`
	ExecuteCommandEnabled bool          `json:"execute_command_enabled" mapstructure:"execute_command_enabled"`
	AllowedCommands       []string      `json:"allowed_commands,omitempty" mapstructure:"allowed_commands,omitempty"`
	Timeout               time.Duration `json:"timeout,omitempty" mapstructure:"timeout,omitempty"`
	MaxOutputSize         int64         `json:"max_output_size,omitempty" mapstructure:"max_output_size,omitempty"`
}

// shellToolSet implements the ToolSet interface for shell operations.
type shellToolSet struct {
	ToolSetConfig
	tools          []tool.CallableTool
	currentWorkDir string // Current working directory for cd command
}

func NewWithDI(injector *do.Injector) (tools.ToolSetFactoryFunc, error) {
	return func(config tools.ConfigPayload) (tool.ToolSet, error) {
		// Extract configuration and convert to options
		var settings shellToolSet
		if err := config.Bind(&settings.ToolSetConfig); err != nil {
			return nil, err
		}

		// Create options from settings
		var opts []Option
		if settings.BaseDir != "" {
			opts = append(opts, WithBaseDir(settings.BaseDir))
		}
		opts = append(opts, WithExecuteCommandEnabled(settings.ExecuteCommandEnabled))
		if len(settings.AllowedCommands) > 0 {
			opts = append(opts, WithAllowedCommands(settings.AllowedCommands))
		}
		if settings.Timeout > 0 {
			opts = append(opts, WithTimeout(settings.Timeout))
		}
		if settings.MaxOutputSize > 0 {
			opts = append(opts, WithMaxOutputSize(settings.MaxOutputSize))
		}

		return New(opts...)
	}, nil
}

// NewToolSet creates a new shell operation tool set with the provided options.
func New(opts ...Option) (tool.ToolSet, error) {
	// Apply default configuration.
	shellToolSet := &shellToolSet{
		ToolSetConfig: ToolSetConfig{
			BaseDir:               defaultBaseDir,
			ExecuteCommandEnabled: true,
			AllowedCommands:       []string{}, // Empty means use default safe list
			Timeout:               30 * time.Second,
			MaxOutputSize:         1024 * 1024, // 1MB default
		},

		currentWorkDir: "", // Will be set to baseDir after validation
	}

	// Apply user-provided options.
	for _, opt := range opts {
		opt(shellToolSet)
	}

	// Clean and validate the base directory.
	shellToolSet.BaseDir = filepath.Clean(shellToolSet.BaseDir)

	// Convert to absolute path for security
	absBaseDir, err := filepath.Abs(shellToolSet.BaseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for base directory: %w", err)
	}
	shellToolSet.BaseDir = absBaseDir

	// Check if the base directory exists.
	stat, err := os.Stat(shellToolSet.BaseDir)
	if err != nil {
		return nil, fmt.Errorf("base directory '%s' does not exist: %w", shellToolSet.BaseDir, err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("base directory '%s' is not a directory", shellToolSet.BaseDir)
	}

	// Validate configuration
	if err := shellToolSet.validateConfiguration(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize current working directory to base directory
	shellToolSet.currentWorkDir = shellToolSet.BaseDir

	// Create function tools based on enabled features.
	var tools []tool.CallableTool
	if shellToolSet.ExecuteCommandEnabled {
		tools = append(tools, shellToolSet.executeCommandTool())
		tools = append(tools, shellToolSet.changeDirectoryTool())
	}
	shellToolSet.tools = tools

	return shellToolSet, nil
}

// validateConfiguration validates the tool set configuration
func (f *shellToolSet) validateConfiguration() error {
	if f.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if f.MaxOutputSize <= 0 {
		return fmt.Errorf("max output size must be positive")
	}
	if f.MaxOutputSize > 100*1024*1024 { // 100MB limit
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
	// Allow "./..." pattern which is common in Go for recursive package operations
	if strings.Contains(cleanPath, "..") && cleanPath != "./..." {
		return "", fmt.Errorf("path traversal detected: %s", relativePath)
	}

	// Join with base directory
	fullPath := filepath.Join(f.BaseDir, cleanPath)

	// Get absolute paths for comparison
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	absBaseDir, err := filepath.Abs(f.BaseDir)
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

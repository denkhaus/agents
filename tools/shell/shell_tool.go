package shelltoolset

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// ShellToolInput represents the input structure for the shell tool
type ShellToolInput struct {
	Command   string   `json:"command" description:"The command to execute (must be in allowed commands list)"`
	Arguments []string `json:"arguments,omitempty" description:"Arguments to pass to the command"`
	WorkDir   string   `json:"work_dir,omitempty" description:"Working directory relative to base path (optional)"`
}

type ShellToolOutput struct {
	StdOut   string `json:"stdout" description:"Standard output from the command"`
	StdError string `json:"stderr" description:"Standard error from the command"`
	ExitCode int    `json:"exit_code" description:"Exit code of the command"`
	Error    string `json:"error,omitempty" description:"Error message if command failed"`
	WorkDir  string `json:"work_dir" description:"Working directory where command was executed"`
}

// validateInput checks if the provided input is valid for this tool
func (t *shellToolSet) validateInput(input ShellToolInput) error {
	if input.Command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Check if command is in allowed list
	if !t.isCommandAllowed(input.Command) {
		return fmt.Errorf("command '%s' is not allowed. Allowed commands: %v", input.Command, t.allowedCommands)
	}

	// Validate working directory if provided
	if input.WorkDir != "" {
		if _, err := t.resolvePath(input.WorkDir); err != nil {
			return fmt.Errorf("invalid work directory: %w", err)
		}
	}

	// Validate arguments for dangerous patterns (skip for cd command)
	if input.Command != "cd" {
		for _, arg := range input.Arguments {
			if err := t.validateArgument(arg); err != nil {
				return fmt.Errorf("invalid argument '%s': %w", arg, err)
			}
		}
	}

	return nil
}

// isCommandAllowed checks if a command is in the allowed list
func (t *shellToolSet) isCommandAllowed(command string) bool {
	if len(t.allowedCommands) == 0 {
		// If no allowed commands specified, use default safe list
		return t.isDefaultSafeCommand(command)
	}

	for _, allowed := range t.allowedCommands {
		if command == allowed {
			return true
		}
	}
	return false
}

// isDefaultSafeCommand checks if a command is in the default safe list
func (t *shellToolSet) isDefaultSafeCommand(command string) bool {
	safeCommands := []string{
		"ls", "cat", "head", "tail", "grep", "find", "wc", "sort", "uniq", "cd",
		"echo", "pwd", "date", "whoami", "id", "uname", "df", "du", "ps",
		"git", "npm", "yarn", "go", "python", "python3", "node", "java",
		"make", "cmake", "curl", "wget", "ping", "nslookup", "dig",
	}

	for _, safe := range safeCommands {
		if command == safe {
			return true
		}
	}
	return false
}

// validateArgument checks if an argument contains dangerous patterns
func (t *shellToolSet) validateArgument(arg string) error {
	// Check for absolute paths and validate they are within workspace
	if filepath.IsAbs(arg) {
		if err := t.validateAbsolutePathWithinWorkspace(arg); err != nil {
			return err
		}
	}

	// Check for dangerous patterns
	dangerousPatterns := []string{
		`\$\(.*\)`,                 // Command substitution $(...)
		"`.*`",                     // Command substitution `...`
		`&&`,                       // Command chaining
		`\|\|`,                     // Command chaining
		`;`,                        // Command separator
		`\|`,                       // Pipe (can be dangerous in some contexts)
		`>`,                        // Redirection
		`<`,                        // Redirection
		`>>`,                       // Append redirection
		`&`,                        // Background execution
		`\$\{.*\}`,                 // Variable expansion ${...}
		`\$[A-Za-z_][A-Za-z0-9_]*`, // Environment variable expansion $VAR
	}

	for _, pattern := range dangerousPatterns {
		matched, err := regexp.MatchString(pattern, arg)
		if err != nil {
			return fmt.Errorf("regex error: %w", err)
		}
		if matched {
			return fmt.Errorf("argument contains dangerous pattern: %s", pattern)
		}
	}

	// Check for path traversal in arguments
	if strings.Contains(arg, "..") {
		return fmt.Errorf("path traversal detected")
	}

	return nil
}

// validateAbsolutePathWithinWorkspace checks if an absolute path is within the workspace
func (t *shellToolSet) validateAbsolutePathWithinWorkspace(absPath string) error {
	// Clean the path
	cleanPath := filepath.Clean(absPath)
	
	// Get absolute path of base directory
	absBaseDir, err := filepath.Abs(t.baseDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute base directory: %w", err)
	}
	
	// Check if the absolute path is within the base directory
	relPath, err := filepath.Rel(absBaseDir, cleanPath)
	if err != nil {
		return fmt.Errorf("absolute path outside workspace: %s", absPath)
	}
	
	// If the relative path starts with "..", it's outside the base directory
	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("absolute path outside workspace boundary: %s", absPath)
	}
	
	return nil
}

// executeCommand executes the shell command with the given input
func (t *shellToolSet) executeCommand(ctx context.Context, input ShellToolInput) (*ShellToolOutput, error) {
	// Special handling for cd command
	if input.Command == "cd" {
		return t.handleChangeDirectory(ctx, input)
	}
	// Handle context timeout
	if t.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t.timeout)
		defer cancel()
	}

	// Validate input first
	if err := t.validateInput(input); err != nil {
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Determine working directory
	workDir := t.currentWorkDir // Use current working directory instead of base directory
	if input.WorkDir != "" {
		resolvedWorkDir, err := t.resolvePath(input.WorkDir)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve work directory: %w", err)
		}
		workDir = resolvedWorkDir
	}

	// Ensure working directory exists and is within base path
	if err := t.validateWorkingDirectory(workDir); err != nil {
		return nil, fmt.Errorf("invalid working directory: %w", err)
	}

	// Create the command
	cmd := exec.CommandContext(ctx, input.Command, input.Arguments...)
	cmd.Dir = workDir

	// Set up environment restrictions
	cmd.Env = t.getRestrictedEnvironment()

	// Set up output capture
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set process group for better cleanup
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Execute the command
	err := cmd.Run()

	// Get exit code
	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	// Prepare result
	result := &ShellToolOutput{
		StdOut:   stdout.String(),
		StdError: stderr.String(),
		ExitCode: exitCode,
		WorkDir:  workDir,
	}

	// If command failed, include error information
	if err != nil {
		result.Error = err.Error()
		// Don't return error for non-zero exit codes, just include in result
		if exitCode != 0 {
			return result, nil
		}
		return result, fmt.Errorf("command execution failed: %w", err)
	}

	return result, nil
}

// validateWorkingDirectory ensures the working directory is safe
func (t *shellToolSet) validateWorkingDirectory(workDir string) error {
	// Check if directory exists
	stat, err := os.Stat(workDir)
	if err != nil {
		return fmt.Errorf("working directory does not exist: %w", err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("working directory is not a directory")
	}

	// Ensure it's within base path
	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absBaseDir, err := filepath.Abs(t.baseDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute base path: %w", err)
	}

	relPath, err := filepath.Rel(absBaseDir, absWorkDir)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("working directory is outside base path")
	}

	return nil
}

// getRestrictedEnvironment returns a restricted environment for command execution
func (t *shellToolSet) getRestrictedEnvironment() []string {
	// Start with a minimal environment
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"USER=" + os.Getenv("USER"),
		"LANG=" + os.Getenv("LANG"),
		"LC_ALL=" + os.Getenv("LC_ALL"),
	}

	// Add additional safe environment variables if needed
	safeVars := []string{"TERM", "SHELL", "PWD"}
	for _, varName := range safeVars {
		if value := os.Getenv(varName); value != "" {
			env = append(env, varName+"="+value)
		}
	}

	return env
}

// executeCommandTool returns a callable tool for executing shell commands.
func (f *shellToolSet) executeCommandTool() tool.CallableTool {
	return function.NewFunctionTool(
		f.executeCommand,
		function.WithName("execute_command"),
		function.WithDescription("Execute a shell command with optional arguments. Only allowed commands can be executed. Commands are executed within the configured base directory for security. Returns stdout, stderr, exit code, and error information."),
	)
}

// ChangeDirectoryInput represents the input structure for the cd command
type ChangeDirectoryInput struct {
	Path string `json:"path" description:"The directory path to change to (relative to current working directory)"`
}

// ChangeDirectoryOutput represents the output structure for the cd command
type ChangeDirectoryOutput struct {
	NewWorkDir string `json:"new_work_dir" description:"The new current working directory"`
	OldWorkDir string `json:"old_work_dir" description:"The previous working directory"`
	Error      string `json:"error,omitempty" description:"Error message if command failed"`
}

// handleChangeDirectory handles the cd command implementation
func (t *shellToolSet) handleChangeDirectory(_ context.Context, input ShellToolInput) (*ShellToolOutput, error) {
	// Basic validation (skip argument validation for cd)
	if input.Command == "" {
		return nil, fmt.Errorf("command cannot be empty")
	}

	if !t.isCommandAllowed(input.Command) {
		return nil, fmt.Errorf("command '%s' is not allowed", input.Command)
	}

	// Get target directory from arguments
	var targetDir string
	if len(input.Arguments) == 0 {
		// No arguments means go to base directory (our "home" in the sandbox)
		targetDir = "."
	} else if len(input.Arguments) == 1 {
		targetDir = input.Arguments[0]
	} else {
		return &ShellToolOutput{
			StdOut:   "",
			StdError: "cd: too many arguments",
			ExitCode: 1,
			Error:    "cd: too many arguments",
			WorkDir:  t.currentWorkDir,
		}, nil
	}

	// Store old working directory
	_ = t.currentWorkDir // oldWorkDir for potential future use

	// Resolve target directory relative to current working directory
	var newWorkDir string

	// Handle tilde (~) expansion - map to base directory (our sandbox "home")
	if targetDir == "~" {
		// Go directly to base directory
		newWorkDir = t.baseDir
	} else if strings.HasPrefix(targetDir, "~/") {
		// Replace ~ with base directory, then append the rest
		targetDir = targetDir[2:] // Remove "~/"
		newWorkDir = filepath.Join(t.baseDir, targetDir)
	}

	if newWorkDir == "" { // Only process if not already set by tilde expansion
		if targetDir == "." {
			// Stay in current directory
			newWorkDir = t.currentWorkDir
		} else if targetDir == ".." {
			// Go to parent directory
			newWorkDir = filepath.Dir(t.currentWorkDir)
		} else if filepath.IsAbs(targetDir) {
			// Absolute paths not allowed
			return &ShellToolOutput{
				StdOut:   "",
				StdError: "cd: absolute paths are not allowed",
				ExitCode: 1,
				Error:    "cd: absolute paths are not allowed",
				WorkDir:  t.currentWorkDir,
			}, nil
		} else {
			// Relative path - resolve from current working directory
			newWorkDir = filepath.Join(t.currentWorkDir, targetDir)
		}
	}

	// Clean the path
	newWorkDir = filepath.Clean(newWorkDir)

	// Check for path traversal attempts that would go outside base directory
	absNewWorkDir, err := filepath.Abs(newWorkDir)
	if err != nil {
		return &ShellToolOutput{
			StdOut:   "",
			StdError: fmt.Sprintf("cd: failed to resolve path: %v", err),
			ExitCode: 1,
			Error:    fmt.Sprintf("cd: failed to resolve path: %v", err),
			WorkDir:  t.currentWorkDir,
		}, nil
	}

	absBaseDir, err := filepath.Abs(t.baseDir)
	if err != nil {
		return &ShellToolOutput{
			StdOut:   "",
			StdError: fmt.Sprintf("cd: failed to resolve base directory: %v", err),
			ExitCode: 1,
			Error:    fmt.Sprintf("cd: failed to resolve base directory: %v", err),
			WorkDir:  t.currentWorkDir,
		}, nil
	}

	relPath, err := filepath.Rel(absBaseDir, absNewWorkDir)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return &ShellToolOutput{
			StdOut:   "",
			StdError: "cd: access denied - cannot navigate outside the allowed workspace boundary",
			ExitCode: 1,
			Error:    "cd: access denied - cannot navigate outside the allowed workspace boundary",
			WorkDir:  t.currentWorkDir,
		}, nil
	}

	newWorkDir = absNewWorkDir

	// Validate that the new directory is within base directory
	if err := t.validateWorkingDirectory(newWorkDir); err != nil {
		return &ShellToolOutput{
			StdOut:   "",
			StdError: fmt.Sprintf("cd: %v", err),
			ExitCode: 1,
			Error:    fmt.Sprintf("cd: %v", err),
			WorkDir:  t.currentWorkDir,
		}, nil
	}

	// Check if directory exists
	if _, err := os.Stat(newWorkDir); err != nil {
		return &ShellToolOutput{
			StdOut:   "",
			StdError: fmt.Sprintf("cd: %s: No such file or directory", targetDir),
			ExitCode: 1,
			Error:    fmt.Sprintf("cd: %s: No such file or directory", targetDir),
			WorkDir:  t.currentWorkDir,
		}, nil
	}

	// Update current working directory
	t.currentWorkDir = newWorkDir

	// Return success
	return &ShellToolOutput{
		StdOut:   "",
		StdError: "",
		ExitCode: 0,
		Error:    "",
		WorkDir:  t.currentWorkDir,
	}, nil
}

// changeDirectoryTool returns a callable tool for changing directories.
func (f *shellToolSet) changeDirectoryTool() tool.CallableTool {
	return function.NewFunctionTool(
		f.changeDirectory,
		function.WithName("change_directory"),
		function.WithDescription("Change the current working directory. Only directories within the base directory are allowed. Supports relative paths, '.', '..', and tilde (~) for home directory."),
	)
}

// changeDirectory is the wrapper function for the change directory tool
func (t *shellToolSet) changeDirectory(ctx context.Context, input ChangeDirectoryInput) (*ChangeDirectoryOutput, error) {
	// Store old working directory before making any changes
	oldWorkDirBeforeChange := t.currentWorkDir

	// Convert to ShellToolInput format
	shellInput := ShellToolInput{
		Command:   "cd",
		Arguments: []string{input.Path},
	}

	// Use the existing handleChangeDirectory function
	result, err := t.handleChangeDirectory(ctx, shellInput)
	if err != nil {
		return &ChangeDirectoryOutput{
			NewWorkDir: t.currentWorkDir,
			OldWorkDir: t.currentWorkDir,
			Error:      err.Error(),
		}, err
	}

	if result.ExitCode == 0 {
		// Success - get relative paths for user-friendly output
		relOldPath, _ := filepath.Rel(t.baseDir, oldWorkDirBeforeChange)
		relNewPath, _ := filepath.Rel(t.baseDir, result.WorkDir)

		return &ChangeDirectoryOutput{
			NewWorkDir: relNewPath,
			OldWorkDir: relOldPath,
			Error:      result.Error,
		}, nil
	} else {
		// Error - working directory didn't change
		relCurrentPath, _ := filepath.Rel(t.baseDir, oldWorkDirBeforeChange)
		return &ChangeDirectoryOutput{
			NewWorkDir: relCurrentPath,
			OldWorkDir: relCurrentPath,
			Error:      result.Error,
		}, nil
	}
}

package shell

import "time"

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

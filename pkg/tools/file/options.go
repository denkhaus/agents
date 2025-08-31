package file

// Option is a configuration function for the file ToolSet.
type Option func(*ToolSetConfig)

func WithReadOnly(readOnly bool) Option {
	return func(t *ToolSetConfig) {
		t.ReadOnly = readOnly
	}
}

func WithWorkspacePath(workspacePath string) Option {
	return func(t *ToolSetConfig) {
		t.WorkspacePath = workspacePath
	}
}

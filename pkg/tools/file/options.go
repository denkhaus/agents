package file

// Option is a configuration function for the file ToolSet.
type Option func(*FileToolSet)

func WithReadOnly(readOnly bool) Option {
	return func(t *FileToolSet) {
		t.ReadOnly = readOnly
	}
}

func WithWorkspacePath(workspacePath string) Option {
	return func(t *FileToolSet) {
		t.WorkspacePath = workspacePath
	}
}

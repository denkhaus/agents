package project

import "go.uber.org/zap"

// Option is a functional option for configuring the project task tool set
type Option func(*projectTaskToolSet)

// WithManager sets a custom manager instance
func WithManager(manager ProjectManager) Option {
	return func(pts *projectTaskToolSet) {
		pts.manager = manager
	}
}

// WithReadOnly allows to select tools that have no write access to project information
func WithReadOnly(readOnly bool) Option {
	return func(pts *projectTaskToolSet) {
		pts.isReadOnly = readOnly
	}
}

// WithRepository sets a custom repository instance
func WithRepository(repo Repository) Option {
	return func(pts *projectTaskToolSet) {
		// Create a new manager with the custom repository
		config := DefaultConfig()
		if pts.manager != nil {
			config = pts.manager.GetConfig()
		}
		pts.manager = NewManagerWithRepository(repo, config)
	}
}

// WithConfig sets a custom configuration
func WithConfig(config *Config) Option {
	return func(pts *projectTaskToolSet) {
		if pts.manager != nil {
			pts.manager.UpdateConfig(config)
		}
	}
}

// WithLogger sets a custom logger
// It' neccessary to switch the logger off while in chat
// since log messages interfere with chat output.
func WithLogger(logger *zap.Logger) Option {
	return func(pts *projectTaskToolSet) {
		pts.logger = logger
	}
}

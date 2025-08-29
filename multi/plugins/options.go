package plugins

import (
	"github.com/denkhaus/agents/multi"
)

// Options contains configuration for multi-agent chat plugins.
type Options struct {
	Processor        multi.ChatProcessor
	ProcessorOptions []multi.ChatProcessorOption
	DisplayWidth     int // Width for chat display borders (default: 120)
}

// MultiAgentChatOption is a function type for configuring multi-agent chat options.
type MultiAgentChatOption func(*Options)

// WithProcessorOptions sets the ChatProcessor options for the multi-agent chat.
func WithProcessorOptions(processorOptions ...multi.ChatProcessorOption) MultiAgentChatOption {
	return func(opts *Options) {
		opts.ProcessorOptions = processorOptions
	}
}

// WithDisplayWidth sets the width for chat display borders.
// If width is less than 40, it will be set to 40 (minimum usable width).
// Default width is 120 characters.
func WithDisplayWidth(width int) MultiAgentChatOption {
	return func(opts *Options) {
		if width < 40 {
			width = 40 // Minimum usable width
		}
		opts.DisplayWidth = width
	}
}

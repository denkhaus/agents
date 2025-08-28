package plugins

import (
	"github.com/denkhaus/agents/multi"
)

// Options contains configuration for multi-agent chat plugins.
type Options struct {
	processor        multi.ChatProcessor
	processorOptions []multi.ChatProcessorOption
}

// MultiAgentChatOption is a function type for configuring multi-agent chat options.
type MultiAgentChatOption func(*Options)

// WithProcessorOptions sets the ChatProcessor options for the multi-agent chat.
func WithProcessorOptions(processorOptions ...multi.ChatProcessorOption) MultiAgentChatOption {
	return func(opts *Options) {
		opts.processorOptions = processorOptions
	}
}

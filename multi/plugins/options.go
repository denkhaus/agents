package plugins

import (
	"github.com/denkhaus/agents/multi"
)

type Options struct {
	processor       multi.ChatProcessor
	applicationName string
}

type MultiAgentChatOption func(*Options)

// WithSessionID sets the SessionID to use.
func WithChatProcessor(processor multi.ChatProcessor) MultiAgentChatOption {
	return func(opts *Options) {
		opts.processor = processor
	}
}

package chat

import (
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

type Options struct {
	userID         string
	sessionID      string
	streaming      bool
	appName        string
	sessionService session.Service
	agent          agent.Agent
}

// Option is a function that configures an LLMAgent.
type Option func(*Options)

// WithModel sets the SessionID to use.
func WithSessionID(sessionID uuid.UUID) Option {
	return func(opts *Options) {
		opts.sessionID = sessionID.String()
	}
}

// WithModel sets the SessionID to use.
func WithUserID(userID uuid.UUID) Option {
	return func(opts *Options) {
		opts.userID = userID.String()
	}
}

// WithStreaming endables or disables streaming support.
func WithStreaming(streaming bool) Option {
	return func(opts *Options) {
		opts.streaming = streaming
	}
}

// WithAppName sets the AppName to use.
func WithAppName(appName string) Option {
	return func(opts *Options) {
		opts.appName = appName
	}
}

// WithSessionService sets the session service to use.
func WithSessionService(service session.Service) Option {
	return func(opts *Options) {
		opts.sessionService = service
	}
}

// WithSessionService sets the session service to use.
func WithAgent(agent agent.Agent) Option {
	return func(opts *Options) {
		opts.agent = agent
	}
}

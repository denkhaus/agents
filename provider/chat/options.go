package chat

import (
	"github.com/denkhaus/agents/provider"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// WithModel sets the SessionID to use.
func WithSessionID(sessionID uuid.UUID) provider.ChatProviderOption {
	return func(opts *provider.ChatProviderOptions) {
		opts.SessionID = sessionID.String()
	}
}

// WithModel sets the SessionID to use.
func WithUserID(userID uuid.UUID) provider.ChatProviderOption {
	return func(opts *provider.ChatProviderOptions) {
		opts.UserID = userID.String()
	}
}

// WithAppName sets the AppName to use.
func WithAppName(appName string) provider.ChatProviderOption {
	return func(opts *provider.ChatProviderOptions) {
		opts.AppName = appName
	}
}

// WithSessionService sets the session service to use.
func WithSessionService(service session.Service) provider.ChatProviderOption {
	return func(opts *provider.ChatProviderOptions) {
		opts.SessionService = service
	}
}

func WithAgentProviderOptions(opt ...provider.AgentProviderOption) provider.ChatProviderOption {
	return func(opts *provider.ChatProviderOptions) {
		opts.AgentProviderOptions = opt
	}
}

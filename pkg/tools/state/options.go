package state

import "trpc.group/trpc-go/trpc-agent-go/session"

// Option defines a function type for configuring StateAccessTool.
type Option func(*StateAccessTool)

// WithSessionService sets the session service for the StateAccessTool.
func WithSessionService(service session.Service) Option {
	return func(t *StateAccessTool) {
		t.sessionService = service
	}
}

// WithAppName sets the application name for the StateAccessTool.
func WithAppName(appName string) Option {
	return func(t *StateAccessTool) {
		t.appName = appName
	}
}

// WithUserID sets the user ID for the StateAccessTool.
func WithUserID(userID string) Option {
	return func(t *StateAccessTool) {
		t.userID = userID
	}
}

// WithSessionID sets the session ID for the StateAccessTool.
func WithSessionID(sessionID string) Option {
	return func(t *StateAccessTool) {
		t.sessionID = sessionID
	}
}

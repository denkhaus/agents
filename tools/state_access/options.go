package stateaccess

import "trpc.group/trpc-go/trpc-agent-go/session"

// StateAccessToolOption defines a function type for configuring StateAccessTool.
type StateAccessToolOption func(*StateAccessTool)

// WithSessionService sets the session service for the StateAccessTool.
func WithSessionService(service session.Service) StateAccessToolOption {
	return func(t *StateAccessTool) {
		t.sessionService = service
	}
}

// WithAppName sets the application name for the StateAccessTool.
func WithAppName(appName string) StateAccessToolOption {
	return func(t *StateAccessTool) {
		t.appName = appName
	}
}

// WithUserID sets the user ID for the StateAccessTool.
func WithUserID(userID string) StateAccessToolOption {
	return func(t *StateAccessTool) {
		t.userID = userID
	}
}

// WithSessionID sets the session ID for the StateAccessTool.
func WithSessionID(sessionID string) StateAccessToolOption {
	return func(t *StateAccessTool) {
		t.sessionID = sessionID
	}
}

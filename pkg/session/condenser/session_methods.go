package condenser

import (
	"context"

	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// CreateSession creates a new session
func (s *Service) CreateSession(ctx context.Context, key session.Key, state session.StateMap, options ...session.Option) (*session.Session, error) {
	return s.sessionService.CreateSession(ctx, key, state, options...)
}

// GetSession gets a session
func (s *Service) GetSession(ctx context.Context, key session.Key, options ...session.Option) (*session.Session, error) {
	return s.sessionService.GetSession(ctx, key, options...)
}

// ListSessions lists all sessions by user scope of session key
func (s *Service) ListSessions(ctx context.Context, userKey session.UserKey, options ...session.Option) ([]*session.Session, error) {
	return s.sessionService.ListSessions(ctx, userKey, options...)
}

// DeleteSession deletes a session
func (s *Service) DeleteSession(ctx context.Context, key session.Key, options ...session.Option) error {
	return s.sessionService.DeleteSession(ctx, key, options...)
}

// UpdateAppState updates the state by target scope and key
func (s *Service) UpdateAppState(ctx context.Context, appName string, state session.StateMap) error {
	return s.sessionService.UpdateAppState(ctx, appName, state)
}

// DeleteAppState deletes the state by target scope and key
func (s *Service) DeleteAppState(ctx context.Context, appName string, key string) error {
	return s.sessionService.DeleteAppState(ctx, appName, key)
}

// ListAppStates lists app states
func (s *Service) ListAppStates(ctx context.Context, appName string) (session.StateMap, error) {
	return s.sessionService.ListAppStates(ctx, appName)
}

// UpdateUserState updates the state by target scope and key
func (s *Service) UpdateUserState(ctx context.Context, userKey session.UserKey, state session.StateMap) error {
	return s.sessionService.UpdateUserState(ctx, userKey, state)
}

// ListUserStates lists user states
func (s *Service) ListUserStates(ctx context.Context, userKey session.UserKey) (session.StateMap, error) {
	return s.sessionService.ListUserStates(ctx, userKey)
}

// DeleteUserState deletes the state by target scope and key
func (s *Service) DeleteUserState(ctx context.Context, userKey session.UserKey, key string) error {
	return s.sessionService.DeleteUserState(ctx, userKey, key)
}

// AppendEvent appends an event to a session and triggers condensation if needed
func (s *Service) AppendEvent(ctx context.Context, sess *session.Session, evt *event.Event, options ...session.Option) error {
	// First, append the event to the underlying session service
	if err := s.sessionService.AppendEvent(ctx, sess, evt, options...); err != nil {
		return err
	}

	// Check if condensation is needed
	return s.checkAndCondense(ctx, sess)
}
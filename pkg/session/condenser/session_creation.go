package condenser

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// createCondensedSession creates a new session with the condensed content
func (s *Service) createCondensedSession(ctx context.Context, sessKey session.Key, summary string, recentEvents []event.Event, originalSess *session.Session) error {
	// Create condensed state with the summary
	condensedState := session.StateMap{
		"condensed_summary": []byte(summary),
		"condensed_at":      []byte(time.Now().Format(time.RFC3339)),
	}

	// Copy any existing state that should be preserved
	for key, value := range originalSess.State {
		// Skip the condensed summary if it already exists
		if key != "condensed_summary" && key != "condensed_at" {
			condensedState[key] = value
		}
	}

	// Delete the original session first
	if err := s.sessionService.DeleteSession(ctx, sessKey); err != nil {
		s.logger.Error("Failed to delete original session",
			zap.String("sessionID", sessKey.SessionID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete original session: %w", err)
	}

	// Create the new condensed session
	newSess, err := s.sessionService.CreateSession(ctx, sessKey, condensedState)
	if err != nil {
		s.logger.Error("Failed to create condensed session",
			zap.String("sessionID", sessKey.SessionID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create condensed session: %w", err)
	}

	// Append the recent events to the new session
	for i, evt := range recentEvents {
		if err := s.sessionService.AppendEvent(ctx, newSess, &evt); err != nil {
			s.logger.Error("Failed to append recent event to condensed session",
				zap.String("sessionID", newSess.ID),
				zap.Int("eventIndex", i),
				zap.Error(err),
			)
			// Continue with other events rather than failing completely
		}
	}

	s.logger.Debug("Created condensed session",
		zap.String("sessionID", newSess.ID),
		zap.Int("stateKeys", len(condensedState)),
		zap.Int("recentEvents", len(recentEvents)),
	)

	return nil
}

// validateCondensedSession verifies that the condensed session is smaller than the original
func (s *Service) validateCondensedSession(ctx context.Context, sessKey session.Key, originalTokens int) error {
	// Get the newly created session
	condensedSess, err := s.sessionService.GetSession(ctx, sessKey)
	if err != nil {
		return fmt.Errorf("failed to get condensed session for validation: %w", err)
	}

	// Calculate the new token count
	condensedTokens, err := s.calculateSessionTokens(ctx, condensedSess)
	if err != nil {
		return fmt.Errorf("failed to calculate condensed session tokens: %w", err)
	}

	// Log the results
	s.logger.Info("Condensation validation",
		zap.String("sessionID", sessKey.SessionID),
		zap.Int("originalTokens", originalTokens),
		zap.Int("condensedTokens", condensedTokens),
		zap.Int("tokenSavings", originalTokens-condensedTokens),
	)

	// Warn if condensation didn't reduce size significantly
	if condensedTokens >= originalTokens {
		s.logger.Warn("Condensed session is not smaller than original",
			zap.String("sessionID", sessKey.SessionID),
			zap.Int("originalTokens", originalTokens),
			zap.Int("condensedTokens", condensedTokens),
		)
	}

	return nil
}

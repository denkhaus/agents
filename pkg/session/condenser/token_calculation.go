package condenser

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// calculateSessionTokens calculates the total token count for a session
func (s *Service) calculateSessionTokens(ctx context.Context, sess *session.Session) (int, error) {
	totalTokens := 0

	// Count tokens in state
	for key, value := range sess.State {
		keyTokens, err := s.tokenCounter.CountTokens(ctx, key)
		if err != nil {
			return 0, fmt.Errorf("failed to count tokens for state key '%s': %w", key, err)
		}

		valueTokens, err := s.tokenCounter.EstimateTokens(ctx, value)
		if err != nil {
			return 0, fmt.Errorf("failed to count tokens for state value: %w", err)
		}

		totalTokens += keyTokens + valueTokens
	}

	// Count tokens in events
	for i, evt := range sess.Events {
		eventTokens, err := s.calculateEventTokens(ctx, &evt)
		if err != nil {
			return 0, fmt.Errorf("failed to count tokens for event %d: %w", i, err)
		}
		totalTokens += eventTokens
	}

	counterInfo := s.tokenCounter.GetInfo()
	s.logger.Debug("Calculated session tokens",
		zap.String("sessionID", sess.ID),
		zap.Int("totalTokens", totalTokens),
		zap.String("method", counterInfo.Method.String()),
		zap.String("accuracy", counterInfo.Accuracy.String()),
		zap.Int("stateKeys", len(sess.State)),
		zap.Int("events", len(sess.Events)),
	)

	return totalTokens, nil
}

// calculateEventTokens calculates the token count for a single event
func (s *Service) calculateEventTokens(ctx context.Context, evt *event.Event) (int, error) {
	tokens := 0

	// Count request tokens - events may not have direct request access
	// For now, we'll focus on response content which is more commonly available

	// Count response tokens
	if evt.Response != nil {
		for _, choice := range evt.Response.Choices {
			contentTokens, err := s.tokenCounter.CountTokens(ctx, choice.Message.Content)
			if err != nil {
				return 0, fmt.Errorf("failed to count tokens for response content: %w", err)
			}
			tokens += contentTokens
		}
	}

	return tokens, nil
}

// checkAndCondense checks if condensation is needed and performs it
func (s *Service) checkAndCondense(ctx context.Context, sess *session.Session) error {
	// Get updated session to ensure we have the latest state
	sessKey := session.Key{
		AppName:   sess.AppName,
		UserID:    sess.UserID,
		SessionID: sess.ID,
	}

	updatedSess, err := s.sessionService.GetSession(ctx, sessKey)
	if err != nil {
		return fmt.Errorf("failed to get updated session: %w", err)
	}

	// Calculate current token usage
	currentTokens, err := s.calculateSessionTokens(ctx, updatedSess)
	if err != nil {
		s.logger.Error("Failed to calculate session tokens", 
			zap.String("sessionID", sess.ID),
			zap.Error(err))
		return nil // Don't fail the append operation
	}

	thresholdTokens := int(float64(s.config.MaxContextTokens) * s.config.TriggerThreshold)

	if currentTokens >= thresholdTokens {
		s.logger.Info("Token threshold exceeded, starting condensation",
			zap.String("sessionID", sess.ID),
			zap.Int("currentTokens", currentTokens),
			zap.Int("threshold", thresholdTokens),
			zap.Float64("thresholdPercentage", s.config.TriggerThreshold*100),
		)

		if err := s.condenseSession(ctx, updatedSess); err != nil {
			s.logger.Error("Condensation failed",
				zap.String("sessionID", sess.ID),
				zap.Error(err),
			)
			s.metrics.FailureCount++
			// Don't return error - allow session to continue
		}
	} else {
		s.logger.Debug("Session tokens below threshold",
			zap.String("sessionID", sess.ID),
			zap.Int("currentTokens", currentTokens),
			zap.Int("threshold", thresholdTokens),
		)
	}

	return nil
}
package condenser

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// condenseSession performs the actual condensation of a session
func (s *Service) condenseSession(ctx context.Context, sess *session.Session) error {
	startTime := time.Now()
	sessKey := session.Key{
		AppName:   sess.AppName,
		UserID:    sess.UserID,
		SessionID: sess.ID,
	}

	// Calculate original token count
	originalTokens, err := s.calculateSessionTokens(ctx, sess)
	if err != nil {
		return fmt.Errorf("failed to calculate original session tokens: %w", err)
	}

	// Generate summary
	summary, summaryTokens, err := s.generateSummary(ctx, sess)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	// Select recent events to keep
	recentEvents := s.selectRecentEvents(sess.Events)
	recentEventTokens, err := s.calculateEventsTokens(ctx, recentEvents)
	if err != nil {
		return fmt.Errorf("failed to calculate recent event tokens: %w", err)
	}

	// Calculate projected token count after condensation
	projectedTokens := summaryTokens + recentEventTokens
	tokenSavings := originalTokens - projectedTokens

	s.logger.Info("Condensation analysis",
		zap.String("sessionID", sess.ID),
		zap.Int("originalTokens", originalTokens),
		zap.Int("summaryTokens", summaryTokens),
		zap.Int("recentEventTokens", recentEventTokens),
		zap.Int("projectedTokens", projectedTokens),
		zap.Int("tokenSavings", tokenSavings),
		zap.Int("recentEventsKept", len(recentEvents)),
	)

	// Create condensed session
	if err := s.createCondensedSession(ctx, sessKey, summary, recentEvents, sess); err != nil {
		return fmt.Errorf("failed to create condensed session: %w", err)
	}

	// Validate that the condensed session is actually smaller
	if err := s.validateCondensedSession(ctx, sessKey, originalTokens); err != nil {
		s.logger.Warn("Condensed session validation failed", 
			zap.String("sessionID", sess.ID),
			zap.Error(err),
		)
		// Continue anyway - validation failure is not fatal
	}

	// Update metrics
	s.updateMetrics(originalTokens, projectedTokens, startTime)

	s.logger.Info("Session condensed successfully",
		zap.String("sessionID", sess.ID),
		zap.Duration("duration", time.Since(startTime)),
		zap.Int("tokenSavings", tokenSavings),
	)

	return nil
}

// generateSummary creates a summary of the session using the LLM
func (s *Service) generateSummary(ctx context.Context, sess *session.Session) (string, int, error) {
	// Build conversation history
	conversationHistory := s.buildConversationHistory(sess)
	if conversationHistory == "" {
		return "", 0, fmt.Errorf("no conversation history to summarize")
	}

	// Prepare messages for summarization
	messages := []model.Message{
		model.NewSystemMessage(s.config.SummaryPrompt),
		model.NewUserMessage(conversationHistory),
	}

	// Generate summary using the LLM
	respChan, err := s.summarizerLLM.GenerateContent(ctx, &model.Request{Messages: messages})
	if err != nil {
		return "", 0, fmt.Errorf("failed to start summary generation: %w", err)
	}

	var finalResp *model.Response
	for r := range respChan {
		if r.Error != nil {
			return "", 0, fmt.Errorf("LLM error during summary generation: %s", r.Error.Message)
		}
		finalResp = r
	}

	if finalResp == nil || len(finalResp.Choices) == 0 || finalResp.Choices[0].Message.Content == "" {
		return "", 0, fmt.Errorf("generated an empty summary")
	}

	summary := finalResp.Choices[0].Message.Content

	// Count tokens in the summary
	summaryTokens, err := s.tokenCounter.CountTokens(ctx, summary)
	if err != nil {
		return "", 0, fmt.Errorf("failed to count summary tokens: %w", err)
	}

	return summary, summaryTokens, nil
}

// buildConversationHistory constructs a text representation of the conversation
func (s *Service) buildConversationHistory(sess *session.Session) string {
	var history string

	for _, evt := range sess.Events {
		if evt.Response != nil && len(evt.Response.Choices) > 0 {
			choice := evt.Response.Choices[0]
			if choice.Message.Content != "" {
				switch choice.Message.Role {
				case model.RoleUser:
					history += fmt.Sprintf("User: %s\n", choice.Message.Content)
				case model.RoleAssistant:
					history += fmt.Sprintf("Assistant: %s\n", choice.Message.Content)
				case model.RoleSystem:
					history += fmt.Sprintf("System: %s\n", choice.Message.Content)
				}
			}
		}
	}

	return history
}

// selectRecentEvents selects the most recent events to keep after condensation
func (s *Service) selectRecentEvents(events []event.Event) []event.Event {
	if s.config.RecentEventsToKeep <= 0 || len(events) <= s.config.RecentEventsToKeep {
		return events
	}

	// Return the last N events
	startIndex := len(events) - s.config.RecentEventsToKeep
	return events[startIndex:]
}

// calculateEventsTokens calculates the total token count for a slice of events
func (s *Service) calculateEventsTokens(ctx context.Context, events []event.Event) (int, error) {
	totalTokens := 0

	for i, evt := range events {
		eventTokens, err := s.calculateEventTokens(ctx, &evt)
		if err != nil {
			return 0, fmt.Errorf("failed to calculate tokens for event %d: %w", i, err)
		}
		totalTokens += eventTokens
	}

	return totalTokens, nil
}

// updateMetrics updates the condensation metrics
func (s *Service) updateMetrics(originalTokens, projectedTokens int, startTime time.Time) {
	s.metrics.CondensationCount++
	s.metrics.LastCondensationTime = startTime

	if originalTokens > 0 {
		tokenSavings := originalTokens - projectedTokens
		s.metrics.TotalTokensSaved += int64(tokenSavings)

		// Calculate average reduction ratio
		reductionRatio := float64(tokenSavings) / float64(originalTokens)
		s.metrics.AverageReductionRatio = (s.metrics.AverageReductionRatio*float64(s.metrics.CondensationCount-1) + reductionRatio) / float64(s.metrics.CondensationCount)
	}
}
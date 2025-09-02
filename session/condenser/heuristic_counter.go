package condenser

import (
	"context"
	"fmt"
	"math"

	"go.uber.org/zap"
)

// HeuristicTokenCounter uses character-based estimation for token counting
type HeuristicTokenCounter struct {
	charsPerToken float64
	logger        *zap.Logger
}

// NewHeuristicTokenCounter creates a new heuristic token counter
func NewHeuristicTokenCounter(charsPerToken float64, logger *zap.Logger) *HeuristicTokenCounter {
	if charsPerToken <= 0 {
		charsPerToken = 4.0 // Default GPT-like ratio
	}
	return &HeuristicTokenCounter{
		charsPerToken: charsPerToken,
		logger:        logger.Named("heuristic-counter"),
	}
}

// CountTokens estimates token count based on character count
func (h *HeuristicTokenCounter) CountTokens(ctx context.Context, text string) (int, error) {
	if text == "" {
		return 0, nil
	}

	// Use simple character count for estimation
	charCount := len(text)
	tokens := int(math.Ceil(float64(charCount) / h.charsPerToken))

	h.logger.Debug("Counted tokens heuristically",
		zap.Int("chars", charCount),
		zap.Int("tokens", tokens),
		zap.Float64("ratio", h.charsPerToken),
	)

	return tokens, nil
}

// EstimateTokens estimates token count for byte content
func (h *HeuristicTokenCounter) EstimateTokens(ctx context.Context, content []byte) (int, error) {
	return h.CountTokens(ctx, string(content))
}

// GetInfo returns information about this token counter
func (h *HeuristicTokenCounter) GetInfo() TokenCounterInfo {
	return TokenCounterInfo{
		Method:      TokenCountingHeuristic,
		Accuracy:    AccuracyEstimated,
		Description: fmt.Sprintf("Character-based estimation (%.1f chars/token)", h.charsPerToken),
	}
}
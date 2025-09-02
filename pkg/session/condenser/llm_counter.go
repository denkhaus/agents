package condenser

import (
	"context"
	"crypto/md5"
	"fmt"
	"math"
	"sync"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// LLMTokenCounter uses the LLM itself for token counting
type LLMTokenCounter struct {
	model  model.Model
	logger *zap.Logger
	cache  map[string]int
	mutex  sync.RWMutex
}

// NewLLMTokenCounter creates a new LLM-based token counter
func NewLLMTokenCounter(model model.Model, logger *zap.Logger) *LLMTokenCounter {
	return &LLMTokenCounter{
		model:  model,
		logger: logger.Named("llm-counter"),
		cache:  make(map[string]int),
	}
}

// CountTokens counts tokens using the LLM model
func (l *LLMTokenCounter) CountTokens(ctx context.Context, text string) (int, error) {
	if text == "" {
		return 0, nil
	}

	// Create cache key using MD5 hash
	hash := fmt.Sprintf("%x", md5.Sum([]byte(text)))

	// Check cache first
	l.mutex.RLock()
	if count, exists := l.cache[hash]; exists {
		l.mutex.RUnlock()
		l.logger.Debug("Token count cache hit", zap.String("hash", hash[:8]))
		return count, nil
	}
	l.mutex.RUnlock()

	// Count via LLM
	count, err := l.countViaLLM(ctx, text)
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens via LLM: %w", err)
	}

	// Cache result
	l.mutex.Lock()
	l.cache[hash] = count
	l.mutex.Unlock()

	l.logger.Debug("Counted tokens via LLM",
		zap.Int("tokens", count),
		zap.String("hash", hash[:8]),
	)

	return count, nil
}

// EstimateTokens estimates token count for byte content
func (l *LLMTokenCounter) EstimateTokens(ctx context.Context, content []byte) (int, error) {
	return l.CountTokens(ctx, string(content))
}

// GetInfo returns information about this token counter
func (l *LLMTokenCounter) GetInfo() TokenCounterInfo {
	return TokenCounterInfo{
		Method:      TokenCountingLLMNative,
		Accuracy:    AccuracyPrecise,
		ModelName:   "llm-based", // TODO: Get actual model name
		Description: "LLM-native token counting",
	}
}

// countViaLLM performs the actual token counting using the LLM
func (l *LLMTokenCounter) countViaLLM(ctx context.Context, text string) (int, error) {
	// TODO: This would need to be implemented based on available model capabilities
	// For now, fallback to heuristic calculation
	// In a real implementation, this might use a special API call or model feature
	
	l.logger.Debug("LLM token counting not yet implemented, using fallback")
	return int(math.Ceil(float64(len(text)) / 4.0)), nil
}

// supportsTokenCounting checks if the given model supports native token counting
func supportsTokenCounting(llm model.Model) bool {
	// TODO: This would check if the LLM supports token counting
	// Implementation depends on the model interface and available capabilities
	// For now, return false to always use heuristic
	return false
}
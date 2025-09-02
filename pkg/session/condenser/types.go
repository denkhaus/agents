package condenser

import (
	"context"
	"time"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// TokenCountingMethod defines the method used for counting tokens
type TokenCountingMethod int

const (
	TokenCountingAuto TokenCountingMethod = iota
	TokenCountingHeuristic
	TokenCountingLLMNative
	TokenCountingTikToken
	TokenCountingCustom
)

func (t TokenCountingMethod) String() string {
	switch t {
	case TokenCountingAuto:
		return "auto"
	case TokenCountingHeuristic:
		return "heuristic"
	case TokenCountingLLMNative:
		return "llm-native"
	case TokenCountingTikToken:
		return "tiktoken"
	case TokenCountingCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// AccuracyLevel indicates the precision of token counting
type AccuracyLevel int

const (
	AccuracyEstimated AccuracyLevel = iota
	AccuracyPrecise
	AccuracyExact
)

func (a AccuracyLevel) String() string {
	switch a {
	case AccuracyEstimated:
		return "estimated"
	case AccuracyPrecise:
		return "precise"
	case AccuracyExact:
		return "exact"
	default:
		return "unknown"
	}
}

// TokenCounter interface for counting tokens in text
type TokenCounter interface {
	CountTokens(ctx context.Context, text string) (int, error)
	EstimateTokens(ctx context.Context, content []byte) (int, error)
	GetInfo() TokenCounterInfo
}

// TokenCounterInfo provides metadata about the token counter
type TokenCounterInfo struct {
	Method      TokenCountingMethod
	Accuracy    AccuracyLevel
	ModelName   string
	Description string
}

// Config holds all configuration for the condenser
type Config struct {
	// Token limits
	MaxContextTokens int     // Maximum tokens before condensation
	TriggerThreshold float64 // Percentage (0.0-1.0) to trigger condensation

	// Condensation behavior
	RecentEventsToKeep int    // Number of recent events to preserve
	SummaryPrompt      string // Custom prompt for summarization

	// Token counting
	TokenCountingMethod TokenCountingMethod
	CharsPerToken       float64 // For heuristic method

	// Performance
	EnableTokenCaching bool
	CacheSize          int

	// Logging
	Logger *zap.Logger
}

// Metrics tracks condensation performance
type Metrics struct {
	CondensationCount     int64
	TotalTokensSaved      int64
	AverageReductionRatio float64
	LastCondensationTime  time.Time
	FailureCount          int64
}

// Service wraps a session service with intelligent condensation
type Service struct {
	sessionService session.Service
	summarizerLLM  model.Model
	tokenCounter   TokenCounter
	config         Config
	logger         *zap.Logger
	metrics        *Metrics
}

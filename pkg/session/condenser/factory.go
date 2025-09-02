package condenser

import (
	"fmt"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// createTokenCounter creates a token counter based on the configuration
func createTokenCounter(llm model.Model, config Config, logger *zap.Logger) (TokenCounter, error) {
	switch config.TokenCountingMethod {
	case TokenCountingHeuristic:
		logger.Info("Creating heuristic token counter",
			zap.Float64("charsPerToken", config.CharsPerToken))
		return NewHeuristicTokenCounter(config.CharsPerToken, logger), nil

	case TokenCountingLLMNative:
		logger.Info("Creating LLM-native token counter")
		return NewLLMTokenCounter(llm, logger), nil

	case TokenCountingAuto:
		// Try LLM-native first, fallback to heuristic
		if supportsTokenCounting(llm) {
			logger.Info("Auto-selected LLM-native token counting")
			return NewLLMTokenCounter(llm, logger), nil
		}
		logger.Info("LLM doesn't support native token counting, using heuristic",
			zap.Float64("charsPerToken", config.CharsPerToken))
		return NewHeuristicTokenCounter(config.CharsPerToken, logger), nil

	case TokenCountingTikToken:
		// TODO: Implement tiktoken support
		return nil, fmt.Errorf("tiktoken support not yet implemented")

	case TokenCountingCustom:
		return nil, fmt.Errorf("custom token counter must be provided via WithCustomTokenCounter")

	default:
		return nil, fmt.Errorf("unsupported token counting method: %v", config.TokenCountingMethod)
	}
}

// WithCustomTokenCounter allows providing a custom token counter implementation
func WithCustomTokenCounter(counter TokenCounter) ConfigOption {
	return func(c *Config) {
		c.TokenCountingMethod = TokenCountingCustom
		// Store the custom counter in a way that can be retrieved
		// This would need to be handled in the service creation
	}
}
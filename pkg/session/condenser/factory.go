package condenser

import (
	"fmt"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// createTokenCounter creates a token counter based on the configuration
func createTokenCounter(llm model.Model, config Config) (TokenCounter, error) {
	switch config.TokenCountingMethod {
	case TokenCountingHeuristic:
		config.Logger.Info("Creating heuristic token counter",
			zap.Float64("charsPerToken", config.CharsPerToken))
		return NewHeuristicTokenCounter(config.CharsPerToken, config.Logger), nil

	case TokenCountingLLMNative:
		config.Logger.Info("Creating LLM-native token counter")
		return NewLLMTokenCounter(llm, config.Logger), nil

	case TokenCountingAuto:
		// Try LLM-native first, fallback to heuristic
		if supportsTokenCounting(llm) {
			config.Logger.Info("Auto-selected LLM-native token counting")
			return NewLLMTokenCounter(llm, config.Logger), nil
		}
		config.Logger.Info("LLM doesn't support native token counting, using heuristic",
			zap.Float64("charsPerToken", config.CharsPerToken))
		return NewHeuristicTokenCounter(config.CharsPerToken, config.Logger), nil

	case TokenCountingTikToken:
		config.Logger.Info("Creating tiktoken-based token counter")
		return NewTikTokenCounter(config.Logger), nil

	case TokenCountingCustom:
		return nil, fmt.Errorf("custom token counter must be provided via WithCustomTokenCounter")

	default:
		return nil, fmt.Errorf("unsupported token counting method: %v", config.TokenCountingMethod)
	}
}

package condenser

import "go.uber.org/zap"

// DefaultConfig returns sensible defaults for the condenser service
func DefaultConfig() Config {
	return Config{
		MaxContextTokens:    8000,
		TriggerThreshold:    0.75,
		RecentEventsToKeep:  3,
		SummaryPrompt:       "Summarize this conversation concisely, preserving key context and decisions:",
		TokenCountingMethod: TokenCountingAuto,
		CharsPerToken:       4.0,
		EnableTokenCaching:  true,
		CacheSize:           1000,
		Logger:              zap.NewNop(),
	}
}

// ConfigOption allows functional configuration of the condenser
type ConfigOption func(*Config)

// WithMaxContextTokens sets the maximum context size in tokens
func WithMaxContextTokens(tokens int) ConfigOption {
	return func(c *Config) {
		c.MaxContextTokens = tokens
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

// WithLogger allows providing a custom logger for debugging reasons
func WithLogger(logger *zap.Logger) ConfigOption {
	return func(c *Config) {
		c.Logger = logger
	}
}

// WithTriggerThreshold sets the trigger threshold for condensation (0.0-1.0)
func WithTriggerThreshold(threshold float64) ConfigOption {
	return func(c *Config) {
		c.TriggerThreshold = threshold
	}
}

// WithRecentEventsToKeep sets the number of recent events to preserve
func WithRecentEventsToKeep(count int) ConfigOption {
	return func(c *Config) {
		c.RecentEventsToKeep = count
	}
}

// WithSummaryPrompt sets the custom prompt for summarization
func WithSummaryPrompt(prompt string) ConfigOption {
	return func(c *Config) {
		c.SummaryPrompt = prompt
	}
}

// WithTokenCountingMethod sets the token counting method
func WithTokenCountingMethod(method TokenCountingMethod) ConfigOption {
	return func(c *Config) {
		c.TokenCountingMethod = method
	}
}

// WithCharsPerToken sets the characters per token ratio for heuristic counting
func WithCharsPerToken(ratio float64) ConfigOption {
	return func(c *Config) {
		c.CharsPerToken = ratio
	}
}

// WithTokenCaching enables or disables token count caching
func WithTokenCaching(enabled bool) ConfigOption {
	return func(c *Config) {
		c.EnableTokenCaching = enabled
	}
}

// WithCacheSize sets the maximum cache size for token counts
func WithCacheSize(size int) ConfigOption {
	return func(c *Config) {
		c.CacheSize = size
	}
}

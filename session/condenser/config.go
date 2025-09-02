package condenser

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
		CacheSize:          1000,
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
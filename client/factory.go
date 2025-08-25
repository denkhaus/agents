package client

import (
	"fmt"
	"time"

	"github.com/samber/mo"
)

// NewUnifiedClient creates a new unified A2A client.
// Returns an error if the configuration is invalid.
func NewUnifiedClient(config ClientConfig) (UnifiedClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	timeout := config.Timeout.OrElse(30 * time.Second)
	userAgent := config.UserAgent.OrElse("unified-a2a-client/1.0")
	streaming := config.Streaming.OrElse(true)
	retryCount := config.RetryCount.OrElse(3)

	client := &unifiedClient{
		baseURL:    config.BaseURL,
		timeout:    timeout,
		userAgent:  userAgent,
		streaming:  streaming,
		retryCount: retryCount,
		auth:       config.Auth,
	}

	return client, nil
}

// NewAPIClient creates a basic API client.
// Returns an error if the configuration is invalid.
func NewAPIClient(baseURL string, opts ...ClientOption) (APIClient, error) {
	config := ClientConfig{BaseURL: baseURL}

	for _, opt := range opts {
		opt(&config)
	}

	return NewUnifiedClient(config)
}

// NewStreamingClient creates a streaming client.
// Returns an error if the configuration is invalid.
func NewStreamingClient(baseURL string, opts ...ClientOption) (StreamingClient, error) {
	config := ClientConfig{
		BaseURL:   baseURL,
		Streaming: mo.Some(true),
	}

	for _, opt := range opts {
		opt(&config)
	}

	return NewUnifiedClient(config)
}

// NewAgentCardFetcher creates an agent card fetcher.
// Returns an error if the timeout is invalid.
func NewAgentCardFetcher(timeout time.Duration) (AgentCardFetcher, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}
	
	return newAgentCardFetcher(timeout), nil
}

// NewTaskTracker creates a new task tracker.
func NewTaskTracker() TaskTracker {
	return newTaskTracker()
}

// ClientOption defines configuration options.
type ClientOption func(*ClientConfig)

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.Timeout = mo.Some(timeout)
	}
}

// WithAuth sets the authenticator.
func WithAuth(auth Authenticator) ClientOption {
	return func(c *ClientConfig) {
		c.Auth = mo.Some(auth)
	}
}

// WithUserAgent sets the user agent.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *ClientConfig) {
		c.UserAgent = mo.Some(userAgent)
	}
}

// WithStreaming enables/disables streaming.
func WithStreaming(enabled bool) ClientOption {
	return func(c *ClientConfig) {
		c.Streaming = mo.Some(enabled)
	}
}

// WithRetryCount sets the retry count.
func WithRetryCount(count int) ClientOption {
	return func(c *ClientConfig) {
		c.RetryCount = mo.Some(count)
	}
}

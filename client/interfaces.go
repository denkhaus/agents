// Package client provides a unified interface for A2A client functionality.
//
// This package defines client interfaces and provides factory functions for creating
// A2A clients. All concrete implementations are unexported and accessed through
// well-defined interfaces.
//
// Key interfaces:
//   - APIClient: Core A2A communication functionality
//   - StreamingClient: Streaming message support
//   - AgentCardFetcher: Agent capability discovery
//   - JWKSVerifier: JWT verification with JWKS
//   - TaskTracker: Task status tracking
//   - UnifiedClient: Combines all client functionality
//
// Example usage:
//
//	// Create a unified client
//	client, err := client.NewUnifiedClient(client.ClientConfig{
//		BaseURL: "https://api.example.com",
//		Timeout: mo.Some(30 * time.Second),
//		Auth:    mo.Some(authProvider),
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Send a message
//	result, err := client.SendMessage(ctx, params)
//	if err != nil {
//		log.Printf("Failed to send message: %v", err)
//	}
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/denkhaus/agents/auth"
	"github.com/samber/mo"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/server"
)

// APIClient defines the core A2A client interface.
// Provides essential functionality for A2A communication.
type APIClient interface {
	// SendMessage sends a message and returns the result.
	// Returns an error if the message cannot be sent or processed.
	SendMessage(ctx context.Context, params protocol.SendMessageParams) (*protocol.MessageResult, error)

	// GetTasks retrieves task information based on query parameters.
	// Returns an error if the task cannot be found or accessed.
	GetTasks(ctx context.Context, params protocol.TaskQueryParams) (*protocol.Task, error)

	// CancelTasks cancels a running task by ID.
	// Returns the updated task state or an error if cancellation fails.
	CancelTasks(ctx context.Context, params protocol.TaskIDParams) (*protocol.Task, error)

	// SetPushNotification configures push notifications for a task.
	// Returns the updated configuration or an error if setup fails.
	SetPushNotification(ctx context.Context, config protocol.TaskPushNotificationConfig) (*protocol.TaskPushNotificationConfig, error)

	// GetPushNotification retrieves push notification configuration for a task.
	// Returns an error if the configuration cannot be retrieved.
	GetPushNotification(ctx context.Context, params protocol.TaskIDParams) (*protocol.TaskPushNotificationConfig, error)
}

// StreamingClient defines the streaming interface.
// Provides real-time streaming capabilities for A2A communication.
type StreamingClient interface {
	// StreamMessage starts a streaming message session.
	// Returns a channel for receiving streaming events and an error if streaming cannot be started.
	// The channel will be closed when the stream ends or context is cancelled.
	StreamMessage(ctx context.Context, params protocol.SendMessageParams) (<-chan protocol.StreamingMessageEvent, error)
}

// Authenticator defines authentication interface.
// This is an alias for auth.Provider to maintain consistency.
type Authenticator = auth.Provider

// AgentCardFetcher defines agent capability discovery interface.
// Provides functionality to discover and retrieve agent capabilities.
type AgentCardFetcher interface {
	// FetchAgentCard retrieves agent capabilities from the specified base URL.
	// Returns an error if the agent card cannot be fetched or parsed.
	FetchAgentCard(ctx context.Context, baseURL string) (*server.AgentCard, error)
}

// JWKSVerifier defines JWT verification interface.
// Provides JWT token verification using JSON Web Key Sets (JWKS).
type JWKSVerifier interface {
	// VerifyJWT verifies a JWT token against the payload.
	// Returns an error if the token is invalid, expired, or payload doesn't match.
	VerifyJWT(ctx context.Context, token string, payload []byte) error

	// RefreshJWKS refreshes the JWKS keyset from the remote endpoint.
	// Should be called periodically or when verification fails.
	RefreshJWKS(ctx context.Context) error
}

// TaskTracker defines task status tracking interface.
// Provides functionality to track and monitor task execution status.
type TaskTracker interface {
	// TrackTask adds a task to tracking with initial "pending" status.
	// Safe to call multiple times for the same task ID.
	TrackTask(taskID string)

	// UpdateTaskStatus updates the status of a tracked task.
	// Creates the task entry if it doesn't exist.
	UpdateTaskStatus(taskID, status string)

	// GetTaskStatus retrieves the current status of a task.
	// Returns "unknown" if the task is not being tracked.
	GetTaskStatus(taskID string) string
}

// ClientConfig defines configuration options using mo.Option for type safety.
type ClientConfig struct {
	// BaseURL is the base URL for the A2A service (required)
	BaseURL string
	// Timeout specifies the request timeout (default: 30 seconds)
	Timeout mo.Option[time.Duration]
	// UserAgent specifies the HTTP User-Agent header (default: "unified-a2a-client/1.0")
	UserAgent mo.Option[string]
	// Auth specifies the authentication provider (optional)
	Auth mo.Option[Authenticator]
	// Streaming enables/disables streaming support (default: true)
	Streaming mo.Option[bool]
	// RetryCount specifies the number of retries for failed requests (default: 3)
	RetryCount mo.Option[int]
}

// Validate validates the client configuration.
func (c ClientConfig) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	
	if timeout, hasTimeout := c.Timeout.Get(); hasTimeout && timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	
	if retryCount, hasRetryCount := c.RetryCount.Get(); hasRetryCount && retryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}
	
	return nil
}

// UnifiedClient combines all client interfaces.
// Provides a single interface for all A2A client functionality.
type UnifiedClient interface {
	APIClient
	StreamingClient
	AgentCardFetcher
}

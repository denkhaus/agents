// Package client provides a unified, thread-safe interface for A2A client functionality.
//
// This package consolidates all client features from the examples directory into
// a single, well-structured API with the following key components:
//
// # Core Interfaces
//
// - APIClient: Basic A2A operations (send message, get/cancel tasks, push notifications)
// - StreamingClient: Real-time streaming message capabilities
// - Authenticator: Authentication provider interface
// - AgentCardFetcher: Agent capability discovery
// - JWKSVerifier: JWT token verification using JWKS
// - TaskTracker: Task status tracking and management
//
// # Authentication
//
// The auth subpackage provides multiple authentication methods:
// - JWT authentication with configurable secrets, audiences, and expiry
// - API key authentication with custom headers
// - OAuth2 client credentials flow
//
// # Usage Examples
//
//	// Create a basic client
//	client, err := client.NewAPIClient("http://localhost:8080",
//		client.WithTimeout(30*time.Second))
//
//	// Create a streaming client with JWT auth
//	auth := auth.NewJWTProvider(auth.JWTConfig{
//		Secret:   []byte("secret"),
//		Audience: "a2a-server",
//		Issuer:   "client",
//		Expiry:   time.Hour,
//	})
//
//	streamClient, err := client.NewStreamingClient("http://localhost:8080",
//		client.WithAuth(auth),
//		client.WithStreaming(true))
//
//	// Create a unified client with all capabilities
//	config := client.ClientConfig{
//		BaseURL: "http://localhost:8080",
//		Timeout: mo.Some(30 * time.Second),
//		Auth:    mo.Some(auth),
//	}
//	unifiedClient, err := client.NewUnifiedClient(config)
//
// # Thread Safety
//
// All client implementations are thread-safe and can be used concurrently
// from multiple goroutines. Internal state is protected by appropriate
// synchronization primitives.
//
// # Error Handling
//
// The package provides structured error types in the helper subpackage:
// - HTTPError: HTTP-specific errors with status codes
// - ClientError: Client-side operation errors
// - AuthError: Authentication-related errors
// - ValidationError: Input validation errors
//
// # Configuration
//
// Configuration uses the samber/mo package for optional values, providing
// type-safe configuration with sensible defaults:
//
//	config := ClientConfig{
//		BaseURL:    "http://localhost:8080",
//		Timeout:    mo.Some(30 * time.Second),  // Optional with default
//		UserAgent:  mo.Some("my-client/1.0"),   // Optional with default
//		Streaming:  mo.Some(true),              // Optional with default
//		RetryCount: mo.Some(3),                 // Optional with default
//	}
//
// # Modular Design
//
// The package is structured in logical modules:
// - client/: Core interfaces and unified client implementation
// - client/auth/: Authentication providers (JWT, API key, OAuth2)
// - client/helper/: Utility functions (HTTP handling, response parsing, errors)
//
// Each module is self-contained and can be used independently while
// maintaining compatibility with the unified interface.
package client

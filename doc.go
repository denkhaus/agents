// Package agents provides a unified, production-ready A2A (Agent-to-Agent) communication library.
//
// This package offers comprehensive client and server implementations for A2A communication,
// including authentication, streaming, multi-agent routing, and push notifications.
//
// Key Features:
//   - Multiple authentication methods (JWT, API Key, OAuth2)
//   - Streaming message support
//   - Agent capability discovery
//   - Task tracking and management
//   - JWKS verification
//   - Multi-agent routing
//   - Push notifications
//   - Production-ready error handling
//
// Architecture:
//
// The library is structured around clear interfaces and factory functions:
//   - All concrete implementations are unexported
//   - Components are accessed through well-defined interfaces
//   - Factory functions provide type-safe configuration
//   - Comprehensive error handling throughout
//
// Basic Usage:
//
//	// Create a client
//	client, err := client.NewUnifiedClient(client.ClientConfig{
//		BaseURL: "https://api.example.com",
//		Auth:    mo.Some(auth.NewAPIKeyProvider(auth.APIKeyConfig{Key: "your-key"})),
//	})
//
//	// Create a server
//	server, err := server.NewServer(agentCard, taskManager,
//		server.WithPort(8080),
//		server.WithCORS(true),
//	)
//
// Package Structure:
//   - auth: Authentication providers and interfaces
//   - client: Client implementations and interfaces
//   - server: Server implementations and interfaces
//   - helper: Utility functions and error types
//
// All components follow Go best practices with proper error handling,
// context support, and thread safety.
package agents

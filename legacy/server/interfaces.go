// Package server provides a unified, production-ready A2A server implementation.
//
// This package defines server interfaces and provides factory functions for creating
// A2A servers. All concrete implementations are unexported and accessed through
// well-defined interfaces.
//
// Key interfaces:
//   - Server: Main server interface for A2A communication
//   - MessageProcessor: Interface for processing messages
//   - AuthenticationManager: Authentication management
//   - StreamingManager: Streaming functionality management
//   - NotificationManager: Push notification management
//
// Example usage:
//
//	// Create a server
//	server, err := server.NewServer(agentCard, taskManager,
//		server.WithPort(8080),
//		server.WithTLS("cert.pem", "key.pem"),
//		server.WithCORS(true),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Start the server
//	if err := server.Start(":8080"); err != nil {
//		log.Printf("Server failed: %v", err)
//	}
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/samber/mo"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

// Server defines the main server interface for A2A communication.
// Provides the core functionality for running an A2A server.
type Server interface {
	// Start starts the server on the specified address.
	// Blocks until the server stops or an error occurs.
	// Returns an error if the server cannot be started.
	Start(addr string) error
	
	// Stop gracefully stops the server with the given context timeout.
	// Returns an error if the server cannot be stopped gracefully.
	Stop(ctx context.Context) error
	
	// Handler returns the HTTP handler for integration with external servers.
	// Can be used to mount the A2A server on a custom HTTP server.
	Handler() http.Handler
}

// MessageProcessor defines the interface for processing messages.
// Provides the core message processing functionality.
type MessageProcessor interface {
	// ProcessMessage processes a message and returns a result.
	// Returns an error if the message cannot be processed.
	ProcessMessage(
		ctx context.Context,
		message protocol.Message,
		options taskmanager.ProcessOptions,
		handle taskmanager.TaskHandler,
	) (*taskmanager.MessageProcessingResult, error)
}

// AuthContext contains authentication information.
type AuthContext struct {
	UserID       string
	ProviderType string
	Claims       map[string]interface{}
	Scopes       []string
}

// ServerConfig defines server configuration using mo.Option for type safety.
type ServerConfig struct {
	// Host specifies the server host (default: "localhost")
	Host mo.Option[string]
	// Port specifies the server port (default: 8080)
	Port mo.Option[int]
	// TLS specifies TLS configuration (optional)
	TLS mo.Option[TLSConfig]
	// Timeout specifies request timeout (default: 30 seconds)
	Timeout mo.Option[time.Duration]
	// MaxRequestSize specifies maximum request size in bytes (default: 10MB)
	MaxRequestSize mo.Option[int64]
	// EnableCORS enables CORS support (default: false)
	EnableCORS mo.Option[bool]
	// EnableMetrics enables metrics collection (default: false)
	EnableMetrics mo.Option[bool]
	// ShutdownTimeout specifies graceful shutdown timeout (default: 30 seconds)
	ShutdownTimeout mo.Option[time.Duration]
}

// Validate validates the server configuration.
func (c ServerConfig) Validate() error {
	if port, hasPort := c.Port.Get(); hasPort && (port <= 0 || port > 65535) {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	
	if timeout, hasTimeout := c.Timeout.Get(); hasTimeout && timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	
	if maxSize, hasMaxSize := c.MaxRequestSize.Get(); hasMaxSize && maxSize <= 0 {
		return fmt.Errorf("max request size must be positive")
	}
	
	if shutdownTimeout, hasShutdownTimeout := c.ShutdownTimeout.Get(); hasShutdownTimeout && shutdownTimeout <= 0 {
		return fmt.Errorf("shutdown timeout must be positive")
	}
	
	return nil
}

// TLSConfig defines TLS configuration.
type TLSConfig struct {
	CertFile string
	KeyFile  string
	AutoTLS  bool
}

// ProcessorCapabilities defines what a message processor can do.
type ProcessorCapabilities struct {
	Streaming            bool
	MultiTurn            bool
	PushNotifications    bool
	CustomAuthentication bool
	InputModes           []string
	OutputModes          []string
}

// ProcessorMetrics contains processor performance metrics.
type ProcessorMetrics struct {
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	AverageLatency     time.Duration
	LastProcessedAt    time.Time
	ErrorRate          float64
}
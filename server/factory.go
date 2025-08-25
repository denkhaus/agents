// Package server provides factory functions for creating server instances.
package server

import "time"

import (
	"fmt"

	"github.com/samber/mo"
	"trpc.group/trpc-go/trpc-a2a-go/server"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

// ServerOption defines configuration options for server creation.
type ServerOption func(*ServerConfig)

// NewServer creates a new A2A server with the given configuration.
// Returns an error if the configuration is invalid or server creation fails.
func NewServer(agentCard server.AgentCard, taskManager taskmanager.TaskManager, opts ...ServerOption) (Server, error) {
	config := &ServerConfig{}
	
	// Apply options
	for _, opt := range opts {
		opt(config)
	}
	
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid server configuration: %w", err)
	}
	
	return newUnifiedServer(agentCard, taskManager, config)
}

// WithHost sets the server host.
func WithHost(host string) ServerOption {
	return func(c *ServerConfig) {
		c.Host = mo.Some(host)
	}
}

// WithPort sets the server port.
func WithPort(port int) ServerOption {
	return func(c *ServerConfig) {
		c.Port = mo.Some(port)
	}
}

// WithTLS enables TLS with the given configuration.
func WithTLS(certFile, keyFile string) ServerOption {
	return func(c *ServerConfig) {
		c.TLS = mo.Some(TLSConfig{
			CertFile: certFile,
			KeyFile:  keyFile,
		})
	}
}

// WithAutoTLS enables automatic TLS certificate management.
func WithAutoTLS() ServerOption {
	return func(c *ServerConfig) {
		c.TLS = mo.Some(TLSConfig{
			AutoTLS: true,
		})
	}
}

// WithTimeout sets the server request timeout.
func WithTimeout(timeout time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.Timeout = mo.Some(timeout)
	}
}

// WithMaxRequestSize sets the maximum request size.
func WithMaxRequestSize(size int64) ServerOption {
	return func(c *ServerConfig) {
		c.MaxRequestSize = mo.Some(size)
	}
}

// WithShutdownTimeout sets the graceful shutdown timeout.
func WithShutdownTimeout(timeout time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.ShutdownTimeout = mo.Some(timeout)
	}
}

// WithCORS enables CORS support.
func WithCORS(enabled bool) ServerOption {
	return func(c *ServerConfig) {
		c.EnableCORS = mo.Some(enabled)
	}
}

// WithMetrics enables metrics collection.
func WithMetrics(enabled bool) ServerOption {
	return func(c *ServerConfig) {
		c.EnableMetrics = mo.Some(enabled)
	}
}
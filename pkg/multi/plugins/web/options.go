package web

import (
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/runner"
)

// Server exposes HTTP endpoints compatible with the ADK Web UI. Internally it
// reuses the trpc-agent-go components for sessions, runners and events.
type Server struct {
	router     *mux.Router
	runnerOpts []runner.Option // Extra options applied when creating a runner.

	traces         map[string]attribute.Set // key: event_id
	memoryExporter *inMemoryExporter

	// Multi-Agent Chat support
	chatProcessor multi.ChatProcessor

	// SSE Connection Pool for Inter-Agent Communication
	ssePool *SSEConnectionPool
	logger  *zap.Logger
}

// Option configures the Server instance.
type Option func(*Server)

// WithRunnerOptions appends additional runner.Option values applied when the
// server lazily constructs a Runner for an agent.
func WithRunnerOptions(opts ...runner.Option) Option {
	return func(s *Server) { s.runnerOpts = append(s.runnerOpts, opts...) }
}

// WithLogger provides a zap.Logger to the server.
// If omitted, a no-op logger is used.
func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) { s.logger = logger }
}

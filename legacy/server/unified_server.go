// Package server provides the main server implementation.
package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"trpc.group/trpc-go/trpc-a2a-go/server"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

// unifiedServer implements Server.
type unifiedServer struct {
	config      ServerConfig
	httpServer  *http.Server
	agentCard   server.AgentCard
	taskManager taskmanager.TaskManager
	mu          sync.RWMutex
}

// newUnifiedServer creates a new unified server with the given configuration.
func newUnifiedServer(agentCard server.AgentCard, taskManager taskmanager.TaskManager, config *ServerConfig) (Server, error) {
	srv := &unifiedServer{
		config:      *config,
		agentCard:   agentCard,
		taskManager: taskManager,
	}

	return srv, nil
}

// Start starts the server on the specified address.
func (s *unifiedServer) Start(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.httpServer != nil {
		return fmt.Errorf("server already started")
	}

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.Handler(),
		ReadTimeout:  s.config.Timeout.OrElse(30 * time.Second),
		WriteTimeout: s.config.Timeout.OrElse(30 * time.Second),
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	if tlsConfig, hasTLS := s.config.TLS.Get(); hasTLS {
		return s.httpServer.ListenAndServeTLS(tlsConfig.CertFile, tlsConfig.KeyFile)
	}

	return s.httpServer.ListenAndServe()
}

// Stop gracefully stops the server.
func (s *unifiedServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.httpServer == nil {
		return nil
	}

	err := s.httpServer.Shutdown(ctx)
	s.httpServer = nil
	return err
}

// Handler returns the HTTP handler for integration.
func (s *unifiedServer) Handler() http.Handler {
	mux := http.NewServeMux()

	// Add basic routes
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/agent-card", s.handleAgentCard)

	// Add CORS if enabled
	if s.config.EnableCORS.OrElse(false) {
		return s.corsMiddleware(mux)
	}

	return mux
}

// handleHealth handles health check requests.
func (s *unifiedServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// handleAgentCard handles agent card requests.
func (s *unifiedServer) handleAgentCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Simple implementation - would use JSON encoder in production
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"name":"%s","description":"%s"}`, s.agentCard.Name, s.agentCard.Description)))
}

// corsMiddleware adds CORS headers.
func (s *unifiedServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
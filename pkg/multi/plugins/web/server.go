//
// Tencent is pleased to support the open source community by making trpc-agent-go available.
//
// Copyright (C) 2025 Tencent.  All rights reserved.
//
// trpc-agent-go is licensed under the Apache License Version 2.0.
//
//

// Package debug provides a HTTP server for debugging and testing.
package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	//itelemetry "trpc.group/trpc-go/trpc-agent-go/internal/telemetry"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins/web/schema"
)

// New creates a new CLI HTTP server with explicit agent registration. The
// behaviour can be tweaked via functional options.
func New(opts ...Option) *Server {
	s := &Server{
		router:         mux.NewRouter(),
		traces:         make(map[string]attribute.Set),
		memoryExporter: newInMemoryExporter(),
		ssePool:        NewSSEConnectionPool(),
	}

	// Apply user-provided options.
	for _, opt := range opts {
		opt(s)
	}

	// Set default logger if none provided
	if s.logger == nil {
		s.logger = zap.NewNop()
	}

	// Add CORS middleware for ADK Web compatibility.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
	})

	s.router.Use(c.Handler)
	s.registerRoutes()

	provider := sdktrace.NewTracerProvider()
	provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(newApiServerSpanExporter(s.traces)))
	provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(s.memoryExporter))
	otel.SetTracerProvider(provider)
	//atrace.Tracer = otel.Tracer(itelemetry.InstrumentName)
	//setTraceInfo()
	return s
}

// SetChatProcessor sets the multi-agent chat processor for the server.
func (s *Server) SetChatProcessor(p multi.ChatProcessor) {
	s.chatProcessor = p
	s.setupInterAgentInterceptor()
}

// Handler returns the http.Handler for the server.
func (s *Server) Handler() http.Handler { return s.router }

// registerRoutes sets up all REST endpoints expected by ADK Web.
func (s *Server) registerRoutes() {
	s.router.HandleFunc("/app-info", s.handleAppInfo).Methods(http.MethodGet)

	// Session APIs.
	s.router.HandleFunc("/apps/{appName}/users/{userId}/sessions",
		s.handleListSessions).Methods(http.MethodGet)
	s.router.HandleFunc("/apps/{appName}/users/{userId}/sessions",
		s.handleCreateSession).Methods(http.MethodPost)
	s.router.HandleFunc("/apps/{appName}/users/{userId}/sessions/{sessionId}",
		s.handleGetSession).Methods(http.MethodGet)
	s.router.HandleFunc("/apps/{appName}/users/{userId}/sessions/{sessionId}",
		s.handleDeleteSession).Methods(http.MethodDelete)

	// Debug APIs
	s.router.HandleFunc("/debug/trace/{event_id}",
		s.handleEventTrace).Methods(http.MethodGet)
	s.router.HandleFunc("/debug/trace/session/{session_id}",
		s.handleSessionTrace).Methods(http.MethodGet)

	// Runner APIs.
	s.router.HandleFunc("/run", s.handleRun).Methods(http.MethodPost)
	s.router.HandleFunc("/run_sse", s.handleRunSSE).Methods(http.MethodPost)

	// OPTIONS handlers to allow CORS pre-flight
	preflight := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	}
	s.router.HandleFunc("/list-agents", preflight).Methods(http.MethodOptions)
	s.router.HandleFunc("/run", preflight).Methods(http.MethodOptions)
	s.router.HandleFunc("/run_sse", preflight).Methods(http.MethodOptions)

	// Session API OPTIONS handlers
	s.router.HandleFunc("/apps/{appName}/users/{userId}/sessions", preflight).Methods(http.MethodOptions)
	s.router.HandleFunc("/apps/{appName}/users/{userId}/sessions/{sessionId}", preflight).Methods(http.MethodOptions)

	// Debug API OPTIONS handlers
	s.router.HandleFunc("/debug/trace/{event_id}", preflight).Methods(http.MethodOptions)
	s.router.HandleFunc("/debug/trace/session/{session_id}", preflight).Methods(http.MethodOptions)
}

// ---- Handlers -----------------------------------------------------------

func (s *Server) handleAppInfo(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("handleAppInfo called", zap.String("path", r.URL.Path))

	// Prefer agents from chatProcessor if available (they have send_message tools)
	if s.chatProcessor == nil {
		http.Error(w, "chat processor is undefined", http.StatusBadRequest)
		return
	}

	aux := map[string]interface{}{
		"applicationName": s.chatProcessor.GetApplicationName(),
		"agents":          s.GetMultiChatAgents(),
	}

	s.writeJSON(w, aux)
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("handleRun called", zap.String("path", r.URL.Path))

	var req schema.AgentRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// If the request is for streaming, delegate to the SSE handler.
	if req.Streaming {
		// As we can't directly pass the decoded body, the SSE handler will re-decode.
		// A more optimized approach might involve passing the decoded struct via context.
		s.handleRunSSE(w, r)
		return
	}

	out, err := s.chatProcessor.SendMessage(
		r.Context(),
		req.FromAgentID,
		req.ToAgentID,
		req.SessionID,
		req.Content.ToMessage(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For non-streaming, we might want to collect all events or just return the final one.
	// ADK web might expect a list of events. Let's send all of them.
	var events []*schema.LLMEvent
	for e := range out {
		if e.Response != nil && e.Response.IsPartial {
			continue // skip streaming chunks in non-streaming endpoint
		}
		if ev, err := schema.NewLLMEvent(e, false); err == nil && ev != nil {
			events = append(events, ev)
		}
		// Note: Silently ignoring events that fail to convert for now
		// In a production system, you might want to log these errors
	}

	s.writeJSON(w, events)
}

func (s *Server) handleRunSSE(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("handleRunSSE called", zap.String("path", r.URL.Path))

	var req schema.AgentRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering

	// Register SSE connection for inter-agent communication
	cleanup := s.RegisterSSEConnectionForRequest(req, w, r)
	defer cleanup()

	out, err := s.chatProcessor.SendMessage(
		r.Context(),
		req.FromAgentID,
		req.ToAgentID,
		req.SessionID,
		req.Content.ToMessage(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Streaming {
		for e := range out {
			sseEvent, err := schema.NewLLMEvent(e, req.Streaming)
			if err != nil || sseEvent == nil {
				continue
			}
			data, err := json.Marshal(sseEvent)
			if err != nil {
				s.logger.Error("Error marshalling SSE event", zap.Error(err))
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	} else {
		// Non-streaming mode: wait for the first complete event and send only that.
		for e := range out {
			sseEvent, err := schema.NewLLMEvent(e, req.Streaming)
			if err != nil {
				s.logger.Error("Error creating LLMEvent", zap.Error(err))
				continue
			}
			if sseEvent == nil {
				continue
			}
			data, err := json.Marshal(sseEvent)
			if err != nil {
				s.logger.Error("Error marshalling SSE event", zap.Error(err))
				break
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}

	// Send a final event to properly close the stream
	fmt.Fprintf(w, "data: {\"done\": true}\n\n")
	flusher.Flush()

	s.logger.Info("handleRunSSE finished", zap.Any("sessionID", req.SessionID))
}

func (s *Server) writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

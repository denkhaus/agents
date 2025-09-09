package web

import (
	"net/http"

	"github.com/denkhaus/agents/pkg/multi/plugins/web/schema"
	"trpc.group/trpc-go/trpc-agent-go/log"
)

// RegisterSSEConnectionForRequest registers an SSE connection for the current request
func (s *Server) RegisterSSEConnectionForRequest(req schema.AgentRunRequest, w http.ResponseWriter, r *http.Request) func() {
	// Register SSE connection for inter-agent communication
	sseConn := s.ssePool.RegisterConnection(r.Context(), req.SessionID, req.ToAgentID, w)

	if sseConn != nil {
		log.Infof("SSE connection registered for agent %s, session %s", req.AppName, req.SessionID)
		// Return cleanup function
		return func() {
			s.ssePool.UnregisterConnection(req.SessionID, req.ToAgentID)
		}
	}

	// Return no-op cleanup function if registration failed
	return func() {}
}

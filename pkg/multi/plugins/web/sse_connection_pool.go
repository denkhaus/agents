package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/log"
)

// SSEConnection represents an active SSE connection for an agent session
type SSEConnection struct {
	Writer    http.ResponseWriter
	Flusher   http.Flusher
	AgentID   uuid.UUID
	SessionID uuid.UUID
	Context   context.Context
	Cancel    context.CancelFunc
}

// SSEConnectionPool manages active SSE connections for inter-agent communication
type SSEConnectionPool struct {
	connections map[string]*SSEConnection // key: connectionID (sessionID:agentName)
	mu          sync.RWMutex
}

// NewSSEConnectionPool creates a new connection pool
func NewSSEConnectionPool() *SSEConnectionPool {
	return &SSEConnectionPool{
		connections: make(map[string]*SSEConnection),
	}
}

// generateConnectionID creates a unique connection ID for SSE connections
func (pool *SSEConnectionPool) generateConnectionID(sessionID, agentID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", sessionID, agentID)
}

// RegisterConnection registers a new SSE connection in the pool
func (pool *SSEConnectionPool) RegisterConnection(
	ctx context.Context,
	sessionID uuid.UUID,
	agentID uuid.UUID,
	w http.ResponseWriter,
) *SSEConnection {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil
	}

	connCtx, cancel := context.WithCancel(ctx)
	conn := &SSEConnection{
		Writer:    w,
		Flusher:   flusher,
		AgentID:   agentID,
		SessionID: sessionID,
		Context:   connCtx,
		Cancel:    cancel,
	}

	connectionID := pool.generateConnectionID(sessionID, agentID)

	pool.mu.Lock()
	pool.connections[connectionID] = conn
	pool.mu.Unlock()

	log.Infof("Registered SSE connection: %s (agent: %s, session: %s)", connectionID, agentID, sessionID)
	return conn
}

// UnregisterConnection removes an SSE connection from the pool
func (pool *SSEConnectionPool) UnregisterConnection(sessionID, agentID uuid.UUID) {
	connectionID := pool.generateConnectionID(sessionID, agentID)

	pool.mu.Lock()
	if conn, exists := pool.connections[connectionID]; exists {
		conn.Cancel()
		delete(pool.connections, connectionID)
		log.Infof("Unregistered SSE connection: %s", connectionID)
	}
	pool.mu.Unlock()
}

// SendEventToConnection sends an event to a specific SSE connection
func (pool *SSEConnectionPool) SendEventToConnection(conn *SSEConnection, event *messaging.LLMEvent) {
	select {
	case <-conn.Context.Done():
		// Connection is closed, skip sending
		return
	default:
		data, err := json.Marshal(event)
		if err != nil {
			log.Errorf("Error marshalling inter-agent event: %v", err)
			return
		}

		_, err = fmt.Fprintf(conn.Writer, "data: %s\n\n", data)
		if err != nil {
			log.Errorf("Error writing to SSE connection: %v", err)
			return
		}

		conn.Flusher.Flush()
		log.Debugf("Sent inter-agent event to connection %s:%s", conn.SessionID, conn.AgentID)
	}
}

// BroadcastToAgent sends an event to all connections of a specific agent
func (pool *SSEConnectionPool) BroadcastToAgent(event *messaging.LLMEvent) int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	sentCount := 0
	for connectionID, conn := range pool.connections {
		if conn.AgentID == event.Routing.ToAgentID {
			go pool.SendEventToConnection(conn, event)
			sentCount++
			log.Debugf("Sent inter-agent event to connection: %s", connectionID)
		}
	}

	return sentCount
}

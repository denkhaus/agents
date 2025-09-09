package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/denkhaus/agents/pkg/multi/plugins/web/schema"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// mockSSEChatProcessor is a mock implementation of the multi.ChatProcessor interface for SSE testing.
type mockSSEChatProcessor struct {
	mockChatProcessor
	sessions map[string]*session.Session
}

func newMockSSEChatProcessor() *mockSSEChatProcessor {
	return &mockSSEChatProcessor{
		sessions: make(map[string]*session.Session),
	}
}

func (m *mockSSEChatProcessor) SendMessage(
	ctx context.Context,
	fromAgentID, toAgentID, sessionID uuid.UUID,
	message model.Message,
) (<-chan *event.Event, error) {
	events := make(chan *event.Event, 2)

	// Send a partial event for streaming
	go func() {
		defer close(events)

		// First event - partial
		events <- &event.Event{
			Response: &model.Response{
				Done:      false,
				IsPartial: true,
				Choices: []model.Choice{{
					Delta: model.Message{
						Role:    model.RoleAssistant,
						Content: "Partial response",
					},
				}},
			},
		}

		// Second event - complete
		time.Sleep(10 * time.Millisecond) // Small delay to simulate streaming
		events <- &event.Event{
			Response: &model.Response{
				Done:      true,
				IsPartial: false,
				Choices: []model.Choice{{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: "Complete response",
					},
				}},
			},
		}
	}()

	return events, nil
}

func (m *mockSSEChatProcessor) CreateSession(ctx context.Context, key session.Key, state session.StateMap, options ...session.Option) (*session.Session, error) {
	sess := &session.Session{
		ID:        uuid.New().String(),
		AppName:   key.AppName,
		UserID:    key.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		State:     state,
	}
	m.sessions[sess.ID] = sess
	return sess, nil
}

func (m *mockSSEChatProcessor) GetSession(ctx context.Context, key session.Key, options ...session.Option) (*session.Session, error) {
	sess, exists := m.sessions[key.SessionID]
	if !exists {
		return nil, nil
	}
	return sess, nil
}

func (m *mockSSEChatProcessor) ListSessions(ctx context.Context, userKey session.UserKey, options ...session.Option) ([]*session.Session, error) {
	var sessions []*session.Session
	for _, sess := range m.sessions {
		if sess.UserID == userKey.UserID && sess.AppName == userKey.AppName {
			sessions = append(sessions, sess)
		}
	}
	return sessions, nil
}

func (m *mockSSEChatProcessor) DeleteSession(ctx context.Context, key session.Key, options ...session.Option) error {
	delete(m.sessions, key.SessionID)
	return nil
}

func TestServer_handleRunSSE(t *testing.T) {
	processor := newMockSSEChatProcessor()
	server := New(nil)
	server.SetChatProcessor(processor)

	fromAgentID := uuid.New()
	toAgentID := uuid.New()
	sessionID := uuid.New()

	// Create a test request for SSE streaming
	requestBody := schema.AgentRunRequest{
		AppName:     "test-app",
		FromAgentID: fromAgentID,
		ToAgentID:   toAgentID,
		SessionID:   sessionID,
		Content: schema.Content{
			Role: "user",
			Parts: []schema.PartIncoming{
				{Text: "Hello, SSE streaming!"},
			},
		},
		Streaming: true,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/run_sse", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder that captures the full response
	w := httptest.NewRecorder()

	server.handleRunSSE(w, req)

	// Check the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}

	// Check that the response has the correct headers for SSE
	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("Expected Content-Type 'text/event-stream', got '%s'", w.Header().Get("Content-Type"))
	}

	// Check that the response body contains SSE events
	responseBody := w.Body.String()
	if responseBody == "" {
		t.Error("Expected non-empty response body")
	}

	// Check that we got streaming events
	if !bytes.Contains(w.Body.Bytes(), []byte("data:")) {
		t.Error("Expected SSE data events in response")
	}

	t.Logf("SSE Response: %s", responseBody)
}

func TestServer_handleRunSSE_NonStreaming(t *testing.T) {
	processor := newMockSSEChatProcessor()
	server := New(nil)
	server.SetChatProcessor(processor)

	fromAgentID := uuid.New()
	toAgentID := uuid.New()
	sessionID := uuid.New()

	// Create a test request for non-streaming
	requestBody := schema.AgentRunRequest{
		AppName:     "test-app",
		FromAgentID: fromAgentID,
		ToAgentID:   toAgentID,
		SessionID:   sessionID,
		Content: schema.Content{
			Role: "user",
			Parts: []schema.PartIncoming{
				{Text: "Hello, non-streaming!"},
			},
		},
		Streaming: false,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/run_sse", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	server.handleRunSSE(w, req)

	// Check the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}

	// Check that the response has the correct headers for SSE
	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("Expected Content-Type 'text/event-stream', got '%s'", w.Header().Get("Content-Type"))
	}

	// Check that the response body contains SSE events
	responseBody := w.Body.String()
	if responseBody == "" {
		t.Error("Expected non-empty response body")
	}

	// Check that we got events
	if !bytes.Contains(w.Body.Bytes(), []byte("data:")) {
		t.Error("Expected SSE data events in response")
	}

	t.Logf("Non-streaming SSE Response: %s", responseBody)
}

func TestSSEConnectionPool_RegisterConnection(t *testing.T) {
	pool := NewSSEConnectionPool()

	if pool == nil {
		t.Fatal("Expected non-nil SSEConnectionPool")
	}

	// Create a mock HTTP response writer
	w := httptest.NewRecorder()

	var testFlusher any = w
	// Simulate an HTTP response that supports flushing
	_, ok := testFlusher.(http.Flusher)
	if !ok {
		t.Fatal("Expected httptest.ResponseRecorder to implement http.Flusher")
	}

	sessionID := uuid.New()
	agentID := uuid.New()

	// Test registering a connection
	conn := pool.RegisterConnection(context.Background(), sessionID, agentID, w)

	if conn == nil {
		t.Fatal("Expected non-nil SSEConnection")
	}

	// Check that the connection has the expected properties
	if conn.AgentID != agentID {
		t.Error("Expected AgentID to match")
	}

	if conn.SessionID != sessionID {
		t.Error("Expected SessionID to match")
	}

	if conn.Context == nil {
		t.Error("Expected non-nil Context")
	}

	if conn.Cancel == nil {
		t.Error("Expected non-nil Cancel function")
	}

	// Check that the connection was added to the pool
	pool.mu.RLock()
	_, exists := pool.connections[pool.generateConnectionID(sessionID, agentID)]
	pool.mu.RUnlock()

	if !exists {
		t.Error("Expected connection to be added to the pool")
	}
}

func TestSSEConnectionPool_UnregisterConnection(t *testing.T) {
	pool := NewSSEConnectionPool()

	sessionID := uuid.New()
	agentID := uuid.New()

	// Create a mock HTTP response writer
	w := httptest.NewRecorder()

	// Register a connection
	conn := pool.RegisterConnection(context.Background(), sessionID, agentID, w)
	if conn == nil {
		t.Fatal("Expected non-nil SSEConnection")
	}

	// Verify the connection was registered
	pool.mu.RLock()
	_, exists := pool.connections[pool.generateConnectionID(sessionID, agentID)]
	pool.mu.RUnlock()

	if !exists {
		t.Error("Expected connection to be registered")
	}

	// Unregister the connection
	pool.UnregisterConnection(sessionID, agentID)

	// Verify the connection was unregistered
	pool.mu.RLock()
	_, exists = pool.connections[pool.generateConnectionID(sessionID, agentID)]
	pool.mu.RUnlock()

	if exists {
		t.Error("Expected connection to be unregistered")
	}
}

func TestSSEConnectionPool_BroadcastToAgent(t *testing.T) {
	pool := NewSSEConnectionPool()

	agentID := uuid.New()
	sessionID1 := uuid.New()
	sessionID2 := uuid.New()

	// Create mock HTTP response writers
	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	// Register two connections for the same agent
	conn1 := pool.RegisterConnection(context.Background(), sessionID1, agentID, w1)
	conn2 := pool.RegisterConnection(context.Background(), sessionID2, agentID, w2)

	if conn1 == nil || conn2 == nil {
		t.Fatal("Expected non-nil SSEConnections")
	}

	// Create a test event
	event := &schema.LLMEvent{
		ID:           uuid.New().String(),
		InvocationID: uuid.New().String(),
		Author:       "test-agent",
		Type:         schema.EventTypeAssistant,
		Done:         true,
		Role:         model.RoleAssistant,
		Parts: []schema.Part{
			&schema.TextPart{Content: "Test broadcast message"},
		},
	}

	// Broadcast the event
	sentCount := pool.BroadcastToAgent(agentID, event)

	// Give some time for goroutines to complete
	time.Sleep(50 * time.Millisecond)

	// Should have sent to 2 connections
	if sentCount != 2 {
		t.Errorf("Expected sentCount to be 2, got %d", sentCount)
	}

	// Check that data was written to both connections
	// Note: Since SendEventToConnection runs in a goroutine, we need to wait a bit
	time.Sleep(100 * time.Millisecond)

	if w1.Body.Len() == 0 {
		t.Error("Expected data to be written to first connection")
	}

	if w2.Body.Len() == 0 {
		t.Error("Expected data to be written to second connection")
	}

	t.Logf("Connection 1 response: %s", w1.Body.String())
	t.Logf("Connection 2 response: %s", w2.Body.String())
}

func TestServer_RegisterSSEConnectionForRequest(t *testing.T) {
	server := New(nil)

	sessionID := uuid.New()
	agentID := uuid.New()

	request := schema.AgentRunRequest{
		AppName:     "test-app",
		FromAgentID: uuid.New(),
		ToAgentID:   agentID,
		SessionID:   sessionID,
		Content: schema.Content{
			Role: "user",
			Parts: []schema.PartIncoming{
				{Text: "Test message"},
			},
		},
		Streaming: true,
	}

	// Create a mock HTTP response writer
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/run_sse", nil)

	// Register SSE connection
	cleanup := server.RegisterSSEConnectionForRequest(request, w, req)

	if cleanup == nil {
		t.Error("Expected non-nil cleanup function")
	}

	// Verify connection was registered
	connectionID := server.ssePool.generateConnectionID(sessionID, agentID)

	server.ssePool.mu.RLock()
	_, exists := server.ssePool.connections[connectionID]
	server.ssePool.mu.RUnlock()

	if !exists {
		t.Error("Expected SSE connection to be registered")
	}

	// Call cleanup function
	cleanup()

	// Verify connection was unregistered
	server.ssePool.mu.RLock()
	_, exists = server.ssePool.connections[connectionID]
	server.ssePool.mu.RUnlock()

	if exists {
		t.Error("Expected SSE connection to be unregistered after cleanup")
	}
}

// Test session management endpoints
func TestServer_handleListSessions(t *testing.T) {
	processor := newMockSSEChatProcessor()
	server := New(nil)
	server.SetChatProcessor(processor)

	// Create a few test sessions
	agentID := uuid.New()
	appName := "test-app"

	for i := 0; i < 3; i++ {
		_, err := processor.CreateSession(context.Background(), session.Key{
			AppName: appName,
			UserID:  agentID.String(),
		}, session.StateMap{})

		if err != nil {
			t.Fatalf("Failed to create test session: %v", err)
		}
	}

	// Create request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/apps/%s/users/%s/sessions", appName, agentID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{
		"appName": appName,
		"agentId": agentID.String(),
	})

	w := httptest.NewRecorder()

	server.handleListSessions(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	// Parse response
	var sessions []*schema.ADKSession
	if err := json.Unmarshal(w.Body.Bytes(), &sessions); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("Expected 3 sessions, got %d", len(sessions))
	}

	t.Logf("List sessions response: %s", w.Body.String())
}

func TestServer_handleGetSession(t *testing.T) {
	processor := newMockSSEChatProcessor()
	server := New(nil)
	server.SetChatProcessor(processor)

	// Create a test session
	agentID := uuid.New()
	appName := "test-app"
	sessionID := uuid.New()

	testSession := &session.Session{
		ID:        sessionID.String(),
		AppName:   appName,
		UserID:    agentID.String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		State:     session.StateMap{},
	}
	processor.sessions[sessionID.String()] = testSession

	// Create request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/apps/%s/users/%s/sessions/%s", appName, agentID.String(), sessionID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{
		"appName":   appName,
		"userId":    agentID.String(),
		"sessionId": sessionID.String(),
	})

	w := httptest.NewRecorder()

	server.handleGetSession(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	// Parse response
	var session schema.ADKSession
	if err := json.Unmarshal(w.Body.Bytes(), &session); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if session.ID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, session.ID)
	}

	t.Logf("Get session response: %s", w.Body.String())
}

func TestServer_handleGetSession_NotFound(t *testing.T) {
	processor := newMockSSEChatProcessor()
	server := New(nil)
	server.SetChatProcessor(processor)

	// Create request for non-existent session
	agentID := uuid.New()
	appName := "test-app"
	sessionID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/apps/%s/users/%s/sessions/%s", appName, agentID.String(), sessionID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{
		"appName":   appName,
		"userId":    agentID.String(),
		"sessionId": sessionID.String(),
	})

	w := httptest.NewRecorder()

	server.handleGetSession(w, req)

	// Check response - should be 500 because of the mock processor behavior
	// In a real implementation, this would be 404
	if w.Code == http.StatusOK {
		t.Error("Expected error status for non-existent session")
	}

	t.Logf("Get session not found response: %s", w.Body.String())
}

func TestServer_handleDeleteSession(t *testing.T) {
	processor := newMockSSEChatProcessor()
	server := New(nil)
	server.SetChatProcessor(processor)

	// Create a test session
	agentID := uuid.New()
	appName := "test-app"
	sessionID := uuid.New()

	testSession := &session.Session{
		ID:        sessionID.String(),
		AppName:   appName,
		UserID:    agentID.String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		State:     session.StateMap{},
	}
	processor.sessions[sessionID.String()] = testSession

	// Verify session exists before deletion
	if _, exists := processor.sessions[sessionID.String()]; !exists {
		t.Fatal("Expected session to exist before deletion")
	}

	// Create request
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/apps/%s/users/%s/sessions/%s", appName, agentID.String(), sessionID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{
		"appName":   appName,
		"agentId":   agentID.String(),
		"sessionId": sessionID.String(),
	})

	w := httptest.NewRecorder()

	server.handleDeleteSession(w, req)

	// Check response
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status NoContent, got %d", w.Code)
	}

	// Verify session was deleted
	if _, exists := processor.sessions[sessionID.String()]; exists {
		t.Error("Expected session to be deleted")
	}

	t.Logf("Delete session response status: %d", w.Code)
}

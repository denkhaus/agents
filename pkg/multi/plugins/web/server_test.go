//
// Tencent is pleased to support the open source community by making trpc-agent-go available.
//
// Copyright (C) 2025 Tencent.  All rights reserved.
//
// trpc-agent-go is licensed under the Apache License Version 2.0.
//
//

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins/web/schema"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/session"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// mockChatProcessor is a mock implementation of the multi.ChatProcessor interface.
type mockChatProcessor struct{}

func (m *mockChatProcessor) SetMessageInterceptor(interceptor messaging.Interceptor) {}
func (m *mockChatProcessor) SendMessage(ctx context.Context, fromAgentID, toAgentID, sessionID uuid.UUID, message model.Message) (<-chan *event.Event, error) {
	events := make(chan *event.Event, 1)
	go func() {
		defer close(events)
		events <- &event.Event{
			Response: &model.Response{
				Done: true, // Set Done to true to indicate this is a complete response
				Choices: []model.Choice{{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: "test response",
					},
				}},
			},
		}
	}()
	return events, nil
}
func (m *mockChatProcessor) SendMessageWithProcessing(ctx context.Context, fromAgentID, toAgentID, sessionID uuid.UUID, message model.Message) error {
	return nil
}
func (m *mockChatProcessor) CreateSession(ctx context.Context, key session.Key, state session.StateMap, options ...session.Option) (*session.Session, error) {
	return &session.Session{
		ID:        uuid.New().String(),
		AppName:   key.AppName,
		UserID:    key.UserID, // Use the UserID from the key
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
func (m *mockChatProcessor) ListSessions(ctx context.Context, userKey session.UserKey, options ...session.Option) ([]*session.Session, error) {
	return nil, nil
}
func (m *mockChatProcessor) DeleteSession(ctx context.Context, key session.Key, options ...session.Option) error {
	return nil
}
func (m *mockChatProcessor) GetSession(ctx context.Context, key session.Key, options ...session.Option) (*session.Session, error) {
	return nil, nil
}
func (m *mockChatProcessor) GetAgentInfoByAuthor(author string) *shared.AgentInfo { return nil }
func (m *mockChatProcessor) GetAgentInfoByID(agentID uuid.UUID) *shared.AgentInfo { return nil }
func (m *mockChatProcessor) GetAllAgentInfos() []*shared.AgentInfo {
	return []*shared.AgentInfo{
		{Info: agent.Info{Name: "agent1"}},
		{Info: agent.Info{Name: "agent2"}},
	}
}
func (m *mockChatProcessor) GetApplicationName() string { return "test-app" }
func (m *mockChatProcessor) GetAgentNameByID(agentID uuid.UUID) string {
	return ""
}
func (m *mockChatProcessor) GetAgentByName(name string) shared.TheAgent        { return nil }
func (m *mockChatProcessor) SetOnMessageCallback(onMessage multi.OnMessage)    {}
func (m *mockChatProcessor) SetOnToolCallCallback(onToolCall multi.OnToolCall) {}

// mockAgent is a simple mock agent for testing.
type mockAgent struct {
	name        string
	description string
}

func (m *mockAgent) Info() agent.Info {
	return agent.Info{
		Name:        m.name,
		Description: m.description,
	}
}

func (m *mockAgent) Tools() []tool.Tool { return nil }

func (m *mockAgent) SubAgents() []agent.Agent { return nil }

func (m *mockAgent) FindSubAgent(name string) agent.Agent { return nil }

func (m *mockAgent) Run(ctx context.Context, inv *agent.Invocation) (<-chan *event.Event, error) {
	// Return a simple event channel for testing.
	events := make(chan *event.Event, 1)
	go func() {
		defer close(events)
		events <- &event.Event{
			Response: &model.Response{
				Choices: []model.Choice{{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: "test response",
					},
				}},
			},
		}
	}()
	return events, nil
}

func TestServer_Handler(t *testing.T) {

	server := New()
	handler := server.Handler()

	if handler == nil {
		t.Fatal("Handler() returned nil")
	}
}

func TestServer_handleListAgents(t *testing.T) {
	server := New(nil)
	server.SetChatProcessor(&mockChatProcessor{})
	req := httptest.NewRequest(http.MethodGet, "/app-info", nil)
	w := httptest.NewRecorder()

	server.handleAppInfo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var appInfo struct {
		ApplicationName string `json:"applicationName"`
		Agents          []struct {
			Name string `json:"name"`
		} `json:"agents"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &appInfo); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if appInfo.ApplicationName != "test-app" {
		t.Errorf("expected applicationName 'test-app', got '%s'", appInfo.ApplicationName)
	}

	if len(appInfo.Agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(appInfo.Agents))
	}

	found := make(map[string]bool)
	for _, agent := range appInfo.Agents {
		found[agent.Name] = true
	}

	if !found["agent1"] || !found["agent2"] {
		t.Error("expected agent names not found in response")
	}
}

func TestServer_handleCreateSession(t *testing.T) {
	server := New(nil)
	server.SetChatProcessor(&mockChatProcessor{})

	agentID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/apps/test-agent/users/"+agentID.String()+"/sessions", nil)
	w := httptest.NewRecorder()

	// Set up the route variables that gorilla/mux would normally set.
	req = mux.SetURLVars(req, map[string]string{
		"appName": "test-agent",
		"agentId": agentID.String(), // Changed from userId to agentId to match the route
	})

	server.handleCreateSession(w, req)

	// Print the response body for debugging
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var session schema.ADKSession
	if err := json.Unmarshal(w.Body.Bytes(), &session); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if session.AppName != "test-agent" {
		t.Errorf("expected appName 'test-agent', got '%s'", session.AppName)
	}

	if session.AgentID != agentID {
		t.Errorf("expected agentId %s, got '%s'", agentID, session.AgentID)
	}

	if session.ID == uuid.Nil {
		t.Error("expected non-empty session ID")
	}
}

func TestServer_handleRun(t *testing.T) {
	server := New(nil)
	server.SetChatProcessor(&mockChatProcessor{})

	fromAgentID := uuid.New()
	toAgentID := uuid.New()
	sessionID := uuid.New()

	// Create a test request.
	requestBody := schema.AgentRunRequest{
		AppName:     "test-agent",
		FromAgentID: fromAgentID,
		ToAgentID:   toAgentID,
		SessionID:   sessionID,
		Content: schema.Content{
			Role: "user",
			Parts: []schema.PartIncoming{
				{Text: "Hello, world!"},
			},
		},
		Streaming: false,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/run", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleRun(w, req)

	// Print response for debugging
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())

	if w.Code != http.StatusOK {
		t.Errorf("unexpected status: got %d, want %d", w.Code, http.StatusOK)
	}

	var events []*schema.LLMEvent
	if err := json.Unmarshal(w.Body.Bytes(), &events); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// The mock now returns events, so we should have at least one
	if len(events) == 0 {
		t.Fatalf("expected at least 1 event, got %d", len(events))
	}

	// Check the content of the first part
	if len(events[0].Parts) > 0 {
		// Since the JSON marshaling/unmarshaling converts the struct to a map,
		// we need to check the map directly
		partMap, ok := events[0].Parts[0].(map[string]interface{})
		if !ok {
			t.Fatalf("expected map[string]interface{}, got %T", events[0].Parts[0])
		}

		content, ok := partMap["content"].(string)
		if !ok {
			t.Fatalf("expected content to be a string, got %T", partMap["content"])
		}

		if content != "test response" {
			t.Errorf("unexpected response content: got '%s', want 'test response'", content)
		}
	}
}

func TestConvertContentToMessage(t *testing.T) {
	content := schema.Content{
		Role: "user",
		Parts: []schema.PartIncoming{
			{Text: "Hello, world!"},
		},
	}

	message := content.ToMessage()

	if message.Role != model.RoleUser {
		t.Errorf("expected role 'user', got '%s'", message.Role)
	}

	if message.Content != "Hello, world!" {
		t.Errorf("expected content 'Hello, world!', got '%s'", message.Content)
	}
}

func TestConvertContentToMessage_Func(t *testing.T) {
	content := schema.Content{
		Role: "assistant",
		Parts: []schema.PartIncoming{
			{
				FunctionCall: &schema.FunctionCall{
					Name: "test_function",
					Args: map[string]interface{}{
						"param1": "value1",
						"param2": 42,
					},
				},
			},
		},
	}

	message := content.ToMessage()

	if message.Role != model.RoleAssistant {
		t.Errorf("expected role 'assistant', got '%s'", message.Role)
	}

	if len(message.ToolCalls) != 1 {
		t.Errorf("expected 1 tool call, got %d", len(message.ToolCalls))
	}

	toolCall := message.ToolCalls[0]
	if toolCall.Type != "function" {
		t.Errorf("expected type 'function', got '%s'", toolCall.Type)
	}

	if toolCall.Function.Name != "test_function" {
		t.Errorf("expected function name 'test_function', got '%s'", toolCall.Function.Name)
	}
}

func TestConvertSessionToADKFormat(t *testing.T) {
	now := time.Now()

	userID := uuid.New()
	sessionID := uuid.New()

	sess := &session.Session{
		ID:        sessionID.String(),
		AppName:   "test-app",
		UserID:    userID.String(),
		CreatedAt: now,
		UpdatedAt: now,
		State:     map[string][]byte{"key1": []byte("value1")},
	}

	adkSession, err := schema.NewADKSession(sess)
	if err != nil {
		t.Errorf("convertSessionToADKFormat: %v", err)
	}

	if adkSession.ID != sessionID {
		t.Errorf("expected ID 'test-session-id', got '%s'", adkSession.ID)
	}

	if adkSession.AppName != "test-app" {
		t.Errorf("expected AppName 'test-app', got '%s'", adkSession.AppName)
	}

	if adkSession.AgentID != userID {
		t.Errorf("expected AgentID %s, got '%s'", userID, adkSession.AgentID)
	}

	if adkSession.CreateTime == 0 {
		t.Error("expected non-zero CreateTime")
	}

	if adkSession.LastUpdateTime == 0 {
		t.Error("expected non-zero LastUpdateTime")
	}

	if len(adkSession.State) != 1 {
		t.Errorf("expected 1 state entry, got %d", len(adkSession.State))
	}
}

// mockSessionService is a simple mock session service for testing.
type mockSessionService struct {
}

func (m *mockSessionService) CreateSession(ctx context.Context, key session.Key, state session.StateMap, options ...session.Option) (*session.Session, error) {
	now := time.Now()
	sess := &session.Session{
		ID:        "mock-session-id",
		AppName:   key.AppName,
		UserID:    key.UserID,
		CreatedAt: now,
		UpdatedAt: now,
		State:     state,
	}
	return sess, nil
}

func (m *mockSessionService) GetSession(ctx context.Context, key session.Key, options ...session.Option) (*session.Session, error) {
	return nil, nil
}

func (m *mockSessionService) ListSessions(ctx context.Context, userKey session.UserKey, options ...session.Option) ([]*session.Session, error) {
	return []*session.Session{}, nil
}

func (m *mockSessionService) DeleteSession(ctx context.Context, key session.Key, options ...session.Option) error {
	return nil
}

func (m *mockSessionService) UpdateAppState(ctx context.Context, appName string, state session.StateMap) error {
	return nil
}

func (m *mockSessionService) DeleteAppState(ctx context.Context, appName string, key string) error {
	return nil
}

func (m *mockSessionService) ListAppStates(ctx context.Context, appName string) (session.StateMap, error) {
	return session.StateMap{}, nil
}

func (m *mockSessionService) UpdateUserState(ctx context.Context, userKey session.UserKey, state session.StateMap) error {
	return nil
}

func (m *mockSessionService) ListUserStates(ctx context.Context, userKey session.UserKey) (session.StateMap, error) {
	return session.StateMap{}, nil
}

func (m *mockSessionService) DeleteUserState(ctx context.Context, userKey session.UserKey, key string) error {
	return nil
}

func (m *mockSessionService) AppendEvent(ctx context.Context, session *session.Session, event *event.Event, options ...session.Option) error {
	return nil
}

func (m *mockSessionService) Close() error {
	return nil
}

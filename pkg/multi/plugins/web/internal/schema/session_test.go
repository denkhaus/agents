package schema

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

func TestNewADKSession(t *testing.T) {
	agentID := uuid.New()
	sessionID := uuid.New()

	now := time.Now()
	baseSession := &session.Session{
		ID:        sessionID.String(),
		AppName:   "test-app",
		UserID:    agentID.String(),
		CreatedAt: now,
		UpdatedAt: now,
		State: map[string][]byte{
			"key1": []byte("value1"),
		},
	}

	adkSession, err := NewADKSession(baseSession)
	if err != nil {
		t.Fatalf("NewADKSession failed: %v", err)
	}

	if adkSession.AppName != "test-app" {
		t.Errorf("expected AppName 'test-app', got '%s'", adkSession.AppName)
	}

	if adkSession.AgentID != agentID {
		t.Errorf("expected AgentID %s, got '%s'", agentID, adkSession.AgentID)
	}

	if adkSession.ID != sessionID {
		t.Errorf("expected ID %s, got '%s'", sessionID, adkSession.ID)
	}

	if adkSession.CreateTime != now.Unix() {
		t.Errorf("expected CreateTime %d, got %d", now.Unix(), adkSession.CreateTime)
	}

	if adkSession.LastUpdateTime != now.Unix() {
		t.Errorf("expected LastUpdateTime %d, got %d", now.Unix(), adkSession.LastUpdateTime)
	}

	if len(adkSession.State) != 1 {
		t.Errorf("expected 1 state entry, got %d", len(adkSession.State))
	}

	if string(adkSession.State["key1"]) != "value1" {
		t.Errorf("expected state value 'value1', got '%s'", string(adkSession.State["key1"]))
	}
}

func TestNewADKSession_InvalidUUID(t *testing.T) {
	// Test with invalid agent ID
	baseSession := &session.Session{
		ID:        uuid.New().String(),
		AppName:   "test-app",
		UserID:    "invalid-uuid",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := NewADKSession(baseSession)
	if err == nil {
		t.Error("expected error for invalid agent ID")
	}

	// Test with invalid session ID
	baseSession.UserID = uuid.New().String()
	baseSession.ID = "invalid-uuid"

	_, err = NewADKSession(baseSession)
	if err == nil {
		t.Error("expected error for invalid session ID")
	}
}
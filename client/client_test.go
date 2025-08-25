package client

import (
	"context"
	"testing"
	"time"

	"github.com/denkhaus/agents/auth"
	"github.com/samber/mo"
)

func TestNewAPIClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		opts    []ClientOption
		wantErr bool
	}{
		{
			name:    "valid basic client",
			baseURL: "http://localhost:8080",
			opts:    []ClientOption{WithTimeout(30 * time.Second)},
			wantErr: false,
		},
		{
			name:    "empty base URL",
			baseURL: "",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "with streaming enabled",
			baseURL: "http://localhost:8080",
			opts:    []ClientOption{WithStreaming(true)},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewAPIClient(tt.baseURL, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAPIClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewAPIClient() returned nil client")
			}
		})
	}
}

func TestNewUnifiedClient(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ClientConfig{
				BaseURL: "http://localhost:8080",
				Timeout: mo.Some(30 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "empty base URL",
			config: ClientConfig{
				BaseURL: "",
			},
			wantErr: true,
		},
		{
			name: "with auth provider",
			config: ClientConfig{
				BaseURL: "http://localhost:8080",
				Auth: mo.Some(func() auth.Provider {
					provider, _ := auth.NewAPIKeyProvider(auth.APIKeyConfig{
						Key: "test-key",
					})
					return provider
				}()),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewUnifiedClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUnifiedClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewUnifiedClient() returned nil client")
			}
		})
	}
}

func TestClientOptions(t *testing.T) {
	config := ClientConfig{BaseURL: "http://localhost:8080"}

	// Test WithTimeout
	WithTimeout(30 * time.Second)(&config)
	if timeout, ok := config.Timeout.Get(); !ok || timeout != 30*time.Second {
		t.Error("WithTimeout option not applied correctly")
	}

	// Test WithStreaming
	WithStreaming(true)(&config)
	if streaming, ok := config.Streaming.Get(); !ok || !streaming {
		t.Error("WithStreaming option not applied correctly")
	}

	// Test WithUserAgent
	WithUserAgent("test-agent")(&config)
	if userAgent, ok := config.UserAgent.Get(); !ok || userAgent != "test-agent" {
		t.Error("WithUserAgent option not applied correctly")
	}

	// Test WithRetryCount
	WithRetryCount(5)(&config)
	if retryCount, ok := config.RetryCount.Get(); !ok || retryCount != 5 {
		t.Error("WithRetryCount option not applied correctly")
	}
}

func TestTaskTracker(t *testing.T) {
	tracker := NewTaskTracker()

	taskID := "test-task-123"

	// Test tracking a task
	tracker.TrackTask(taskID)
	status := tracker.GetTaskStatus(taskID)
	if status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", status)
	}

	// Test updating task status
	tracker.UpdateTaskStatus(taskID, "running")
	status = tracker.GetTaskStatus(taskID)
	if status != "running" {
		t.Errorf("Expected status 'running', got '%s'", status)
	}

	// Test unknown task
	unknownStatus := tracker.GetTaskStatus("unknown-task")
	if unknownStatus != "unknown" {
		t.Errorf("Expected status 'unknown', got '%s'", unknownStatus)
	}
}

func TestAgentCardFetcher(t *testing.T) {
	fetcher := newAgentCardFetcher(5 * time.Second)

	// Test with invalid URL (should fail gracefully)
	ctx := context.Background()
	_, err := fetcher.FetchAgentCard(ctx, "http://invalid-url-that-does-not-exist")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestJWTAuth(t *testing.T) {
	config := auth.JWTConfig{
		Secret:   []byte("test-secret"),
		Audience: "test-audience",
		Issuer:   "test-issuer",
		Expiry:   time.Hour,
	}

	provider, err := auth.NewJWTProvider(config)
	if err != nil {
		t.Fatalf("Failed to create JWT provider: %v", err)
	}

	if provider.Type() != "jwt" {
		t.Errorf("Expected type 'jwt', got '%s'", provider.Type())
	}

	// Test refresh
	ctx := context.Background()
	if refreshErr := provider.Refresh(ctx); refreshErr != nil {
		t.Errorf("JWT refresh failed: %v", refreshErr)
	}

	if !provider.IsValid() {
		t.Error("JWT provider should be valid after refresh")
	}
}

func TestAPIKeyAuth(t *testing.T) {
	config := auth.APIKeyConfig{
		Key:    "test-api-key",
		Header: mo.Some("X-Test-Key"),
	}

	provider, err := auth.NewAPIKeyProvider(config)
	if err != nil {
		t.Fatalf("Failed to create API key provider: %v", err)
	}

	if provider.Type() != "apikey" {
		t.Errorf("Expected type 'apikey', got '%s'", provider.Type())
	}

	if !provider.IsValid() {
		t.Error("API key provider should be valid with non-empty key")
	}

	// Test refresh (should be no-op)
	ctx := context.Background()
	if refreshErr := provider.Refresh(ctx); refreshErr != nil {
		t.Errorf("API key refresh failed: %v", refreshErr)
	}
}

func TestOAuth2Auth(t *testing.T) {
	config := auth.OAuth2Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     "http://localhost:8080/oauth2/token",
		Scopes:       []string{"read", "write"},
	}

	provider, err := auth.NewOAuth2Provider(config)
	if err != nil {
		t.Fatalf("Failed to create OAuth2 provider: %v", err)
	}

	if provider.Type() != "oauth2" {
		t.Errorf("Expected type 'oauth2', got '%s'", provider.Type())
	}

	// Initially should not be valid (no token)
	if provider.IsValid() {
		t.Error("OAuth2 provider should not be valid initially")
	}
}

func TestStreamingClientInterface(t *testing.T) {
	// Test that unified client implements StreamingClient interface
	config := ClientConfig{
		BaseURL:   "http://localhost:8080",
		Streaming: mo.Some(true),
	}

	client, err := NewUnifiedClient(config)
	if err != nil {
		t.Fatalf("Failed to create unified client: %v", err)
	}

	// Test interface assertion
	if _, ok := client.(StreamingClient); !ok {
		t.Error("Unified client should implement StreamingClient interface")
	}

	if _, ok := client.(APIClient); !ok {
		t.Error("Unified client should implement APIClient interface")
	}

	if _, ok := client.(AgentCardFetcher); !ok {
		t.Error("Unified client should implement AgentCardFetcher interface")
	}
}

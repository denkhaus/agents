package condenser

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

// mockModel implements a simple mock for model.Model
type mockModel struct{}

func (m *mockModel) GenerateContent(ctx context.Context, req *model.Request) (<-chan *model.Response, error) {
	respChan := make(chan *model.Response, 1)

	// Create a simple mock response
	resp := &model.Response{
		Choices: []model.Choice{
			{
				Message: model.Message{
					Role:    model.RoleAssistant,
					Content: "This is a test summary of the conversation.",
				},
			},
		},
	}

	respChan <- resp
	close(respChan)

	return respChan, nil
}

func (m *mockModel) Info() model.Info {
	return model.Info{
		Name: "mock-model",
	}
}

func TestNewCondenserService(t *testing.T) {
	// Create a mock session service
	sessionService := inmemory.NewSessionService()

	// Create a mock LLM
	mockLLM := &mockModel{}

	// Test creating condenser service with default config
	service, err := New(sessionService, mockLLM, DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create condenser service: %v", err)
	}

	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	// Test that we can get config and metrics
	config := service.GetConfig()
	if config.MaxContextTokens != 8000 {
		t.Errorf("Expected MaxContextTokens to be 8000, got %d", config.MaxContextTokens)
	}

	metrics := service.GetMetrics()
	if metrics.CondensationCount != 0 {
		t.Errorf("Expected initial CondensationCount to be 0, got %d", metrics.CondensationCount)
	}
}

func TestNewWithOptions(t *testing.T) {
	// Create a mock session service
	sessionService := inmemory.NewSessionService()

	// Create a mock LLM
	mockLLM := &mockModel{}

	// Test creating condenser service with options
	service, err := NewWithOptions(
		sessionService,
		mockLLM,
		WithMaxContextTokens(4000),
		WithTriggerThreshold(0.5),
		WithTokenCountingMethod(TokenCountingHeuristic),
	)
	if err != nil {
		t.Fatalf("Failed to create condenser service with options: %v", err)
	}

	config := service.GetConfig()
	if config.MaxContextTokens != 4000 {
		t.Errorf("Expected MaxContextTokens to be 4000, got %d", config.MaxContextTokens)
	}

	if config.TriggerThreshold != 0.5 {
		t.Errorf("Expected TriggerThreshold to be 0.5, got %f", config.TriggerThreshold)
	}

	if config.TokenCountingMethod != TokenCountingHeuristic {
		t.Errorf("Expected TokenCountingMethod to be TokenCountingHeuristic, got %v", config.TokenCountingMethod)
	}
}

func TestHeuristicTokenCounter(t *testing.T) {
	logger := zap.NewNop()
	counter := NewHeuristicTokenCounter(4.0, logger)

	ctx := context.Background()

	// Test empty string
	tokens, err := counter.CountTokens(ctx, "")
	if err != nil {
		t.Fatalf("Failed to count tokens for empty string: %v", err)
	}
	if tokens != 0 {
		t.Errorf("Expected 0 tokens for empty string, got %d", tokens)
	}

	// Test simple string
	tokens, err = counter.CountTokens(ctx, "hello")
	if err != nil {
		t.Fatalf("Failed to count tokens for 'hello': %v", err)
	}
	if tokens != 2 { // 5 chars / 4.0 chars per token = 1.25, rounded up to 2
		t.Errorf("Expected 2 tokens for 'hello', got %d", tokens)
	}

	// Test token counter info
	info := counter.GetInfo()
	if info.Method != TokenCountingHeuristic {
		t.Errorf("Expected method to be TokenCountingHeuristic, got %v", info.Method)
	}
	if info.Accuracy != AccuracyEstimated {
		t.Errorf("Expected accuracy to be AccuracyEstimated, got %v", info.Accuracy)
	}
}

func TestTokenCountingMethodString(t *testing.T) {
	tests := []struct {
		method   TokenCountingMethod
		expected string
	}{
		{TokenCountingAuto, "auto"},
		{TokenCountingHeuristic, "heuristic"},
		{TokenCountingLLMNative, "llm-native"},
		{TokenCountingTikToken, "tiktoken"},
		{TokenCountingCustom, "custom"},
	}

	for _, test := range tests {
		if test.method.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.method.String())
		}
	}
}

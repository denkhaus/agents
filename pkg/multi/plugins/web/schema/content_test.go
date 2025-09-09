package schema

import (
	"testing"

	"trpc.group/trpc-go/trpc-agent-go/model"
)

func TestContent_ToMessage(t *testing.T) {
	// Test with text part
	content := Content{
		Role: "user",
		Parts: []PartIncoming{
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

	if len(message.ToolCalls) != 0 {
		t.Errorf("expected 0 tool calls, got %d", len(message.ToolCalls))
	}
}

func TestContent_ToMessage_WithFunctionCall(t *testing.T) {
	// Test with function call part
	content := Content{
		Role: "assistant",
		Parts: []PartIncoming{
			{
				FunctionCall: &FunctionCall{
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

	if message.Content != "" {
		t.Errorf("expected empty content, got '%s'", message.Content)
	}

	if len(message.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(message.ToolCalls))
	}

	toolCall := message.ToolCalls[0]
	if toolCall.Type != "function" {
		t.Errorf("expected type 'function', got '%s'", toolCall.Type)
	}

	if toolCall.Function.Name != "test_function" {
		t.Errorf("expected function name 'test_function', got '%s'", toolCall.Function.Name)
	}
}

func TestContent_ToMessage_WithInlineData(t *testing.T) {
	// Test with inline data part
	content := Content{
		Role: "user",
		Parts: []PartIncoming{
			{
				InlineData: &InlineData{
					Data:        "test data",
					MimeType:    "image/png",
					DisplayName: "test.png",
				},
			},
		},
	}

	message := content.ToMessage()

	if message.Role != model.RoleUser {
		t.Errorf("expected role 'user', got '%s'", message.Role)
	}

	expectedContent := "[image: test.png (image/png)]"
	if message.Content != expectedContent {
		t.Errorf("expected content '%s', got '%s'", expectedContent, message.Content)
	}
}

func TestContent_ToMessage_WithFunctionResponse(t *testing.T) {
	// Test with function response part
	content := Content{
		Role: "tool",
		Parts: []PartIncoming{
			{
				FunctionResponse: &FunctionResponse{
					Name:     "test_function",
					Response: "test response",
					ID:       "test-id",
				},
			},
		},
	}

	message := content.ToMessage()

	if message.Role != model.RoleTool {
		t.Errorf("expected role 'tool', got '%s'", message.Role)
	}

	expectedContent := "[Function test_function responded: \"test response\"]"
	if message.Content != expectedContent {
		t.Errorf("expected content '%s', got '%s'", expectedContent, message.Content)
	}
}

func TestContent_ToMessage_MultipleParts(t *testing.T) {
	// Test with multiple parts
	content := Content{
		Role: "user",
		Parts: []PartIncoming{
			{Text: "Hello, world!"},
			{
				InlineData: &InlineData{
					Data:        "test data",
					MimeType:    "image/png",
					DisplayName: "test.png",
				},
			},
		},
	}

	message := content.ToMessage()

	if message.Role != model.RoleUser {
		t.Errorf("expected role 'user', got '%s'", message.Role)
	}

	expectedContent := "Hello, world!\n[image: test.png (image/png)]"
	if message.Content != expectedContent {
		t.Errorf("expected content '%s', got '%s'", expectedContent, message.Content)
	}
}
package schema

import (
	"testing"

	"github.com/denkhaus/agents/pkg/messaging"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

func TestNewLLMEvent(t *testing.T) {
	baseEvent := &event.Event{
		ID:           "test-event-id",
		InvocationID: "test-invocation-id",
		Author:       "test-author",
		Response: &model.Response{
			ID:      "test-response-id",
			Object:  "test-object",
			Created: 1234567890,
			Model:   "test-model",
			Choices: []model.Choice{
				{
					Index: 0,
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: "test response",
					},
				},
			},
			Done: true,
		},
	}

	llmEvent, err := messaging.NewLLMEvent(baseEvent, false)
	if err != nil {
		t.Fatalf("NewLLMEvent failed: %v", err)
	}

	if llmEvent == nil {
		t.Fatal("expected non-nil LLMEvent")
	}

	if llmEvent.Type == "" {
		t.Error("expected LLMEvent to have a non-empty type")
	}

	if llmEvent.ID != "test-response-id" {
		t.Errorf("expected ID 'test-response-id', got '%s'", llmEvent.ID)
	}

	if llmEvent.InvocationID != "test-invocation-id" {
		t.Errorf("expected InvocationID 'test-invocation-id', got '%s'", llmEvent.InvocationID)
	}

	if llmEvent.Author != "test-author" {
		t.Errorf("expected Author 'test-author', got '%s'", llmEvent.Author)
	}

	if llmEvent.Role != model.RoleAssistant {
		t.Errorf("expected Role 'assistant', got '%s'", llmEvent.Role)
	}

	if !llmEvent.Done {
		t.Error("expected Done to be true")
	}

	if llmEvent.Type != "test-object" {
		t.Errorf("expected Type 'test-object', got '%s'", llmEvent.Type)
	}

	if llmEvent.Created != 1234567890 {
		t.Errorf("expected Created 1234567890, got %d", llmEvent.Created)
	}

	if llmEvent.Model != "test-model" {
		t.Errorf("expected Model 'test-model', got '%s'", llmEvent.Model)
	}
}

func TestLLMEvent_DetermineEventRole(t *testing.T) {
	// Test with regular response
	baseEvent := &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Role: model.RoleAssistant,
					},
				},
			},
		},
	}

	llmEvent := &LLMEvent{base: baseEvent}
	role := llmEvent.determineEventRole()

	if role != model.RoleAssistant {
		t.Errorf("expected role 'assistant', got '%s'", role)
	}

	// Test with tool response
	baseEvent.Response.Object = model.ObjectTypeToolResponse
	role = llmEvent.determineEventRole()

	if role != model.RoleTool {
		t.Errorf("expected role 'tool', got '%s'", role)
	}

	// Test with nil response
	llmEvent.base.Response = nil
	role = llmEvent.determineEventRole()

	if role != "" {
		t.Errorf("expected empty role, got '%s'", role)
	}
}

func TestLLMEvent_BuildEventParts(t *testing.T) {
	// Test with text content
	baseEvent := &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Content: "test response",
					},
				},
			},
		},
	}

	llmEvent := &LLMEvent{base: baseEvent}
	parts := llmEvent.buildEventParts()

	if len(parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(parts))
	}

	textPart, ok := parts[0].(*TextPart)
	if !ok {
		t.Fatalf("expected TextPart, got %T", parts[0])
	}

	if textPart.Content != "test response" {
		t.Errorf("expected content 'test response', got '%s'", textPart.Content)
	}
}

func TestLLMEvent_BuildEventPartsWithToolCalls(t *testing.T) {
	// Test with tool calls
	baseEvent := &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						ToolCalls: []model.ToolCall{
							{
								ID:   "test-tool-id",
								Type: "function",
								Function: model.FunctionDefinitionParam{
									Name:      "test_function",
									Arguments: []byte(`{"param1": "value1"}`),
								},
							},
						},
					},
				},
			},
		},
	}

	llmEvent := &LLMEvent{base: baseEvent}
	parts := llmEvent.buildEventParts()

	if len(parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(parts))
	}

	functionCallPart, ok := parts[0].(*FunctionCallPart)
	if !ok {
		t.Fatalf("expected FunctionCallPart, got %T", parts[0])
	}

	if functionCallPart.Name != "test_function" {
		t.Errorf("expected function name 'test_function', got '%s'", functionCallPart.Name)
	}

	if functionCallPart.ID != "test-tool-id" {
		t.Errorf("expected ID 'test-tool-id', got '%s'", functionCallPart.ID)
	}
}

func TestLLMEvent_IsToolResponse(t *testing.T) {
	// Test with tool response
	baseEvent := &event.Event{
		Response: &model.Response{
			Object: model.ObjectTypeToolResponse,
		},
	}

	llmEvent := &LLMEvent{base: baseEvent}
	isToolResponse := llmEvent.isToolResponse()

	if !isToolResponse {
		t.Error("expected isToolResponse to be true")
	}

	// Test without tool response
	baseEvent.Response.Object = "chat.completion"
	isToolResponse = llmEvent.isToolResponse()

	if isToolResponse {
		t.Error("expected isToolResponse to be false")
	}

	// Test with nil response
	llmEvent.base.Response = nil
	isToolResponse = llmEvent.isToolResponse()

	if isToolResponse {
		t.Error("expected isToolResponse to be false")
	}
}

func TestLLMEvent_HasToolCalls(t *testing.T) {
	// Test with tool calls
	baseEvent := &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						ToolCalls: []model.ToolCall{
							{
								Type: "function",
							},
						},
					},
				},
			},
		},
	}

	llmEvent := &LLMEvent{base: baseEvent}
	hasToolCalls := llmEvent.hasToolCalls()

	if !hasToolCalls {
		t.Error("expected hasToolCalls to be true")
	}

	// Test without tool calls
	baseEvent.Response.Choices[0].Message.ToolCalls = nil
	hasToolCalls = llmEvent.hasToolCalls()

	if hasToolCalls {
		t.Error("expected hasToolCalls to be false")
	}
}

func TestLLMEvent_AddUsageMetadata(t *testing.T) {
	// Test with usage data
	baseEvent := &event.Event{
		Response: &model.Response{
			Usage: &model.Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		},
	}

	llmEvent := &LLMEvent{base: baseEvent}
	llmEvent.addUsageMetadata()

	if llmEvent.Usage == nil {
		t.Fatal("expected non-nil Usage")
	}

	if llmEvent.Usage.PromptTokenCount != 10 {
		t.Errorf("expected PromptTokenCount 10, got %d", llmEvent.Usage.PromptTokenCount)
	}

	if llmEvent.Usage.CandidatesTokenCount != 20 {
		t.Errorf("expected CandidatesTokenCount 20, got %d", llmEvent.Usage.CandidatesTokenCount)
	}

	if llmEvent.Usage.TotalTokenCount != 30 {
		t.Errorf("expected TotalTokenCount 30, got %d", llmEvent.Usage.TotalTokenCount)
	}

	// Test without usage data
	baseEvent.Response.Usage = nil
	llmEvent = &LLMEvent{base: baseEvent}
	llmEvent.addUsageMetadata()

	if llmEvent.Usage != nil {
		t.Error("expected nil Usage")
	}
}

func TestLLMEvent_AddResponseMetadata(t *testing.T) {
	// Test with response data
	baseEvent := &event.Event{
		Response: &model.Response{
			Done:   true,
			Object: "test-object",
			Model:  "test-model",
		},
	}

	llmEvent := &LLMEvent{base: baseEvent}
	llmEvent.addResponseMetadata()

	if !llmEvent.Done {
		t.Error("expected Done to be true")
	}

	if llmEvent.Type != "test-object" {
		t.Errorf("expected Type 'test-object', got '%s'", llmEvent.Type)
	}

	if llmEvent.Model != "test-model" {
		t.Errorf("expected Model 'test-model', got '%s'", llmEvent.Model)
	}

	// Test without response data
	baseEvent.Response = nil
	llmEvent = &LLMEvent{base: baseEvent}
	llmEvent.addResponseMetadata()

	if llmEvent.Done {
		t.Error("expected Done to be false")
	}
}

func TestEventType_Validate(t *testing.T) {
	// Test valid event types
	validTypes := []EventType{EventTypeAssistant, EventTypeToolCall, EventTypeToolResponse, EventTypeReasoning}

	for _, eventType := range validTypes {
		if !eventType.Validate() {
			t.Errorf("expected event type '%s' to be valid", eventType)
		}
	}

	// Test that any other specific type is also valid
	customType := EventType("some-other-type")
	if !customType.Validate() {
		t.Errorf("expected custom event type '%s' to be valid", customType)
	}

	// Test that empty type is invalid
	emptyType := EventType("")
	if emptyType.Validate() {
		t.Error("expected empty event type to be invalid")
	}
}

func TestLLMEvent_Validate(t *testing.T) {
	// Test valid event types
	validTypes := []EventType{EventTypeAssistant, EventTypeToolCall, EventTypeToolResponse, EventTypeReasoning}

	for _, eventType := range validTypes {
		llmEvent := &LLMEvent{
			Type: eventType,
		}

		if err := llmEvent.Validate(); err != nil {
			t.Errorf("expected event with type '%s' to be valid, got error: %v", eventType, err)
		}
	}

	// Test that any other specific type is also valid
	llmEvent := &LLMEvent{
		Type: "some-other-type",
	}

	if err := llmEvent.Validate(); err != nil {
		t.Errorf("expected event with custom type to be valid, got error: %v", err)
	}

	// Test that empty type is invalid
	llmEventEmpty := &LLMEvent{
		Type: "",
	}

	if err := llmEventEmpty.Validate(); err == nil {
		t.Error("expected event with empty type to be invalid")
	}
}

func TestLLMEvent_DetectReasoning(t *testing.T) {
	llmEvent := &LLMEvent{}

	// Test content with reasoning indicators
	reasoningContent := "Let me think through this step by step. /REASONING/ First, I need to analyze the problem. /*PLANNING*/ Then I'll create a plan."
	if !llmEvent.detectReasoning(reasoningContent) {
		t.Error("expected reasoning to be detected")
	}

	// Test content without reasoning indicators
	normalContent := "This is a normal response without any reasoning indicators."
	if llmEvent.detectReasoning(normalContent) {
		t.Error("expected reasoning to not be detected")
	}
}

func TestLLMEvent_ExtractWithReasoning(t *testing.T) {
	// Test with reasoning content
	baseEvent := &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Content: "Let me think through this step by step. /REASONING/ First, I need to analyze the problem.",
						Role:    model.RoleAssistant,
					},
				},
			},
			Done: true,
		},
	}

	llmEvent := &LLMEvent{base: baseEvent}
	result, err := llmEvent.extract(false)
	if err != nil {
		t.Fatalf("extract failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil LLMEvent")
	}

	if result.Type != EventTypeReasoning {
		t.Errorf("expected Type to be EventTypeReasoning, got '%s'", result.Type)
	}

	// Test with normal content
	baseEvent2 := &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Content: "This is a normal response.",
						Role:    model.RoleAssistant,
					},
				},
			},
			Done: true,
		},
	}

	llmEvent2 := &LLMEvent{base: baseEvent2}
	result2, err := llmEvent2.extract(false)
	if err != nil {
		t.Fatalf("extract failed: %v", err)
	}

	if result2 == nil {
		t.Fatal("expected non-nil LLMEvent")
	}

	if result2.Type != EventTypeAssistant {
		t.Errorf("expected Type to be EventTypeAssistant for normal content, got '%s'", result2.Type)
	}

	// Test with tool call content
	baseEvent3 := &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Content: "I need to call a tool.",
						Role:    model.RoleAssistant,
						ToolCalls: []model.ToolCall{
							{
								ID:   "test-tool-id",
								Type: "function",
								Function: model.FunctionDefinitionParam{
									Name:      "test_function",
									Arguments: []byte(`{"param": "value"}`),
								},
							},
						},
					},
				},
			},
			Done: true,
		},
	}

	llmEvent3 := &LLMEvent{base: baseEvent3}
	result3, err := llmEvent3.extract(false)
	if err != nil {
		t.Fatalf("extract failed: %v", err)
	}

	if result3 == nil {
		t.Fatal("expected non-nil LLMEvent")
	}

	if result3.Type != EventTypeToolCall {
		t.Errorf("expected Type to be EventTypeToolCall for tool call content, got '%s'", result3.Type)
	}

	// Test with tool response content
	baseEvent4 := &event.Event{
		Response: &model.Response{
			Object: model.ObjectTypeToolResponse,
			Choices: []model.Choice{
				{
					Message: model.Message{
						Content: "Tool response content",
						Role:    model.RoleTool,
					},
				},
			},
			Done: true,
		},
	}

	llmEvent4 := &LLMEvent{base: baseEvent4}
	result4, err := llmEvent4.extract(false)
	if err != nil {
		t.Fatalf("extract failed: %v", err)
	}

	if result4 == nil {
		t.Fatal("expected non-nil LLMEvent")
	}

	if result4.Type != EventTypeToolResponse {
		t.Errorf("expected Type to be EventTypeToolResponse for tool response content, got '%s'", result4.Type)
	}
}

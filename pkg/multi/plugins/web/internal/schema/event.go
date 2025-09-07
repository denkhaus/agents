package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type EventType string

const (
	EventTypeAssistant    EventType = "assistant"
	EventTypeToolCall     EventType = "tool.call"
	EventTypeToolResponse EventType = "tool.response"
	EventTypeReasoning    EventType = "reasoning"
	//EventTypeError        EventType = "error"
	//EventTypeSystem       EventType = "system"
)

// validate checks if the event type is valid
func (et EventType) Validate() bool {
	switch et {
	case EventTypeAssistant, EventTypeToolCall, EventTypeToolResponse, EventTypeReasoning:
		return true
	default:
		// Any other specific type is also valid as long as it's not empty
		return et != ""
	}
}

func NewLLMEvent(base *event.Event, isStreaming bool) (*LLMEvent, error) {
	ev := &LLMEvent{
		base: base,
	}

	return ev.extract(isStreaming)
}

func (p *LLMEvent) extract(isStreaming bool) (*LLMEvent, error) {
	p.ID = p.eventID()
	p.InvocationID = p.base.InvocationID
	p.Timestamp = p.base.Timestamp.Unix()
	p.Author = p.base.Author

	p.Role = p.determineEventRole()
	// Build parts.
	parts := p.buildEventParts()
	// Filter parts based on streaming mode.
	parts, isToolEvent := p.filterEventParts(parts, isStreaming)
	// Skip event if no meaningful parts, unless it's a tool-related event
	if len(parts) == 0 && !isToolEvent {
		return nil, errors.New("empty event")
	}

	p.Parts = parts

	// Detect reasoning in the content
	for _, part := range parts {
		if textPart, ok := part.(*TextPart); ok {
			if p.detectReasoning(textPart.Content) {
				p.Type = EventTypeReasoning
				break
			}
		}
	}

	p.addUsageMetadata()
	p.addResponseMetadata()

	// Assign default type if none was set
	if p.Type == "" {
		// If it's a tool event but wasn't categorized, set appropriate type
		if p.isToolResponse() {
			p.Type = EventTypeToolResponse
		} else if p.hasToolCalls() {
			p.Type = EventTypeToolCall
		} else {
			// Default to assistant for regular messages
			p.Type = EventTypeAssistant
		}
	}

	// Validate the LLMEvent
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate LLMEven: %w", err)
	}

	return p, nil
}

// determineEventRole determines the role for the event content.
func (p *LLMEvent) determineEventRole() model.Role {
	var role model.Role
	if p.base.Response != nil {
		if p.base.Response.Object == model.ObjectTypeToolResponse {
			role = model.RoleTool
		} else if len(p.base.Response.Choices) > 0 {
			role = p.base.Response.Choices[0].Message.Role
		}
	}

	return role
}

func (p *LLMEvent) addUsageMetadata() {
	if p.base.Usage == nil {
		return
	}

	p.Usage = &UsageMetaData{
		PromptTokenCount:     p.base.Usage.PromptTokens,
		CandidatesTokenCount: p.base.Usage.CompletionTokens,
		TotalTokenCount:      p.base.Usage.TotalTokens,
	}
}

func (p *LLMEvent) detectReasoning(content string) bool {
	// Check for React planner tags that indicate reasoning/planning content
	if strings.Contains(content, "/PLANNING/") ||
		strings.Contains(content, "/REASONING/") ||
		strings.Contains(content, "/REPLANNING/") ||
		strings.Contains(content, "/*PLANNING*/") ||
		strings.Contains(content, "/*REASONING*/") ||
		strings.Contains(content, "/*REPLANNING*/") {
		return true
	}

	// Check for other reasoning indicators
	if strings.Contains(content, "/ACTION/") ||
		strings.Contains(content, "/*ACTION*/") {
		return true
	}

	return false
}

func (p *LLMEvent) addResponseMetadata() {
	if p.base.Response == nil {
		return
	}

	p.Done = p.base.Response.Done
	p.Partial = p.base.Response.IsPartial

	// Ensure partial flag is correctly set for streaming
	if p.base.Response.IsPartial {
		p.Partial = true
		p.Done = false
	} else if p.base.Response.Done {
		p.Partial = false
		p.Done = true
	}

	if p.base.Response.Object != "" {
		p.Type = EventType(p.base.Response.Object)
	}
	if p.base.Response.Created != 0 {
		p.Created = p.base.Response.Created
	}
	if p.base.Response.Model != "" {
		p.Model = p.base.Response.Model
	}
}

// eventID returns the canonical identifier for an event.
// If the underlying model.Response already contains a non-empty ID we
// prefer it; otherwise we fall back to the envelope‐level event ID.
func (p *LLMEvent) eventID() string {
	if p.base.Response != nil && p.base.Response.ID != "" {
		return p.base.Response.ID
	}
	return p.base.ID
}

// isToolResponse reports whether the supplied event represents a tool
// response produced by the LLM flow.
func (p *LLMEvent) isToolResponse() bool {
	return p.base.Response != nil && p.base.Response.Object == model.ObjectTypeToolResponse
}

func (p *LLMEvent) hasToolCalls() bool {
	if len(p.base.Response.Choices) > 0 && len(p.base.Response.Choices[0].Message.ToolCalls) > 0 {
		return true
	}
	return false
}

// buildFunctionCallPart converts a model.ToolCall into the ADK Web part map.
// The returned map follows the schema expected by the Web UI.
func (p *LLMEvent) buildFunctionCallPart(tc model.ToolCall) Part {
	var args interface{}
	if err := json.Unmarshal(tc.Function.Arguments, &args); err != nil {
		// Preserve raw string if not valid JSON.
		args = map[string]interface{}{"raw": string(tc.Function.Arguments)}
	}

	return &FunctionCallPart{
		Name: tc.Function.Name,
		Args: args,
		ID:   tc.ID,
	}
}

// buildFunctionResponsePart builds a single functionResponse part.
// respObj can be either a structured object (decoded JSON) or the original
// raw string when JSON decoding fails. The name field is currently unknown
// from the upstream payload, so we intentionally leave it blank.
func (p *LLMEvent) buildFunctionResponsePart(respObj interface{}, id string, name string) Part {
	return &FunctionResponsePart{
		Name: name,
		Args: respObj,
		ID:   id,
	}
}

func (p *LLMEvent) buildEventParts() []Part {
	var parts []Part

	if p.base.Response == nil {
		return parts
	}

	// Handle normal / streaming assistant or model messages.
	for _, choice := range p.base.Response.Choices {
		// Regular text (full message).
		if choice.Message.Content != "" {
			// For tool response events, we do NOT include the raw JSON string as a
			// separate text part, otherwise the ADK Web UI will render duplicated
			// information (both as plain text and as function_response). Keeping
			// only the structured function_response part provides a cleaner view.
			if p.base.Response.Object != model.ObjectTypeToolResponse {
				parts = append(parts, &TextPart{Content: choice.Message.Content})
			}
		}

		// Tool calls in full message.
		for _, tc := range choice.Message.ToolCalls {
			parts = append(parts, p.buildFunctionCallPart(tc))
		}

		// Streaming delta text.
		if choice.Delta.Content != "" {
			parts = append(parts, &TextPart{Content: choice.Delta.Content})
		}
		// Tool calls in streaming delta.
		for _, tc := range choice.Delta.ToolCalls {
			parts = append(parts, p.buildFunctionCallPart(tc))
		}
	}

	// Tool response events.
	if p.base.Response.Object == model.ObjectTypeToolResponse {
		for _, choice := range p.base.Response.Choices {
			var respObj interface{}
			if choice.Message.Content != "" {
				if err := json.Unmarshal([]byte(choice.Message.Content), &respObj); err != nil {
					respObj = choice.Message.Content // raw string fallback
				}
			}

			parts = append(parts, p.buildFunctionResponsePart(respObj, choice.Message.ToolID, choice.Message.ToolName))
		}
	}

	return parts
}

// filterEventParts filters parts based on streaming mode and event type.
func (p *LLMEvent) filterEventParts(parts []Part, isStreaming bool) ([]Part, bool) {
	if p.base.Response == nil {
		return parts, false
	}

	// Always include tool calls and tool responses regardless of streaming mode
	toolResp := p.isToolResponse()
	hasToolCall := p.hasToolCalls()

	// Set object type for tool calls and responses only if not already set
	if p.Type == "" {
		if hasToolCall {
			p.Type = EventTypeToolCall
		} else if toolResp {
			p.Type = EventTypeToolResponse
		}
	}

	if toolResp || hasToolCall {
		return parts, true
	}

	if isStreaming {
		// In streaming mode, include all partial events and the final done event
		// Don't drop the final event as it may contain important completion info
		return parts, false
	} else {
		// Non-streaming endpoint should include final assistant messages
		if !p.base.Response.Done {
			return nil, false
		}
	}

	return parts, false
}

// validate checks if the LLMEvent is valid
func (p *LLMEvent) Validate() error {
	// Event must have a non-empty type
	if p.Type == "" {
		return fmt.Errorf("event must have a non-empty type")
	}

	// Validate the event type
	if !p.Type.Validate() {
		return fmt.Errorf("invalid event type: %s", p.Type)
	}

	return nil
}

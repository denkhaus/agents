package web

// ---------------------------------------------------------------------
// Internal helpers for event conversion --------------------------------
// ---------------------------------------------------------------------

// ADK Web payload JSON keys. Keeping them as constants helps avoid
// typographical errors and makes refactoring easier.
// const (
// 	keyText             = "text"             // Plain textual content part.
// 	keyFunctionCall     = "functionCall"     // Function call part key.
// 	keyFunctionResponse = "functionResponse" // Function response part key.
// )

// // eventID returns the canonical identifier for an event.
// // If the underlying model.Response already contains a non-empty ID we
// // prefer it; otherwise we fall back to the envelope‐level event ID.
// func eventID(e *event.Event) string {
// 	if e.Response != nil && e.Response.ID != "" {
// 		return e.Response.ID
// 	}
// 	return e.ID
// }

// // isToolResponse reports whether the supplied event represents a tool
// // response produced by the LLM flow.
// func isToolResponse(e *event.Event) bool {
// 	return e.Response != nil && e.Response.Object == model.ObjectTypeToolResponse
// }

// // buildFunctionCallPart converts a model.ToolCall into the ADK Web part map.
// // The returned map follows the schema expected by the Web UI.
// func buildFunctionCallPart(tc model.ToolCall) map[string]interface{} {
// 	var args interface{}
// 	if err := json.Unmarshal(tc.Function.Arguments, &args); err != nil {
// 		// Preserve raw string if not valid JSON.
// 		args = map[string]interface{}{"raw": string(tc.Function.Arguments)}
// 	}
// 	return map[string]interface{}{
// 		keyFunctionCall: map[string]interface{}{
// 			"name": tc.Function.Name,
// 			"args": args,
// 			"id":   tc.ID,
// 		},
// 	}
// }

// // convertContentToMessage converts Google GenAI Content to trpc-agent model.Message
// func convertContentToMessage(content schema.Content) model.Message {
// 	log.Debugf("convertContentToMessage: role=%s parts=%+v", content.Role, content.Parts)
// 	var textParts []string
// 	var toolCalls []model.ToolCall
// 	for _, part := range content.Parts {
// 		if part.Text != "" {
// 			textParts = append(textParts, part.Text)
// 		}

// 		if part.FunctionCall != nil {
// 			argsBytes, _ := json.Marshal(part.FunctionCall.Args)
// 			toolCall := model.ToolCall{
// 				Type: "function",
// 				Function: model.FunctionDefinitionParam{
// 					Name:      part.FunctionCall.Name,
// 					Arguments: argsBytes,
// 				},
// 			}
// 			toolCalls = append(toolCalls, toolCall)
// 		}

// 		if part.InlineData != nil {
// 			dataType := "file"
// 			if part.InlineData.MimeType != "" {
// 				if strings.HasPrefix(part.InlineData.MimeType, "image") {
// 					dataType = "image"
// 				} else if strings.HasPrefix(part.InlineData.MimeType, "audio") {
// 					dataType = "audio"
// 				} else if strings.HasPrefix(part.InlineData.MimeType, "video") {
// 					dataType = "video"
// 				}
// 			}
// 			fileName := part.InlineData.DisplayName
// 			if fileName == "" {
// 				fileName = "attachment"
// 			}
// 			attachmentText := fmt.Sprintf("[%s: %s (%s)]", dataType, fileName, part.InlineData.MimeType)
// 			textParts = append(textParts, attachmentText)
// 		}

// 		if part.FunctionResponse != nil {
// 			responseJSON, _ := json.Marshal(part.FunctionResponse.Response)
// 			responseText := fmt.Sprintf("[Function %s responded: %s]", part.FunctionResponse.Name, string(responseJSON))
// 			textParts = append(textParts, responseText)
// 		}
// 	}
// 	var combinedText string
// 	if len(textParts) > 0 {
// 		combinedText = strings.Join(textParts, "\n")
// 	}
// 	msg := model.Message{
// 		Role:    model.Role(content.Role),
// 		Content: combinedText,
// 	}

// 	if len(toolCalls) > 0 {
// 		msg.ToolCalls = toolCalls
// 	}
// 	return msg
// }

// // buildADKEventEnvelope creates the basic ADK event envelope.
// func buildADKEventEnvelope(e *event.Event) map[string]interface{} {
// 	id := eventID(e)
// 	return map[string]interface{}{
// 		"invocationId": e.InvocationID,
// 		"author":       e.Author,
// 		"actions": map[string]interface{}{
// 			"stateDelta":           map[string]interface{}{},
// 			"artifactDelta":        map[string]interface{}{},
// 			"requestedAuthConfigs": map[string]interface{}{},
// 		},
// 		"id":        id,
// 		"timestamp": e.Timestamp.Unix(),
// 	}
// }

// determineEventRole determines the role for the event content.
// func determineEventRole(e *event.Event) string {
// 	role := e.Author // fallback
// 	if e.Response != nil {
// 		if e.Response.Object == model.ObjectTypeToolResponse {
// 			role = string(model.RoleTool)
// 		} else if len(e.Response.Choices) > 0 {
// 			role = string(e.Response.Choices[0].Message.Role)
// 		}
// 	}
// 	return role
// }

// // buildEventParts constructs the parts array for the event content.
// func buildEventParts(e *event.Event) []map[string]interface{} {
// 	var parts []map[string]interface{}

// 	if e.Response == nil {
// 		return parts
// 	}

// 	// Handle normal / streaming assistant or model messages.
// 	for _, choice := range e.Response.Choices {
// 		// Regular text (full message).
// 		if choice.Message.Content != "" {
// 			// For tool response events, we do NOT include the raw JSON string as a
// 			// separate text part, otherwise the ADK Web UI will render duplicated
// 			// information (both as plain text and as function_response). Keeping
// 			// only the structured function_response part provides a cleaner view.
// 			if e.Response.Object != model.ObjectTypeToolResponse {
// 				parts = append(parts, map[string]interface{}{keyText: choice.Message.Content})
// 			}
// 		}

// 		// Tool calls in full message.
// 		for _, tc := range choice.Message.ToolCalls {
// 			parts = append(parts, buildFunctionCallPart(tc))
// 		}

// 		// Streaming delta text.
// 		if choice.Delta.Content != "" {
// 			parts = append(parts, map[string]interface{}{keyText: choice.Delta.Content})
// 		}
// 		// Tool calls in streaming delta.
// 		for _, tc := range choice.Delta.ToolCalls {
// 			parts = append(parts, buildFunctionCallPart(tc))
// 		}
// 	}

// 	// Tool response events.
// 	if e.Response.Object == model.ObjectTypeToolResponse {
// 		for _, choice := range e.Response.Choices {
// 			var respObj interface{}
// 			if choice.Message.Content != "" {
// 				if err := json.Unmarshal([]byte(choice.Message.Content), &respObj); err != nil {
// 					respObj = choice.Message.Content // raw string fallback
// 				}
// 			}
// 			parts = append(parts, buildFunctionResponsePart(respObj, choice.Message.ToolID, choice.Message.ToolName))
// 		}
// 	}

// 	return parts
// }

// // filterEventParts filters parts based on streaming mode and event type.
// func filterEventParts(e *event.Event, parts []map[string]interface{}, isStreaming bool) []map[string]interface{} {
// 	if e.Response == nil {
// 		return parts
// 	}

// 	// Always include tool calls and tool responses regardless of streaming mode
// 	toolResp := isToolResponse(e)
// 	hasToolCall := false
// 	if len(e.Response.Choices) > 0 && len(e.Response.Choices[0].Message.ToolCalls) > 0 {
// 		hasToolCall = true
// 	}

// 	if toolResp || hasToolCall {
// 		return parts
// 	}

// 	if isStreaming {
// 		// In streaming mode, include all partial events and the final done event
// 		// Don't drop the final event as it may contain important completion info
// 		return parts
// 	} else {
// 		// Non-streaming endpoint should include final assistant messages
// 		if !e.Response.Done {
// 			return nil
// 		}
// 	}

// 	return parts
// }

// // addResponseMetadata adds response-level metadata to the ADK event.
// func addResponseMetadata(adkEvent map[string]interface{}, e *event.Event) {
// 	if e.Response == nil {
// 		return
// 	}

// 	adkEvent["done"] = e.Response.Done
// 	adkEvent["partial"] = e.Response.IsPartial

// 	// Ensure partial flag is correctly set for streaming
// 	if e.Response.IsPartial {
// 		adkEvent["partial"] = true
// 		adkEvent["done"] = false
// 	} else if e.Response.Done {
// 		adkEvent["partial"] = false
// 		adkEvent["done"] = true
// 	}

// 	if e.Response.Object != "" {
// 		adkEvent["object"] = e.Response.Object
// 	}
// 	if e.Response.Created != 0 {
// 		adkEvent["created"] = e.Response.Created
// 	}
// 	if e.Response.Model != "" {
// 		adkEvent["model"] = e.Response.Model
// 	}
// }

// // addUsageMetadata adds usage metadata to the ADK event.
// func addUsageMetadata(adkEvent map[string]interface{}, e *event.Event) {
// 	if e.Usage == nil {
// 		return
// 	}

// 	adkEvent["usageMetadata"] = map[string]interface{}{
// 		"promptTokenCount":     e.Usage.PromptTokens,
// 		"candidatesTokenCount": e.Usage.CompletionTokens,
// 		"totalTokenCount":      e.Usage.TotalTokens,
// 	}
// }

// // convertEventToADKFormat converts trpc-agent Event to ADK Web UI expected
// // format. The isStreaming flag indicates whether the UI is currently
// // displaying token-level streaming (true) or expecting a single complete
// // response (false).
// func convertEventToADKFormat(e *event.Event, isStreaming bool) map[string]interface{} {
// 	// Build basic envelope.
// 	adkEvent := buildADKEventEnvelope(e)

// 	// Determine role and build content.
// 	role := determineEventRole(e)
// 	content := map[string]interface{}{
// 		"role": role,
// 	}

// 	// Build parts.
// 	parts := buildEventParts(e)

// 	// Filter parts based on streaming mode.
// 	parts = filterEventParts(e, parts, isStreaming)

// 	// For tool calls and tool responses, always include even if no text parts
// 	toolResp := isToolResponse(e)
// 	hasToolCall := false
// 	if e.Response != nil && len(e.Response.Choices) > 0 && len(e.Response.Choices[0].Message.ToolCalls) > 0 {
// 		hasToolCall = true
// 	}

// 	// Skip event if no meaningful parts, unless it's a tool-related event
// 	if len(parts) == 0 && !toolResp && !hasToolCall {
// 		return nil
// 	}

// 	// Set object type for tool calls and responses
// 	if hasToolCall {
// 		adkEvent["object"] = "tool_call"
// 	} else if toolResp {
// 		adkEvent["object"] = "tool_response"
// 	}

// 	content["parts"] = parts
// 	adkEvent["content"] = content

// 	// Add metadata.
// 	addResponseMetadata(adkEvent, e)
// 	addUsageMetadata(adkEvent, e)

// 	return adkEvent
// }

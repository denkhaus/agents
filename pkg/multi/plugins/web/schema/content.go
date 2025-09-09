package schema

import (
	"encoding/json"
	"fmt"
	"strings"

	"trpc.group/trpc-go/trpc-a2a-go/log"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// Content matches the GenAI content contract used by ADK Web.
type Content struct {
	Role  string         `json:"role"`
	Parts []PartIncoming `json:"parts"`
}

// convertContentToMessage converts Google GenAI Content to trpc-agent model.Message
func (p *Content) ToMessage() model.Message {
	log.Debugf("convertContentToMessage: role=%s parts=%+v", p.Role, p.Parts)
	var textParts []string
	var toolCalls []model.ToolCall
	for _, part := range p.Parts {
		if part.Text != "" {
			textParts = append(textParts, part.Text)
		}

		if part.FunctionCall != nil {
			argsBytes, _ := json.Marshal(part.FunctionCall.Args)
			toolCall := model.ToolCall{
				Type: "function",
				Function: model.FunctionDefinitionParam{
					Name:      part.FunctionCall.Name,
					Arguments: argsBytes,
				},
			}
			toolCalls = append(toolCalls, toolCall)
		}

		if part.InlineData != nil {
			dataType := "file"
			if part.InlineData.MimeType != "" {
				if strings.HasPrefix(part.InlineData.MimeType, "image") {
					dataType = "image"
				} else if strings.HasPrefix(part.InlineData.MimeType, "audio") {
					dataType = "audio"
				} else if strings.HasPrefix(part.InlineData.MimeType, "video") {
					dataType = "video"
				}
			}
			fileName := part.InlineData.DisplayName
			if fileName == "" {
				fileName = "attachment"
			}
			attachmentText := fmt.Sprintf("[%s: %s (%s)]", dataType, fileName, part.InlineData.MimeType)
			textParts = append(textParts, attachmentText)
		}

		if part.FunctionResponse != nil {
			responseJSON, _ := json.Marshal(part.FunctionResponse.Response)
			responseText := fmt.Sprintf("[Function %s responded: %s]", part.FunctionResponse.Name, string(responseJSON))
			textParts = append(textParts, responseText)
		}
	}
	var combinedText string
	if len(textParts) > 0 {
		combinedText = strings.Join(textParts, "\n")
	}

	msg := model.Message{
		Role:    model.Role(p.Role),
		Content: combinedText,
	}

	if len(toolCalls) > 0 {
		msg.ToolCalls = toolCalls
	}
	return msg
}

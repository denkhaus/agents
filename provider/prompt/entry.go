package prompt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/denkhaus/agents/utils"
)

func (p *promptEntry) GetDescription() string {
	return p.metadata.Description
}

func (p *promptEntry) GetName() string {
	return p.metadata.Name
}

func (p *promptEntry) GetGlobalInstruction() string {
	return p.metadata.GlobalInstruction
}

func (p *promptEntry) GetInstruction(data interface{}) (string, error) {
	processedData, result, err := utils.ValidateJSON(data, p.schema)
	if err != nil {
		return "", &PromptError{
			Message: "prompt data validation failed",
			AgentID: p.metadata.AgentID,
			Err:     err,
		}
	}
	if !result.Valid() {
		var validationErrors []string
		for _, desc := range result.Errors() {
			validationErrors = append(validationErrors, fmt.Sprintf("- %s", desc))
		}
		return "", &PromptError{
			Message: "prompt data validation failed: " + strings.Join(validationErrors, "; "),
			AgentID: p.metadata.AgentID,
		}
	}

	var buf bytes.Buffer
	if err := p.template.Execute(&buf, processedData); err != nil {
		return "", &PromptError{
			Message: "failed to render prompt",
			AgentID: p.metadata.AgentID,
			Err:     err,
		}
	}
	return buf.String(), nil
}

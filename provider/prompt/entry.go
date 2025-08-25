package prompt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
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
	dataLoader := gojsonschema.NewGoLoader(data)
	result, err := p.schema.Validate(dataLoader)
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
	if err := p.template.Execute(&buf, data); err != nil {
		return "", &PromptError{
			Message: "failed to render prompt",
			AgentID: p.metadata.AgentID,
			Err:     err,
		}
	}
	return buf.String(), nil
}

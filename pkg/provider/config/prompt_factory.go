package config

import (
	"bytes"
	"fmt"
	"html/template"
)

func (f *UnifiedAgentFactory) renderPromptContent(templateContent string, data interface{}) (string, error) {
	tmpl, err := template.New("prompt").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return buf.String(), nil
}

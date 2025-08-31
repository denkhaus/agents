package prompt

import (
	"text/template"

	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
)

// PromptMetadata represents the metadata extracted from the Markdown front matter.
type PromptMetadata struct {
	Name              string     `yaml:"name"`
	GlobalInstruction string     `yaml:"global_instruction"`
	Description       string     `yaml:"description"`
	AgentID           uuid.UUID  `yaml:"agent_id"`
	Schema            JSONSchema `yaml:"schema"`
}

// promptEntry holds a compiled template and its associated metadata and JSON schema.
type promptEntry struct {
	template *template.Template
	metadata PromptMetadata
	schema   *gojsonschema.Schema
}

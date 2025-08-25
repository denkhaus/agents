package prompt

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/denkhaus/agents/provider"
	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

//go:embed templates/*.md
var promptFS embed.FS

// NewPromptManager creates a new instance of PromptManager.
// It takes an embed.FS for loading prompt templates.
func NewPromptManager(fsys embed.FS, rootPath string) (provider.PromptManager, error) {
	prompts := make(map[uuid.UUID]*promptEntry)

	// Use fs.WalkDir to correctly traverse the embedded directory
	err := fs.WalkDir(fsys, rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := fsys.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded prompt file %s: %w", path, err)
		}

		parts := strings.SplitN(string(content), "---", 3)
		if len(parts) < 3 {
			return fmt.Errorf("invalid Markdown front matter in %s", path)
		}

		var metadata PromptMetadata
		if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
			return fmt.Errorf("failed to parse YAML front matter in %s: %w", path, err)
		}

		if metadata.Name == "" {
			return fmt.Errorf("prompt name cannot be empty in %s", path)
		}

		if metadata.AgentID == uuid.Nil {
			return fmt.Errorf("agent ID cannot be empty in %s", path)
		}

		templateContent := strings.TrimSpace(parts[2])
		tpl, err := template.New(metadata.Name).Parse(templateContent)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", metadata.Name, err)
		}

		schemaLoader := gojsonschema.NewGoLoader(metadata.Schema)
		schema, err := gojsonschema.NewSchema(schemaLoader)
		if err != nil {
			return fmt.Errorf("failed to compile schema for [%s]-[%s]: %w", metadata.Name, metadata.AgentID, err)
		}

		if _, exists := prompts[metadata.AgentID]; exists {
			return fmt.Errorf("duplicate agent id on prompt %s: prompt with agent id %s already exists", metadata.Name, metadata.AgentID)
		}

		prompts[metadata.AgentID] = &promptEntry{
			template: tpl,
			metadata: metadata,
			schema:   schema,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return NewManager(prompts), nil
}

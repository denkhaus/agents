package settings

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

//go:embed templates/*.yaml
var SettingsFS embed.FS

// SettingsManager defines the interface for managing agent settings.
type SettingsManager interface {
	GetSettings(agentID uuid.UUID) (*Settings, error)
}

// settingsManagerImpl is an unexported implementation of SettingsManager.
type settingsManagerImpl struct {
	settings map[uuid.UUID]*Settings
}

// NewSettingsManager creates a new instance of SettingsManager.
// It takes an embed.FS for loading settings templates.
func NewSettingsManager(fsys embed.FS, rootPath string) (SettingsManager, error) {
	settings := make(map[uuid.UUID]*Settings)

	// Use fs.WalkDir to correctly traverse the embedded directory
	err := fs.WalkDir(fsys, rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}

		content, err := fsys.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded settings file %s: %w", path, err)
		}

		// Parse the YAML content
		var settingsData struct {
			Settings
		}

		if err := yaml.Unmarshal(content, &settingsData.Settings); err != nil {
			return fmt.Errorf("failed to parse YAML settings in %s: %w", path, err)
		}

		if settingsData.AgentID == uuid.Nil {
			return fmt.Errorf("agent ID cannot be empty in %s", path)
		}

		// Validate the agent role
		if err := settingsData.Agent.Role.Validate(); err != nil {
			return fmt.Errorf("invalid agent role in %s: %w", path, err)
		}

		if _, exists := settings[settingsData.AgentID]; exists {
			return fmt.Errorf("duplicate agent id in settings %s: settings with agent id %s already exists", path, settingsData.AgentID)
		}

		// Create Settings struct from parsed data
		settingsStruct := &Settings{
			Model: settingsData.Model,
			Agent: settingsData.Agent,
		}

		settings[settingsData.AgentID] = settingsStruct

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &settingsManagerImpl{settings: settings}, nil
}

func (sm *settingsManagerImpl) GetSettings(agentID uuid.UUID) (*Settings, error) {
	settings, ok := sm.settings[agentID]
	if !ok {
		return nil, fmt.Errorf("settings not found for agent ID: %s", agentID)
	}

	return settings, nil
}

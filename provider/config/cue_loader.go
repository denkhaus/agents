package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

// cueConfigProviderImpl loads agent configurations from CUE files
type cueConfigProviderImpl struct {
	ctx        *cue.Context
	configPath string
}

// NewCUEConfigProvider creates a new CUE configuration provider
func NewCUEConfigProvider(configPath string) ConfigProvider {
	return &cueConfigProviderImpl{
		ctx:        cuecontext.New(),
		configPath: configPath,
	}
}

// LoadAgentComposition loads a complete agent configuration from environment
func (p *cueConfigProviderImpl) LoadAgentComposition(environment, agentName string) (*AgentConfig, error) {
	// Load environment-specific composition
	envPath := filepath.Join(p.configPath, "compositions", "environments", fmt.Sprintf("%s.cue", environment))

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("environment configuration not found: %s", environment)
	}

	instances := load.Instances([]string{envPath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return nil, fmt.Errorf("no CUE instances found for environment: %s", environment)
	}

	values, err := p.ctx.BuildInstances(instances)
	if err != nil {
		return nil, fmt.Errorf("failed to build instances: %w", err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no CUE instances built")
	}
	value := values[0]
	if value.Err() != nil {
		return nil, fmt.Errorf("failed to build CUE instance: %w", value.Err())
	}

	// Extract specific agent configuration
	agentPath := fmt.Sprintf("%s.agents.%s", environment, agentName)
	agentValue := value.LookupPath(cue.ParsePath(agentPath))

	if !agentValue.Exists() {
		return nil, fmt.Errorf("agent %s not found in environment %s", agentName, environment)
	}

	// Decode the agent configuration
	var config AgentConfig
	if err := agentValue.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode agent config: %w", err)
	}

	// Resolve environment variables in the configuration
	if err := p.resolveEnvironmentVariables(&config); err != nil {
		return nil, fmt.Errorf("failed to resolve environment variables: %w", err)
	}

	return &config, nil
}

// LoadPrompt loads a specific prompt configuration
func (p *cueConfigProviderImpl) LoadPrompt(agentName, version string) (*PromptConfig, error) {
	promptPath := filepath.Join(p.configPath, "prompts", version, fmt.Sprintf("%s.cue", agentName))

	instances := load.Instances([]string{promptPath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return nil, fmt.Errorf("no CUE instances found for prompt: %s/%s", version, agentName)
	}

	values, err := p.ctx.BuildInstances(instances)
	if err != nil {
		return nil, fmt.Errorf("failed to build instances: %w", err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no prompt CUE instances built")
	}
	value := values[0]
	if value.Err() != nil {
		return nil, fmt.Errorf("failed to build prompt CUE instance: %w", value.Err())
	}

	// Extract prompt configuration
	promptValue := value.LookupPath(cue.ParsePath(agentName))
	if !promptValue.Exists() {
		return nil, fmt.Errorf("prompt %s not found", agentName)
	}

	var prompt PromptConfig
	if err := promptValue.Decode(&prompt); err != nil {
		return nil, fmt.Errorf("failed to decode prompt config: %w", err)
	}

	return &prompt, nil
}

// LoadSettings loads agent settings
func (p *cueConfigProviderImpl) LoadSettings(agentName, profile string) (*SettingsConfig, error) {
	settingsPath := filepath.Join(p.configPath, "settings", profile, fmt.Sprintf("%s.cue", agentName))

	instances := load.Instances([]string{settingsPath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return nil, fmt.Errorf("no CUE instances found for settings: %s/%s", profile, agentName)
	}

	values, err := p.ctx.BuildInstances(instances)
	if err != nil {
		return nil, fmt.Errorf("failed to build instances: %w", err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no settings CUE instances built")
	}
	value := values[0]
	if value.Err() != nil {
		return nil, fmt.Errorf("failed to build settings CUE instance: %w", value.Err())
	}

	// Extract settings configuration
	settingsValue := value.LookupPath(cue.ParsePath(agentName))
	if !settingsValue.Exists() {
		return nil, fmt.Errorf("settings %s not found", agentName)
	}

	var settings SettingsConfig
	if err := settingsValue.Decode(&settings); err != nil {
		return nil, fmt.Errorf("failed to decode settings config: %w", err)
	}

	return &settings, nil
}

// LoadToolProfile loads tool profile configuration
func (p *cueConfigProviderImpl) LoadToolProfile(profileName string) (*ToolsConfig, error) {
	toolsPath := filepath.Join(p.configPath, "tools", "profiles", fmt.Sprintf("%s.cue", profileName))

	instances := load.Instances([]string{toolsPath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return nil, fmt.Errorf("no CUE instances found for tool profile: %s", profileName)
	}

	values, err := p.ctx.BuildInstances(instances)
	if err != nil {
		return nil, fmt.Errorf("failed to build instances: %w", err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no tools CUE instances built")
	}
	value := values[0]
	if value.Err() != nil {
		return nil, fmt.Errorf("failed to build tools CUE instance: %w", value.Err())
	}

	// Extract tools configuration
	toolsValue := value.LookupPath(cue.ParsePath(profileName))
	if !toolsValue.Exists() {
		return nil, fmt.Errorf("tool profile %s not found", profileName)
	}

	var tools ToolsConfig
	if err := toolsValue.Decode(&tools); err != nil {
		return nil, fmt.Errorf("failed to decode tools config: %w", err)
	}

	// Resolve environment variables in tool configurations
	if err := p.resolveToolEnvironmentVariables(&tools); err != nil {
		return nil, fmt.Errorf("failed to resolve tool environment variables: %w", err)
	}

	return &tools, nil
}

// ValidateConfiguration validates all CUE configurations
func (p *cueConfigProviderImpl) ValidateConfiguration() error {
	// Load all CUE files and validate them
	instances := load.Instances([]string{"./..."}, &load.Config{
		Dir: p.configPath,
	})

	values, err := p.ctx.BuildInstances(instances)
	if err != nil {
		return fmt.Errorf("failed to build instances: %w", err)
	}

	for _, instance := range values {
		if instance.Err() != nil {
			return fmt.Errorf("CUE validation failed: %w", instance.Err())
		}
	}

	return nil
}

// resolveEnvironmentVariables resolves environment variables in the agent configuration
func (p *cueConfigProviderImpl) resolveEnvironmentVariables(config *AgentConfig) error {
	return p.resolveToolEnvironmentVariables(&config.Tools)
}

// resolveToolEnvironmentVariables resolves environment variables in tool configurations
func (p *cueConfigProviderImpl) resolveToolEnvironmentVariables(tools *ToolsConfig) error {
	// Resolve environment variables in tool configs
	for toolName, toolConfig := range tools.Tools {
		if err := p.resolveConfigMap(toolConfig.Config); err != nil {
			return fmt.Errorf("failed to resolve environment variables for tool %s: %w", toolName, err)
		}
	}

	// Resolve environment variables in toolset configs
	for toolSetName, toolSetConfig := range tools.ToolSets {
		if err := p.resolveConfigMap(toolSetConfig.Config); err != nil {
			return fmt.Errorf("failed to resolve environment variables for toolset %s: %w", toolSetName, err)
		}
	}

	return nil
}

// resolveConfigMap resolves environment variables in a configuration map
func (p *cueConfigProviderImpl) resolveConfigMap(config map[string]interface{}) error {
	for key, value := range config {
		resolved, err := p.resolveValue(value)
		if err != nil {
			return fmt.Errorf("failed to resolve config key %s: %w", key, err)
		}
		config[key] = resolved
	}
	return nil
}

// resolveValue resolves environment variables in a single value
func (p *cueConfigProviderImpl) resolveValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return p.resolveStringValue(v)
	case map[string]interface{}:
		return p.resolveMapValue(v)
	case []interface{}:
		return p.resolveSliceValue(v)
	default:
		return value, nil
	}
}

// resolveStringValue resolves environment variables in string values
func (p *cueConfigProviderImpl) resolveStringValue(value string) (interface{}, error) {
	if !strings.HasPrefix(value, "env:") {
		return value, nil
	}

	// Parse env:VAR_NAME or env:VAR_NAME:default
	parts := strings.SplitN(value, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid environment variable format: %s", value)
	}

	envVarName := parts[1]
	if envVarName == "" {
		return nil, fmt.Errorf("empty environment variable name")
	}

	envValue := os.Getenv(envVarName)

	// If no value and we have a default
	if envValue == "" && len(parts) == 3 {
		return parts[2], nil // Return default value
	}

	// If no value and no default
	if envValue == "" {
		return nil, fmt.Errorf("environment variable %s is not set", envVarName)
	}

	return envValue, nil
}

// resolveMapValue recursively resolves environment variables in a map
func (p *cueConfigProviderImpl) resolveMapValue(value map[string]interface{}) (interface{}, error) {
	resolved := make(map[string]interface{})
	for key, val := range value {
		resolvedVal, err := p.resolveValue(val)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve map key %s: %w", key, err)
		}
		resolved[key] = resolvedVal
	}
	return resolved, nil
}

// resolveSliceValue recursively resolves environment variables in a slice
func (p *cueConfigProviderImpl) resolveSliceValue(value []interface{}) (interface{}, error) {
	resolved := make([]interface{}, len(value))
	for i, val := range value {
		resolvedVal, err := p.resolveValue(val)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve slice index %d: %w", i, err)
		}
		resolved[i] = resolvedVal
	}
	return resolved, nil
}

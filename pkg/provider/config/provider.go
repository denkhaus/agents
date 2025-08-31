package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/config"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"github.com/samber/do"
	"go.uber.org/zap"
)

// cueConfigProviderImpl loads agent configurations from CUE files
type cueConfigProviderImpl struct {
	ctx              *cue.Context
	appConfigService config.Service
	configPath       string
}

// NewCUEConfigProvider creates a new CUE configuration provider
func NewCUEConfigProvider(injector *do.Injector) (ConfigProvider, error) {
	appConfigService := do.MustInvoke[config.Service](injector)
	configPath, err := appConfigService.GetConfigPath()
	if err != nil {
		return nil, err
	}

	return &cueConfigProviderImpl{
		appConfigService: appConfigService,
		ctx:              cuecontext.New(),
		configPath:       configPath,
	}, nil
}

// LoadAgentComposition loads a complete agent configuration from environment
func (p *cueConfigProviderImpl) LoadAgentComposition(environment string, agentRole shared.AgentRole) (*AgentConfig, error) {
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

	// Iterate through agents in the environment to find the one with the matching role
	agentsValue := value.LookupPath(cue.ParsePath(fmt.Sprintf("%s.agents", environment)))
	if !agentsValue.Exists() {
		return nil, fmt.Errorf("no agents defined in environment %s", environment)
	}

	iter, err := agentsValue.Fields()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over agents in environment %s: %w", environment, err)
	}

	var foundAgentValue cue.Value

	for iter.Next() {
		agentName := iter.Selector().String() // Keep agentName for logging/error messages if needed
		currentAgentValue := iter.Value()

		var basicConfig struct {
			Role shared.AgentRole `json:"role"`
		}
		if err := currentAgentValue.Decode(&basicConfig); err != nil {
			logger.Log.Warn("Failed to decode basic agent config for role lookup", zap.String("agentName", agentName), zap.Error(err))
			continue
		}

		if basicConfig.Role == agentRole {
			foundAgentValue = currentAgentValue
			// foundAgentName = agentName // Removed unused variable assignment
			break
		}
	}

	if !foundAgentValue.Exists() {
		return nil, fmt.Errorf("agent with role %s not found in environment %s", agentRole, environment)
	}

	// Now decode the full agent configuration
	var basicConfig struct {
		AgentID     uuid.UUID        `json:"agent_id"`
		Name        string           `json:"name"`
		Description string           `json:"description,omitempty"`
		Type        string           `json:"type"`
		Prompt      struct {
			Source cue.Value `json:"source"`
		}`json:"prompt"`
		Setting struct {
			Source cue.Value `json:"source"`
		}`json:"setting"`
		Tool struct {
			Source cue.Value `json:"source"`
		}`json:"tool"`
	}

	// Decode the basic configuration
	if err := foundAgentValue.Decode(&basicConfig); err != nil {
		return nil, fmt.Errorf("failed to decode basic agent config for role %s: %w", agentRole, err)
	}

	// Explicitly get the role value
	roleValue := foundAgentValue.LookupPath(cue.ParsePath("role"))
	if !roleValue.Exists() {
		return nil, fmt.Errorf("role field not found for agent %s", agentRole)
	}
	var decodedRole shared.AgentRole
	if err := roleValue.Decode(&decodedRole); err != nil {
		return nil, fmt.Errorf("failed to decode role for agent %s: %w", agentRole, err)
	}

	// Create the final config with resolved references
	config := &AgentConfig{
		AgentID:     basicConfig.AgentID,
		Name:        basicConfig.Name,
		Role:        decodedRole, // Use the explicitly decoded role
		Description: basicConfig.Description,
		Type:        shared.AgentType(basicConfig.Type),
	}

	// Resolve prompt reference
	if basicConfig.Prompt.Source.Exists() {
		promptConfig, err := p.LoadPrompt(decodedRole) // Use decodedRole for prompt lookup
		if err != nil {
			return nil, fmt.Errorf("failed to load prompt for agent role %s: %w", decodedRole, err)
		}
		config.Prompt = *promptConfig
	}

	// Resolve settings reference
	if basicConfig.Setting.Source.Exists() {
		settingsConfig, err := p.LoadSettings(decodedRole) // Use decodedRole for settings lookup
		if err != nil {
			return nil, fmt.Errorf("failed to load settings for agent role %s: %w", decodedRole, err)
		}
		config.Setting = *settingsConfig
	}

	// Resolve tools reference
	if basicConfig.Tool.Source.Exists() {
		toolsConfig, err := p.LoadToolProfile(decodedRole) // Use decodedRole for tool profile lookup
		if err != nil {
			return nil, fmt.Errorf("failed to load tools for agent role %s: %w", decodedRole, err)
		}
		config.Tool = *toolsConfig
	}

	// Resolve environment variables in the configuration
	if err := p.resolveEnvironmentVariables(config); err != nil {
		return nil, fmt.Errorf("failed to resolve environment variables: %w", err)
	}

	return config, nil
}

// LoadPrompt loads a specific prompt configuration
func (p *cueConfigProviderImpl) LoadPrompt(agentRole shared.AgentRole) (*PromptConfig, error) {
	promptPath := filepath.Join(p.configPath, "prompts", fmt.Sprintf("%s.cue", agentRole))

	instances := load.Instances([]string{promptPath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return nil, fmt.Errorf("no CUE instances found for prompt: %s", agentRole)
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
	promptValue := value.LookupPath(cue.ParsePath(string(agentRole)))
	if !promptValue.Exists() {
		return nil, fmt.Errorf("prompt %s not found in file", agentRole)
	}

	var prompt PromptConfig
	if err := promptValue.Decode(&prompt); err != nil {
		return nil, fmt.Errorf("failed to decode prompt config: %w", err)
	}

	return &prompt, nil
}

// LoadSettings loads agent settings
func (p *cueConfigProviderImpl) LoadSettings(agentRole shared.AgentRole) (*SettingsConfig, error) {
	settingsPath := filepath.Join(p.configPath, "settings", fmt.Sprintf("%s.cue", agentRole))

	instances := load.Instances([]string{settingsPath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return nil, fmt.Errorf("no CUE instances found for settings: %s", agentRole)
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
	settingsValue := value.LookupPath(cue.ParsePath(string(agentRole)))
	if !settingsValue.Exists() {
		return nil, fmt.Errorf("settings %s not found in file", agentRole)
	}

	var settings SettingsConfig
	if err := settingsValue.Decode(&settings); err != nil {
		return nil, fmt.Errorf("failed to decode settings config: %w", err)
	}

	return &settings, nil
}

// LoadToolProfile loads tool profile configuration
func (p *cueConfigProviderImpl) LoadToolProfile(agentRole shared.AgentRole) (*ToolsConfig, error) {
	toolsPath := filepath.Join(p.configPath, "tools", fmt.Sprintf("%s.cue", agentRole))

	instances := load.Instances([]string{toolsPath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return nil, fmt.Errorf("no CUE instances found for tool profile: %s", agentRole)
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
	toolsValue := value.LookupPath(cue.ParsePath(string(agentRole)))
	if !toolsValue.Exists() {
		return nil, fmt.Errorf("tool profile %s not found in file", agentRole)
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
	return p.resolveToolEnvironmentVariables(&config.Tool)
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

// getAllEnvironments returns a list of all available environment names
func (p *cueConfigProviderImpl) getAllEnvironments() ([]string, error) {
	envDir := filepath.Join(p.configPath, "compositions", "environments")
	files, err := os.ReadDir(envDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read environments directory %s: %w", envDir, err)
	}

	var environments []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".cue") {
			envName := strings.TrimSuffix(file.Name(), ".cue")
			environments = append(environments, envName)
		}
	}
	return environments, nil
}

// GetAgentInfoByID retrieves agent information by its UUID.
func (p *cueConfigProviderImpl) GetAgentInfoByID(agentID uuid.UUID) (*shared.AgentInfo, error) {
	environments, err := p.getAllEnvironments()
	if err != nil {
		return nil, fmt.Errorf("failed to get all environments: %w", err)
	}

	for _, env := range environments {
		agents, err := p.GetAgentsInEnvironment(env)
		if err != nil {
			// Log the error but continue to other environments
			logger.Log.Warn("Failed to get agents in environment", zap.String("environment", env), zap.Error(err))
			continue
		}

		for _, agentInfo := range agents {
			if agentInfo.ID() == agentID {
				return agentInfo, nil
			}
		}
	}

	return nil, fmt.Errorf("agent with ID %s not found in any environment", agentID)
}

func (p *cueConfigProviderImpl) loadSubAgentInfo(environment string, agentConfig *AgentConfig) (map[uuid.UUID]*shared.AgentInfo, error) {
	subAgents := agentConfig.Setting.Agent.SubAgents
	result := make(map[uuid.UUID]*shared.AgentInfo)

	if len(subAgents) == 0 {
		return result, nil
	}

	for _, subAgentRole := range subAgents { // Changed from subAgentID to subAgentRole
		subAgentConfig, err := p.LoadAgentComposition(environment, subAgentRole) // Use role for lookup
		if err != nil {
			return nil, fmt.Errorf("failed to load agent composition for sub-agent role %q: %w", subAgentRole, err)
		}

		result[subAgentConfig.AgentID] = subAgentConfig.ToAgentInfo()
	}

	return result, nil
}

// GetAgentsInEnvironment retrieves information about all agents defined within a specific environment.
func (p *cueConfigProviderImpl) GetAgentsInEnvironment(environment string) ([]*shared.AgentInfo, error) {

	// Load the environment-specific composition file
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
		return nil, fmt.Errorf("failed to build instances for environment %s: %w", environment, err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no CUE instances built for environment %s", environment)
	}
	value := values[0]
	if value.Err() != nil {
		return nil, fmt.Errorf("failed to build CUE instance for environment %s: %w", environment, value.Err())
	}

	// Extract the 'agents' field from the environment composition
	agentsValue := value.LookupPath(cue.ParsePath(fmt.Sprintf("%s.agents", environment)))
	if !agentsValue.Exists() {
		return nil, fmt.Errorf("no agents defined in environment %s", environment)
	}

	// Iterate over the agents defined in the environment
	iter, err := agentsValue.Fields()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over agents in environment %s: %w", environment, err)
	}

	agentMap := make(map[uuid.UUID]*shared.AgentInfo)
	for iter.Next() {

		agentName := iter.Selector().String() // Keep agentName for logging/error messages if needed
		currentAgentValue := iter.Value()

		var basicConfig struct {
			Role shared.AgentRole `json:"role"`
		}
		if err := currentAgentValue.Decode(&basicConfig); err != nil {
			logger.Log.Warn("Failed to decode basic agent config for role lookup in GetAgentsInEnvironment", zap.String("agentName", agentName), zap.Error(err))
			continue
		}

		agentConfig, err := p.LoadAgentComposition(environment, basicConfig.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to load agent composition for role %s in environment %s: %w", basicConfig.Role, environment, err)
		}

		info, err := p.loadSubAgentInfo(environment, agentConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load subagent info for agent %s: %w", agentName, err)
		}

		for agentID, agentInfo := range info {
			if _, exists := agentMap[agentID]; !exists {
				agentMap[agentID] = agentInfo
			}
		}

		if _, exists := agentMap[agentConfig.AgentID]; !exists {
			agentMap[agentConfig.AgentID] = agentConfig.ToAgentInfo()
		}
	}

	agentsInfo := make([]*shared.AgentInfo, 0, len(agentMap))
	for _, info := range agentMap {
		agentsInfo = append(agentsInfo, info)
	}

	return agentsInfo, nil
}

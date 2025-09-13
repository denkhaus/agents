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

// simplifiedCueConfigProvider loads agent configurations from simplified CUE files
type simplifiedCueConfigProvider struct {
	ctx              *cue.Context
	appConfigService config.Service
	configPath       string
}

// NewCUEConfigProvider creates a new simplified CUE configuration provider
func NewCUEConfigProvider(injector *do.Injector) (ConfigProvider, error) {
	appConfigService := do.MustInvoke[config.Service](injector)
	configPath, err := appConfigService.GetConfigPath()
	if err != nil {
		return nil, err
	}

	return &simplifiedCueConfigProvider{
		appConfigService: appConfigService,
		ctx:              cuecontext.New(),
		configPath:       configPath,
	}, nil
}

// LoadSystemConfig loads the complete system configuration
func (p *simplifiedCueConfigProvider) LoadSystemConfig() (*SystemConfig, error) {
	systemPath := filepath.Join(p.configPath, "system.cue")

	if _, err := os.Stat(systemPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("system configuration file not found: %s", systemPath)
	}

	instances := load.Instances([]string{systemPath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return nil, fmt.Errorf("no CUE instances found for system config")
	}

	values, err := p.ctx.BuildInstances(instances)
	if err != nil {
		return nil, fmt.Errorf("failed to build CUE instances for system config: %w", err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no CUE instances built for system config")
	}
	value := values[0]
	if value.Err() != nil {
		return nil, fmt.Errorf("CUE validation error in system config: %w", value.Err())
	}

	// Extract system configuration
	systemValue := value.LookupPath(cue.ParsePath("system"))
	if !systemValue.Exists() {
		return nil, fmt.Errorf("system configuration not found in system.cue")
	}

	var systemConfig SystemConfig
	if err := systemValue.Decode(&systemConfig); err != nil {
		return nil, fmt.Errorf("failed to decode system config: %w", err)
	}

	// Resolve environment variables in the configuration
	if err := p.resolveSystemEnvironmentVariables(&systemConfig); err != nil {
		return nil, fmt.Errorf("failed to resolve environment variables: %w", err)
	}

	return &systemConfig, nil
}

// LoadEnvironmentConfig loads a specific environment configuration
func (p *simplifiedCueConfigProvider) LoadEnvironmentConfig(envName EnvironmentName) (*EnvironmentConfig, error) {
	systemConfig, err := p.LoadSystemConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load system config: %w", err)
	}

	envConfig, exists := systemConfig.Environments[envName]
	if !exists {
		return nil, fmt.Errorf("environment %s not found in system config", envName)
	}

	return &envConfig, nil
}

// ResolveAgentByRole resolves an agent configuration by role within an environment
func (p *simplifiedCueConfigProvider) ResolveAgentByRole(envName EnvironmentName, role shared.AgentRole) (*AgentConfig, error) {
	envConfig, err := p.LoadEnvironmentConfig(envName)
	if err != nil {
		return nil, fmt.Errorf("failed to load environment config: %w", err)
	}

	// Resolve role to agent name
	agentName, exists := envConfig.Roles[role]
	if !exists {
		return nil, fmt.Errorf("role %s not mapped in environment %s", role, envName)
	}

	// Get agent config
	agentConfig, exists := envConfig.Agents[agentName]
	if !exists {
		return nil, fmt.Errorf("agent %s not found in environment %s", agentName, envName)
	}

	return &agentConfig, nil
}

// LoadAgentComposition loads a complete agent configuration from environment (legacy method)
func (p *simplifiedCueConfigProvider) LoadAgentComposition(environment EnvironmentName, agentRole shared.AgentRole) (*AgentConfig, error) {
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
		return nil, fmt.Errorf("failed to build CUE instances for environment '%s': %w", environment, err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no CUE instances built for environment '%s'", environment)
	}
	value := values[0]
	if value.Err() != nil {
		return nil, fmt.Errorf("CUE validation error in environment '%s': %w", environment, value.Err())
	}

	// Find the agent with the matching role in the environment
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
		agentName := iter.Selector().String()
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
			break
		}
	}

	if !foundAgentValue.Exists() {
		return nil, fmt.Errorf("agent with role %s not found in environment %s", agentRole, environment)
	}

	// Decode the complete agent configuration directly from the merged CUE value
	var config AgentConfig
	if err := foundAgentValue.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode agent config for role %s: %w", agentRole, err)
	}

	// Resolve environment variables in the configuration
	if err := p.resolveEnvironmentVariables(&config); err != nil {
		return nil, fmt.Errorf("failed to resolve environment variables: %w", err)
	}

	return &config, nil
}

// LoadPrompt loads a specific prompt configuration (kept for backward compatibility)
func (p *simplifiedCueConfigProvider) LoadPrompt(agentRole shared.AgentRole) (*PromptConfig, error) {
	promptPath := filepath.Join(p.configPath, "prompts", fmt.Sprintf("%s.cue", agentRole))

	if _, err := os.Stat(promptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("prompt configuration file not found for role '%s': %s", agentRole, promptPath)
	}

	var prompt PromptConfig
	err := p.loadAndDecode(promptPath, string(agentRole), &prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt configuration for role '%s' from file '%s': %w", agentRole, promptPath, err)
	}
	return &prompt, nil
}

// LoadSettings loads agent settings (kept for backward compatibility)
func (p *simplifiedCueConfigProvider) LoadSettings(agentRole shared.AgentRole) (*SettingsConfig, error) {
	settingsPath := filepath.Join(p.configPath, "settings", fmt.Sprintf("%s.cue", agentRole))

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("settings configuration file not found for role '%s': %s", agentRole, settingsPath)
	}

	var settings SettingsConfig
	err := p.loadAndDecode(settingsPath, string(agentRole), &settings)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings configuration for role '%s' from file '%s': %w", agentRole, settingsPath, err)
	}
	return &settings, nil
}

// LoadToolProfile loads tool profile configuration (kept for backward compatibility)
func (p *simplifiedCueConfigProvider) LoadToolProfile(agentRole shared.AgentRole) (*ToolsConfig, error) {
	toolsPath := filepath.Join(p.configPath, "tools", fmt.Sprintf("%s.cue", agentRole))

	if _, err := os.Stat(toolsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("tools configuration file not found for role '%s': %s", agentRole, toolsPath)
	}

	var tools ToolsConfig
	err := p.loadAndDecode(toolsPath, string(agentRole), &tools)
	if err != nil {
		return nil, fmt.Errorf("failed to load tools configuration for role '%s' from file '%s': %w", agentRole, toolsPath, err)
	}

	return &tools, nil
}

// ValidateConfiguration validates all CUE configurations
func (p *simplifiedCueConfigProvider) ValidateConfiguration() error {
	instances := load.Instances([]string{"./..."}, &load.Config{
		Dir: p.configPath,
	})

	values, err := p.ctx.BuildInstances(instances)
	if err != nil {
		return fmt.Errorf("failed to build CUE instances during validation in directory '%s': %w", p.configPath, err)
	}

	for i, instance := range values {
		if instance.Err() != nil {
			var fileName string
			if len(instances) > i && len(instances[i].Files) > 0 {
				fileName = instances[i].Files[0].Filename
			} else {
				fileName = fmt.Sprintf("instance_%d", i)
			}
			return fmt.Errorf("CUE validation failed in file '%s': %w", fileName, instance.Err())
		}
	}

	return nil
}

// GetAgentsInEnvironment retrieves information about all agents defined within a specific environment
func (p *simplifiedCueConfigProvider) GetAgentsInEnvironment(environment EnvironmentName, includeHuman bool) ([]*shared.AgentInfo, error) {
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

	agentsValue := value.LookupPath(cue.ParsePath(fmt.Sprintf("%s.agents", environment)))
	if !agentsValue.Exists() {
		return nil, fmt.Errorf("no agents defined in environment %s", environment)
	}

	iter, err := agentsValue.Fields()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over agents in environment %s: %w", environment, err)
	}

	agentMap := make(map[uuid.UUID]*shared.AgentInfo)
	for iter.Next() {
		agentName := iter.Selector().String()
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

		// Load sub-agent info
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

	if includeHuman {
		agentsInfo = append(agentsInfo, shared.AgentInfoHuman)
	}

	return agentsInfo, nil
}

// GetAgentInfoByID retrieves agent information by its UUID
func (p *simplifiedCueConfigProvider) GetAgentInfoByID(agentID uuid.UUID) (*shared.AgentInfo, error) {
	environments, err := p.getAllEnvironments()
	if err != nil {
		return nil, fmt.Errorf("failed to get all environments: %w", err)
	}

	for _, env := range environments {
		agents, err := p.GetAgentsInEnvironment(env, true)
		if err != nil {
			logger.Log.Warn("Failed to get agents in environment", zap.Any("environment", env), zap.Error(err))
			continue
		}

		for _, agentInfo := range agents {
			if agentInfo.ID == agentID {
				return agentInfo, nil
			}
		}
	}

	return nil, fmt.Errorf("agent with ID %s not found in any environment", agentID)
}

// Helper methods (simplified versions)

func (p *simplifiedCueConfigProvider) loadAndDecode(filePath, lookupPath string, target interface{}) error {
	instances := load.Instances([]string{filePath}, &load.Config{
		Dir: p.configPath,
	})

	if len(instances) == 0 {
		return fmt.Errorf("no CUE instances found for path: %s", filePath)
	}

	values, err := p.ctx.BuildInstances(instances)
	if err != nil {
		return fmt.Errorf("failed to build CUE instances for file '%s': %w", filePath, err)
	}
	if len(values) == 0 {
		return fmt.Errorf("no CUE instances built for file '%s'", filePath)
	}
	value := values[0]
	if value.Err() != nil {
		return fmt.Errorf("CUE validation error in file '%s': %w", filePath, value.Err())
	}

	configValue := value.LookupPath(cue.ParsePath(lookupPath))
	if !configValue.Exists() {
		return fmt.Errorf("path '%s' not found in file %s", lookupPath, filePath)
	}

	if err := configValue.Decode(target); err != nil {
		return fmt.Errorf("failed to decode config for path '%s' in file '%s': %w", lookupPath, filePath, err)
	}

	return nil
}

func (p *simplifiedCueConfigProvider) getAllEnvironments() ([]EnvironmentName, error) {
	envDir := filepath.Join(p.configPath, "compositions", "environments")
	files, err := os.ReadDir(envDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read environments directory %s: %w", envDir, err)
	}

	var environments []EnvironmentName
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".cue") {
			envName := strings.TrimSuffix(file.Name(), ".cue")
			environments = append(environments, EnvironmentName(envName))
		}
	}
	return environments, nil
}

func (p *simplifiedCueConfigProvider) loadSubAgentInfo(environment EnvironmentName, agentConfig *AgentConfig) (map[uuid.UUID]*shared.AgentInfo, error) {
	subAgents := agentConfig.Setting.Agent.SubAgents
	result := make(map[uuid.UUID]*shared.AgentInfo)

	if len(subAgents) == 0 {
		return result, nil
	}

	for _, subAgentRole := range subAgents {
		subAgentConfig, err := p.LoadAgentComposition(environment, subAgentRole)
		if err != nil {
			return nil, fmt.Errorf("failed to load agent composition for sub-agent role %q: %w", subAgentRole, err)
		}

		result[subAgentConfig.AgentID] = subAgentConfig.ToAgentInfo()
	}

	return result, nil
}

// Environment variable resolution (simplified)
func (p *simplifiedCueConfigProvider) resolveEnvironmentVariables(config *AgentConfig) error {
	return p.resolveToolEnvironmentVariables(&config.Tool)
}

func (p *simplifiedCueConfigProvider) resolveToolEnvironmentVariables(tools *ToolsConfig) error {
	for toolName, toolConfig := range tools.Tools {
		if err := p.resolveConfigMap(toolConfig.Config); err != nil {
			return fmt.Errorf("failed to resolve environment variables for tool %s: %w", toolName, err)
		}
	}

	for toolSetName, toolSetConfig := range tools.ToolSets {
		if err := p.resolveConfigMap(toolSetConfig.Config); err != nil {
			return fmt.Errorf("failed to resolve environment variables for toolset %s: %w", toolSetName, err)
		}
	}

	return nil
}

func (p *simplifiedCueConfigProvider) resolveConfigMap(config map[string]interface{}) error {
	for key, value := range config {
		resolved, err := p.resolveValue(value)
		if err != nil {
			return fmt.Errorf("failed to resolve config key %s: %w", key, err)
		}
		config[key] = resolved
	}
	return nil
}

func (p *simplifiedCueConfigProvider) resolveValue(value interface{}) (interface{}, error) {
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

func (p *simplifiedCueConfigProvider) resolveStringValue(value string) (interface{}, error) {
	if !strings.HasPrefix(value, "env:") {
		return value, nil
	}

	parts := strings.SplitN(value, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid environment variable format: %s", value)
	}

	envVarName := parts[1]
	if envVarName == "" {
		return nil, fmt.Errorf("empty environment variable name")
	}

	envValue := os.Getenv(envVarName)

	if envValue == "" && len(parts) == 3 {
		return parts[2], nil // Return default value
	}

	if envValue == "" {
		return nil, fmt.Errorf("environment variable %s is not set", envVarName)
	}

	return envValue, nil
}

func (p *simplifiedCueConfigProvider) resolveMapValue(value map[string]interface{}) (interface{}, error) {
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

func (p *simplifiedCueConfigProvider) resolveSliceValue(value []interface{}) (interface{}, error) {
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

// resolveSystemEnvironmentVariables resolves environment variables in the system configuration
func (p *simplifiedCueConfigProvider) resolveSystemEnvironmentVariables(systemConfig *SystemConfig) error {
	// Resolve environment variables in each environment
	for envName, envConfig := range systemConfig.Environments {
		// Resolve variables in agents
		for agentName, agentConfig := range envConfig.Agents {
			if err := p.resolveToolEnvironmentVariables(&agentConfig.Tool); err != nil {
				return fmt.Errorf("failed to resolve environment variables for agent %s in environment %s: %w", agentName, envName, err)
			}
			// Update the agent config back to the map
			envConfig.Agents[agentName] = agentConfig
		}
		// Update the environment config back to the map
		systemConfig.Environments[envName] = envConfig
	}
	return nil
}

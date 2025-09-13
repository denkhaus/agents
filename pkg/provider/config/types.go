package config

import (
	"context"

	"github.com/denkhaus/agents/pkg/session/condenser"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

type EnvironmentName string
type AgentName string

// AgentConfig represents a complete agent configuration
type AgentConfig struct {
	AgentID     uuid.UUID        `json:"agent_id"`
	Name        string           `json:"name"`
	Role        shared.AgentRole `json:"role"`
	Description string           `json:"description,omitempty"`
	Type        shared.AgentType `json:"type"`
	Prompt      PromptConfig     `json:"prompt"`
	Setting     SettingsConfig   `json:"setting"`
	Tool        ToolsConfig      `json:"tool"`
}

func (p *AgentConfig) ToAgentInfo() *shared.AgentInfo {
	agent := p.Setting.Agent

	agentInfo := shared.NewAgentInfo(
		p.AgentID,
		p.Role,
		agent.StreamingEnabled,
		p.Name,
		p.Description,
	)

	agentInfo.InputSchema = agent.InputSchema
	agentInfo.OutputSchema = agent.OutputSchema

	return &agentInfo
}

// PromptConfig represents prompt configuration
type PromptConfig struct {
	AgentID           uuid.UUID              `json:"agent_id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description,omitempty"`
	GlobalInstruction string                 `json:"global_instruction,omitempty"`
	Content           string                 `json:"content"`
	Schema            map[string]interface{} `json:"schema"`
}

// SettingsConfig represents agent settings
type SettingsConfig struct {
	AgentID     uuid.UUID     `json:"agent_id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Agent       AgentSettings `json:"agent"`
}

// AgentSettings represents the agent runtime settings
type AgentSettings struct {
	ApplicationName   string                 `json:"application_name"`
	PlanningEnabled   bool                   `json:"planning_enabled"`
	MaxIterations     int                    `json:"max_iterations"`
	Timeout           int                    `json:"timeout"`
	StreamingEnabled  bool                   `json:"streaming_enabled"`
	ChannelBufferSize int                    `json:"channel_buffer_size"`
	LLM               LLMSettings            `json:"llm"`
	SubAgents         []shared.AgentRole     `json:"sub_agents,omitempty"`
	InputSchema       map[string]interface{} `json:"input_schema,omitempty"`
	OutputSchema      map[string]interface{} `json:"output_schema,omitempty"`
	OutputKey         string                 `json:"output_key,omitempty"`
	TimeAwareness     *TimeAwarenessSettings `json:"time_awareness,omitempty"`
}

// LLMSettings represents LLM configuration
type LLMSettings struct {
	Model             string               `json:"model"`
	Temperature       float64              `json:"temperature"`
	MaxTokens         int                  `json:"max_tokens"`
	TopP              float64              `json:"top_p"`
	FrequencyPenalty  float64              `json:"frequency_penalty"`
	PresencePenalty   float64              `json:"presence_penalty"`
	Provider          shared.ModelProvider `json:"provider"`
	BaseURL           string               `json:"base_url,omitempty"`
	APIKey            string               `json:"api_key,omitempty"`
	ChannelBufferSize int                  `json:"channel_buffer_size,omitempty"`
}

// ToolsConfig represents tool configuration
type ToolsConfig struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Tools       map[string]ToolConfig    `json:"tools"`
	ToolSets    map[string]ToolSetConfig `json:"toolsets"`
}

// ToolConfig represents individual tool configuration
type ToolConfig struct {
	Enabled bool                   `json:"enabled"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

// ToolSetConfig represents tool set configuration
type ToolSetConfig struct {
	Enabled bool                   `json:"enabled"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

// AgentFactory creates agents using configuration-based approach
type AgentFactory interface {
	CreateAgent(ctx context.Context, environment EnvironmentName, agentRole shared.AgentRole) (shared.TheAgent, error)
	CreateAgentByID(ctx context.Context, environment EnvironmentName, agentID uuid.UUID) (shared.TheAgent, error)
	CreateAllAgentsInEnvironment(ctx context.Context, envName EnvironmentName) ([]shared.TheAgent, error)
	ValidateConfiguration() error
	GetAgentConfig(environment EnvironmentName, agentRole shared.AgentRole) (*AgentConfig, error)
}

// ToolFactory creates tools from configuration
type ToolFactory interface {
	CreateTools(toolsConfig ToolsConfig, availableAgents []*shared.AgentInfo) ([]tool.Tool, []tool.ToolSet, error)
}

// ConfigProvider loads configurations from various sources
type ConfigProvider interface {
	// System-level operations
	LoadSystemConfig() (*SystemConfig, error)
	LoadEnvironmentConfig(envName EnvironmentName) (*EnvironmentConfig, error)
	ResolveAgentByRole(envName EnvironmentName, role shared.AgentRole) (*AgentConfig, error)

	// Legacy methods (for backward compatibility)
	LoadAgentComposition(environment EnvironmentName, agentRole shared.AgentRole) (*AgentConfig, error)
	LoadPrompt(agentRole shared.AgentRole) (*PromptConfig, error)
	LoadSettings(agentRole shared.AgentRole) (*SettingsConfig, error)
	LoadToolProfile(agentRole shared.AgentRole) (*ToolsConfig, error)
	ValidateConfiguration() error
	GetAgentsInEnvironment(environment EnvironmentName, includeHuman bool) ([]*shared.AgentInfo, error)
	GetAgentInfoByID(agentID uuid.UUID) (*shared.AgentInfo, error)
}

type CondenserServiceSettings struct {
	LoggingEnabled      bool                          `json:"logging_enabled,omitempty"`
	TriggerThreshold    float64                       `json:"trigger_threshold,omitempty"`
	SummaryPrompt       string                        `json:"summary_prompt,omitempty"`
	TokenCountingMethod condenser.TokenCountingMethod `json:"token_counting_method,omitempty"`
	RecentEventsToKeep  int                           `json:"recent_events_to_keep,omitempty"`
	MaxContextTokens    int                           `json:"max_context_tokens,omitempty"`
	ModelProvider       shared.ModelProvider          `json:"model_provider,omitempty"`
	ModelName           string                        `json:"model_name,omitempty"`
}

type EnvironmentConfig struct {
	Description string                         `json:"description,omitempty"`
	Name        string                         `json:"name,omitempty"`
	Agents      map[AgentName]AgentConfig      `json:"agents,omitempty"`
	Roles       map[shared.AgentRole]AgentName `json:"roles,omitempty"`
	Condenser   CondenserServiceSettings       `json:"condenser,omitempty"`
}

type CommonConfig struct {
}

type SystemConfig struct {
	Environments map[EnvironmentName]EnvironmentConfig `json:"environments,omitempty"`
	Common       CommonConfig                          `json:"common_settings,omitempty"`
}

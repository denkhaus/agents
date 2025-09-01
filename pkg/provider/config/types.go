package config

import (
	"context"

	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

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
	ReactEnabled      bool                   `json:"react_enabled"`
	MaxIterations     int                    `json:"max_iterations"`
	Timeout           int                    `json:"timeout"`
	StreamingEnabled  bool                   `json:"streaming_enabled"`
	ChannelBufferSize int                    `json:"channel_buffer_size"`
	LLM               LLMSettings            `json:"llm"`
	SubAgents         []shared.AgentRole     `json:"sub_agents,omitempty"`
	InputSchema       map[string]interface{} `json:"input_schema,omitempty"`
	OutputSchema      map[string]interface{} `json:"output_schema,omitempty"`
	OutputKey         string                 `json:"output_key,omitempty"`
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
	CreateAgent(ctx context.Context, environment string, agentRole shared.AgentRole) (shared.TheAgent, error)
	CreateAgentByID(ctx context.Context, environment string, agentID uuid.UUID) (shared.TheAgent, error)
	ValidateConfiguration() error
	GetAgentConfig(environment string, agentRole shared.AgentRole) (*AgentConfig, error)
}

// ToolFactory creates tools from configuration
type ToolFactory interface {
	CreateTools(toolsConfig ToolsConfig) ([]tool.Tool, []tool.ToolSet, error)
}

// ConfigProvider loads configurations from various sources
type ConfigProvider interface {
	LoadAgentComposition(environment string, agentRole shared.AgentRole) (*AgentConfig, error)
	LoadPrompt(agentRole shared.AgentRole) (*PromptConfig, error)
	LoadSettings(agentRole shared.AgentRole) (*SettingsConfig, error)
	LoadToolProfile(agentRole shared.AgentRole) (*ToolsConfig, error)
	ValidateConfiguration() error
	GetAgentsInEnvironment(environment string, includeHuman bool) ([]*shared.AgentInfo, error)
	GetAgentInfoByID(agentID uuid.UUID) (*shared.AgentInfo, error) // New method
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

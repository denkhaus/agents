package config

import (
	"context"

	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

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

package provider

import (
	"context"

	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/chainagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/cycleagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
)

type AgentConfiguration interface {
	GetName() string
	GetType() shared.AgentType
	IsStreamingEnabled() bool
	GetDefaultOptions(ctx context.Context, provider AgentProvider) ([]llmagent.Option, error)
	GetCycleOptions(ctx context.Context, provider AgentProvider) ([]cycleagent.Option, error)
	GetChainOptions(ctx context.Context, provider AgentProvider) ([]chainagent.Option, error)
}

type AgentProvider interface {
	GetAgent(ctx context.Context, agentID uuid.UUID) (agent.Agent, bool, error)
}

type SettingsProvider interface {
	GetAgentConfiguration(agentID uuid.UUID) (AgentConfiguration, error)
}

type Workspace interface {
	GetWorkspacePath() (string, error)
}

type WorkspaceProvider interface {
	GetWorkspace(agentID uuid.UUID) (Workspace, error)
}

type Prompt interface {
	GetName() string
	GetDescription() string
	GetGlobalInstruction() string
	GetInstruction(data interface{}) (string, error)
}

// PromptManager defines the interface for managing and rendering prompts.
type PromptManager interface {
	GetPrompt(agentID uuid.UUID) (Prompt, error)
}

type PromptProvider interface {
	GetPrompt(agentID uuid.UUID, data interface{}) (Prompt, error)
}

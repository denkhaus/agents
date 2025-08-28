package provider

import (
	"context"

	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent/chainagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/cycleagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/parallelagent"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

type AgentProviderOptions struct {
	LLMOpt      []llmagent.Option
	CycleOpt    []cycleagent.Option
	ParallelOpt []parallelagent.Option
	ChainOpt    []chainagent.Option
}

// Option is a function that configures an Agent.
type AgentProviderOption func(*AgentProviderOptions)

type AgentProvider interface {
	GetAgent(ctx context.Context, agentID uuid.UUID, opt ...AgentProviderOption) (shared.TheAgent, error)
}

type AgentConfiguration interface {
	GetName() string
	GetType() shared.AgentType
	IsStreamingEnabled() bool
	GetDefaultOptions(ctx context.Context, provider AgentProvider, opt ...llmagent.Option) ([]llmagent.Option, error)
	GetCycleOptions(ctx context.Context, provider AgentProvider, opt ...cycleagent.Option) ([]cycleagent.Option, error)
	GetChainOptions(ctx context.Context, provider AgentProvider, opt ...chainagent.Option) ([]chainagent.Option, error)
	GetParallelOptions(ctx context.Context, provider AgentProvider, opt ...parallelagent.Option) ([]parallelagent.Option, error)
}

type SettingsProvider interface {
	GetActiveAgents(includeHumanAgent bool) ([]shared.AgentInfo, error)
	GetAgentConfiguration(agentID uuid.UUID) (AgentConfiguration, error)
}

type Workspace interface {
	GetPath() (string, error)
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

type ChatProviderOptions struct {
	UserID               string
	SessionID            string
	AppName              string
	SessionService       session.Service
	AgentProviderOptions []AgentProviderOption
}

// Option is a function that configures an LLMAgent.
type ChatProviderOption func(*ChatProviderOptions)

type Chat interface {
	ProcessMessage(ctx context.Context, userMessage string) error
}

type ChatProvider interface {
	GetChat(ctx context.Context, agentID uuid.UUID, opts ...ChatProviderOption) (Chat, error)
}

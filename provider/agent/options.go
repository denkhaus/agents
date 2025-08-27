// Package agent provides option functions for configuring different types of agents.
package agent

import (
	"github.com/denkhaus/agents/provider"
	"trpc.group/trpc-go/trpc-agent-go/agent/chainagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/cycleagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/parallelagent"
)

// WithLLMAgentOptions configures options for LLM agents.
// It allows passing custom options to the underlying llmagent.New function.
func WithLLMAgentOptions(opt ...llmagent.Option) provider.AgentProviderOption {
	return func(opts *provider.AgentProviderOptions) {
		opts.LLMOpt = opt
	}
}

// WithCycleAgentOptions configures options for cycle agents.
// It allows passing custom options to the underlying cycleagent.New function.
func WithCycleAgentOptions(opt ...cycleagent.Option) provider.AgentProviderOption {
	return func(opts *provider.AgentProviderOptions) {
		opts.CycleOpt = opt
	}
}

// WithParallelAgentOptions configures options for parallel agents.
// It allows passing custom options to the underlying parallelagent.New function.
func WithParallelAgentOptions(opt ...parallelagent.Option) provider.AgentProviderOption {
	return func(opts *provider.AgentProviderOptions) {
		opts.ParallelOpt = opt
	}
}

// WithChainAgentOptions configures options for chain agents.
// It allows passing custom options to the underlying chainagent.New function.
func WithChainAgentOptions(opt ...chainagent.Option) provider.AgentProviderOption {
	return func(opts *provider.AgentProviderOptions) {
		opts.ChainOpt = opt
	}
}

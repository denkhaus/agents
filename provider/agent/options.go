package agent

import (
	"github.com/denkhaus/agents/provider"
	"trpc.group/trpc-go/trpc-agent-go/agent/chainagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/cycleagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/parallelagent"
)

func WithLLMAgentOptions(opt ...llmagent.Option) provider.AgentProviderOption {
	return func(opts *provider.AgentProviderOptions) {
		opts.LLMOpt = opt
	}
}

func WithCycleAgentOptions(opt ...cycleagent.Option) provider.AgentProviderOption {
	return func(opts *provider.AgentProviderOptions) {
		opts.CycleOpt = opt
	}
}

func WithParallelAgentOptions(opt ...parallelagent.Option) provider.AgentProviderOption {
	return func(opts *provider.AgentProviderOptions) {
		opts.ParallelOpt = opt
	}
}

func WithChainAgentOptions(opt ...chainagent.Option) provider.AgentProviderOption {
	return func(opts *provider.AgentProviderOptions) {
		opts.ChainOpt = opt
	}
}

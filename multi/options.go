package multi

import (
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// OnError is a callback function type for handling errors from agents.
type OnError func(info *shared.AgentInfo, err error)

// OnProgress is a callback function type for reporting progress updates.
type OnProgress func(format string, a ...any)

// OnMessage is a callback function type for handling messages from agents.
type OnMessage func(info *shared.AgentInfo, content string)

// OnToolCall is a callback function type for handling tool calls made by agents.
type OnToolCall func(info *shared.AgentInfo, functionDef model.FunctionDefinitionParam)

// Options contains configuration settings for the ChatProcessor.
type Options struct {
	sessionID       uuid.UUID
	humanAgent      shared.TheAgent
	agents          []shared.TheAgent
	applicationName string
	onToolCall      OnToolCall
	onMessage       OnMessage
	onProgress      OnProgress
	onError         OnError
}

// ChatProcessorOption is a function type for configuring ChatProcessor options.
type ChatProcessorOption func(*Options)

// WithSessionID sets the SessionID to use.
func WithSessionID(sessionID uuid.UUID) ChatProcessorOption {
	return func(opts *Options) {
		opts.sessionID = sessionID
	}
}

// WithApplicationName sets the application name for the ChatProcessor.
func WithApplicationName(applicationName string) ChatProcessorOption {
	return func(opts *Options) {
		opts.applicationName = applicationName
	}
}

// WithHumanAgent sets the human agent for the ChatProcessor.
func WithHumanAgent(humanAgent shared.TheAgent) ChatProcessorOption {
	return func(opts *Options) {
		opts.humanAgent = humanAgent
	}
}

// WithAgents sets the AI agents for the ChatProcessor.
func WithAgents(agents ...shared.TheAgent) ChatProcessorOption {
	return func(opts *Options) {
		opts.agents = agents
	}
}

// WithOnError sets the error callback function for the ChatProcessor.
func WithOnError(onError OnError) ChatProcessorOption {
	return func(opts *Options) {
		opts.onError = onError
	}
}

// WithOnProgress sets the progress callback function for the ChatProcessor.
func WithOnProgress(onProgress OnProgress) ChatProcessorOption {
	return func(opts *Options) {
		opts.onProgress = onProgress
	}
}

// WithOnMessage sets the message callback function for the ChatProcessor.
func WithOnMessage(onMessage OnMessage) ChatProcessorOption {
	return func(opts *Options) {
		opts.onMessage = onMessage
	}
}

// WithOnToolCall sets the tool call callback function for the ChatProcessor.
func WithOnToolCall(onToolCall OnToolCall) ChatProcessorOption {
	return func(opts *Options) {
		opts.onToolCall = onToolCall
	}
}

package multi

import (
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type OnError func(info *shared.AgentInfo, err error)
type OnProgress func(format string, a ...any)
type OnMessage func(info *shared.AgentInfo, content string)
type OnToolCall func(info *shared.AgentInfo, functionDef model.FunctionDefinitionParam)

type Options struct {
	sessionID  uuid.UUID
	humanAgent shared.TheAgent
	onToolCall OnToolCall
	onMessage  OnMessage
	onProgress OnProgress
	onError    OnError
}

type ChatProcessorOption func(*Options)

// WithSessionID sets the SessionID to use.
func WithSessionID(sessionID uuid.UUID) ChatProcessorOption {
	return func(opts *Options) {
		opts.sessionID = sessionID
	}
}

func WithHumanAgent(humanAgent shared.TheAgent) ChatProcessorOption {
	return func(opts *Options) {
		opts.humanAgent = humanAgent
	}
}

func WithOnError(onError OnError) ChatProcessorOption {
	return func(opts *Options) {
		opts.onError = onError
	}
}

func WithOnProgress(onProgress OnProgress) ChatProcessorOption {
	return func(opts *Options) {
		opts.onProgress = onProgress
	}
}

func WithOnMessage(onMessage OnMessage) ChatProcessorOption {
	return func(opts *Options) {
		opts.onMessage = onMessage
	}
}

func WithOnToolCall(onToolCall OnToolCall) ChatProcessorOption {
	return func(opts *Options) {
		opts.onToolCall = onToolCall
	}
}

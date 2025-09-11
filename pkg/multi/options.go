package multi

import (
	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/shared"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

// SystemMessageType defines the type of system message.
//
//go:generate stringer -type=SystemMessageType
type SystemMessageType int

const (
	SystemMessageDefault SystemMessageType = iota
	SystemMessageSending
	SystemMessageDelivered
	SystemMessageProcessed
)

// OnError is a callback function type for handling errors from agents.
type OnError func(info *messaging.RoutingInfo, err error)

// OnProgress is a callback function type for reporting progress updates.
type OnProgress func(info *messaging.RoutingInfo, messageType SystemMessageType, format string, a ...any)

// OnMessage is a callback function type for handling messages from agents.
type OnMessage func(info *messaging.RoutingInfo, content string)

// OnRawEvent is a callback function type for handling every event without filtering.
type OnRawEvent func(info *messaging.RoutingInfo, event *event.Event)

// OnReasoningMessage is a callback function type for handling reasoning/thinking messages from agents.
type OnReasoningMessage func(info *messaging.RoutingInfo, content string)

// OnToolCall is a callback function type for handling tool calls made by agents.
type OnToolCall func(info *messaging.RoutingInfo, functionDef model.FunctionDefinitionParam)

// Options contains configuration settings for the ChatProcessor.
type Options struct {
	availableAgents    []shared.TheAgent
	applicationName    string
	sessionService     session.Service
	onToolCall         OnToolCall
	onMessage          OnMessage
	onRawEvent         OnRawEvent
	onReasoningMessage OnReasoningMessage
	onProgress         OnProgress
	onError            OnError
}

// ChatProcessorOption is a function type for configuring ChatProcessor options.
type ChatProcessorOption func(*Options)

// WithApplicationName sets the application name for the ChatProcessor.
func WithApplicationName(applicationName string) ChatProcessorOption {
	return func(opts *Options) {
		opts.applicationName = applicationName
	}
}

// WithSessionService sets the session service to use.
func WithSessionService(service session.Service) ChatProcessorOption {
	return func(opts *Options) {
		opts.sessionService = service
	}
}

// WithAgents sets the AI agents for the ChatProcessor.
func WithAgents(agents ...shared.TheAgent) ChatProcessorOption {
	return func(opts *Options) {
		opts.availableAgents = agents
	}
}

// WithOnError sets the error callback function for the ChatProcessor.
func WithOnError(onError OnError) ChatProcessorOption {
	return func(opts *Options) {
		opts.onError = onError
	}
}

// WithOnProgress sets the progress callback function for the ChatProcessor.
func WithOnRawEvent(onRawEvent OnRawEvent) ChatProcessorOption {
	return func(opts *Options) {
		opts.onRawEvent = onRawEvent
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

// WithOnReasoningMessage sets the reasoning message callback function for the ChatProcessor.
func WithOnReasoningMessage(onReasoningMessage OnReasoningMessage) ChatProcessorOption {
	return func(opts *Options) {
		opts.onReasoningMessage = onReasoningMessage
	}
}

// WithOnToolCall sets the tool call callback function for the ChatProcessor.
func WithOnToolCall(onToolCall OnToolCall) ChatProcessorOption {
	return func(opts *Options) {
		opts.onToolCall = onToolCall
	}
}

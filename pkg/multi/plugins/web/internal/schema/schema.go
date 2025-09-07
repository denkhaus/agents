//
// Tencent is pleased to support the open source community by making trpc-agent-go available.
//
// Copyright (C) 2025 Tencent.  All rights reserved.
//
// trpc-agent-go is licensed under the Apache License Version 2.0.
//
//

// Package schema defines JSON schema structs used by the CLI HTTP server.
// These types are internal – they are not intended to be imported by other
// packages. They only exist to facilitate request/response marshalling.
package schema

import (
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type Part interface {
}

type TextPart struct {
	Content string `json:"content,omitempty"`
}

type FunctionCallPart struct {
	Name string      `json:"name,omitempty"`
	Args interface{} `json:"args,omitempty"`
	ID   string      `json:"id,omitempty"`
}

type FunctionResponsePart struct {
	Name     string      `json:"name,omitempty"`
	Args     interface{} `json:"args,omitempty"`
	ID       string      `json:"id,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

type UsageMetaData struct {
	PromptTokenCount     int `json:"prompt_token_count,omitempty"`
	CandidatesTokenCount int `json:"candidates_token_count,omitempty"`
	TotalTokenCount      int `json:"total_token_count,omitempty"`
}

type LLMEvent struct {
	base         *event.Event   `json:"-"`
	Usage        *UsageMetaData `json:"usage,omitempty"`
	Done         bool           `json:"done,omitempty"`
	Partial      bool           `json:"partial,omitempty"`
	Object       string         `json:"object,omitempty"`
	Created      int64          `json:"created,omitempty"`
	Model        string         `json:"model,omitempty"`
	Role         model.Role     `json:"role,omitempty"`
	Parts        []Part         `json:"parts,omitempty"`
	Timestamp    int64          `json:"timestamp,omitempty"`
	ID           string         `json:"id,omitempty"`
	InvocationID string         `json:"invocation_id,omitempty"`
	Author       string         `json:"author,omitempty"`
}

// ADKSession mirrors the structure expected by ADK Web UI for a session.
// Field names follow the camel-case convention required by the UI.
type ADKSession struct {
	AppName        string            `json:"appName"`
	AgentID        uuid.UUID         `json:"agentId"`
	ID             uuid.UUID         `json:"id"`
	CreateTime     int64             `json:"createTime"`
	LastUpdateTime int64             `json:"lastUpdateTime"`
	State          map[string][]byte `json:"state"`
	Events         []*LLMEvent       `json:"events"`
}

// Span represents a single span in the trace.
type Span struct {
	Name         string         `json:"name"`
	SpanID       string         `json:"span_id"`
	TraceID      string         `json:"trace_id"`
	StartTime    int64          `json:"start_time"`
	EndTime      int64          `json:"end_time"`
	Attributes   map[string]any `json:"attributes"`
	ParentSpanID string         `json:"parent_span_id"`
}

// -----------------------------------------------------------------------------
// Incoming request payloads ----------------------------------------------------
// -----------------------------------------------------------------------------

// Part represents a single message segment used by ADK Web.
type PartIncoming struct {
	Text             string            `json:"text,omitempty"`
	InlineData       *InlineData       `json:"inlineData,omitempty"`
	FunctionCall     *FunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *FunctionResponse `json:"functionResponse,omitempty"`
}

// InlineData encapsulates binary data (image/audio/video/file).
type InlineData struct {
	Data        string `json:"data"`
	MimeType    string `json:"mimeType"`
	DisplayName string `json:"displayName,omitempty"`
}

// FunctionCall matches GenAI functionCall part.
type FunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args,omitempty"`
}

// FunctionResponse matches GenAI functionResponse part.
type FunctionResponse struct {
	Name     string      `json:"name"`
	Response interface{} `json:"response"`
	ID       string      `json:"id,omitempty"`
}

type AgentRunRequest struct {
	AppName     string    `json:"appName"`
	FromAgentID uuid.UUID `json:"fromAgentId"`
	ToAgentID   uuid.UUID `json:"toAgentId"`
	SessionID   uuid.UUID `json:"sessionId"`
	Content     Content   `json:"content"`
	Streaming   bool      `json:"streaming"`
}

// TraceLLMRequest represents a trace request for LLM operations.
type TraceLLMRequest struct {
	Contents []Content `json:"contents"`
}

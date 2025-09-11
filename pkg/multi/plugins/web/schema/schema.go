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
	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/google/uuid"
)

// ADKSession mirrors the structure expected by ADK Web UI for a session.
// Field names follow the camel-case convention required by the UI.
type ADKSession struct {
	AppName        string                `json:"appName"`
	AgentID        uuid.UUID             `json:"agentId"`
	ID             uuid.UUID             `json:"id"`
	CreateTime     int64                 `json:"createTime"`
	LastUpdateTime int64                 `json:"lastUpdateTime"`
	State          map[string][]byte     `json:"state"`
	Events         []*messaging.LLMEvent `json:"events"`
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
	messaging.RoutingInfo `json:",inline"`
	AppName               string  `json:"appName"`
	Content               Content `json:"content"`
}

// TraceLLMRequest represents a trace request for LLM operations.
type TraceLLMRequest struct {
	Contents []Content `json:"contents"`
}

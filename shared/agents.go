package shared

import "github.com/google/uuid"

var (
	AgentIDCoder    = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	AgentIDDebugger = uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
)

type AgentRole string

func (p AgentRole) String() string {
	return string(p)
}

const (
	AgentRoleCoder    AgentRole = "coder"
	AgentRoleDebugger AgentRole = "debugger"
)

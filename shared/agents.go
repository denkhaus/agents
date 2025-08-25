package shared

import (
	"fmt"
	"github.com/google/uuid"
)

var (
	AgentIDCoder    = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	AgentIDDebugger = uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
)

type AgentRole string

func (p AgentRole) String() string {
	return string(p)
}

// Validate checks if the AgentRole is a valid defined role
func (p AgentRole) Validate() error {
	switch p {
	case AgentRoleCoder, AgentRoleDebugger:
		return nil
	default:
		return fmt.Errorf("invalid agent role: %s. Valid roles are: %s, %s", 
			p, AgentRoleCoder, AgentRoleDebugger)
	}
}

const (
	AgentRoleCoder    AgentRole = "coder"
	AgentRoleDebugger AgentRole = "debugger"
)

package shared

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	AgentIDSupervisor     = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	AgentIDCoder          = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	AgentIDDebugger       = uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
	AgentIDProjectManager = uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")
)

type AgentType string

func (p AgentType) String() string {
	return string(p)
}

// Validate checks if the AgentRole is a valid defined role
func (p AgentType) Validate() error {
	switch p {
	case AgentTypeDefault, AgentTypeChain, AgentTypeCycle, AgentTypeParallel:
		return nil
	default:
		return fmt.Errorf("invalid agent role: %s. Valid roles are: %s, %s, %s, %s",
			p, AgentTypeDefault, AgentTypeChain, AgentTypeCycle, AgentTypeParallel)
	}
}

const (
	AgentTypeDefault  AgentType = "default"
	AgentTypeChain    AgentType = "chain"
	AgentTypeCycle    AgentType = "cycle"
	AgentTypeParallel AgentType = "parallel"
)

type AgentRole string

func (p AgentRole) String() string {
	return string(p)
}

// Validate checks if the AgentRole is a valid defined role
func (p AgentRole) Validate() error {
	switch p {
	case AgentRoleSupervisor, AgentRoleCoder, AgentRoleDebugger, AgentRoleProjectManager:
		return nil
	default:
		return fmt.Errorf("invalid agent role: %s. Valid roles are: %s, %s, %s, %s",
			p, AgentRoleSupervisor, AgentRoleCoder, AgentRoleDebugger, AgentRoleProjectManager)
	}
}

const (
	AgentRoleSupervisor     AgentRole = "supervisor"
	AgentRoleCoder          AgentRole = "coder"
	AgentRoleDebugger       AgentRole = "debugger"
	AgentRoleProjectManager AgentRole = "project-manager"
)

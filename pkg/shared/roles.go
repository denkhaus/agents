package shared

import "fmt"

type AgentRole string

func (p AgentRole) String() string {
	return string(p)
}

// Validate checks if the AgentRole is a valid defined role
func (p AgentRole) Validate() error {
	switch p {
	case AgentRoleSupervisor,
		AgentRoleCoder,
		AgentRoleDebugger,
		AgentRoleProjectManager,
		AgentRoleResearcher,
		AgentRoleHuman:
		return nil
	default:
		return fmt.Errorf("invalid agent role: %s. Valid roles are: %s, %s, %s, %s, %s, %s",
			p, AgentRoleSupervisor, AgentRoleCoder, AgentRoleDebugger, AgentRoleProjectManager, AgentRoleHuman, AgentRoleResearcher)
	}
}

const (
	AgentRoleSupervisor     AgentRole = "supervisor"
	AgentRoleCoder          AgentRole = "coder"
	AgentRoleDebugger       AgentRole = "debugger"
	AgentRoleProjectManager AgentRole = "project_manager" // Changed hyphen to underscore
	AgentRoleHuman          AgentRole = "human"
	AgentRoleResearcher     AgentRole = "researcher"
)
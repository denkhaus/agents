package agents

type AgentRole string

func (p AgentRole) String() string {
	return string(p)
}

const (
	AgentRoleCoder AgentRole = "coder"
)

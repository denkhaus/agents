package messaging

import (
	"context"
	"testing"

	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/utils"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// mockAgent is a simple mock agent for testing
type mockAgent struct {
	name string
	id   uuid.UUID
}

func (ma *mockAgent) Run(ctx context.Context, invocation *agent.Invocation) (<-chan *event.Event, error) {
	return nil, nil
}

func (ma *mockAgent) Tools() []tool.Tool {
	return []tool.Tool{}
}

func (ma *mockAgent) Info() agent.Info {
	return agent.Info{
		Name: ma.name,
	}
}

func (ma *mockAgent) SubAgents() []agent.Agent {
	return []agent.Agent{}
}

func (ma *mockAgent) FindSubAgent(name string) agent.Agent {
	return nil
}

func (ma *mockAgent) ID() uuid.UUID {
	return ma.id
}

func (ma *mockAgent) IsStreaming() bool {
	return false
}

func (ma *mockAgent) GetInfo() *shared.AgentInfo {
	return &shared.AgentInfo{
		ID:          ma.id,
		Role:        "test",
		IsStreaming: utils.BoolPtr(false),
		Name:        ma.name,
		Description: "Mock agent for testing",
	}
}

func (ma *mockAgent) GetRole() shared.AgentRole {
	return "test"
}

func TestResourceManagerIntegration(t *testing.T) {
	sessionID := uuid.New()
	broker := NewMessageBroker(sessionID)

	// Create test agents
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	agent1 := &mockAgent{name: "Agent1", id: uuid1}
	agent2 := &mockAgent{name: "Agent2", id: uuid2}

	// Convert to shared.TheAgent using the shared.NewAgent function
	wrapper1 := shared.NewAgent(agent1, uuid1, false, shared.AgentRole("test"), []uuid.UUID{})
	wrapper2 := shared.NewAgent(agent2, uuid2, false, shared.AgentRole("test"), []uuid.UUID{})

	// Register agents
	broker.RegisterAgent(uuid1, wrapper1)
	broker.RegisterAgent(uuid2, wrapper2)

	// Check that agents are registered
	agentIDs := broker.ListAgentIDs()
	if len(agentIDs) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(agentIDs))
	}

	// Send a message
	err := broker.SendMessage(uuid1, uuid2, "Test message")
	if err != nil {
		t.Logf("Expected error in test message (agent not wrapped), but got: %v", err)
	}

	// Unregister an agent
	broker.UnregisterAgent(uuid1)

	// Check that the agent was unregistered
	agentIDs = broker.ListAgentIDs()
	if len(agentIDs) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agentIDs))
	}

	// Verify the remaining agent is the correct one
	if agentIDs[0] != uuid2 {
		t.Errorf("Expected remaining agent to be uuid2, got %s", agentIDs[0])
	}
}

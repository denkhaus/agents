package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// TestAgent is a simple agent for testing
type TestAgent struct {
	name string
	uuid uuid.UUID
}

func (ta *TestAgent) Run(ctx context.Context, invocation *agent.Invocation) (<-chan *event.Event, error) {
	eventChan := make(chan *event.Event, 1)

	go func() {
		defer close(eventChan)

		// Create a message in the content
		message := model.NewAssistantMessage("Test message from " + ta.name)

		response := &model.Response{
			Object:    model.ObjectTypeChatCompletion,
			Done:      true,
			Created:   time.Now().Unix(),
			Choices:   []model.Choice{{Message: message}},
			Timestamp: time.Now(),
		}

		event := &event.Event{
			Response:     response,
			InvocationID: uuid.New().String(),
			Author:       ta.name,
			ID:           uuid.New().String(),
			Timestamp:    time.Now(),
		}

		eventChan <- event
	}()

	return eventChan, nil
}

func (ta *TestAgent) Tools() []tool.Tool {
	return []tool.Tool{}
}

func (ta *TestAgent) Info() agent.Info {
	return agent.Info{
		Name:        ta.name,
		Description: "Test agent",
	}
}

func (ta *TestAgent) SubAgents() []agent.Agent {
	return []agent.Agent{}
}

func (ta *TestAgent) FindSubAgent(name string) agent.Agent {
	return nil
}

// ID returns the agent's UUID
func (ta *TestAgent) ID() uuid.UUID {
	return ta.uuid
}

// IsStreaming returns whether the agent is streaming
func (ta *TestAgent) IsStreaming() bool {
	return false
}

// GetInfo returns the agent information
func (ta *TestAgent) GetInfo() *shared.AgentInfo {
	info := shared.NewAgentInfo(
		ta.uuid,
		"test", // Using a simple string for test role
		false,
		ta.name,
		"Test agent",
	)
	return &info
}

// GetRole returns the agent role
func (ta *TestAgent) GetRole() shared.AgentRole {
	return "test"
}

func TestMessageBroker(t *testing.T) {
	broker := NewMessageBroker()

	agent1 := &TestAgent{name: "Agent1", uuid: uuid.New()}
	agent2 := &TestAgent{name: "Agent2", uuid: uuid.New()}

	// Convert to shared.TheAgent using the shared.NewAgent function
	wrapper1 := NewMessagingWrapper(shared.NewAgent(agent1, agent1.uuid, false), broker)
	wrapper2 := NewMessagingWrapper(shared.NewAgent(agent2, agent2.uuid, false), broker)

	// Type assert to access messaging-specific methods
	messagingWrapper1, ok := wrapper1.(*messagingWrapper)
	if !ok {
		t.Fatal("wrapper1 is not a messagingWrapper")
	}
	messagingWrapper2, ok := wrapper2.(*messagingWrapper)
	if !ok {
		t.Fatal("wrapper2 is not a messagingWrapper")
	}

	// Test sending a message
	err := messagingWrapper1.SendMessage(messagingWrapper2.ID(), "Hello from Agent1")
	if err != nil {
		t.Errorf("Failed to send message: %v", err)
	}

	// Test receiving a message
	msgChan, err := messagingWrapper2.GetMessageChannel()
	if err != nil {
		t.Errorf("Failed to get message channel: %v", err)
	}

	select {
	case msg := <-msgChan:
		if msg.Content != "Hello from Agent1" {
			t.Errorf("Expected 'Hello from Agent1', got '%s'", msg.Content)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestMessagingWrapper(t *testing.T) {
	broker := NewMessageBroker()

	uuid1 := uuid.New()
	testAgent := &TestAgent{name: "TestAgent", uuid: uuid1}
	// Convert to shared.TheAgent using the shared.NewAgent function
	wrapper := NewMessagingWrapper(shared.NewAgent(testAgent, uuid1, false), broker)

	// Type assert to access messaging-specific methods
	messagingWrapper, ok := wrapper.(*messagingWrapper)
	if !ok {
		t.Fatal("wrapper is not a messagingWrapper")
	}

	// Test that the wrapper has an ID
	if messagingWrapper.ID() == uuid.Nil {
		t.Error("Wrapper should have a valid ID")
	}

	// Test that the wrapper info is correct
	info := messagingWrapper.Info()
	if info.Name != "TestAgent" {
		t.Errorf("Expected name 'TestAgent', got '%s'", info.Name)
	}

	// Test running the agent
	ctx := context.Background()
	invocation := &agent.Invocation{
		InvocationID: uuid.New().String(),
	}

	eventChan, err := messagingWrapper.Run(ctx, invocation)
	if err != nil {
		t.Errorf("Failed to run agent: %v", err)
	}

	select {
	case event := <-eventChan:
		if event.Author != "TestAgent" {
			t.Errorf("Expected author 'TestAgent', got '%s'", event.Author)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for event")
	}
}

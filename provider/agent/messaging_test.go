package agent

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// TestAgent is a simple agent for testing
type TestAgent struct {
	name string
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

func TestMessageBroker(t *testing.T) {
	broker := NewMessageBroker()
	
	agent1 := &TestAgent{name: "Agent1"}
	agent2 := &TestAgent{name: "Agent2"}
	
	wrapper1 := NewMessagingWrapper(agent1, broker)
	wrapper2 := NewMessagingWrapper(agent2, broker)
	
	// Test sending a message
	err := wrapper1.SendMessage(wrapper2.ID(), "Hello from Agent1")
	if err != nil {
		t.Errorf("Failed to send message: %v", err)
	}
	
	// Test receiving a message
	msgChan, err := wrapper2.GetMessageChannel()
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
	
	testAgent := &TestAgent{name: "TestAgent"}
	wrapper := NewMessagingWrapper(testAgent, broker)
	
	// Test that the wrapper has an ID
	if wrapper.ID() == uuid.Nil {
		t.Error("Wrapper should have a valid ID")
	}
	
	// Test that the wrapper info is correct
	info := wrapper.Info()
	if info.Name != "TestAgent" {
		t.Errorf("Expected name 'TestAgent', got '%s'", info.Name)
	}
	
	// Test running the agent
	ctx := context.Background()
	invocation := &agent.Invocation{
		InvocationID: uuid.New().String(),
	}
	
	eventChan, err := wrapper.Run(ctx, invocation)
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
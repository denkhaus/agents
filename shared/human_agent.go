package shared

import (
	"context"
	"time"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// humanAgentImpl implements the agent.Agent interface for the human
type humanAgentImpl struct {
	AgentInfo
}

func NewHumanAgent(info AgentInfo) TheAgent {
	return &humanAgentImpl{
		AgentInfo: info,
	}
}

func (d *humanAgentImpl) GetInfo() *AgentInfo {
	return &d.AgentInfo
}

func (d *humanAgentImpl) IsStreaming() bool {
	return false
}

func (d *humanAgentImpl) GetRole() AgentRole {
	return AgentRoleHuman
}

func (d *humanAgentImpl) Run(ctx context.Context, invocation *agent.Invocation) (<-chan *event.Event, error) {
	// Create an event channel for the human agent
	eventChan := make(chan *event.Event, 10)

	go func() {
		defer close(eventChan)

		// If there's a message in the invocation, convert it to an event
		if invocation != nil && invocation.Message.Content != "" {
			// Create an assistant message event to display the received message
			response := &model.Response{
				Object:  model.ObjectTypeChatCompletion,
				Done:    true,
				Created: time.Now().Unix(),
				Choices: []model.Choice{{
					Message: model.NewAssistantMessage(invocation.Message.Content),
				}},
			}

			event := &event.Event{
				Response:     response,
				InvocationID: uuid.New().String(),
				Author:       d.ID().String(),
				ID:           uuid.New().String(),
				Timestamp:    time.Now(),
			}

			// Send the event
			select {
			case eventChan <- event:
			case <-ctx.Done():
				return
			}
		}

		// Wait for context cancellation
		<-ctx.Done()
	}()

	return eventChan, nil
}

func (d *humanAgentImpl) Info() agent.Info {
	return d.AgentInfo.Info
}

func (d *humanAgentImpl) Tools() []tool.Tool {
	return []tool.Tool{}
}

func (d *humanAgentImpl) FindSubAgent(name string) agent.Agent {
	return nil
}

func (d *humanAgentImpl) SubAgents() []agent.Agent {
	return []agent.Agent{}
}

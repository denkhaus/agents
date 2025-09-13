package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type Invocation struct {
	ID         uuid.UUID
	EventsByID map[uuid.UUID][]*messaging.LLMEvent `json:"events"`
}

func (p *Invocation) AddEvent(event *messaging.LLMEvent) {
	if p.EventsByID == nil {
		p.EventsByID = make(map[uuid.UUID][]*messaging.LLMEvent)
	}

	if p.EventsByID[event.ID] == nil {
		p.EventsByID[event.ID] = []*messaging.LLMEvent{event}
	} else {
		p.EventsByID[event.ID] = append(p.EventsByID[event.ID], event)
	}
}

func (p *Invocation) EventCount() int {
	total := 0
	for _, evts := range p.EventsByID {
		total += len(evts)
	}

	return total
}

type EventCollector struct {
	Invokations map[uuid.UUID]*Invocation
}

func (p *EventCollector) AddEvent(event *messaging.LLMEvent) {
	if p.Invokations == nil {
		p.Invokations = make(map[uuid.UUID]*Invocation)
	}

	if inv, ok := p.Invokations[event.InvocationID]; ok {
		inv.AddEvent(event)
	} else {
		inv := &Invocation{
			ID: event.InvocationID,
		}
		p.Invokations[event.InvocationID] = inv
		inv.AddEvent(event)
	}
}

func (p *EventCollector) EventCount() int {
	total := 0
	for _, inv := range p.Invokations {
		total += inv.EventCount()
	}

	return total
}

func (p *EventCollector) Persist() error {

	// Serialize output struct to JSON and save to file
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output to JSON: %w", err)
	}

	// Create filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("llm_events_%s.json", timestamp)

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	logger.Log.Info("LLM events saved to file",
		zap.String("filename", filename),
		zap.Int("event_count", p.EventCount()),
	)

	return nil
}

func Test_Processor(t *testing.T) {

	os.Setenv("AGENTS_WORKSPACE_PATH", "/home/denkhaus/dev/gomodules/agents/test_workspace")
	os.Setenv("AGENTS_CONFIG_PATH", "/home/denkhaus/dev/gomodules/agents/config")
	os.Setenv("AGENTS_DATABASE_URL", "postgres://agents:agents@localhost:6888/agents?sslmode=disable")

	app, err := NewApp()
	assert.NoError(t, err)

	ctx := context.Background()
	environmentName := "production"
	routing := &messaging.RoutingInfo{
		FromAgentID: shared.AgentIDHuman,
		ToAgentID:   shared.AgentIDCoder,
		SessionID:   uuid.New(),
		Streaming:   utils.BoolPtr(false),
	}

	envName := config.EnvironmentName(environmentName)

	logger.Log.Info("Starting agents system",
		zap.String("version", appVersion),
		zap.String("environment", string(envName)),
	)

	// Validate environment exists
	envConfig, err := app.configProvider.LoadEnvironmentConfig(envName)
	if err != nil {
		assert.NoError(t, err)
	}

	logger.Log.Info("Environment loaded successfully",
		zap.String("name", envConfig.Name),
		zap.String("description", envConfig.Description),
		zap.Int("agents", len(envConfig.Agents)),
		zap.Int("roles", len(envConfig.Roles)),
		zap.Bool("condenser_enabled", envConfig.Condenser.LoggingEnabled),
	)

	// Create all agents automatically
	agents, err := app.agentFactory.CreateAllAgentsInEnvironment(ctx, envName)
	if err != nil {
		assert.NoError(t, err)
	}

	if len(agents) == 0 {
		return
	}

	// Log created agents
	for _, ag := range agents {
		logger.Log.Info("Agent ready",
			zap.String("name", ag.Info().Name),
			zap.String("role", string(ag.GetRole())),
			zap.Any("id", ag.GetID()),
		)
	}

	condenserService, err := app.createCondenser(ctx, envConfig)
	if err != nil {
		assert.NoError(t, err)
	}

	collector := EventCollector{}
	collectEvents := func(info *messaging.RoutingInfo, ev *event.Event) {
		llmEvent, err := messaging.NewLLMEvent(info, ev)
		if err != nil {
			logger.Log.Error("failed to create llm event", zap.Error(err))
			return
		}

		if llmEvent != nil {
			collector.AddEvent(llmEvent)
			logger.Log.Debug("Added event to collector",
				zap.String("event_id", llmEvent.ID.String()),
				zap.String("invocation_id", llmEvent.InvocationID.String()),
				zap.Int("total_events", collector.EventCount()),
			)
			spew.Dump(llmEvent)
		}
	}

	processor := multi.NewChatProcessor(
		routing.SessionID,
		multi.WithSessionService(condenserService),
		multi.WithApplicationName(fmt.Sprintf("%s-%s", appName, envConfig.Name)),
		multi.WithOnRawEvent(func(info *messaging.RoutingInfo, ev *event.Event) {
			collectEvents(info, ev)
		}),
		multi.WithAgents(agents...),
	)

	evts, err := processor.SendMessage(
		ctx,
		routing,
		model.NewUserMessage("Get current time with tool. Send only the time value to researcher. Researcher should respond with just 'received'."),
	)

	if err != nil {
		assert.NoError(t, err)
	}

	timer := time.NewTimer(time.Minute * 3)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			logger.Log.Info("Timer expired, persisting events")
			err := collector.Persist()
			assert.NoError(t, err)
			return
		case ev, ok := <-evts:
			if !ok {
				// Channel closed, persist remaining events
				//logger.Log.Info("Event channel closed")
				// err := collector.Persist()
				// assert.NoError(t, err)
				time.Sleep(time.Second)
				continue
			}

			collectEvents(routing, ev)
		}
	}
}

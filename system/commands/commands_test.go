package commands

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/utils"
	sys_shared "github.com/denkhaus/agents/system/shared"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

func Test_Processor(t *testing.T) {

	os.Setenv("AGENTS_WORKSPACE_PATH", "/home/denkhaus/dev/gomodules/agents/test_workspace")
	os.Setenv("AGENTS_CONFIG_PATH", "/home/denkhaus/dev/gomodules/agents/config")
	os.Setenv("AGENTS_DATABASE_URL", "postgres://agents:agents@localhost:6888/agents?sslmode=disable")

	injector := di.NewContainer()

	configProvider := do.MustInvoke[config.ConfigProvider](injector)
	agentFactory := do.MustInvoke[config.AgentFactory](injector)

	ctx := context.Background()
	appName := "test-app"
	environmentName := "production"
	routing := &messaging.RoutingInfo{
		FromAgentID: shared.AgentIDHuman,
		ToAgentID:   shared.AgentIDCoder,
		SessionID:   uuid.New(),
		Streaming:   utils.BoolPtr(true),
	}

	envName := config.EnvironmentName(environmentName)

	// Validate environment exists
	envConfig, err := configProvider.LoadEnvironmentConfig(envName)
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
	agents, err := agentFactory.CreateAllAgentsInEnvironment(ctx, envName)
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

	condenserService, err := sys_shared.CreateCondenser(ctx, envConfig)
	if err != nil {
		assert.NoError(t, err)
	}

	collector := NewEventCollector(ctx, time.Second*30)
	defer collector.Close()

	processor, err := multi.NewChatProcessor(
		routing.SessionID,
		multi.WithSessionService(condenserService),
		multi.WithApplicationName(fmt.Sprintf("%s-%s", appName, envConfig.Name)),
		multi.WithOnRawEvent(func(info *messaging.RoutingInfo, ev *event.Event) {
			collector.Collect(info, ev)
		}),
		multi.WithAgents(agents...),
	)

	if err != nil {
		assert.NoError(t, err)
	}

	evts, err := processor.SendMessage(
		ctx,
		routing,
		//model.NewUserMessage("Get current time with tool. Send only the time value to researcher. Researcher should respond with just 'received'."),
		model.NewUserMessage("Tell me the current time. My timezone is 'Europe/Berlin' "),
	)

	if err != nil {
		assert.NoError(t, err)
	}

	for {
		select {
		case <-collector.Done():
			logger.Log.Info("collection timeout -> persisting events")
			err := collector.Persist()
			assert.NoError(t, err)
			return
		case ev, ok := <-evts:
			if !ok {
				// the main channel is closed
				// wait until potential inter-agent communication has finished
				// then the timeout in the collector will trigger
				time.Sleep(time.Second)
				continue
			}

			collector.Collect(routing, ev)
		}
	}
}

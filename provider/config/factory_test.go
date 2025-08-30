package config

import (
	"context"
	"testing"

	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// Mock implementations for testing
type mockConfigProvider struct {
	mock.Mock
}

func (m *mockConfigProvider) LoadAgentComposition(environment, agentName string) (*AgentConfig, error) {
	args := m.Called(environment, agentName)
	return args.Get(0).(*AgentConfig), args.Error(1)
}

func (m *mockConfigProvider) LoadPrompt(agentName, version string) (*PromptConfig, error) {
	args := m.Called(agentName, version)
	return args.Get(0).(*PromptConfig), args.Error(1)
}

func (m *mockConfigProvider) LoadSettings(agentName, profile string) (*SettingsConfig, error) {
	args := m.Called(agentName, profile)
	return args.Get(0).(*SettingsConfig), args.Error(1)
}

func (m *mockConfigProvider) LoadToolProfile(profileName string) (*ToolsConfig, error) {
	args := m.Called(profileName)
	return args.Get(0).(*ToolsConfig), args.Error(1)
}

func (m *mockConfigProvider) ValidateConfiguration() error {
	args := m.Called()
	return args.Error(0)
}

type mockToolFactory struct {
	mock.Mock
}

func (m *mockToolFactory) CreateTools(toolsConfig ToolsConfig) ([]tool.Tool, []tool.ToolSet, error) {
	args := m.Called(toolsConfig)
	return args.Get(0).([]tool.Tool), args.Get(1).([]tool.ToolSet), args.Error(2)
}

func TestUnifiedAgentFactory_CreateAgent(t *testing.T) {
	// Create mocks
	mockConfigProvider := new(mockConfigProvider)
	mockToolFactory := new(mockToolFactory)

	// Create factory with mocks
	factory := &UnifiedAgentFactory{
		configProvider: mockConfigProvider,
		toolFactory:    mockToolFactory,
	}

	// Set up test data
	agentID := uuid.New()
	agentConfig := &AgentConfig{
		AgentID: agentID,
		Name:    "test-agent",
		Type:    shared.AgentTypeDefault,
		Prompt: PromptConfig{
			Content: "Test prompt",
		},
		Settings: SettingsConfig{
			Agent: AgentSettings{
				LLM: LLMSettings{
					Model:    "gpt-3.5-turbo",
					Provider: shared.ModelProviderOpenAI,
				},
			},
		},
	}

	// Set up mock expectations
	mockConfigProvider.On("LoadAgentComposition", "production", "test-agent").Return(agentConfig, nil)
	mockToolFactory.On("CreateTools", agentConfig.Tools).Return([]tool.Tool{}, []tool.ToolSet{}, nil)

	// Test CreateAgent
	ctx := context.Background()
	agent, err := factory.CreateAgent(ctx, "production", "test-agent")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, agent)
	
	// Verify mock expectations
	mockConfigProvider.AssertExpectations(t)
	mockToolFactory.AssertExpectations(t)
}

func TestUnifiedAgentFactory_CreateAgentByID(t *testing.T) {
	// Create mocks
	mockConfigProvider := new(mockConfigProvider)
	mockToolFactory := new(mockToolFactory)

	// Create factory with mocks
	factory := &UnifiedAgentFactory{
		configProvider: mockConfigProvider,
		toolFactory:    mockToolFactory,
	}

	// Set up test data
	agentID := shared.AgentIDCoder
	agentConfig := &AgentConfig{
		AgentID: agentID,
		Name:    "coder",
		Type:    shared.AgentTypeDefault,
		Prompt: PromptConfig{
			Content: "Coder prompt",
		},
		Settings: SettingsConfig{
			Agent: AgentSettings{
				LLM: LLMSettings{
					Model:    "gpt-3.5-turbo",
					Provider: shared.ModelProviderOpenAI,
				},
			},
		},
	}

	// Set up mock expectations
	mockConfigProvider.On("LoadAgentComposition", "production", "coder").Return(agentConfig, nil)
	mockToolFactory.On("CreateTools", agentConfig.Tools).Return([]tool.Tool{}, []tool.ToolSet{}, nil)

	// Test CreateAgentByID
	ctx := context.Background()
	agent, err := factory.CreateAgentByID(ctx, agentID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, agentID, agent.ID())
	
	// Verify mock expectations
	mockConfigProvider.AssertExpectations(t)
	mockToolFactory.AssertExpectations(t)
}
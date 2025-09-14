package project

import (
	"fmt"
	"time"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/tools"
	"github.com/denkhaus/agents/pkg/tools/project/repository/postgres"
	"github.com/samber/do"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

const (
	ToolSetName = "project_toolset"
)

// projectTaskToolSet implements the ToolSet interface for project task management
type projectTaskToolSet struct {
	isReadOnly      bool
	manager         ProjectManager
	logger          *zap.Logger
	tools           []tool.CallableTool
	availableAgents []*shared.AgentInfo
}

func NewWithDI(injector *do.Injector) (tools.ToolSetFactoryFunc, error) {
	return func(config tools.ConfigPayload, availableAgents []*shared.AgentInfo) (tool.ToolSet, error) {
		// Extract configuration and convert to options
		var tsk ToolSetConfig
		if err := config.Bind(&tsk); err != nil {
			return nil, err
		}

		if err := tsk.RepositoryType.Validate(); err != nil {
			return nil, err
		}

		// Create options from settings
		var opts []Option
		if len(availableAgents) > 0 {
			opts = append(opts, WithAvailableAgents(availableAgents))
		}

		if tsk.RepositoryType == ProjectRepositoryTypePostgres {

			if tsk.DatabaseURL == nil || *tsk.DatabaseURL == "" {
				return nil, fmt.Errorf("database url must be defined if repository type is %s", tsk.RepositoryType)
			}

			repo, err := postgres.NewRepository(
				postgres.WithLogger(logger.Log),
				postgres.WithDatabaseURL(*tsk.DatabaseURL),
				postgres.WithAutoMigrate(true),
				postgres.WithConnectionPool(25, 5),
				postgres.WithConnectionLifetime(time.Hour, time.Minute*15),
			)

			if err != nil {
				return nil, fmt.Errorf("failed to create postgres repository for project toolset: %w", err)
			}

			opts = append(opts, WithRepository(repo))
		}

		opts = append(opts, WithReadOnly(tsk.IsReadOnly))

		return New(opts...)
	}, nil
}

// NewToolSet creates a new project task management tool set
func New(opts ...Option) (tool.ToolSet, error) {
	toolSet := &projectTaskToolSet{
		manager: NewManager(DefaultConfig()),
		logger:  zap.NewNop(), // Use null logger by default to avoid interfering with chat output
	}

	// Apply options
	for _, opt := range opts {
		opt(toolSet)
	}

	if toolSet.manager == nil {
		return nil, fmt.Errorf("manager cannot be nil")
	}

	if toolSet.isReadOnly {
		// Initialize readonly tools
		toolSet.tools = []tool.CallableTool{
			toolSet.getProjectTool(),
			toolSet.listProjectsTool(),
			toolSet.getTaskTool(),
			toolSet.getProjectProgressTool(),
			toolSet.getChildTasksTool(),
			toolSet.getParentTaskTool(),
			toolSet.findNextActionableTaskTool(),
			toolSet.findTasksNeedingBreakdownTool(),
			toolSet.getRootTasksTool(),
			toolSet.listTasksByStateTool(),
			toolSet.listTasksForProjectTool(),
			toolSet.getTaskDependenciesTool(),
			toolSet.getDependentTasksTool(),
			// Agent management tools (readonly)
			toolSet.listTasksByAgentTool(),
			toolSet.listUnassignedTasksTool(),
		}
	} else {
		// Initialize all tools
		toolSet.tools = []tool.CallableTool{
			toolSet.createProjectTool(),
			toolSet.getProjectTool(),
			toolSet.updateProjectDescriptionTool(), // Add this line
			toolSet.listProjectsTool(),
			toolSet.createTaskTool(),
			toolSet.getTaskTool(),
			toolSet.updateTaskDescriptionTool(), // Add this line
			toolSet.updateTaskStateTool(),
			toolSet.getProjectProgressTool(),
			toolSet.getChildTasksTool(),
			toolSet.getParentTaskTool(),
			toolSet.findNextActionableTaskTool(),
			toolSet.findTasksNeedingBreakdownTool(),
			toolSet.getRootTasksTool(),
			toolSet.listTasksByStateTool(),
			toolSet.deleteTaskSubtreeTool(),
			toolSet.updateTaskTool(),
			toolSet.deleteTaskTool(),
			toolSet.updateProjectTool(),
			toolSet.deleteProjectTool(),
			toolSet.listTasksForProjectTool(),
			toolSet.bulkUpdateTasksTool(),
			toolSet.duplicateTaskTool(),
			toolSet.setTaskEstimateTool(),
			toolSet.addTaskDependencyTool(),
			toolSet.removeTaskDependencyTool(),
			toolSet.getTaskDependenciesTool(),
			toolSet.getDependentTasksTool(),
			// Agent management tools
			toolSet.assignTaskToAgentTool(),
			toolSet.unassignTaskFromAgentTool(),
			toolSet.listTasksByAgentTool(),
			toolSet.listUnassignedTasksTool(),
		}
	}

	return toolSet, nil
}

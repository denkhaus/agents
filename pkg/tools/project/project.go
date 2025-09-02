package project

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/tools/project/shared"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// Project management tools

// createProject performs project creation
func (pts *projectTaskToolSet) createProject(ctx context.Context, args createProjectArgs) (createProjectResult, error) {
	pts.logger.Info("Creating project", zap.String("title", args.Title))

	project, err := pts.manager.CreateProject(ctx, args.Title, args.Details)
	if err != nil {
		pts.logger.Error("Failed to create project", zap.Error(err))
		return createProjectResult{}, err
	}

	pts.logger.Info("Created project successfully", zap.String("projectID", project.ID.String()))
	return createProjectResult{
		Project: project,
		Message: "Project created successfully",
	}, nil
}

func (pts *projectTaskToolSet) createProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.createProject,
		function.WithName("create_project"),
		function.WithDescription("Create a new project for task management"),
	)
}

// getProject performs project retrieval
func (pts *projectTaskToolSet) getProject(ctx context.Context, args getProjectArgs) (*shared.Project, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Getting project", zap.String("projectID", projectID.String()))

	project, err := pts.manager.GetProject(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to get project", zap.Error(err))
		return nil, err
	}

	return project, nil
}

func (pts *projectTaskToolSet) getProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getProject,
		function.WithName("get_project"),
		function.WithDescription("Get project details by ID"),
	)
}

// updateProjectDescription performs project description update
func (pts *projectTaskToolSet) updateProjectDescription(ctx context.Context, args updateProjectDescriptionArgs) (*shared.Project, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Updating project description", zap.String("projectID", projectID.String()))

	project, err := pts.manager.UpdateProjectDescription(ctx, projectID, args.Description)
	if err != nil {
		pts.logger.Error("Failed to update project description", zap.Error(err))
		return nil, err
	}

	return project, nil
}

func (pts *projectTaskToolSet) updateProjectDescriptionTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.updateProjectDescription,
		function.WithName("update_project_description"),
		function.WithDescription("Update only the project description"),
	)
}

// listProjects performs project listing
func (pts *projectTaskToolSet) listProjects(ctx context.Context, args listProjectsArgs) (listProjectsResult, error) {
	pts.logger.Info("Listing all projects")

	projects, err := pts.manager.ListProjects(ctx)
	if err != nil {
		pts.logger.Error("Failed to list projects", zap.Error(err))
		return listProjectsResult{}, err
	}

	return listProjectsResult{
		Projects: projects,
		Count:    len(projects),
	}, nil
}

func (pts *projectTaskToolSet) listProjectsTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.listProjects,
		function.WithName("list_projects"),
		function.WithDescription("List all projects"),
	)
}

// updateProject updates a project with all fields
func (pts *projectTaskToolSet) updateProject(ctx context.Context, args updateProjectArgs) (updateProjectResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return updateProjectResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Updating project", zap.String("projectID", projectID.String()))

	project, err := pts.manager.UpdateProject(ctx, projectID, args.Title, args.Description)
	if err != nil {
		pts.logger.Error("Failed to update project", zap.Error(err))
		return updateProjectResult{}, err
	}

	pts.logger.Info("Successfully updated project", zap.String("projectID", project.ID.String()))
	return updateProjectResult{
		Project: project,
		Message: fmt.Sprintf("Successfully updated project: %s", project.ID),
	}, nil
}

// updateProjectTool creates a tool for updating a project
func (pts *projectTaskToolSet) updateProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.updateProject,
		function.WithName("update_project"),
		function.WithDescription("Update a project with all fields"),
	)
}

// deleteProject deletes a project
func (pts *projectTaskToolSet) deleteProject(ctx context.Context, args deleteProjectArgs) (deleteProjectResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return deleteProjectResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Deleting project", zap.String("projectID", projectID.String()))

	err = pts.manager.DeleteProject(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to delete project", zap.Error(err))
		return deleteProjectResult{}, err
	}

	pts.logger.Info("Successfully deleted project", zap.String("projectID", projectID.String()))
	return deleteProjectResult{
		Message: fmt.Sprintf("Successfully deleted project: %s", projectID),
	}, nil
}

// deleteProjectTool creates a tool for deleting a project
func (pts *projectTaskToolSet) deleteProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.deleteProject,
		function.WithName("delete_project"),
		function.WithDescription("Delete a project"),
	)
}

func (pts *projectTaskToolSet) getProjectProgress(ctx context.Context, args getProjectProgressArgs) (*shared.ProjectProgress, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Getting project progress", zap.String("projectID", projectID.String()))

	progress, err := pts.manager.GetProjectProgress(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to get project progress", zap.Error(err))
		return nil, err
	}

	return progress, nil
}

func (pts *projectTaskToolSet) getProjectProgressTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getProjectProgress,
		function.WithName("get_project_progress"),
		function.WithDescription("Get detailed progress metrics for a project"),
	)
}

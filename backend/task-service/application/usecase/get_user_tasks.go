package usecase

import (
	"context"
	apperrors "github.com/todoist/backend/pkg/errors"

	"github.com/todoist/backend/task-service/application/dto"
	"github.com/todoist/backend/task-service/application/mapper"
	"github.com/todoist/backend/task-service/domain"
)

type GetUserTasksUseCase struct {
	taskRepo domain.TaskRepository
}

func NewGetUserTasksUseCase(taskRepo domain.TaskRepository) *GetUserTasksUseCase {
	return &GetUserTasksUseCase{
		taskRepo: taskRepo,
	}
}

func (uc *GetUserTasksUseCase) Execute(ctx context.Context, userID string, status string, priority *int, projectID *string) ([]*dto.TaskResponse, error) {
	// Validate input
	if userID == "" {
		return nil, apperrors.NewBadRequestError("user ID is required")
	}

	// Get tasks from repository
	tasks, err := uc.taskRepo.GetByUserID(userID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get tasks", err)
	}

	// Filter tasks based on criteria
	var filteredTasks []*domain.Task
	for _, task := range tasks {
		// Filter by status
		if status != "" && task.Status != status {
			continue
		}

		// Filter by priority
		if priority != nil && task.Priority != *priority {
			continue
		}

		// Filter by project ID
		if projectID != nil {
			if task.ProjectID == nil || *task.ProjectID != *projectID {
				continue
			}
		}

		filteredTasks = append(filteredTasks, task)
	}

	// Convert to response DTOs
	var taskResponses []*dto.TaskResponse
	for _, task := range filteredTasks {
		taskResponses = append(taskResponses, mapper.ToTaskResponse(task))
	}

	return taskResponses, nil
}

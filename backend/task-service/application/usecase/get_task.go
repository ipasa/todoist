package usecase

import (
	"context"
	apperrors "github.com/todoist/backend/pkg/errors"

	"github.com/todoist/backend/task-service/application/dto"
	"github.com/todoist/backend/task-service/application/mapper"
	"github.com/todoist/backend/task-service/domain"
)

type GetTaskUseCase struct {
	taskRepo domain.TaskRepository
}

func NewGetTaskUseCase(taskRepo domain.TaskRepository) *GetTaskUseCase {
	return &GetTaskUseCase{
		taskRepo: taskRepo,
	}
}

func (uc *GetTaskUseCase) Execute(ctx context.Context, taskID, userID string) (*dto.TaskResponse, error) {
	// Validate inputs
	if taskID == "" {
		return nil, apperrors.NewBadRequestError("task ID is required")
	}
	if userID == "" {
		return nil, apperrors.NewBadRequestError("user ID is required")
	}

	// Get task from repository
	task, err := uc.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("task not found")
	}

	// Verify the task belongs to the user
	if task.UserID != userID {
		return nil, apperrors.NewForbiddenError("access denied to this task")
	}

	// Return response DTO
	return mapper.ToTaskResponse(task), nil
}

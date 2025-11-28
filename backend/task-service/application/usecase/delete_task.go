package usecase

import (
	"context"
	apperrors "github.com/todoist/backend/pkg/errors"

	"github.com/todoist/backend/task-service/domain"
)

type DeleteTaskUseCase struct {
	taskRepo domain.TaskRepository
}

func NewDeleteTaskUseCase(taskRepo domain.TaskRepository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{
		taskRepo: taskRepo,
	}
}

func (uc *DeleteTaskUseCase) Execute(ctx context.Context, taskID, userID string) error {
	// Validate inputs
	if taskID == "" {
		return apperrors.NewBadRequestError("task ID is required")
	}
	if userID == "" {
		return apperrors.NewBadRequestError("user ID is required")
	}

	// Get existing task
	task, err := uc.taskRepo.GetByID(taskID)
	if err != nil {
		return apperrors.NewNotFoundError("task not found")
	}

	// Verify task belongs to user
	if task.UserID != userID {
		return apperrors.NewForbiddenError("access denied to this task")
	}

	// Delete task
	if err := uc.taskRepo.Delete(taskID); err != nil {
		return apperrors.NewInternalError("failed to delete task", err)
	}

	return nil
}

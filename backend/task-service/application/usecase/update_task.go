package usecase

import (
	"context"
	"time"

	"github.com/todoist/backend/task-service/application/dto"
	"github.com/todoist/backend/task-service/application/mapper"
	"github.com/todoist/backend/task-service/domain"
	apperrors "github.com/todoist/backend/pkg/errors"
)

type UpdateTaskUseCase struct {
	taskRepo domain.TaskRepository
}

func NewUpdateTaskUseCase(taskRepo domain.TaskRepository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		taskRepo: taskRepo,
	}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, taskID, userID string, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	// Validate inputs
	if taskID == "" {
		return nil, apperrors.NewBadRequestError("task ID is required")
	}
	if userID == "" {
		return nil, apperrors.NewBadRequestError("user ID is required")
	}

	// Get existing task
	task, err := uc.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("task not found")
	}

	// Verify task belongs to user
	if task.UserID != userID {
		return nil, apperrors.NewForbiddenError("access denied to this task")
	}

	// Update fields if provided
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != 0 {
		task.Priority = req.Priority
	}
	if req.ProjectID != nil {
		task.ProjectID = req.ProjectID
	}
	if req.DueDate != nil {
		if *req.DueDate != "" {
			dueDate, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				return nil, apperrors.NewBadRequestError("invalid due_date format, should be RFC3339")
			}
			task.DueDate = &dueDate
		} else {
			task.DueDate = nil
		}
	}

	// Update timestamp
	task.UpdatedAt = time.Now()

	// Update in repository
	if err := uc.taskRepo.Update(task); err != nil {
		return nil, apperrors.NewInternalError("failed to update task", err)
	}

	// Return response DTO
	return mapper.ToTaskResponse(task), nil
}

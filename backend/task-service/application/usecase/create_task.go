package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/todoist/backend/task-service/application/dto"
	"github.com/todoist/backend/task-service/application/mapper"
	"github.com/todoist/backend/task-service/domain"
	apperrors "github.com/todoist/backend/pkg/errors"
)

type CreateTaskUseCase struct {
	taskRepo domain.TaskRepository
}

func NewCreateTaskUseCase(taskRepo domain.TaskRepository) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		taskRepo: taskRepo,
	}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, req dto.CreateTaskRequest, userID string) (*dto.TaskResponse, error) {
	// Validate required fields
	if req.Title == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}

	// Create task domain model
	now := time.Now()
	task := &domain.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      "pending",
		Priority:    req.Priority,
		UserID:      userID,
		ProjectID:   req.ProjectID,
		DueDate:     nil, // Parse due_date if provided
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Parse due date if provided
	if req.DueDate != nil && *req.DueDate != "" {
		dueDate, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return nil, apperrors.NewBadRequestError("invalid due_date format, should be RFC3339")
		}
		task.DueDate = &dueDate
	}

	// Create task in repository
	if err := uc.taskRepo.Create(task); err != nil {
		return nil, apperrors.NewInternalError("failed to create task", err)
	}

	// Return response DTO
	return mapper.ToTaskResponse(task), nil
}

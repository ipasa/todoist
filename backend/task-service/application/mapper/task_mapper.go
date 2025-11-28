package mapper

import (
	"time"

	"github.com/todoist/backend/task-service/application/dto"
	"github.com/todoist/backend/task-service/domain"
)

func ToTaskResponse(task *domain.Task) *dto.TaskResponse {
	response := &dto.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}

	if task.ProjectID != nil {
		response.ProjectID = task.ProjectID
	}

	if task.DueDate != nil {
		dueDateStr := task.DueDate.Format(time.RFC3339)
		response.DueDate = &dueDateStr
	}

	return response
}

package events

import (
	"time"

	"github.com/google/uuid"
)

// TaskCreated event published when a new task is created
type TaskCreated struct {
	BaseEvent
	TaskID      uuid.UUID  `json:"task_id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// NewTaskCreated creates a new TaskCreated event
func NewTaskCreated(userID, taskID, projectID uuid.UUID, title, description string, priority int, dueDate *time.Time) TaskCreated {
	return TaskCreated{
		BaseEvent:   NewBaseEvent("task.task.created", userID),
		TaskID:      taskID,
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     dueDate,
	}
}

// TaskUpdated event published when a task is updated
type TaskUpdated struct {
	BaseEvent
	TaskID    uuid.UUID              `json:"task_id"`
	ProjectID uuid.UUID              `json:"project_id"`
	Changes   map[string]interface{} `json:"changes"`
}

// NewTaskUpdated creates a new TaskUpdated event
func NewTaskUpdated(userID, taskID, projectID uuid.UUID, changes map[string]interface{}) TaskUpdated {
	return TaskUpdated{
		BaseEvent: NewBaseEvent("task.task.updated", userID),
		TaskID:    taskID,
		ProjectID: projectID,
		Changes:   changes,
	}
}

// TaskCompleted event published when a task is completed
type TaskCompleted struct {
	BaseEvent
	TaskID      uuid.UUID `json:"task_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// NewTaskCompleted creates a new TaskCompleted event
func NewTaskCompleted(userID, taskID, projectID uuid.UUID, completedAt time.Time) TaskCompleted {
	return TaskCompleted{
		BaseEvent:   NewBaseEvent("task.task.completed", userID),
		TaskID:      taskID,
		ProjectID:   projectID,
		CompletedAt: completedAt,
	}
}

// TaskDeleted event published when a task is deleted
type TaskDeleted struct {
	BaseEvent
	TaskID    uuid.UUID `json:"task_id"`
	ProjectID uuid.UUID `json:"project_id"`
}

// NewTaskDeleted creates a new TaskDeleted event
func NewTaskDeleted(userID, taskID, projectID uuid.UUID) TaskDeleted {
	return TaskDeleted{
		BaseEvent: NewBaseEvent("task.task.deleted", userID),
		TaskID:    taskID,
		ProjectID: projectID,
	}
}

// CommentAdded event published when a comment is added to a task
type CommentAdded struct {
	BaseEvent
	CommentID uuid.UUID `json:"comment_id"`
	TaskID    uuid.UUID `json:"task_id"`
	Content   string    `json:"content"`
}

// NewCommentAdded creates a new CommentAdded event
func NewCommentAdded(userID, commentID, taskID uuid.UUID, content string) CommentAdded {
	return CommentAdded{
		BaseEvent: NewBaseEvent("task.comment.added", userID),
		CommentID: commentID,
		TaskID:    taskID,
		Content:   content,
	}
}

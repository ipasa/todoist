package domain

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	Status      string
	Priority    int
	UserID      string
	ProjectID   *string
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskRepository interface {
	Create(task *Task) error
	GetByID(id string) (*Task, error)
	GetByUserID(userID string) ([]*Task, error)
	Update(task *Task) error
	Delete(id string) error
}

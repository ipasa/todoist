package domain

import "time"

type Project struct {
	ID          string
	Name        string
	Description string
	Color       string
	UserID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProjectRepository interface {
	Create(project *Project) error
	GetByID(id string) (*Project, error)
	GetByUserID(userID string) ([]*Project, error)
	Update(project *Project) error
	Delete(id string) error
}

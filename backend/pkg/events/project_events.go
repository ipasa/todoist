package events

import "github.com/google/uuid"

// ProjectCreated event published when a new project is created
type ProjectCreated struct {
	BaseEvent
	ProjectID   uuid.UUID `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
}

// NewProjectCreated creates a new ProjectCreated event
func NewProjectCreated(userID, projectID uuid.UUID, name, description, color string) ProjectCreated {
	return ProjectCreated{
		BaseEvent:   NewBaseEvent("project.project.created", userID),
		ProjectID:   projectID,
		Name:        name,
		Description: description,
		Color:       color,
	}
}

// ProjectUpdated event published when a project is updated
type ProjectUpdated struct {
	BaseEvent
	ProjectID uuid.UUID              `json:"project_id"`
	Changes   map[string]interface{} `json:"changes"`
}

// NewProjectUpdated creates a new ProjectUpdated event
func NewProjectUpdated(userID, projectID uuid.UUID, changes map[string]interface{}) ProjectUpdated {
	return ProjectUpdated{
		BaseEvent: NewBaseEvent("project.project.updated", userID),
		ProjectID: projectID,
		Changes:   changes,
	}
}

// ProjectShared event published when a project is shared with another user
type ProjectShared struct {
	BaseEvent
	ProjectID  uuid.UUID `json:"project_id"`
	SharedWith uuid.UUID `json:"shared_with"`
	Permission string    `json:"permission"` // view, edit, admin
}

// NewProjectShared creates a new ProjectShared event
func NewProjectShared(userID, projectID, sharedWith uuid.UUID, permission string) ProjectShared {
	return ProjectShared{
		BaseEvent:  NewBaseEvent("project.project.shared", userID),
		ProjectID:  projectID,
		SharedWith: sharedWith,
		Permission: permission,
	}
}

// ProjectDeleted event published when a project is deleted
type ProjectDeleted struct {
	BaseEvent
	ProjectID uuid.UUID `json:"project_id"`
}

// NewProjectDeleted creates a new ProjectDeleted event
func NewProjectDeleted(userID, projectID uuid.UUID) ProjectDeleted {
	return ProjectDeleted{
		BaseEvent: NewBaseEvent("project.project.deleted", userID),
		ProjectID: projectID,
	}
}

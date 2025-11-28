package postgres

import (
	"database/sql"
	"time"

	"github.com/todoist/backend/project-service/domain"
)

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) domain.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *domain.Project) error {
	query := `
		INSERT INTO projects (id, name, description, color, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, project.ID, project.Name, project.Description, project.Color,
		project.UserID, project.CreatedAt, project.UpdatedAt)
	return err
}

func (r *projectRepository) GetByID(id string) (*domain.Project, error) {
	query := `
		SELECT id, name, description, color, user_id, created_at, updated_at
		FROM projects WHERE id = $1
	`
	project := &domain.Project{}
	err := r.db.QueryRow(query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Color,
		&project.UserID, &project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectRepository) GetByUserID(userID string) ([]*domain.Project, error) {
	query := `
		SELECT id, name, description, color, user_id, created_at, updated_at
		FROM projects WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		project := &domain.Project{}
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Color,
			&project.UserID, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *projectRepository) Update(project *domain.Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2, color = $3, updated_at = $4
		WHERE id = $5
	`
	project.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, project.Name, project.Description, project.Color,
		project.UpdatedAt, project.ID)
	return err
}

func (r *projectRepository) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

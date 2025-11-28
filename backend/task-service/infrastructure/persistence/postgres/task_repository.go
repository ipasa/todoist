package postgres

import (
	"database/sql"
	"time"

	"github.com/todoist/backend/task-service/domain"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, priority, user_id, project_id, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query, task.ID, task.Title, task.Description, task.Status, task.Priority,
		task.UserID, task.ProjectID, task.DueDate, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *taskRepository) GetByID(id string) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, user_id, project_id, due_date, created_at, updated_at
		FROM tasks WHERE id = $1
	`
	task := &domain.Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&task.UserID, &task.ProjectID, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *taskRepository) GetByUserID(userID string) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, user_id, project_id, due_date, created_at, updated_at
		FROM tasks WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
			&task.UserID, &task.ProjectID, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *taskRepository) Update(task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, priority = $4, project_id = $5, due_date = $6, updated_at = $7
		WHERE id = $8
	`
	task.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.Priority,
		task.ProjectID, task.DueDate, task.UpdatedAt, task.ID)
	return err
}

func (r *taskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

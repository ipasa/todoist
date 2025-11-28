package dto

type CreateTaskRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	Priority    int     `json:"priority"`
	ProjectID   *string `json:"project_id"`
	DueDate     *string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Priority    int     `json:"priority"`
	ProjectID   *string `json:"project_id"`
	DueDate     *string `json:"due_date"`
}

type TaskResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Priority    int     `json:"priority"`
	UserID      string  `json:"user_id"`
	ProjectID   *string `json:"project_id"`
	DueDate     *string `json:"due_date"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

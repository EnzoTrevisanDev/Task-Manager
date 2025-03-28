package dto

// CreateTaskInput for creating a task
type CreateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ProjectID   uint   `json:"project_id" binding:"required"`
	UserID      uint   `json:"user_id"`
}

// CreateProjectInput for creating a project
type CreateProjectInput struct {
	Name string `json:"name" binding:"required"`
}

// AddUserToProjectInput for adding a user to a project
type AddUserToProjectInput struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

// UpdateUserRoleInput for updating a user's role
type UpdateUserRoleInput struct {
	Role string `json:"role" binding:"required"`
}

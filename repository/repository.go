package repository

//database ops

import (
	"work-management/models"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// User operations
func (r *Repository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *Repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// Task operations
func (r *Repository) CreateTask(task *models.Task) error {
	return r.DB.Create(task).Error
}

func (r *Repository) AssignTaskToUser(taskID, userID uint) error {
	return r.DB.Model(&models.Task{Model: gorm.Model{ID: taskID}}).Update("user_id", userID).Error
}

// Project operations
func (r *Repository) CreateProject(project *models.Project) error {
	return r.DB.Create(project).Error
}

func (r *Repository) GetProjects() ([]models.Project, error) {
	var projects []models.Project
	err := r.DB.Preload("Creator").Preload("Tasks").Preload("UserRoles.User").Find(&projects).Error
	return projects, err
}

func (r *Repository) AddUserToProject(userID, projectID uint, role string) error {
	userRole := &models.UserRole{
		UserID:    userID,
		ProjectID: projectID,
		Role:      role,
	}
	return r.DB.Create(userRole).Error
}

func (r *Repository) GetUserRole(userID, projectID uint) (string, error) {
	var userRole models.UserRole
	err := r.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&userRole).Error
	if err != nil {
		return "", err
	}
	return userRole.Role, nil
}

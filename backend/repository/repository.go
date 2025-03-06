package repository

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

func (r *Repository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *Repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *Repository) CreateTask(task *models.Task) error {
	return r.DB.Create(task).Error
}

func (r *Repository) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Preload("User").Preload("Project").Find(&tasks).Error
	return tasks, err
}

func (r *Repository) GetTaskByID(taskID uint) (*models.Task, error) {
	var task models.Task
	err := r.DB.Preload("User").Preload("Project").First(&task, taskID).Error
	return &task, err
}

func (r *Repository) UpdateTask(task *models.Task) error {
	return r.DB.Save(task).Error
}

func (r *Repository) DeleteTask(taskID uint) error {
	return r.DB.Delete(&models.Task{}, taskID).Error
}

func (r *Repository) AssignTaskToUser(taskID, userID uint) error {
	return r.DB.Model(&models.Task{Model: gorm.Model{ID: taskID}}).Update("user_id", userID).Error
}

func (r *Repository) CreateProject(project *models.Project) error {
	return r.DB.Create(project).Error
}

func (r *Repository) GetProjects() ([]models.Project, error) {
	var projects []models.Project
	err := r.DB.Preload("Creator").Preload("Tasks").Preload("UserRoles.User").Find(&projects).Error
	return projects, err
}

func (r *Repository) GetProjectByID(projectID uint) (*models.Project, error) {
	var project models.Project
	err := r.DB.Preload("Creator").Preload("Tasks").Preload("UserRoles.User").First(&project, projectID).Error
	return &project, err
}

func (r *Repository) UpdateProject(project *models.Project) error {
	return r.DB.Save(project).Error
}

func (r *Repository) DeleteProject(projectID uint) error {
	return r.DB.Delete(&models.Project{}, projectID).Error
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

func (r *Repository) UpdateUserRole(userID, projectID uint, role string) error {
	return r.DB.Model(&models.UserRole{}).Where("user_id = ? AND project_id = ?", userID, projectID).Update("role", role).Error
}

func (r *Repository) RemoveUserFromProject(userID, projectID uint) error {
	return r.DB.Where("user_id = ? AND project_id = ?", userID, projectID).Delete(&models.UserRole{}).Error
}

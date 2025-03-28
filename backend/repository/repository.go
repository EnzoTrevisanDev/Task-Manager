package repository

import (
	"time"
	"work-management/models"

	"github.com/sirupsen/logrus"
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

func (r *Repository) GetProjects(userID uint) ([]models.Project, error) {
	var projects []models.Project
	query := r.DB.
		Debug().
		Unscoped(). // Ignore soft deletes
		Joins("LEFT JOIN user_roles ON user_roles.project_id = projects.id").
		Where("user_roles.user_id = ? OR projects.creator_id = ?", userID, userID).
		Group("projects.id")

	err := query.Find(&projects).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err,
		}).Error("Failed to execute GetProjects query")
		return nil, err
	}

	for i := range projects {
		err := r.DB.
			Unscoped(). // Ignore soft deletes for preloading
			Preload("Creator").
			Preload("Tasks").
			Preload("Users.User").
			First(&projects[i], projects[i].ID).Error
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"projectID": projects[i].ID,
				"error":     err,
			}).Warn("Failed to preload data for project")
		}
	}

	logrus.WithFields(logrus.Fields{
		"userID":   userID,
		"projects": projects,
	}).Debug("Projects fetched from database")
	return projects, nil
}

func (r *Repository) GetProjectByID(projectID uint) (*models.Project, error) {
	var project models.Project
	err := r.DB.
		Unscoped().                 // Ignore soft deletes, consistent with GetProjects
		Where("id = ?", projectID). // Explicitly specify the ID condition
		Preload("Creator").
		Preload("Tasks").
		Preload("Users.User").
		First(&project).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found in GetProjectByID")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
	}).Debug("Project retrieved in GetProjectByID")
	return &project, nil
}

func (r *Repository) UpdateProject(project *models.Project) error {
	return r.DB.Save(project).Error
}

func (r *Repository) DeleteProject(projectID uint) error {
	return r.DB.Delete(&models.Project{}, projectID).Error
}

func (r *Repository) AddUserToProject(userID, projectID uint, role string) error {
	userRole := models.UserRole{
		UserID:    userID,
		ProjectID: projectID,
		Role:      role,
	}
	err := r.DB.Create(&userRole).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
			"role":      role,
			"error":     err,
		}).Error("Failed to create user_roles entry")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"userID":    userID,
		"projectID": projectID,
		"role":      role,
	}).Debug("User_roles entry created successfully")
	return nil
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

func (r *Repository) GetTasksByProjectID(projectID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.
		Where("project_id = ?", projectID).
		Preload("User").
		Preload("Project").
		Find(&tasks).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to retrieve tasks for project")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"taskCount": len(tasks),
	}).Debug("Tasks fetched for project")
	return tasks, nil
}
func (r *Repository) LogActivity(projectID, userID uint, action string) error {
	activity := models.Activity{
		ProjectID: projectID,
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}
	err := r.DB.Create(&activity).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    userID,
			"action":    action,
			"error":     err,
		}).Error("Failed to log activity")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"userID":    userID,
		"action":    action,
	}).Debug("Activity logged successfully")
	return nil
}

func (r *Repository) GetActivitiesByProjectID(projectID uint) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.DB.
		Where("project_id = ?", projectID).
		Preload("User").
		Order("timestamp DESC").
		Limit(10). // Limit to 10 recent activities
		Find(&activities).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to retrieve activities for project")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"projectID":     projectID,
		"activityCount": len(activities),
	}).Debug("Activities fetched for project")
	return activities, nil
}

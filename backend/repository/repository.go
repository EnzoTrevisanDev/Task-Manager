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
	return r.DB.Unscoped().Delete(&models.Task{}, taskID).Error
}

func (r *Repository) AssignTaskToUser(taskID, userID uint) error {
	return r.DB.Model(&models.Task{Model: gorm.Model{ID: taskID}}).Update("user_id", userID).Error
}

func (r *Repository) CreateProject(project *models.Project) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		logrus.WithFields(logrus.Fields{
			"error": tx.Error,
		}).Error("Failed to start transaction")
		return tx.Error
	}

	// Create the project
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to create project")
		return err
	}

	// Add the creator as an admin in the user_roles table
	userRole := models.UserRole{
		UserID:    project.CreatorID,
		ProjectID: project.ID,
		Role:      "admin",
	}
	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to create user role for creator")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to commit transaction")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"projectID": project.ID,
		"creatorID": project.CreatorID,
	}).Info("Project created successfully")
	return nil
}

func (r *Repository) GetProjects(userID uint) ([]models.Project, error) {
	var projects []models.Project
	query := `
        SELECT DISTINCT p.*
        FROM projects p
        LEFT JOIN user_roles ur ON ur.project_id = p.id
        WHERE (ur.user_id = ? OR p.creator_id = ?)
        AND p.deleted_at IS NULL
    `
	err := r.DB.Raw(query, userID, userID).Scan(&projects).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err,
		}).Error("Failed to fetch projects with raw SQL")
		return nil, err
	}

	if len(projects) > 0 {
		// Preload associations without refetching the projects
		for i := range projects {
			if err := r.DB.Model(&projects[i]).Association("Creator").Find(&projects[i].Creator); err != nil {
				logrus.WithFields(logrus.Fields{
					"userID":    userID,
					"projectID": projects[i].ID,
					"error":     err,
				}).Error("Failed to preload Creator for project")
				return nil, err
			}
			if err := r.DB.Model(&projects[i]).Association("Users").Find(&projects[i].Users); err != nil {
				logrus.WithFields(logrus.Fields{
					"userID":    userID,
					"projectID": projects[i].ID,
					"error":     err,
				}).Error("Failed to preload Users for project")
				return nil, err
			}
		}
	}

	logrus.WithFields(logrus.Fields{
		"userID":       userID,
		"projectCount": len(projects),
		"projects":     projects,
	}).Debug("Projects fetched from database with raw SQL")
	return projects, nil
}

func (r *Repository) GetProjectByID(projectID uint) (*models.Project, error) {
	var projects []models.Project
	err := r.DB.
		Preload("Creator").
		Preload("Tasks", "deleted_at IS NULL").
		Preload("Tasks.User").
		Preload("Tasks.Project").
		Preload("Users.User").
		Where("id = ?", projectID).
		Find(&projects).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Error fetching project")
		return nil, err
	}
	if len(projects) == 0 {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
		}).Warn("Project not found in GetProjectByID")
		return nil, gorm.ErrRecordNotFound
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"project":   projects[0],
	}).Debug("Project retrieved in GetProjectByID")
	return &projects[0], nil
}

func (r *Repository) UpdateProject(project *models.Project) error {
	return r.DB.Save(project).Error
}

func (r *Repository) DeleteProject(projectID uint) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     tx.Error,
		}).Error("Failed to start transaction")
		return tx.Error
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&models.UserRole{}).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to delete user_roles entries")
		tx.Rollback()
		return err
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&models.Task{}).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to delete tasks")
		tx.Rollback()
		return err
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&models.Activity{}).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to delete activities")
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&models.Project{}, projectID).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to delete project")
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to commit transaction")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
	}).Info("Project and related records deleted successfully")
	return nil
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
		Where("project_id = ? AND deleted_at IS NULL", projectID).
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
		Limit(10).
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

func (r *Repository) UpdateProjectOwner(projectID, newCreatorID uint) error {
	return r.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("creator_id", newCreatorID).Error
}

func (r *Repository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}

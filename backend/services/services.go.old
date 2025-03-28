package services

import (
	"errors"
	"time"
	"work-management/models"
	"work-management/repository"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

const secretKey = "Banana" // Change when going to production

type Service struct {
	Repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) CreateUser(name, email, password string) (*models.User, error) {
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	if err := s.Repo.CreateUser(user); err != nil {
		log.WithFields(log.Fields{
			"email": email,
			"error": err,
		}).Error("Failed to create user")
		return nil, err
	}
	log.WithFields(log.Fields{
		"userID": user.ID,
		"email":  email,
	}).Info("User created successfully")
	return user, nil
}

func (s *Service) Login(email, password string) (string, string, error) {
	user, err := s.Repo.FindUserByEmail(email)
	if err != nil {
		log.WithFields(log.Fields{
			"email": email,
		}).Warn("User not found")
		return "", "", errors.New("user not found")
	}
	if !user.CheckPassword(password) {
		log.WithFields(log.Fields{
			"email": email,
		}).Warn("Invalid credentials")
		return "", "", errors.New("invalid credentials")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		log.WithFields(log.Fields{
			"userID": user.ID,
			"error":  err,
		}).Error("Failed to sign access token")
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		log.WithFields(log.Fields{
			"userID": user.ID,
			"error":  err,
		}).Error("Failed to sign refresh token")
		return "", "", err
	}

	log.WithFields(log.Fields{
		"userID": user.ID,
	}).Info("User logged in successfully")
	return accessTokenString, refreshTokenString, nil
}

func (s *Service) CreateTask(title, description string, projectID, userID uint) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		UserID:      userID,
	}
	if err := s.Repo.CreateTask(task); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    userID,
			"error":     err,
		}).Error("Failed to create task")
		return nil, err
	}
	log.WithFields(log.Fields{
		"taskID":    task.ID,
		"projectID": projectID,
		"userID":    userID,
	}).Info("Task created successfully")
	return task, nil
}

func (s *Service) GetTasks() ([]models.Task, error) {
	tasks, err := s.Repo.GetTasks()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to retrieve tasks")
		return nil, err
	}
	log.Info("Tasks retrieved successfully")
	return tasks, nil
}

func (s *Service) GetTaskByID(taskID uint) (*models.Task, error) {
	task, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"error":  err,
		}).Warn("Task not found")
		return nil, err
	}
	log.WithFields(log.Fields{
		"taskID": taskID,
	}).Info("Task retrieved successfully")
	return task, nil
}

func (s *Service) UpdateTask(taskID uint, title, description string, projectID, userID uint) (*models.Task, error) {
	task, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"error":  err,
		}).Warn("Task not found for update")
		return nil, err
	}
	task.Title = title
	task.Description = description
	task.ProjectID = projectID
	task.UserID = userID
	if err := s.Repo.UpdateTask(task); err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to update task")
		return nil, err
	}
	log.WithFields(log.Fields{
		"taskID": taskID,
	}).Info("Task updated successfully")
	return task, nil
}

func (s *Service) DeleteTask(taskID uint) error {
	if err := s.Repo.DeleteTask(taskID); err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to delete task")
		return err
	}
	log.WithFields(log.Fields{
		"taskID": taskID,
	}).Info("Task deleted successfully")
	return nil
}

func (s *Service) AssignTaskToUser(taskID, userID uint) error {
	if err := s.Repo.AssignTaskToUser(taskID, userID); err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"userID": userID,
			"error":  err,
		}).Error("Failed to assign task to user")
		return err
	}
	log.WithFields(log.Fields{
		"taskID": taskID,
		"userID": userID,
	}).Info("Task assigned to user successfully")
	return nil
}

func (s *Service) CreateProject(name string, creatorID uint) (*models.Project, error) {
	project := &models.Project{
		Name:      name,
		CreatorID: creatorID,
	}
	if err := s.Repo.CreateProject(project); err != nil {
		log.WithFields(log.Fields{
			"creatorID": creatorID,
			"error":     err,
		}).Error("Failed to create project")
		return nil, err
	}
	if err := s.Repo.AddUserToProject(creatorID, project.ID, "admin"); err != nil {
		log.WithFields(log.Fields{
			"projectID": project.ID,
			"creatorID": creatorID,
			"error":     err,
		}).Error("Failed to add creator to project")
		return nil, err
	}
	log.WithFields(log.Fields{
		"projectID": project.ID,
		"creatorID": creatorID,
	}).Info("Project created successfully")
	return project, nil
}

func (s *Service) GetProjects() ([]models.Project, error) {
	projects, err := s.Repo.GetProjects()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to retrieve projects")
		return nil, err
	}
	log.Info("Projects retrieved successfully")
	return projects, nil
}

func (s *Service) GetProjectByID(projectID uint) (*models.Project, error) {
	project, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found")
		return nil, err
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
	}).Info("Project retrieved successfully")
	return project, nil
}

func (s *Service) UpdateProject(projectID uint, name string) (*models.Project, error) {
	project, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found for update")
		return nil, err
	}
	project.Name = name
	if err := s.Repo.UpdateProject(project); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to update project")
		return nil, err
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
	}).Info("Project updated successfully")
	return project, nil
}

func (s *Service) DeleteProject(projectID uint) error {
	if err := s.Repo.DeleteProject(projectID); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to delete project")
		return err
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
	}).Info("Project deleted successfully")
	return nil
}

func (s *Service) AddUserToProject(userID, projectID uint, role string) error {
	validRoles := map[string]bool{"admin": true, "editor": true, "viewer": true}
	if !validRoles[role] {
		log.WithFields(log.Fields{
			"role": role,
		}).Warn("Invalid role provided")
		return errors.New("invalid role")
	}
	if err := s.Repo.AddUserToProject(userID, projectID, role); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    userID,
			"role":      role,
			"error":     err,
		}).Error("Failed to add user to project")
		return err
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
		"userID":    userID,
		"role":      role,
	}).Info("User added to project successfully")
	return nil
}

func (s *Service) UpdateUserRole(userID, projectID uint, role string) error {
	validRoles := map[string]bool{"admin": true, "editor": true, "viewer": true}
	if !validRoles[role] {
		log.WithFields(log.Fields{
			"role": role,
		}).Warn("Invalid role provided")
		return errors.New("invalid role")
	}
	if err := s.Repo.UpdateUserRole(userID, projectID, role); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    userID,
			"role":      role,
			"error":     err,
		}).Error("Failed to update user role")
		return err
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
		"userID":    userID,
		"role":      role,
	}).Info("User role updated successfully")
	return nil
}

func (s *Service) RemoveUserFromProject(userID, projectID uint) error {
	if err := s.Repo.RemoveUserFromProject(userID, projectID); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    userID,
			"error":     err,
		}).Error("Failed to remove user from project")
		return err
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
		"userID":    userID,
	}).Info("User removed from project successfully")
	return nil
}

func (s *Service) CanModifyProject(userID, projectID uint) (bool, error) {
	role, err := s.Repo.GetUserRole(userID, projectID)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    userID,
			"error":     err,
		}).Error("Failed to check user role")
		return false, err
	}
	canModify := role == "admin" || role == "editor"
	log.WithFields(log.Fields{
		"projectID": projectID,
		"userID":    userID,
		"role":      role,
		"canModify": canModify,
	}).Info("Checked user modification permission")
	return canModify, nil
}

func (s *Service) AdminOnly(userID, projectID uint) (bool, error) {
	role, err := s.Repo.GetUserRole(userID, projectID)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    userID,
			"error":     err,
		}).Error("Failed to check admin role")
		return false, err
	}
	isAdmin := role == "admin"
	log.WithFields(log.Fields{
		"projectID": projectID,
		"userID":    userID,
		"role":      role,
		"isAdmin":   isAdmin,
	}).Info("Checked admin-only permission")
	return isAdmin, nil
}

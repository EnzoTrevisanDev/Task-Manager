package services

//business logic

import (
	"errors"
	"work-management/models"
	"work-management/repository"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo *repository.Repository
}

func NewServices(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}

// User logic
func (s *Service) CreateUser(name, email, password string) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

// task logic
func (s *Service) CreateTask(title, description string, projectID, userID uint) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		UserID:      userID, //assigned user
	}
	if err := s.Repo.CreateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

//func (s *Service) GetTasks() ([]models.Task, error) {
//	return s.Repo.GetTasks()
//}

func (s *Service) AssignTaskToUser(taskID, userID uint) error {
	return s.Repo.AssignTaskToUser(taskID, userID)
}

// Project logic
func (s *Service) CreateProject(name string, creatorID uint) (*models.Project, error) {
	project := &models.Project{
		Name:      name,
		CreatorID: creatorID,
	}
	if err := s.Repo.CreateProject(project); err != nil {
		return nil, err
	}
	// Creator gets admin role
	if err := s.Repo.AddUserToProject(creatorID, project.ID, "admin"); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *Service) GetProjects() ([]models.Project, error) {
	return s.Repo.GetProjects()
}

func (s *Service) AddUserToProject(userID, projectID uint, role string) error {
	// validate role
	validRoles := map[string]bool{"admin": true, "editor": true, "viewer": true}
	if !validRoles[role] {
		return errors.New("invalid role")
	}
	return s.Repo.AddUserToProject(userID, projectID, role)
}

func (s *Service) CanModifyProject(userID, projectID uint) (bool, error) {
	role, err := s.Repo.GetUserRole(userID, projectID)
	if err != nil {
		return false, err
	}
	return role == "admin" || role == "editor", nil
}

//creator auto-gets "admin" role on project creation

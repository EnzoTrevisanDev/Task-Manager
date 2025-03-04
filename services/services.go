package services

import (
	"errors"
	"time"
	"work-management/models"
	"work-management/repository"

	"github.com/dgrijalva/jwt-go"
)

const secretKey = "Banana"

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
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.Repo.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}
	if !user.CheckPassword(password) {
		return "", errors.New("invalid credentials")
	}

	//Create a JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), //expires in 24h
	})

	//Sigin the token w the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Service) CreateTask(title, description string, projectID, userID uint) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		UserID:      userID,
	}
	if err := s.Repo.CreateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) GetTasks() ([]models.Task, error) {
	return s.Repo.GetTasks()
}

func (s *Service) GetTaskByID(taskID uint) (*models.Task, error) {
	return s.Repo.GetTaskByID(taskID)
}

func (s *Service) UpdateTask(taskID uint, title, description string, projectID, userID uint) (*models.Task, error) {
	task, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	task.Title = title
	task.Description = description
	task.ProjectID = projectID
	task.UserID = userID
	if err := s.Repo.UpdateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) DeleteTask(taskID uint) error {
	return s.Repo.DeleteTask(taskID)
}

func (s *Service) AssignTaskToUser(taskID, userID uint) error {
	return s.Repo.AssignTaskToUser(taskID, userID)
}

func (s *Service) CreateProject(name string, creatorID uint) (*models.Project, error) {
	project := &models.Project{
		Name:      name,
		CreatorID: creatorID,
	}
	if err := s.Repo.CreateProject(project); err != nil {
		return nil, err
	}
	if err := s.Repo.AddUserToProject(creatorID, project.ID, "admin"); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *Service) GetProjects() ([]models.Project, error) {
	return s.Repo.GetProjects()
}

func (s *Service) GetProjectByID(projectID uint) (*models.Project, error) {
	return s.Repo.GetProjectByID(projectID)
}

func (s *Service) UpdateProject(projectID uint, name string) (*models.Project, error) {
	project, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	project.Name = name
	if err := s.Repo.UpdateProject(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *Service) DeleteProject(projectID uint) error {
	return s.Repo.DeleteProject(projectID)
}

func (s *Service) AddUserToProject(userID, projectID uint, role string) error {
	validRoles := map[string]bool{"admin": true, "editor": true, "viewer": true}
	if !validRoles[role] {
		return errors.New("invalid role")
	}
	return s.Repo.AddUserToProject(userID, projectID, role)
}

func (s *Service) UpdateUserRole(userID, projectID uint, role string) error {
	validRoles := map[string]bool{"admin": true, "editor": true, "viewer": true}
	if !validRoles[role] {
		return errors.New("invalid role")
	}
	return s.Repo.UpdateUserRole(userID, projectID, role)
}

func (s *Service) RemoveUserFromProject(userID, projectID uint) error {
	return s.Repo.RemoveUserFromProject(userID, projectID)
}

func (s *Service) CanModifyProject(userID, projectID uint) (bool, error) {
	role, err := s.Repo.GetUserRole(userID, projectID)
	if err != nil {
		return false, err
	}
	return role == "admin" || role == "editor", nil
}

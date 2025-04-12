// Task-related services (GetTasks, CreateTask, etc.)
package services

import (
	"time"
	"work-management/models"

	"github.com/sirupsen/logrus"
)

func (s *Service) CreateTask(title, description string, projectID, userID uint, status string, dueDate time.Time) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		UserID:      userID,
		Status:      status,
		DueDate:     dueDate,
	}
	if err := s.Repo.CreateTask(task); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    userID,
			"error":     err,
		}).Error("Failed to create task")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"taskID":    task.ID,
		"projectID": projectID,
		"userID":    userID,
	}).Info("Task created successfully")
	// Enhancement: Emit a WebSocket event for real-time updates
	// s.notifyClients("task_created", task)
	return task, nil
}

func (s *Service) GetTasks() ([]models.Task, error) {
	tasks, err := s.Repo.GetTasks()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to retrieve tasks")
		return nil, err
	}
	logrus.Info("Tasks retrieved successfully")
	return tasks, nil
}

func (s *Service) GetTaskByID(taskID uint) (*models.Task, error) {
	task, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Warn("Task not found")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"taskID": taskID,
	}).Info("Task retrieved successfully")
	return task, nil
}

func (s *Service) UpdateTask(taskID uint, title, description string, projectID, userID uint, status string, dueDate time.Time) (*models.Task, error) {
	task, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Warn("Task not found for update")
		return nil, err
	}
	task.Title = title
	task.Description = description
	task.ProjectID = projectID
	task.UserID = userID
	task.Status = status
	task.DueDate = dueDate
	if err := s.Repo.UpdateTask(task); err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to update task")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"taskID": taskID,
	}).Info("Task updated successfully")
	// Enhancement: Emit a WebSocket event for real-time updates
	// s.notifyClients("task_updated", task)
	return task, nil
}

func (s *Service) DeleteTask(taskID uint) error {
	if err := s.Repo.DeleteTask(taskID); err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to delete task")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"taskID": taskID,
	}).Info("Task deleted successfully")
	// Enhancement: Emit a WebSocket event for real-time updates
	// s.notifyClients("task_deleted", taskID)
	return nil
}

func (s *Service) AssignTaskToUser(taskID, userID uint) error {
	if err := s.Repo.AssignTaskToUser(taskID, userID); err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"userID": userID,
			"error":  err,
		}).Error("Failed to assign task to user")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"taskID": taskID,
		"userID": userID,
	}).Info("Task assigned to user successfully")
	return nil
}

func (s *Service) GetTasksByProjectID(projectID uint) ([]models.Task, error) {
	tasks, err := s.Repo.GetTasksByProjectID(projectID)
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
	}).Info("Tasks retrieved for project successfully")
	return tasks, nil
}

// Enhancement: Add a method to fetch task dependencies (future feature)
// func (s *Service) GetTaskDependencies(taskID uint) ([]models.Task, error) {
//     // Implementation for fetching dependent tasks
// }

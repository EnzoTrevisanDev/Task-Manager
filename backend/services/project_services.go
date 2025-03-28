package services

import (
	"errors"

	"work-management/models"

	"github.com/sirupsen/logrus"
)

func (s *Service) CreateProject(name, description, category, status string, creatorID uint) (*models.Project, error) {
	logrus.WithFields(logrus.Fields{
		"creatorID": creatorID,
		"name":      name,
	}).Debug("Starting CreateProject")

	// Start a transaction
	tx := s.Repo.DB.Begin()
	if tx.Error != nil {
		logrus.WithFields(logrus.Fields{
			"creatorID": creatorID,
			"error":     tx.Error,
		}).Error("Failed to start transaction")
		return nil, tx.Error
	}
	logrus.WithFields(logrus.Fields{
		"creatorID": creatorID,
	}).Debug("Transaction started")

	// Create the project
	project := models.Project{
		Name:        name,
		Description: description,
		Category:    category,
		Status:      status,
		CreatorID:   creatorID,
	}
	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"creatorID": creatorID,
			"error":     err,
		}).Error("Failed to create project in database")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"creatorID": creatorID,
		"projectID": project.ID,
	}).Debug("Project created in database")

	// Add the creator as an admin in the user_roles table
	userRole := models.UserRole{
		UserID:    creatorID,
		ProjectID: project.ID,
		Role:      "admin",
	}
	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"creatorID": creatorID,
			"projectID": project.ID,
			"error":     err,
		}).Error("Failed to add creator to user_roles")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"creatorID": creatorID,
		"projectID": project.ID,
	}).Debug("Creator added to user_roles")

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"creatorID": creatorID,
			"projectID": project.ID,
			"error":     err,
		}).Error("Failed to commit transaction")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"creatorID": creatorID,
		"projectID": project.ID,
	}).Info("Creator added to user_roles successfully")
	return &project, nil
}

func (s *Service) GetProjects(userID uint) ([]models.Project, error) {
	projects, err := s.Repo.GetProjects(userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err,
		}).Error("Failed to retrieve projects")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"userID": userID,
	}).Info("Projects retrieved successfully")
	return projects, nil
}

func (s *Service) GetProjectByID(projectID uint) (*models.Project, error) {
	project, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
	}).Info("Project retrieved successfully")
	return project, nil
}

func (s *Service) UpdateProject(projectID uint, name, description, category, status string) (*models.Project, error) {
	project, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found for update")
		return nil, err
	}
	project.Name = name
	project.Description = description
	project.Category = category
	project.Status = status
	if err := s.Repo.UpdateProject(project); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to update project")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
	}).Info("Project updated successfully")
	return project, nil
}

func (s *Service) DeleteProject(projectID uint) error {
	if err := s.Repo.DeleteProject(projectID); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to delete project")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
	}).Info("Project deleted successfully")
	return nil
}

func (s *Service) ToggleFavorite(projectID uint, isFavorite bool) (*models.Project, error) {
	project, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found for toggling favorite")
		return nil, err
	}
	project.IsFavorite = isFavorite
	if err := s.Repo.UpdateProject(project); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to toggle project favorite status")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"projectID":  projectID,
		"isFavorite": isFavorite,
	}).Info("Project favorite status updated successfully")
	return project, nil
}

func (s *Service) AddUserToProject(userID, projectID uint, role string) error {
	validRoles := map[string]bool{"admin": true, "editor": true, "viewer": true}
	if !validRoles[role] {
		logrus.WithFields(logrus.Fields{
			"role": role,
		}).Warn("Invalid role provided")
		return errors.New("invalid role")
	}
	if err := s.Repo.AddUserToProject(userID, projectID, role); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    userID,
			"role":      role,
			"error":     err,
		}).Error("Failed to add user to project")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"userID":    userID,
		"role":      role,
	}).Info("User added to project successfully")
	return nil
}

func (s *Service) UpdateUserRole(userID, projectID uint, role string) error {
	validRoles := map[string]bool{"admin": true, "editor": true, "viewer": true}
	if !validRoles[role] {
		logrus.WithFields(logrus.Fields{
			"role": role,
		}).Warn("Invalid role provided")
		return errors.New("invalid role")
	}
	if err := s.Repo.UpdateUserRole(userID, projectID, role); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    userID,
			"role":      role,
			"error":     err,
		}).Error("Failed to update user role")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"userID":    userID,
		"role":      role,
	}).Info("User role updated successfully")
	return nil
}

func (s *Service) RemoveUserFromProject(userID, projectID uint) error {
	if err := s.Repo.RemoveUserFromProject(userID, projectID); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    userID,
			"error":     err,
		}).Error("Failed to remove user from project")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"userID":    userID,
	}).Info("User removed from project successfully")
	return nil
}

func (s *Service) CanModifyProject(userID, projectID uint) (bool, error) {
	role, err := s.Repo.GetUserRole(userID, projectID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    userID,
			"error":     err,
		}).Error("Failed to check user role")
		return false, err
	}
	canModify := role == "admin" || role == "editor"
	logrus.WithFields(logrus.Fields{
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
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    userID,
			"error":     err,
		}).Error("Failed to check admin role")
		return false, err
	}
	isAdmin := role == "admin"
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"userID":    userID,
		"role":      role,
		"isAdmin":   isAdmin,
	}).Info("Checked admin-only permission")
	return isAdmin, nil
}
func (s *Service) ChangeProjectOwner(projectID, newOwnerID uint) (*models.Project, error) {
	project, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found for changing owner")
		return nil, err
	}

	// Verify the new owner exists and is part of the project
	role, err := s.Repo.GetUserRole(newOwnerID, projectID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID":  projectID,
			"newOwnerID": newOwnerID,
			"error":      err,
		}).Warn("New owner not found in project")
		return nil, errors.New("new owner must be a member of the project")
	}

	// Update the creator_id
	project.CreatorID = newOwnerID
	if err := s.Repo.UpdateProject(project); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID":  projectID,
			"newOwnerID": newOwnerID,
			"error":      err,
		}).Error("Failed to update project owner")
		return nil, err
	}

	// Ensure the new owner has admin role
	if role != "admin" {
		if err := s.Repo.UpdateUserRole(newOwnerID, projectID, "admin"); err != nil {
			logrus.WithFields(logrus.Fields{
				"projectID":  projectID,
				"newOwnerID": newOwnerID,
				"error":      err,
			}).Error("Failed to update new owner's role to admin")
			return nil, err
		}
	}

	logrus.WithFields(logrus.Fields{
		"projectID":  projectID,
		"newOwnerID": newOwnerID,
	}).Info("Project owner updated successfully")
	return project, nil
}

func (s *Service) LogActivity(projectID, userID uint, action string) error {
	return s.Repo.LogActivity(projectID, userID, action)
}
func (s *Service) GetActivitiesByProjectID(projectID uint) ([]models.Activity, error) {
	activities, err := s.Repo.GetActivitiesByProjectID(projectID)
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
	}).Info("Activities retrieved for project successfully")
	return activities, nil
}

func (s *Service) GetProjectAnalytics(projectID uint) (map[string]interface{}, error) {
	// Mock analytics data for now
	analytics := map[string]interface{}{
		"cycleTime": map[string]interface{}{
			"labels": []string{"Task 1", "Task 2", "Task 3", "Task 4", "Task 5"},
			"datasets": []map[string]interface{}{
				{
					"label":           "Cycle Time (Days)",
					"data":            []int{5, 7, 3, 6, 4},
					"backgroundColor": "#9b87f6",
				},
			},
		},
		"velocity": map[string]interface{}{
			"labels": []string{"Sprint 1", "Sprint 2", "Sprint 3", "Sprint 4"},
			"datasets": []map[string]interface{}{
				{
					"label":           "Velocity (Story Points)",
					"data":            []int{18, 27, 36, 30},
					"backgroundColor": "#2ecc71",
				},
			},
		},
		"burndown": map[string]interface{}{
			"labels": []string{"Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7"},
			"datasets": []map[string]interface{}{
				{
					"label":       "Burndown",
					"data":        []int{12, 10, 8, 6, 4, 2, 0},
					"fill":        false,
					"borderColor": "#9b87f6",
					"tension":     0.1,
				},
			},
		},
		"cumulativeFlow": map[string]interface{}{
			"labels": []string{"Week 1", "Week 2", "Week 3", "Week 4", "Week 5"},
			"datasets": []map[string]interface{}{
				{
					"label":           "To Do",
					"data":            []int{15, 12, 10, 8, 5},
					"backgroundColor": "#e74c3c",
				},
				{
					"label":           "In Progress",
					"data":            []int{0, 3, 5, 7, 8},
					"backgroundColor": "#f1c40f",
				},
				{
					"label":           "Completed",
					"data":            []int{0, 0, 0, 0, 2},
					"backgroundColor": "#2ecc71",
				},
			},
		},
		"metrics": map[string]interface{}{
			"cycleTime":    "4.2 days",
			"velocity":     "28 points",
			"defects":      3,
			"codeCoverage": "87%",
		},
	}
	return analytics, nil
}

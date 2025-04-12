package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateProject(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "POST",
		"path":   "/projects",
	}).Info("Incoming request")
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	project, err := h.Service.CreateProject(input.Name, input.Description, input.Category, input.Status, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"creatorID": userID,
			"error":     err,
		}).Error("Failed to create project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": project.ID,
		"creatorID": userID,
	}).Info("Project created successfully")
	c.JSON(http.StatusCreated, project)
}

func (h *Handler) GetProjects(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "GET",
		"path":   "/projects",
	}).Info("Incoming request")
	projects, err := h.Service.GetProjects(userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err,
		}).Error("Failed to get projects")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"userID":   userID,
		"projects": projects,
	}).Debug("Projects returned in response")
	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetProject(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "GET",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}
	project, err := h.Service.GetProjectByID(uint(projectID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found")
		SendError(c, http.StatusNotFound, "project not found")
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"creator":   project.Creator,
	}).Debug("Project details")
	c.JSON(http.StatusOK, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "PUT",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	project, err := h.Service.UpdateProject(uint(projectID), input.Name, input.Description, input.Category, input.Status)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to update project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
	}).Info("Project updated successfully")
	c.JSON(http.StatusOK, project)
}

func (h *Handler) ToggleFavorite(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "PUT",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	var input struct {
		IsFavorite bool `json:"is_favorite"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	project, err := h.Service.ToggleFavorite(uint(projectID), input.IsFavorite)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to toggle project favorite status")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID":  projectID,
		"isFavorite": input.IsFavorite,
	}).Info("Project favorite status updated successfully")
	c.JSON(http.StatusOK, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "DELETE",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.DeleteProject(uint(projectID)); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to delete project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
	}).Info("Project deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "project deleted"})
}

func (h *Handler) AddUserToProject(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "POST",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	var input struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.AddUserToProject(input.UserID, uint(projectID), input.Role); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    input.UserID,
			"role":      input.Role,
			"error":     err,
		}).Error("Failed to add user to project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"userID":    input.UserID,
		"role":      input.Role,
	}).Info("User added to project successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user added to project"})
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "PUT",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	userIDStr := c.Param("user_id")
	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userIDStr,
			"error":  err,
		}).Warn("Invalid user ID")
		SendError(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.UpdateUserRole(uint(targetUserID), uint(projectID), input.Role); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    targetUserID,
			"role":      input.Role,
			"error":     err,
		}).Error("Failed to update user role")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"userID":    targetUserID,
		"role":      input.Role,
	}).Info("User role updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user role updated"})
}

func (h *Handler) RemoveUserFromProject(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "DELETE",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	userIDStr := c.Param("user_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userIDStr,
			"error":  err,
		}).Warn("Invalid user ID")
		SendError(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.RemoveUserFromProject(uint(targetUserID), uint(projectID)); err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"userID":    targetUserID,
			"error":     err,
		}).Error("Failed to remove user from project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"userID":    targetUserID,
	}).Info("User removed from project successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user removed from project"})
}

func (h *Handler) ChangeProjectOwner(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "PUT",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}

	// Check if the user is an admin
	isAdmin, err := h.Service.AdminOnly(userID, uint(projectID))
	if err != nil || !isAdmin {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission - admin only")
		SendError(c, http.StatusForbidden, "insufficient permission - admin only")
		return
	}

	var input struct {
		NewOwnerID uint `json:"new_owner_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	project, err := h.Service.ChangeProjectOwner(uint(projectID), input.NewOwnerID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID":  projectID,
			"newOwnerID": input.NewOwnerID,
			"error":      err,
		}).Error("Failed to change project owner")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID":  projectID,
		"newOwnerID": input.NewOwnerID,
	}).Info("Project owner changed successfully")
	c.JSON(http.StatusOK, project)
}
func (h *Handler) GetProjectActivities(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "GET",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}
	activities, err := h.Service.GetActivitiesByProjectID(uint(projectID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to get activities for project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID":     projectID,
		"activityCount": len(activities),
	}).Info("Activities retrieved for project successfully")
	c.JSON(http.StatusOK, activities)
}

func (h *Handler) GetProjectAnalytics(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "GET",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}
	analytics, err := h.Service.GetProjectAnalytics(uint(projectID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to get project analytics")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
	}).Info("Project analytics retrieved successfully")
	c.JSON(http.StatusOK, analytics)
}

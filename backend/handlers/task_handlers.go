// Task-related handlers (GetTasks, CreateTask, etc.)
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetTasks(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "GET",
		"path":   "/tasks",
	}).Info("Incoming request")
	tasks, err := h.Service.GetTasks()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to get tasks")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetTask(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "GET",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskIDStr,
			"error":  err,
		}).Warn("Invalid task ID")
		SendError(c, http.StatusBadRequest, "invalid task_id")
		return
	}
	task, err := h.Service.GetTaskByID(uint(taskID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Warn("Task not found")
		SendError(c, http.StatusNotFound, "task not found")
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handler) CreateTask(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "POST",
		"path":   "/tasks",
	}).Info("Incoming request")
	var input struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		ProjectID   uint      `json:"project_id" binding:"required"`
		UserID      uint      `json:"user_id" binding:"required"`
		Status      string    `json:"status" binding:"required"`
		DueDate     time.Time `json:"due_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !CheckProjectPermission(c, h, userID, input.ProjectID) {
		return
	}

	task, err := h.Service.CreateTask(input.Title, input.Description, input.ProjectID, input.UserID, input.Status, input.DueDate)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to create task")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"taskID":    task.ID,
		"userID":    userID,
		"projectID": input.ProjectID,
	}).Info("Task created successfully")
	c.JSON(http.StatusCreated, task)
}
func (h *Handler) UpdateTask(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "PUT",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	taskIDStr := c.Param("task_id")
	taskID, ok := ParseID(c, taskIDStr, "task_id")
	if !ok {
		return
	}
	var input struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		ProjectID   uint      `json:"project_id" binding:"required"`
		UserID      uint      `json:"user_id" binding:"required"`
		Status      string    `json:"status" binding:"required"`
		DueDate     time.Time `json:"due_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !CheckProjectPermission(c, h, userID, input.ProjectID) {
		return
	}

	task, err := h.Service.UpdateTask(taskID, input.Title, input.Description, input.ProjectID, input.UserID, input.Status, input.DueDate)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to update task")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"taskID": taskID,
	}).Info("Task updated successfully")
	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	userID := c.GetUint("userID")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"method": "DELETE",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskIDStr,
			"error":  err,
		}).Warn("Invalid task ID")
		SendError(c, http.StatusBadRequest, "invalid task_id")
		return
	}

	task, err := h.Service.GetTaskByID(uint(taskID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Warn("Task not found")
		SendError(c, http.StatusNotFound, "task not found")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, task.ProjectID)
	if err != nil || !canModify {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": task.ProjectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.DeleteTask(uint(taskID)); err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to delete task")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"taskID": taskID,
	}).Info("Task deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

func (h *Handler) AssignTask(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "POST",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	taskIDStr := c.Param("task_id")
	userIDStr := c.Query("user_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskIDStr,
			"error":  err,
		}).Warn("Invalid task ID")
		SendError(c, http.StatusBadRequest, "invalid task_id")
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userIDStr,
			"error":  err,
		}).Warn("Invalid user ID")
		SendError(c, http.StatusBadRequest, "invalid user_id")
		return
	}
	if err := h.Service.AssignTaskToUser(uint(taskID), uint(userID)); err != nil {
		logrus.WithFields(logrus.Fields{
			"taskID": taskID,
			"userID": userID,
			"error":  err,
		}).Error("Failed to assign task")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"taskID": taskID,
		"userID": userID,
	}).Info("Task assigned successfully")
	c.JSON(http.StatusOK, gin.H{"message": "task assigned"})
}

func (h *Handler) GetTasksByProjectID(c *gin.Context) {
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
	tasks, err := h.Service.GetTasksByProjectID(uint(projectID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to get tasks for project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"projectID": projectID,
		"taskCount": len(tasks),
	}).Info("Tasks retrieved for project successfully")
	c.JSON(http.StatusOK, tasks)
}

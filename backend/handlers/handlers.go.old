package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"work-management/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const secretKey = "Banana"

type Handler struct {
	Service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{Service: service}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Incoming request")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("Authorization header missing")
			SendError(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn("Invalid authorization header format")
			SendError(c, http.StatusUnauthorized, "authorization header format must be Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.WithFields(log.Fields{
				"token": tokenString,
				"error": err,
			}).Warn("Invalid token")
			SendError(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["userID"].(float64))
			c.Set("userID", userID)
			log.WithFields(log.Fields{
				"userID": userID,
			}).Info("Token validated successfully")
			c.Next()
		} else {
			log.Warn("Invalid token claims")
			SendError(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
		}
	}
}

func (h *Handler) GetTasks(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": "GET",
		"path":   "/tasks",
	}).Info("Incoming request")
	tasks, err := h.Service.GetTasks()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get tasks")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetTask(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": "GET",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskIDStr,
			"error":  err,
		}).Warn("Invalid task ID")
		SendError(c, http.StatusBadRequest, "invalid task_id")
		return
	}
	task, err := h.Service.GetTaskByID(uint(taskID))
	if err != nil {
		log.WithFields(log.Fields{
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
	log.WithFields(log.Fields{
		"userID": userID,
		"method": "POST",
		"path":   "/tasks",
	}).Info("Incoming request")
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ProjectID   uint   `json:"project_id" binding:"required"`
		UserID      uint   `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, input.ProjectID)
	if err != nil || !canModify {
		log.WithFields(log.Fields{
			"userID":    userID,
			"projectID": input.ProjectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	task, err := h.Service.CreateTask(input.Title, input.Description, input.ProjectID, input.UserID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to create task")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"taskID":    task.ID,
		"userID":    userID,
		"projectID": input.ProjectID,
	}).Info("Task created successfully")
	c.JSON(http.StatusCreated, task)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	userID := c.GetUint("userID")
	log.WithFields(log.Fields{
		"userID": userID,
		"method": "PUT",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskIDStr,
			"error":  err,
		}).Warn("Invalid task ID")
		SendError(c, http.StatusBadRequest, "invalid task_id")
		return
	}
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ProjectID   uint   `json:"project_id" binding:"required"`
		UserID      uint   `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, input.ProjectID)
	if err != nil || !canModify {
		log.WithFields(log.Fields{
			"userID":    userID,
			"projectID": input.ProjectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	task, err := h.Service.UpdateTask(uint(taskID), input.Title, input.Description, input.ProjectID, input.UserID)
	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to update task")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"taskID": taskID,
	}).Info("Task updated successfully")
	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	userID := c.GetUint("userID")
	log.WithFields(log.Fields{
		"userID": userID,
		"method": "DELETE",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskIDStr,
			"error":  err,
		}).Warn("Invalid task ID")
		SendError(c, http.StatusBadRequest, "invalid task_id")
		return
	}

	task, err := h.Service.GetTaskByID(uint(taskID))
	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"error":  err,
		}).Warn("Task not found")
		SendError(c, http.StatusNotFound, "task not found")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, task.ProjectID)
	if err != nil || !canModify {
		log.WithFields(log.Fields{
			"userID":    userID,
			"projectID": task.ProjectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.DeleteTask(uint(taskID)); err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to delete task")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"taskID": taskID,
	}).Info("Task deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

func (h *Handler) AssignTask(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": "POST",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	taskIDStr := c.Param("task_id")
	userIDStr := c.Query("user_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskIDStr,
			"error":  err,
		}).Warn("Invalid task ID")
		SendError(c, http.StatusBadRequest, "invalid task_id")
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"userID": userIDStr,
			"error":  err,
		}).Warn("Invalid user ID")
		SendError(c, http.StatusBadRequest, "invalid user_id")
		return
	}
	if err := h.Service.AssignTaskToUser(uint(taskID), uint(userID)); err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"userID": userID,
			"error":  err,
		}).Error("Failed to assign task")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"taskID": taskID,
		"userID": userID,
	}).Info("Task assigned successfully")
	c.JSON(http.StatusOK, gin.H{"message": "task assigned"})
}

func (h *Handler) CreateUser(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": "POST",
		"path":   "/users",
	}).Info("Incoming request")
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.Service.CreateUser(input.Name, input.Email, input.Password)
	if err != nil {
		log.WithFields(log.Fields{
			"email": input.Email,
			"error": err,
		}).Error("Failed to create user")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"userID": user.ID,
		"email":  input.Email,
	}).Info("User created successfully")
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": "POST",
		"path":   "/login",
	}).Info("Incoming request")
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.Service.Login(input.Email, input.Password)
	if err != nil {
		log.WithFields(log.Fields{
			"email": input.Email,
			"error": err,
		}).Warn("Login failed")
		SendError(c, http.StatusUnauthorized, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"email": input.Email,
	}).Info("Login successful")
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) CreateProject(c *gin.Context) {
	userID := c.GetUint("userID")
	log.WithFields(log.Fields{
		"userID": userID,
		"method": "POST",
		"path":   "/projects",
	}).Info("Incoming request")
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	project, err := h.Service.CreateProject(input.Name, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"creatorID": userID,
			"error":     err,
		}).Error("Failed to create project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"projectID": project.ID,
		"creatorID": userID,
	}).Info("Project created successfully")
	c.JSON(http.StatusCreated, project)
}

func (h *Handler) GetProjects(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": "GET",
		"path":   "/projects",
	}).Info("Incoming request")
	projects, err := h.Service.GetProjects()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get projects")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetProject(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": "GET",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}
	project, err := h.Service.GetProjectByID(uint(projectID))
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"error":     err,
		}).Warn("Project not found")
		SendError(c, http.StatusNotFound, "project not found")
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	userID := c.GetUint("userID")
	log.WithFields(log.Fields{
		"userID": userID,
		"method": "PUT",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		log.WithFields(log.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	project, err := h.Service.UpdateProject(uint(projectID), input.Name)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to update project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
	}).Info("Project updated successfully")
	c.JSON(http.StatusOK, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	userID := c.GetUint("userID")
	log.WithFields(log.Fields{
		"userID": userID,
		"method": "DELETE",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		log.WithFields(log.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.DeleteProject(uint(projectID)); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"error":     err,
		}).Error("Failed to delete project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
	}).Info("Project deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "project deleted"})
}

func (h *Handler) AddUserToProject(c *gin.Context) {
	userID := c.GetUint("userID")
	log.WithFields(log.Fields{
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
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		log.WithFields(log.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.AddUserToProject(input.UserID, uint(projectID), input.Role); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    input.UserID,
			"role":      input.Role,
			"error":     err,
		}).Error("Failed to add user to project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
		"userID":    input.UserID,
		"role":      input.Role,
	}).Info("User added to project successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user added to project"})
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	userID := c.GetUint("userID")
	log.WithFields(log.Fields{
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
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"userID": userIDStr,
			"error":  err,
		}).Warn("Invalid user ID")
		SendError(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		log.WithFields(log.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.UpdateUserRole(uint(targetUserID), uint(projectID), input.Role); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    targetUserID,
			"role":      input.Role,
			"error":     err,
		}).Error("Failed to update user role")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
		"userID":    targetUserID,
		"role":      input.Role,
	}).Info("User role updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user role updated"})
}

func (h *Handler) RemoveUserFromProject(c *gin.Context) {
	userID := c.GetUint("userID")
	log.WithFields(log.Fields{
		"userID": userID,
		"method": "DELETE",
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")
	projectIDStr := c.Param("project_id")
	userIDStr := c.Param("user_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"projectID": projectIDStr,
			"error":     err,
		}).Warn("Invalid project ID")
		SendError(c, http.StatusBadRequest, "invalid project_id")
		return
	}
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"userID": userIDStr,
			"error":  err,
		}).Warn("Invalid user ID")
		SendError(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		log.WithFields(log.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return
	}

	if err := h.Service.RemoveUserFromProject(uint(targetUserID), uint(projectID)); err != nil {
		log.WithFields(log.Fields{
			"projectID": projectID,
			"userID":    targetUserID,
			"error":     err,
		}).Error("Failed to remove user from project")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.WithFields(log.Fields{
		"projectID": projectID,
		"userID":    targetUserID,
	}).Info("User removed from project successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user removed from project"})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": "POST",
		"path":   "/refresh",
	}).Info("Incoming request")
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := jwt.Parse(input.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		log.WithFields(log.Fields{
			"token": input.RefreshToken,
			"error": err,
		}).Warn("Invalid refresh token")
		SendError(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["userID"].(float64))
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": userID,
			"exp":    time.Now().Add(time.Hour * 1).Unix(),
		})
		newAccessTokenString, err := newAccessToken.SignedString([]byte(secretKey))
		if err != nil {
			log.WithFields(log.Fields{
				"userID": userID,
				"error":  err,
			}).Error("Failed to sign new access token")
			SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		log.WithFields(log.Fields{
			"userID": userID,
		}).Info("Access token refreshed successfully")
		c.JSON(http.StatusOK, gin.H{"access_token": newAccessTokenString})
	} else {
		log.Warn("Invalid token claims")
		SendError(c, http.StatusUnauthorized, "invalid token claims")
	}
}

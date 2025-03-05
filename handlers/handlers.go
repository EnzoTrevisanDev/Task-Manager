package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"work-management/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer <token>"})
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["userID"].(float64)) //JWT stores numbers as float64
			c.Set("userID", userID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
		}
	}
}

func (h *Handler) GetTasks(c *gin.Context) {
	tasks, err := h.Service.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task_id"})
		return
	}
	task, err := h.Service.GetTaskByID(uint(taskID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handler) CreateTask(c *gin.Context) {
	userID := c.GetUint("userID") //from authmiddleware
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ProjectID   uint   `json:"project_id" binding:"required"`
		UserID      uint   `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if the user has permission
	canModify, err := h.Service.CanModifyProject(userID, input.ProjectID)
	if err != nil || !canModify {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permission"})
		return
	}

	task, err := h.Service.CreateTask(input.Title, input.Description, input.ProjectID, input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	userID := c.GetUint("userID") //from authmiddleware
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task_id"})
		return
	}
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ProjectID   uint   `json:"project_id" binding:"required"`
		UserID      uint   `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if the user has permission
	canModify, err := h.Service.CanModifyProject(userID, input.ProjectID)
	if err != nil || !canModify {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permission"})
		return
	}

	task, err := h.Service.UpdateTask(uint(taskID), input.Title, input.Description, input.ProjectID, input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	userID := c.GetUint("userID") // From AuthMiddleware
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task_id"})
		return
	}

	// Get the project ID associated with the task to check permissions
	task, err := h.Service.GetTaskByID(uint(taskID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	// Check if the user has permission to modify the project
	canModify, err := h.Service.CanModifyProject(userID, task.ProjectID)
	if err != nil || !canModify {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permission"})
		return
	}

	if err := h.Service.DeleteTask(uint(taskID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

func (h *Handler) AssignTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	userIDStr := c.Query("user_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task_id"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	if err := h.Service.AssignTaskToUser(uint(taskID), uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task assigned"})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Service.CreateUser(input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.Service.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) CreateProject(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := uint(1) // Hardcoded—replace with JWT later
	project, err := h.Service.CreateProject(input.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, project)
}

func (h *Handler) GetProjects(c *gin.Context) {
	projects, err := h.Service.GetProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetProject(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}
	project, err := h.Service.GetProjectByID(uint(projectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	userID := c.GetUint("userID") // From AuthMiddleware
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	// Check if the user has permission to modify the project
	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permission"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := h.Service.UpdateProject(uint(projectID), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	userID := c.GetUint("userID") // From AuthMiddleware
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	// Check if the user has permission to delete the project (admin only)
	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permission"})
		return
	}

	if err := h.Service.DeleteProject(uint(projectID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project deleted"})
}

func (h *Handler) AddUserToProject(c *gin.Context) {
	userID := c.GetUint("userID") // From AuthMiddleware
	projectIDStr := c.Param("project_id")
	var input struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	// Check if the user has permission to modify project membership (admin only)
	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permission"})
		return
	}

	if err := h.Service.AddUserToProject(input.UserID, uint(projectID), input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user added to project"})
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	userID := c.GetUint("userID") // From AuthMiddleware
	projectIDStr := c.Param("project_id")
	userIDStr := c.Param("user_id")
	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	// Check if the user has permission to modify roles (admin only)
	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permission"})
		return
	}
	if err := h.Service.UpdateUserRole(uint(targetUserID), uint(projectID), input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user role updated"})
}

func (h *Handler) RemoveUserFromProject(c *gin.Context) {
	userID := c.GetUint("userID") // From AuthMiddleware
	projectIDStr := c.Param("project_id")
	userIDStr := c.Param("user_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Check if the user has permission to remove users (admin only)
	canModify, err := h.Service.CanModifyProject(userID, uint(projectID))
	if err != nil || !canModify {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permission"})
		return
	}

	if err := h.Service.RemoveUserFromProject(uint(targetUserID), uint(projectID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user removed from project"})
}

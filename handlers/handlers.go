package handlers

import (
	"net/http"
	"work-management/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{Service: service}
}

//func (h *Handler) GetTasks(c *gin.Context){
//	tasks, err := h.Service.GetTasks()
//
//}

func (h *Handler) CreateTask(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ProjectID   uint   `json:"project_id" binding:"required"`
		UserID      uint   `json:"user_id"`
	}
	if err := c.ShouldBingJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.Service.CreateTask(input.Title, input.Description, input.ProjectID, input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func
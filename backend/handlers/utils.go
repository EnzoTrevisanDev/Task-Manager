// Shared utilities (SendError, Handler struct, etc.)
package handlers

import (
	"net/http"
	"strconv"

	"work-management/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler struct to hold the service dependency
type Handler struct {
	Service *services.Service
}

// NewHandler creates a new Handler instance
func NewHandler(service *services.Service) *Handler {
	return &Handler{Service: service}
}

// SendError sends a standardized error response
func SendError(c *gin.Context, status int, message string) {
	logrus.WithFields(logrus.Fields{
		"status":  status,
		"message": message,
	}).Error("Request failed")
	c.JSON(status, gin.H{"error": message})
}

// ParseID parses a string ID into a uint
func ParseID(c *gin.Context, param, name string) (uint, bool) {
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			name:    param,
			"error": err,
		}).Warn("Invalid ID")
		SendError(c, http.StatusBadRequest, "invalid "+name)
		return 0, false
	}
	return uint(id), true
}

// CheckProjectPermission checks if the user can modify a project
func CheckProjectPermission(c *gin.Context, h *Handler, userID, projectID uint) bool {
	canModify, err := h.Service.CanModifyProject(userID, projectID)
	if err != nil || !canModify {
		logrus.WithFields(logrus.Fields{
			"userID":    userID,
			"projectID": projectID,
		}).Warn("Insufficient permission")
		SendError(c, http.StatusForbidden, "insufficient permission")
		return false
	}
	return true
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateUser(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "POST",
		"path":   "/users",
	}).Info("Incoming request")
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.Service.CreateUser(input.Name, input.Email, input.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err,
		}).Error("Failed to create user")
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"userID": user.ID,
		"email":  input.Email,
	}).Info("User created successfully")
	c.JSON(http.StatusCreated, user)
}

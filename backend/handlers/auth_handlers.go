package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"work-management/services" // Import services to access SecretKey

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Removed local secretKey constant since we're using services.SecretKey

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Incoming request")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn("Authorization header missing")
			SendError(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logrus.Warn("Invalid authorization header format")
			SendError(c, http.StatusUnauthorized, "authorization header format must be Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(services.SecretKey), nil // Use the centralized SecretKey
		})

		if err != nil || !token.Valid {
			logrus.WithFields(logrus.Fields{
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
			logrus.WithFields(logrus.Fields{
				"userID": userID,
			}).Info("Token validated successfully")
			c.Next()
		} else {
			logrus.Warn("Invalid token claims")
			SendError(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
		}
	}
}

func (h *Handler) Login(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "POST",
		"path":   "/login",
	}).Info("Incoming request")
	var input struct {
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

	accessToken, refreshToken, err := h.Service.Login(input.Email, input.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err,
		}).Warn("Login failed")
		SendError(c, http.StatusUnauthorized, err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"email": input.Email,
	}).Info("Login successful")
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": "POST",
		"path":   "/refresh",
	}).Info("Incoming request")
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Invalid input")
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := jwt.Parse(input.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(services.SecretKey), nil // Use the centralized SecretKey
	})
	if err != nil || !token.Valid {
		logrus.WithFields(logrus.Fields{
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
			"exp":    time.Now().Add(time.Hour * 24).Unix(), // Increase to 24 hours
		})
		newAccessTokenString, err := newAccessToken.SignedString([]byte(services.SecretKey)) // Use the centralized SecretKey
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"userID": userID,
				"error":  err,
			}).Error("Failed to sign new access token")
			SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Info("Access token refreshed successfully")
		c.JSON(http.StatusOK, gin.H{"access_token": newAccessTokenString})
	} else {
		logrus.Warn("Invalid token claims")
		SendError(c, http.StatusUnauthorized, "invalid token claims")
	}
}

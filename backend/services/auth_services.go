package services

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

func (s *Service) Login(email, password string) (string, string, error) {
	user, err := s.Repo.FindUserByEmail(email)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email": email,
		}).Warn("User not found")
		return "", "", errors.New("user not found")
	}
	if !user.CheckPassword(password) {
		logrus.WithFields(logrus.Fields{
			"email": email,
		}).Warn("Invalid credentials")
		return "", "", errors.New("invalid credentials")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(SecretKey)) // Use the centralized SecretKey
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": user.ID,
			"error":  err,
		}).Error("Failed to sign access token")
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(SecretKey)) // Use the centralized SecretKey
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": user.ID,
			"error":  err,
		}).Error("Failed to sign refresh token")
		return "", "", err
	}

	logrus.WithFields(logrus.Fields{
		"userID": user.ID,
	}).Info("User logged in successfully")
	return accessTokenString, refreshTokenString, nil
}

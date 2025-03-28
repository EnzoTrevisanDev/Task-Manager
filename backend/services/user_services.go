// User-related services (CreateUser)
package services

import (
	"work-management/models"

	"github.com/sirupsen/logrus"
)

func (s *Service) CreateUser(name, email, password string) (*models.User, error) {
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	if err := s.Repo.CreateUser(user); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": email,
			"error": err,
		}).Error("Failed to create user")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"userID": user.ID,
		"email":  email,
	}).Info("User created successfully")
	return user, nil
}

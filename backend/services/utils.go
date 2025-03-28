package services

import (
	"work-management/repository"
)

// Service struct to hold the repository dependency
type Service struct {
	Repo *repository.Repository
}

// NewService creates a new Service instance
func NewService(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}

// SecretKey for JWT signing (move to config in production)
const SecretKey = "Banana"

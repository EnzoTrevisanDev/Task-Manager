package graph

import (
	"context"
	"task-manager/database"
	"task-manager/graph/model"
)

// Resolver struct (pode conter dependências, como DB)
type Resolver struct{}

// QueryResolver (implementação)
func (r *Resolver) Users(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	result := database.DB.Preload("Tasks").Find(&users)
	return users, result.Error
}

func (r *Resolver) User(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	result := database.DB.Preload("Tasks").First(&user, id)
	return &user, result.Error
}

// MutationResolver (implementação)
func (r *Resolver) CreateUser(ctx context.Context, fullname string, email string, password string) (*model.User, error) {
	user := &model.User{FullName: fullname, Email: email, Password: password}
	result := database.DB.Create(user)
	return user, result.Error
}

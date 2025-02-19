package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Statys      string `gorm:"default:'pending'"`
	UserID      uint   // Foreign key
}

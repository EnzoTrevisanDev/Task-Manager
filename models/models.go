package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique;index"`
	Password string //hashed with bcrypt
	Tasks    []Task
	Projects []Project
}

type Task struct {
	gorm.Model
	Title       string
	Description string
	UserID      uint `gorm:"index"` //assigned to
	User        User
	ProjectID   uint `gorm:"index"`
	Project     Project
}

type Project struct {
	gorm.Model
	Name      string
	CreatorID uint `gorm:"index"`
	Creator   User
	Tasks     []Task
	UserRoles []UserRole
}

type UserRole struct {
	gorm.Model
	UserID    uint `gorm:"index"`
	ProjectID uint `gorm:"index"`
	User      User
	Project   Project
	Role      string `gorm:"type:enum('admin', 'editor', 'viewer');default:'viewer'"`
}

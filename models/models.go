package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique;index"`
	Password string `gorm:"type:varchar(255)"`
	Tasks    []Task `gorm:"foreignKey:UserID"`
}

type Task struct {
	gorm.Model
	Title       string
	Description string
	UserID      uint    `gorm:"index"`
	User        User    `gorm:"foreignKey:UserID"`
	ProjectID   uint    `gorm:"index"`
	Project     Project `gorm:"foreignKey:ProjectID"`
}

type Project struct {
	gorm.Model
	Name      string
	CreatorID uint       `gorm:"index"`
	Creator   User       `gorm:"foreignKey:CreatorID"`
	Tasks     []Task     `gorm:"foreignKey:ProjectID"`
	UserRoles []UserRole `gorm:"foreignKey:ProjectID"`
}

type UserRole struct {
	gorm.Model
	UserID    uint    `gorm:"index"`
	ProjectID uint    `gorm:"index"`
	User      User    `gorm:"foreignKey:UserID"`
	Project   Project `gorm:"foreignKey:ProjectID"`
	Role      string  `gorm:"type:varchar(20);default:'viewer'"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

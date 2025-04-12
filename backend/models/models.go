package models

import (
	"time"

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
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ProjectID   uint      `json:"project_id"`
	UserID      uint      `json:"user_id"`
	Status      string    `json:"status"`   // New field for task status (e.g., "To Do", "In Progress", "Completed")
	DueDate     time.Time `json:"due_date"` // New field for due date
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Project     Project   `json:"project" gorm:"foreignKey:ProjectID"`
}
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Status      string         `json:"status"`
	IsFavorite  bool           `json:"is_favorite"`
	CreatorID   uint           `json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator"`
	Tasks       []Task         `json:"tasks"`
	Users       []UserRole     `gorm:"foreignKey:ProjectID"`
	CreatedAt   gorm.DeletedAt `json:"created_at"`
	UpdatedAt   gorm.DeletedAt `json:"updated_at"`
}

type Activity struct {
	gorm.Model
	ProjectID uint      `json:"project_id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
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

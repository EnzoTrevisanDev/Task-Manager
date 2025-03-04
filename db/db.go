package db

import (
	"work-management/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "host=localhost port=5432 user=admin password=1234 dbname=task_manager sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	err = db.AutoMigrate(&models.User{}, &models.Task{}, &models.Project{}, &models.UserRole{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
	return db
}

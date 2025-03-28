package db

import (
	"work-management/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect establishes a connection to the database and returns a *gorm.DB instance.
func Connect() *gorm.DB {
	dsn := "host=localhost port=5432 user=admin password=1234 dbname=task_manager sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Run database migrations
	err = db.AutoMigrate(&models.User{}, &models.Task{}, &models.Project{}, &models.UserRole{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	return db
}

// CloseDB closes the underlying database connection for a given *gorm.DB instance.
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

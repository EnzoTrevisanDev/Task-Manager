package db

import (
	"fmt"
	"work-management/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=admin password=1234 dbname=task_manager"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	} else {
		fmt.Println("Sucesso db con")
	}
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Project{})
	fmt.Println(db)
	return db
}

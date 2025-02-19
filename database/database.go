package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=admin password=1234 dbname=task_manager port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	fmt.Println("Connect to Postgresql")

	// Auto migrate tables
	err = db.AutoMigrate(&User{}, &Task{})
	if err != nil {
		log.Fatal("Error migrating tables: ", err)
	}

	DB = db
}

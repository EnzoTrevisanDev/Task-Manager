package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx"
)

var DB *pgx.Conn

func ConnectPostgres() {
	dbURL := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	fmt.Println("Connected to PostgreSQL")
	DB = conn
}

package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	//load environment variables
	if err := gotoenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	//initialize PostgreSQL and Redis
	database.ConnectPostgres()
	database.ConnectRedis()

	//Fiber instance
	app := fiber.New()

	// REST API Routes
	api := app.Group("/api/v1")
	api.Get("/tasks", handlers.GetTasks)
	api.Post("/tasks", handlers.CreateTasks)

	//Websockets
	app.Get("/ws", handlers.HandleWebsocket)

	//GraphQL
	app.Post("/graphql", handlers.GraphqlHandler)

	//Start server
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))

}

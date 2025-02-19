package main

import (
	"log"
	"task-manager/database"
	"task-manager/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func main() {
	// connect to database
	database.ConnectDB()

	app := fiber.New()

	// GraphQL Server
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	app.Post("/query", adaptor.HTTPHandler(h))

	// GraphqL Playground (using fiber adapter)
	app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query")))

	log.Println("🚀 Server running on http://localhost:8080")
	app.Listen(":8080")
}

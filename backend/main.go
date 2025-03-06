package main

import (
	"os"
	"work-management/db"
	"work-management/handlers"
	"work-management/repository"
	"work-management/services"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Initialize logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	dbConn := db.Connect()
	repo := repository.NewRepository(dbConn)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

	r := gin.Default()

	// Public routes
	r.POST("/users", handler.CreateUser)
	r.POST("/login", handler.Login)
	r.POST("/refresh", handler.RefreshToken)

	// Protected routes
	protected := r.Group("/", handlers.AuthMiddleware())
	// Tasks
	protected.GET("/tasks", handler.GetTasks)
	protected.GET("/tasks/:task_id", handler.GetTask)
	protected.POST("/tasks", handler.CreateTask)
	protected.PUT("/tasks/:task_id", handler.UpdateTask)
	protected.DELETE("/tasks/:task_id", handler.DeleteTask)
	protected.POST("/tasks/:task_id/assign", handler.AssignTask)
	// Projects
	protected.GET("/projects", handler.GetProjects)
	protected.GET("/projects/:project_id", handler.GetProject)
	protected.POST("/projects", handler.CreateProject)
	protected.PUT("/projects/:project_id", handler.UpdateProject)
	protected.DELETE("/projects/:project_id", handler.DeleteProject)
	protected.POST("/projects/:project_id/users", handler.AddUserToProject)
	protected.PUT("/projects/:project_id/users/:user_id", handler.UpdateUserRole)
	protected.DELETE("/projects/:project_id/users/:user_id", handler.RemoveUserFromProject)

	log.Println("Server running at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

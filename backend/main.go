package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"work-management/db"
	"work-management/handlers"
	"work-management/repository"
	"work-management/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	//log.SetLevel(log.InfoLevel)
	log.SetLevel(log.DebugLevel)
	// Initialize database connection
	dbConn := db.Connect()
	if dbConn == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer func() {
		if err := db.CloseDB(dbConn); err != nil {
			log.Error("Failed to close database connection: ", err)
		} else {
			log.Info("Database connection closed successfully")
		}
	}()

	// Initialize repository, service, and handler
	repo := repository.NewRepository(dbConn)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

	// Set up Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Public routes (no authentication required)
	r.POST("/users", handler.CreateUser)
	r.POST("/login", handler.Login)
	r.POST("/refresh", handler.RefreshToken)

	// Protected routes (require authentication via AuthMiddleware)
	protected := r.Group("/", handlers.AuthMiddleware())
	// Task routes
	protected.GET("/tasks", handler.GetTasks)
	protected.GET("/tasks/:task_id", handler.GetTask)
	protected.POST("/tasks", handler.CreateTask)
	protected.PUT("/tasks/:task_id", handler.UpdateTask)
	protected.DELETE("/tasks/:task_id", handler.DeleteTask)
	protected.POST("/tasks/:task_id/assign", handler.AssignTask)
	protected.GET("/projects/:project_id/tasks", handler.GetTasksByProjectID)
	protected.GET("/projects/:project_id/activities", handler.GetProjectActivities)
	protected.GET("/projects/:project_id/analytics", handler.GetProjectAnalytics)
	// Project routes
	protected.GET("/projects", handler.GetProjects)
	protected.GET("/projects/:project_id", handler.GetProject)
	protected.POST("/projects", handler.CreateProject)
	protected.PUT("/projects/:project_id", handler.UpdateProject)
	protected.PUT("/projects/:project_id/favorite", handler.ToggleFavorite) // Added new route
	protected.DELETE("/projects/:project_id", handler.DeleteProject)
	protected.POST("/projects/:project_id/users", handler.AddUserToProject)
	protected.PUT("/projects/:project_id/users/:user_id", handler.UpdateUserRole)
	protected.DELETE("/projects/:project_id/users/:user_id", handler.RemoveUserFromProject)
	protected.PUT("/projects/:project_id/owner", handler.ChangeProjectOwner)

	// Create an HTTP server with the Gin router
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		log.Info("Server running at :" + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed: ", err)
	}

	log.Info("Server exited gracefully")
}

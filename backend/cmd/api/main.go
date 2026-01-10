package main

import (
	"log"
	"os"

	"fullstack-backend/internal/config"
	"fullstack-backend/internal/database"
	"fullstack-backend/internal/handlers"
	"fullstack-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"database": "connected",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			authHandler := handlers.NewAuthHandler(db)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// Task routes
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			taskHandler := handlers.NewTaskHandler(db)
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.POST("", taskHandler.CreateTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

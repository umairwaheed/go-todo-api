package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/umair/go-todo-api/database"
	"github.com/umair/go-todo-api/handlers"
	"github.com/umair/go-todo-api/models"
)

func main() {
	// Initialize database
	dbPath := "todo.db"
	if os.Getenv("DB_PATH") != "" {
		dbPath = os.Getenv("DB_PATH")
	}

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB(db)

	// Initialize models and handlers
	todoModel := models.NewTodoModel(db)
	todoHandler := handlers.NewTodoHandler(todoModel)

	// Set up Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Todo routes
		todos := api.Group("/todos")
		{
			todos.GET("", todoHandler.GetTodos)
			todos.GET("/:id", todoHandler.GetTodo)
			todos.POST("", todoHandler.CreateTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
			todos.PATCH("/:id/complete", todoHandler.CompleteTodo)
			todos.PATCH("/:id/uncomplete", todoHandler.UncompleteTodo)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Todo API is running",
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Go Todo API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health": "/health",
				"todos":  "/api/v1/todos",
			},
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("API available at http://localhost:%s", port)
	log.Printf("Health check at http://localhost:%s/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

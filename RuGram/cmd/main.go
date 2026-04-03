package main

import (
	"database/sql"
	"log"
	"os"

	"rugram-api/internal/config"
	"rugram-api/internal/database"
	"rugram-api/internal/handlers"
	"rugram-api/internal/repository"
	"rugram-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    
    // Connect to database
    db, err := database.Connect(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Initialize repository, service, and handler
    postRepo := repository.NewPostRepository(db)
    postService := service.NewPostService(postRepo)
    postHandler := handlers.NewPostHandler(postService)
    
    // Setup Gin router
    router := gin.Default()
    
    // Middleware
    if cfg.AppEnv == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // Routes
    api := router.Group("/api/v1")
    {
        posts := api.Group("/posts")
        {
            posts.GET("", postHandler.GetAll)
            posts.GET("/:id", postHandler.GetByID)
            posts.POST("", postHandler.Create)
            posts.PUT("/:id", postHandler.Update)
            posts.PATCH("/:id", postHandler.Update)
            posts.DELETE("/:id", postHandler.Delete)
            posts.GET("/user/:userId", postHandler.GetByUserID)
        }
    }
    
    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "service": "rugram-api",
        })
    })
    
    // Start server
    log.Printf("Server starting on port %s", cfg.AppPort)
    if err := router.Run(":" + cfg.AppPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func runMigrations(db *sql.DB) error {
    // Read migration files
    migrationFiles := []string{
        "internal/database/migrations/001_create_posts_table.sql",
    }
    
    for _, migrationFile := range migrationFiles {
        content, err := os.ReadFile(migrationFile)
        if err != nil {
            return err
        }
        
        _, err = db.Exec(string(content))
        if err != nil {
            return err
        }
        
        log.Printf("Successfully ran migration: %s", migrationFile)
    }
    
    return nil
}
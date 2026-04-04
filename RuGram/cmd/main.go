package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

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
		log.Printf("WARNING: Failed to run migrations: %v", err)
		log.Printf("You may need to create tables manually")
	}

	// Initialize repository, service, and handler
	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	// Setup Gin router
	router := gin.Default()

	// Recovery middleware to catch panics
	router.Use(gin.Recovery())

	// Middleware для логирования ошибок
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			log.Printf("Error in request %s: %v", c.Request.URL, c.Errors)
		}
	})

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
		// Проверяем подключение к БД
		if err := db.Ping(); err != nil {
			c.JSON(503, gin.H{
				"status":  "unhealthy",
				"service": "rugram-api",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  "ok",
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
	migrationFiles := []string{
		"internal/database/migrations/001_create_posts_table.sql",
	}

	for _, migrationFile := range migrationFiles {
		// Проверяем существует ли файл
		if _, err := os.Stat(migrationFile); os.IsNotExist(err) {
			log.Printf("Migration file not found: %s", migrationFile)
			continue
		}

		content, err := os.ReadFile(migrationFile)
		if err != nil {
			return err
		}

		// Разделяем SQL на отдельные statements
		statements := strings.Split(string(content), ";")

		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err = db.Exec(stmt)
			if err != nil {
				log.Printf("Error executing migration statement: %v", err)
				// Продолжаем выполнение, не прерываем
			}
		}

		log.Printf("Successfully ran migration: %s", migrationFile)
	}

	return nil
}

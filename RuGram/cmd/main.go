package main

import (
	"database/sql"
	"log"
	"os"

	"rugram-api/internal/cache"
	"rugram-api/internal/config"
	"rugram-api/internal/database"
	"rugram-api/internal/handlers"
	"rugram-api/internal/middleware"
	"rugram-api/internal/repository"
	"rugram-api/internal/service"
	"rugram-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    
    // Set Gin mode
    if cfg.AppEnv == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // Connect to database
    db, err := database.Connect(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Connect to Redis
    cacheSvc, err := cache.NewCacheService()
    if err != nil {
        log.Printf("Warning: Redis connection failed: %v. Cache will be disabled.", err)
    } else {
        defer cacheSvc.Close()
        log.Println("Redis cache service initialized")
    }
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    tokenRepo := repository.NewTokenRepository(db)
    postRepo := repository.NewPostRepository(db)
    
    // Initialize services
    authService := service.NewAuthService(userRepo, tokenRepo, cacheSvc)
    oauthService := service.NewOAuthService(userRepo, tokenRepo)
    userService := service.NewUserService(userRepo)
    postService := service.NewPostService(postRepo, cacheSvc)
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService, oauthService)
    userHandler := handlers.NewUserHandler(userService)
    postHandler := handlers.NewPostHandler(postService)
    
    // Setup Gin router
    router := gin.Default()
    
    // Trust proxy - для продакшена можно указать конкретные прокси
    // router.SetTrustedProxies([]string{"127.0.0.1"})
    
    // Routes
    api := router.Group("/api/v1")
    {
        // Auth routes (public)
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.Refresh)
            auth.GET("/whoami", middleware.AuthMiddleware(authService), authHandler.Whoami)
            auth.POST("/logout", middleware.AuthMiddleware(authService), authHandler.Logout)
            auth.POST("/logout-all", middleware.AuthMiddleware(authService), authHandler.LogoutAll)
            
            // OAuth routes
            auth.GET("/oauth/yandex", authHandler.OAuthYandex)
            auth.GET("/oauth/yandex/callback", authHandler.OAuthYandexCallback)
            auth.GET("/oauth/vk", authHandler.OAuthVK)
            auth.GET("/oauth/vk/callback", authHandler.OAuthVKCallback)
        }
        
        // User routes (protected)
        users := api.Group("/users")
        users.Use(middleware.AuthMiddleware(authService))
        {
            users.GET("", userHandler.GetAll)
            users.GET("/:id", userHandler.GetByID)
            users.GET("/email/:email", userHandler.GetByEmail)
            users.PUT("/:id", userHandler.Update)
            users.PATCH("/:id", userHandler.Update)
            users.DELETE("/:id", userHandler.Delete)
        }
        
        // Post routes (protected)
        posts := api.Group("/posts")
        posts.Use(middleware.AuthMiddleware(authService))
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
    router.Any("/health", func(c *gin.Context) {
        utils.SuccessResponse(c, 200, gin.H{
            "status":  "ok",
            "service": "rugram-api",
        })
    })
    
    // Start server
    log.Printf("Server starting on port %s in %s mode", cfg.AppPort, cfg.AppEnv)
    if err := router.Run(":" + cfg.AppPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func runMigrations(db *sql.DB) error {
    migrationFiles := []string{
        "internal/database/migrations/001_create_users_table.sql",
        "internal/database/migrations/002_create_user_tokens_table.sql",
        "internal/database/migrations/003_create_posts_table.sql",
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
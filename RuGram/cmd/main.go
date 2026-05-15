package main

import (
	"log"

	"rugram-api/internal/cache"
	"rugram-api/internal/config"
	"rugram-api/internal/database"
	"rugram-api/internal/handlers"
	"rugram-api/internal/middleware"
	"rugram-api/internal/queue"
	"rugram-api/internal/repository"
	"rugram-api/internal/service"
	"rugram-api/pkg/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "rugram-api/docs"
)

// @title RuGram API
// @version 4.0
// @description API для социальной сети с поддержкой JWT аутентификации, MinIO файлового хранилища, RabbitMQ асинхронной обработки и MongoDB
// @termsOfService http://swagger.io/terms/

// @contact.name Gleb, Pavel
// @contact.email support@rugram.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:4200
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in cookie
// @name access_token
// @description JWT access token stored in cookie

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Set Gin mode
    if cfg.AppEnv == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Connect to MongoDB
    mongoDB, err := database.Connect(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer mongoDB.Close()

    // Connect to Redis
    cacheSvc, err := cache.NewCacheService()
    if err != nil {
        log.Printf("Warning: Redis connection failed: %v. Cache will be disabled.", err)
        cacheSvc = &cache.CacheService{}
    } else {
        defer cacheSvc.Close()
        log.Println("Redis cache service initialized")
    }

    // Connect to MinIO
    minioSvc, err := service.NewMinioService()
    if err != nil {
        log.Printf("Warning: MinIO connection failed: %v. File upload will be disabled.", err)
    }

    // Connect to RabbitMQ
    rabbitMQSvc, err := service.NewRabbitMQService()
    if err != nil {
        log.Printf("Warning: RabbitMQ connection failed: %v. Async events will be disabled.", err)
    } else {
        defer rabbitMQSvc.Close()
        log.Println("RabbitMQ service initialized")
    }

    // Initialize email service
    emailSvc, err := service.NewEmailService()
    if err != nil {
        log.Printf("Warning: SMTP not configured: %v. Welcome emails will be disabled.", err)
    }

    // Initialize repositories
    userRepo := repository.NewUserRepository(mongoDB.Database)
    tokenRepo := repository.NewTokenRepository(mongoDB.Database)
    postRepo := repository.NewPostRepository(mongoDB.Database)
    fileRepo := repository.NewFileRepository(mongoDB.Database)

    // Initialize services
    authService := service.NewAuthService(userRepo, tokenRepo, cacheSvc, rabbitMQSvc)
    oauthService := service.NewOAuthService(userRepo, tokenRepo)
    userService := service.NewUserService(userRepo, fileRepo, minioSvc)
    postService := service.NewPostService(postRepo, cacheSvc)
    
    var fileService *service.FileService
    if minioSvc != nil {
        fileService, _ = service.NewFileService(fileRepo, minioSvc, cacheSvc)
    }

    // Initialize event handler and consumer
    eventHandler := handlers.NewEventHandler(emailSvc)
    consumer := queue.NewConsumer(rabbitMQSvc, eventHandler)

    // Start RabbitMQ consumer in background
    if rabbitMQSvc != nil && rabbitMQSvc.IsConnected() {
        go func() {
            if err := consumer.Start(); err != nil {
                log.Printf("Failed to start RabbitMQ consumer: %v", err)
            }
        }()
    }

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService, oauthService)
    userHandler := handlers.NewUserHandler(userService)
    postHandler := handlers.NewPostHandler(postService)
    
    var fileHandler *handlers.FileHandler
    if fileService != nil {
        fileHandler = handlers.NewFileHandler(fileService)
    }

    // Setup Gin router
    router := gin.Default()

    // Swagger documentation (only in development)
    if cfg.AppEnv != "production" {
        router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
        log.Println("Swagger UI available at http://localhost:" + cfg.AppPort + "/swagger/index.html")
    }

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

        // Profile routes (protected)
        profile := api.Group("/profile")
        profile.Use(middleware.AuthMiddleware(authService))
        {
            profile.GET("", userHandler.GetProfile)
            profile.POST("", userHandler.UpdateProfile)
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

        // File routes (protected)
        if fileHandler != nil {
            files := api.Group("/files")
            files.Use(middleware.AuthMiddleware(authService))
            {
                files.POST("", fileHandler.UploadFile)
                files.GET("", fileHandler.GetMyFiles)
                files.GET("/:id", fileHandler.GetFile)
                files.DELETE("/:id", fileHandler.DeleteFile)
            }
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
package main

import (
	"RuGramm/internal/config"
	"RuGramm/internal/handler"
	"RuGramm/internal/repository"
	"RuGramm/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnectionString()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация репозитория, сервиса и хендлера
	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postService)

	// Настройка роутера
	router := gin.Default()

	// Routes
	api := router.Group("/api/v1")
	{
		posts := api.Group("/posts")
		{
			posts.POST("/", postHandler.Create)
			posts.GET("/", postHandler.GetAll)
			posts.GET("/:id", postHandler.GetByID)
			posts.PUT("/:id", postHandler.Update)
			posts.PATCH("/:id", postHandler.Update)
			posts.DELETE("/:id", postHandler.Delete)
		}
	}

	// Запуск сервера
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

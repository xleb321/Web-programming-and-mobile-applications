package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// MongoDB настройки
	MongoURI      string
	MongoDatabase string

	AppPort      string
	AppEnv       string
	DefaultPage  int
	DefaultLimit int
	MaxLimit     int

	// Redis настройки
	RedisHost     string
	RedisPort     string
	RedisPassword string
	CacheTTL      int

	// JWT настройки
	JWTAccessSecret  string
	JWTRefreshSecret string

	// OAuth настройки
	YandexClientID     string
	YandexClientSecret string
	YandexRedirectURI  string
	VKClientID         string
	VKClientSecret     string
	VKRedirectURI      string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	defaultPage, _ := strconv.Atoi(getEnv("DEFAULT_PAGE", "1"))
	defaultLimit, _ := strconv.Atoi(getEnv("DEFAULT_LIMIT", "10"))
	maxLimit, _ := strconv.Atoi(getEnv("MAX_LIMIT", "100"))
	cacheTTL, _ := strconv.Atoi(getEnv("CACHE_TTL_DEFAULT", "300"))

	return &Config{
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("DB_NAME", "rugram_db"),
		AppPort:       getEnv("APP_PORT", "4200"),
		AppEnv:        getEnv("APP_ENV", "development"),
		DefaultPage:   defaultPage,
		DefaultLimit:  defaultLimit,
		MaxLimit:      maxLimit,

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		CacheTTL:      cacheTTL,

		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "default_access_secret_change_in_prod"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "default_refresh_secret_change_in_prod"),

		YandexClientID:     getEnv("YANDEX_CLIENT_ID", ""),
		YandexClientSecret: getEnv("YANDEX_CLIENT_SECRET", ""),
		YandexRedirectURI:  getEnv("YANDEX_REDIRECT_URI", ""),
		VKClientID:         getEnv("VK_CLIENT_ID", ""),
		VKClientSecret:     getEnv("VK_CLIENT_SECRET", ""),
		VKRedirectURI:      getEnv("VK_REDIRECT_URI", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

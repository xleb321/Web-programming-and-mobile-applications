package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    AppPort    string
    AppEnv     string
    DefaultPage int
    DefaultLimit int
    MaxLimit    int
}

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    defaultPage, _ := strconv.Atoi(getEnv("DEFAULT_PAGE", "1"))
    defaultLimit, _ := strconv.Atoi(getEnv("DEFAULT_LIMIT", "10"))
    maxLimit, _ := strconv.Atoi(getEnv("MAX_LIMIT", "100"))

    return &Config{
        DBHost:      getEnv("DB_HOST", "localhost"),
        DBPort:      getEnv("DB_PORT", "5432"),
        DBUser:      getEnv("DB_USER", "rugram_user"),
        DBPassword:  getEnv("DB_PASSWORD", "rugram_password"),
        DBName:      getEnv("DB_NAME", "rugram_db"),
        AppPort:     getEnv("APP_PORT", "4200"),
        AppEnv:      getEnv("APP_ENV", "development"),
        DefaultPage: defaultPage,
        DefaultLimit: defaultLimit,
        MaxLimit:    maxLimit,
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
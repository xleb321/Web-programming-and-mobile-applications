package database

import (
	"database/sql"
	"fmt"
	"log"

	"rugram-api/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBName,
    )
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Database connected successfully")
    return db, nil
}
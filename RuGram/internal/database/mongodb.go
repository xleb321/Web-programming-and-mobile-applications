package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"rugram-api/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(cfg *config.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Проверка подключения
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(cfg.MongoDatabase)

	// Создание индексов
	if err := createIndexes(database); err != nil {
		log.Printf("Warning: failed to create indexes: %v", err)
	}

	log.Println("MongoDB connected successfully")
	return &MongoDB{
		Client:   client,
		Database: database,
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func createIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Индексы для users collection
	usersCollection := db.Collection("users")
	usersIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_email_unique"),
		},
		{
			Keys:    bson.D{{Key: "yandex_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("idx_yandex_id"),
		},
		{
			Keys:    bson.D{{Key: "vk_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("idx_vk_id"),
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
	}
	if _, err := usersCollection.Indexes().CreateMany(ctx, usersIndexes); err != nil {
		return err
	}

	// Индексы для user_tokens collection
	tokensCollection := db.Collection("user_tokens")
	tokensIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "token_hash", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "token_type", Value: 1},
				{Key: "revoked", Value: 1},
			},
		},
	}
	if _, err := tokensCollection.Indexes().CreateMany(ctx, tokensIndexes); err != nil {
		return err
	}

	// Индексы для posts collection
	postsCollection := db.Collection("posts")
	postsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	if _, err := postsCollection.Indexes().CreateMany(ctx, postsIndexes); err != nil {
		return err
	}

	log.Println("MongoDB indexes created successfully")
	return nil
}

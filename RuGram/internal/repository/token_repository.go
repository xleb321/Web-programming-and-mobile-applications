package repository

import (
	"context"
	"errors"
	"time"

	"rugram-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenRepository struct {
	collection *mongo.Collection
}

func NewTokenRepository(db *mongo.Database) *TokenRepository {
	return &TokenRepository{
		collection: db.Collection("user_tokens"),
	}
}

func (r *TokenRepository) Create(token *models.UserToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token.CreatedAt = time.Now()
	token.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *TokenRepository) FindValidAccessToken(userID primitive.ObjectID, tokenHash string) (*models.UserToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"user_id":    userID,
		"token_type": "access",
		"revoked":    false,
		"expires_at": bson.M{"$gt": now},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tokens []models.UserToken
	if err := cursor.All(ctx, &tokens); err != nil {
		return nil, err
	}

	for _, token := range tokens {
		if token.TokenHash == tokenHash {
			return &token, nil
		}
	}
	return nil, mongo.ErrNoDocuments
}

func (r *TokenRepository) FindValidRefreshToken(userID primitive.ObjectID, tokenHash string) (*models.UserToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"user_id":    userID,
		"token_type": "refresh",
		"revoked":    false,
		"expires_at": bson.M{"$gt": now},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tokens []models.UserToken
	if err := cursor.All(ctx, &tokens); err != nil {
		return nil, err
	}

	for _, token := range tokens {
		if token.TokenHash == tokenHash {
			return &token, nil
		}
	}
	return nil, mongo.ErrNoDocuments
}

func (r *TokenRepository) RevokeToken(tokenID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": tokenID}
	update := bson.M{"$set": bson.M{"revoked": true}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("token not found")
	}
	return nil
}

func (r *TokenRepository) RevokeAllUserTokens(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "revoked": false}
	update := bson.M{"$set": bson.M{"revoked": true}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *TokenRepository) RevokeAllUserTokensByType(userID primitive.ObjectID, tokenType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":    userID,
		"token_type": tokenType,
		"revoked":    false,
	}
	update := bson.M{"$set": bson.M{"revoked": true}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

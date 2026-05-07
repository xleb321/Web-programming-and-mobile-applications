package repository

import (
	"context"
	"errors"
	"time"

	"rugram-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *PostRepository {
	return &PostRepository{
		collection: db.Collection("posts"),
	}
}

func (r *PostRepository) Create(post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, post)
	return err
}

func (r *PostRepository) FindByID(id primitive.ObjectID) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	var post models.Post
	err := r.collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) FindByIDString(id string) (*models.Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}
	return r.FindByID(objectID)
}

func (r *PostRepository) Update(post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post.UpdatedAt = time.Now()

	filter := bson.M{
		"_id":        post.ID,
		"deleted_at": nil,
	}
	update := bson.M{
		"$set": bson.M{
			"title":       post.Title,
			"description": post.Description,
			"image_url":   post.ImageURL,
			"status":      post.Status,
			"updated_at":  post.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (r *PostRepository) SoftDelete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}
	update := bson.M{
		"$set": bson.M{"deleted_at": now},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (r *PostRepository) FindAll(limit, offset int) ([]models.Post, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"deleted_at": nil}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostRepository) FindByUserID(userID string, limit, offset int) ([]models.Post, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":    userID,
		"deleted_at": nil,
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

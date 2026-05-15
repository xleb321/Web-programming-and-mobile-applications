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

type FileRepository struct {
    collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) *FileRepository {
    return &FileRepository{
        collection: db.Collection("files"),
    }
}

func (r *FileRepository) Create(file *models.File) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    file.CreatedAt = time.Now()
    file.UpdatedAt = time.Now()
    file.ID = primitive.NewObjectID()

    _, err := r.collection.InsertOne(ctx, file)
    return err
}

func (r *FileRepository) FindByID(id primitive.ObjectID) (*models.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "_id":        id,
        "deleted_at": nil,
    }

    var file models.File
    err := r.collection.FindOne(ctx, filter).Decode(&file)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil
        }
        return nil, err
    }
    return &file, nil
}

func (r *FileRepository) FindByIDString(id string) (*models.File, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid file ID format")
    }
    return r.FindByID(objectID)
}

func (r *FileRepository) SoftDelete(id primitive.ObjectID) error {
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
        return errors.New("file not found")
    }
    return nil
}

func (r *FileRepository) GetUserFiles(userID primitive.ObjectID, limit, offset int) ([]models.File, int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    filter := bson.M{
        "user_id":    userID,
        "deleted_at": nil,
    }

    total, err := r.collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }

    opts := options.Find().
        SetSort(bson.M{"created_at": -1}).
        SetLimit(int64(limit)).
        SetSkip(int64(offset))

    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)

    var files []models.File
    if err := cursor.All(ctx, &files); err != nil {
        return nil, 0, err
    }

    return files, total, nil
}
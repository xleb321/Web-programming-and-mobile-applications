package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	OriginalName string             `json:"original_name" bson:"original_name"`
	ObjectKey    string             `json:"-" bson:"object_key"`
	Size         int64              `json:"size" bson:"size"`
	MimeType     string             `json:"mime_type" bson:"mime_type"`
	Bucket       string             `json:"bucket" bson:"bucket"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt    *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func (f *File) GetID() string {
	return f.ID.Hex()
}
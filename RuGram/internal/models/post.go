package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      string             `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Status      string             `json:"status" bson:"status"`
	LikesCount  int                `json:"likes_count" bson:"likes_count"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func (p *Post) GetID() string {
	return p.ID.Hex()
}

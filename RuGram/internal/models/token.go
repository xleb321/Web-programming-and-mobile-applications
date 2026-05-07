package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserToken struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	TokenHash string             `json:"-" bson:"token_hash"`
	TokenSalt string             `json:"-" bson:"token_salt"`
	TokenType string             `json:"token_type" bson:"token_type"` // access or refresh
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	Revoked   bool               `json:"revoked" bson:"revoked"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (UserToken) CollectionName() string {
	return "user_tokens"
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
    Email        string              `json:"email" bson:"email"`
    Phone        *string             `json:"phone,omitempty" bson:"phone,omitempty"`
    PasswordHash string              `json:"-" bson:"password_hash"`
    PasswordSalt string              `json:"-" bson:"password_salt"`
    YandexID     *string             `json:"-" bson:"yandex_id,omitempty"`
    VkID         *string             `json:"-" bson:"vk_id,omitempty"`
    AvatarFileID *primitive.ObjectID `json:"avatar_file_id,omitempty" bson:"avatar_file_id,omitempty"`
    DisplayName  *string             `json:"display_name,omitempty" bson:"display_name,omitempty"`
    Bio          *string             `json:"bio,omitempty" bson:"bio,omitempty"`
    CreatedAt    time.Time           `json:"created_at" bson:"created_at"`
    UpdatedAt    time.Time           `json:"updated_at" bson:"updated_at"`
    DeletedAt    *time.Time          `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func (u *User) GetID() string {
    return u.ID.Hex()
}
package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID  `json:"id" db:"id"`
    Email        string     `json:"email" db:"email"`
    Phone        *string    `json:"phone,omitempty" db:"phone"`
    PasswordHash string     `json:"-" db:"password_hash"`
    PasswordSalt string     `json:"-" db:"password_salt"`
    YandexID     *string    `json:"-" db:"yandex_id"`
    VkID         *string    `json:"-" db:"vk_id"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (User) TableName() string {
    return "users"
}

type UserResponse struct {
    ID        string     `json:"id"`
    Email     string     `json:"email"`
    Phone     *string    `json:"phone,omitempty"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
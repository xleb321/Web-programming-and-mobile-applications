package models

import (
    "time"
    "github.com/google/uuid"
)

type UserToken struct {
    ID         uuid.UUID `json:"id" db:"id"`
    UserID     uuid.UUID `json:"user_id" db:"user_id"`
    TokenHash  string    `json:"-" db:"token_hash"`
    TokenSalt  string    `json:"-" db:"token_salt"`
    TokenType  string    `json:"token_type" db:"token_type"` // access or refresh
    ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
    Revoked    bool      `json:"revoked" db:"revoked"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func (UserToken) TableName() string {
    return "user_tokens"
}
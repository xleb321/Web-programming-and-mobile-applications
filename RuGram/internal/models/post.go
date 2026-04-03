package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
    ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID      string     `json:"user_id" gorm:"not null;index"`
    Title       string     `json:"title" gorm:"not null;size:200"`
    Description string     `json:"description" gorm:"type:text"`
    ImageURL    string     `json:"image_url" gorm:"size:500"`
    Status      string     `json:"status" gorm:"default:'active';size:20"`
    LikesCount  int        `json:"likes_count" gorm:"default:0"`
    CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

func (Post) TableName() string {
    return "posts"
}
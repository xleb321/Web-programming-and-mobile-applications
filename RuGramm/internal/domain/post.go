package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID          string         `gorm:"type:uuid;primary_key" json:"id"`
	Title       string         `gorm:"not null" json:"title" validate:"required,min=3,max=100"`
	ImageURL    string         `gorm:"not null" json:"image_url" validate:"required,url"`
	Description string         `json:"description" validate:"max=500"`
	Likes       int            `gorm:"default:0" json:"likes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New().String()
	return nil
}
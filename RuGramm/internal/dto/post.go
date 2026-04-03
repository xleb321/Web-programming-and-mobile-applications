package dto

import (
	"time"
)

type CreatePostRequest struct {
    UserID      string `json:"user_id" binding:"required"`
    Title       string `json:"title" binding:"required,min=1,max=200"`
    Description string `json:"description" binding:"max=1000"`
    ImageURL    string `json:"image_url" binding:"url"`
    Status      string `json:"status" binding:"oneof=active draft archived"`
}

type UpdatePostRequest struct {
    Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
    Description *string `json:"description" binding:"omitempty,max=1000"`
    ImageURL    *string `json:"image_url" binding:"omitempty,url"`
    Status      *string `json:"status" binding:"omitempty,oneof=active draft archived"`
}

type PostResponse struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    ImageURL    string    `json:"image_url"`
    Status      string    `json:"status"`
    LikesCount  int       `json:"likes_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type PaginationRequest struct {
    Page  int `form:"page" binding:"omitempty,min=1"`
    Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type PaginationResponse struct {
    Data       interface{} `json:"data"`
    Meta       MetaData    `json:"meta"`
}

type MetaData struct {
    Total       int64 `json:"total"`
    Page        int   `json:"page"`
    Limit       int   `json:"limit"`
    TotalPages  int64 `json:"total_pages"`
}
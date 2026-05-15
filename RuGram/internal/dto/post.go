package dto

import (
	"time"
)

// CreatePostRequest represents post creation request body
// @Description Данные для создания поста
type CreatePostRequest struct {
	UserID      string `json:"user_id" binding:"required" example:"507f1f77bcf86cd799439011" description:"ID пользователя (ObjectId)"`
	Title       string `json:"title" binding:"required,min=1,max=200" example:"Мой первый пост" description:"Заголовок поста"`
	Description string `json:"description" binding:"max=1000" example:"Это описание моего первого поста" description:"Описание поста"`
	ImageURL    string `json:"image_url" binding:"url" example:"https://example.com/image.jpg" description:"URL изображения"`
	Status      string `json:"status" binding:"oneof=active draft archived" example:"active" description:"Статус поста"`
}

// UpdatePostRequest represents post update request body
// @Description Данные для обновления поста
type UpdatePostRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=1,max=200" example:"Обновленный заголовок" description:"Новый заголовок"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000" example:"Обновленное описание" description:"Новое описание"`
	ImageURL    *string `json:"image_url,omitempty" binding:"omitempty,url" example:"https://example.com/new-image.jpg" description:"Новый URL изображения"`
	Status      *string `json:"status,omitempty" binding:"omitempty,oneof=active draft archived" example:"draft" description:"Новый статус"`
}

// PostResponse represents post data response
// @Description Информация о посте
type PostResponse struct {
	ID          string    `json:"id" example:"507f1f77bcf86cd799439011" description:"ID поста (ObjectId)"`
	UserID      string    `json:"user_id" example:"507f1f77bcf86cd799439011" description:"ID автора"`
	Title       string    `json:"title" example:"Мой первый пост" description:"Заголовок"`
	Description string    `json:"description" example:"Это описание поста" description:"Описание"`
	ImageURL    string    `json:"image_url" example:"https://example.com/image.jpg" description:"URL изображения"`
	Status      string    `json:"status" example:"active" description:"Статус поста"`
	LikesCount  int       `json:"likes_count" example:"5" description:"Количество лайков"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T12:00:00Z" description:"Дата создания"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z" description:"Дата обновления"`
}

// PaginationRequest represents pagination query parameters
type PaginationRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1" example:"1" description:"Номер страницы"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100" example:"10" description:"Элементов на странице"`
}

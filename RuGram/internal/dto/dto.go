package dto

import (
	"time"
)

// Auth DTOs
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"SecurePass123"`
	Phone    string `json:"phone,omitempty" example:"+79991234567"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePass123"`
}

type LoginResponse struct {
	Message string       `json:"message" example:"login successful"`
	User    UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string    `json:"id" example:"507f1f77bcf86cd799439011"`
	Email     string    `json:"email" example:"user@example.com"`
	Phone     *string   `json:"phone,omitempty" example:"+79991234567"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}

type WhoamiResponse struct {
	ID    string  `json:"id" example:"507f1f77bcf86cd799439011"`
	Email string  `json:"email" example:"user@example.com"`
	Phone *string `json:"phone,omitempty" example:"+79991234567"`
}

// OAuth DTOs
type OAuthRedirectResponse struct {
	URL string `json:"url" example:"https://oauth.yandex.ru/authorize?..."`
}

// Yandex OAuth DTOs
type YandexTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type YandexUserInfo struct {
	ID           string   `json:"id"`
	Login        string   `json:"login"`
	ClientID     string   `json:"client_id"`
	DisplayName  string   `json:"display_name"`
	RealName     string   `json:"real_name"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Sex          string   `json:"sex"`
	DefaultEmail string   `json:"default_email"`
	Emails       []string `json:"emails"`
}

// VK OAuth DTOs
type VKTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
	Email       string `json:"email"`
}

type VKUserInfo struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// User update DTO
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" example:"newemail@example.com"`
	Phone    *string `json:"phone,omitempty" example:"+79998887766"`
	Password *string `json:"password,omitempty" example:"NewSecurePass456"`
}

// Pagination DTOs
type PaginationResponse struct {
	Data interface{} `json:"data"`
	Meta MetaData    `json:"meta"`
}

type MetaData struct {
	Total      int64 `json:"total" example:"100"`
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	TotalPages int64 `json:"totalPages" example:"10"`
}

// Error and Success responses
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request data"`
}

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

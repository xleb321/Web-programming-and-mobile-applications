package dto

import (
    "time"
)

// Auth DTOs
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Phone    string `json:"phone,omitempty"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Message string      `json:"message"`
    User    UserResponse `json:"user"`
}

type UserResponse struct {
    ID        string     `json:"id"`
    Email     string     `json:"email"`
    Phone     *string    `json:"phone,omitempty"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}

type WhoamiResponse struct {
    ID    string  `json:"id"`
    Email string  `json:"email"`
    Phone *string `json:"phone,omitempty"`
}

// OAuth DTOs
type OAuthRedirectResponse struct {
    URL string `json:"url"`
}

type YandexUserInfo struct {
    ID            string `json:"id"`
    Login         string `json:"login"`
    ClientID      string `json:"client_id"`
    DisplayName   string `json:"display_name"`
    RealName      string `json:"real_name"`
    FirstName     string `json:"first_name"`
    LastName      string `json:"last_name"`
    Sex           string `json:"sex"`
    DefaultEmail  string `json:"default_email"`
    Emails        []string `json:"emails"`
}

type YandexTokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
}

type VKUserInfo struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}

type VKTokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    UserID      int    `json:"user_id"`
    Email       string `json:"email"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

type SuccessResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

// Добавить в существующий файл dto.go:

// UpdateUserRequest DTO
type UpdateUserRequest struct {
    Email    *string `json:"email,omitempty"`
    Phone    *string `json:"phone,omitempty"`
    Password *string `json:"password,omitempty"`
}

// PaginationResponse DTO
type PaginationResponse struct {
    Data interface{} `json:"data"`
    Meta MetaData    `json:"meta"`
}

type MetaData struct {
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    TotalPages int64 `json:"totalPages"`
}
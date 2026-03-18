package dto

import "github.com/go-playground/validator/v10"

type CreatePostRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	ImageURL    string `json:"image_url" validate:"required,url"`
	Description string `json:"description" validate:"max=500"`
}

type UpdatePostRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=100"`
	ImageURL    *string `json:"image_url" validate:"omitempty,url"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

type PaginationQuery struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (q *PaginationQuery) SetDefaults() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 10
	}
}

func (q *PaginationQuery) GetOffset() int {
	return (q.Page - 1) * q.Limit
}

type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

func (v *CreatePostRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(v)
}

func (v *UpdatePostRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(v)
}

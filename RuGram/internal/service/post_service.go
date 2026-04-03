package service

import (
	"database/sql"
	"errors"
	"fmt"

	"rugram-api/internal/dto"
	"rugram-api/internal/models"
	"rugram-api/internal/repository"

	"github.com/google/uuid"
)

type PostService struct {
    repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
    return &PostService{repo: repo}
}

func (s *PostService) Create(req *dto.CreatePostRequest) (*dto.PostResponse, error) {
    post := &models.Post{
        UserID:      req.UserID,
        Title:       req.Title,
        Description: req.Description,
        ImageURL:    req.ImageURL,
        Status:      req.Status,
        LikesCount:  0,
    }
    
    if post.Status == "" {
        post.Status = "active"
    }
    
    err := s.repo.Create(post)
    if err != nil {
        return nil, fmt.Errorf("failed to create post: %w", err)
    }
    
    return s.toResponse(post), nil
}

func (s *PostService) GetByID(id string) (*dto.PostResponse, error) {
    postID, err := uuid.Parse(id)
    if err != nil {
        return nil, errors.New("invalid post ID format")
    }
    
    post, err := s.repo.FindByID(postID)
    if err != nil {
        return nil, fmt.Errorf("failed to find post: %w", err)
    }
    
    if post == nil {
        return nil, errors.New("post not found")
    }
    
    return s.toResponse(post), nil
}

func (s *PostService) Update(id string, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
    postID, err := uuid.Parse(id)
    if err != nil {
        return nil, errors.New("invalid post ID format")
    }
    
    post, err := s.repo.FindByID(postID)
    if err != nil {
        return nil, fmt.Errorf("failed to find post: %w", err)
    }
    
    if post == nil {
        return nil, errors.New("post not found")
    }
    
    if req.Title != nil {
        post.Title = *req.Title
    }
    if req.Description != nil {
        post.Description = *req.Description
    }
    if req.ImageURL != nil {
        post.ImageURL = *req.ImageURL
    }
    if req.Status != nil {
        post.Status = *req.Status
    }
    
    err = s.repo.Update(post)
    if err != nil {
        return nil, fmt.Errorf("failed to update post: %w", err)
    }
    
    return s.toResponse(post), nil
}

func (s *PostService) Delete(id string) error {
    postID, err := uuid.Parse(id)
    if err != nil {
        return errors.New("invalid post ID format")
    }
    
    err = s.repo.SoftDelete(postID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errors.New("post not found")
        }
        return fmt.Errorf("failed to delete post: %w", err)
    }
    
    return nil
}

func (s *PostService) GetAll(page, limit int) (*dto.PaginationResponse, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    
    offset := (page - 1) * limit
    
    posts, total, err := s.repo.FindAll(limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to get posts: %w", err)
    }
    
    responses := make([]dto.PostResponse, len(posts))
    for i, post := range posts {
        responses[i] = *s.toResponse(&post)
    }
    
    totalPages := (total + int64(limit) - 1) / int64(limit)
    
    return &dto.PaginationResponse{
        Data: responses,
        Meta: dto.MetaData{
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
        },
    }, nil
}

func (s *PostService) GetByUserID(userID string, page, limit int) (*dto.PaginationResponse, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    
    offset := (page - 1) * limit
    
    posts, total, err := s.repo.FindByUserID(userID, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to get user posts: %w", err)
    }
    
    responses := make([]dto.PostResponse, len(posts))
    for i, post := range posts {
        responses[i] = *s.toResponse(&post)
    }
    
    totalPages := (total + int64(limit) - 1) / int64(limit)
    
    return &dto.PaginationResponse{
        Data: responses,
        Meta: dto.MetaData{
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
        },
    }, nil
}

func (s *PostService) toResponse(post *models.Post) *dto.PostResponse {
    return &dto.PostResponse{
        ID:          post.ID.String(),
        UserID:      post.UserID,
        Title:       post.Title,
        Description: post.Description,
        ImageURL:    post.ImageURL,
        Status:      post.Status,
        LikesCount:  post.LikesCount,
        CreatedAt:   post.CreatedAt,
        UpdatedAt:   post.UpdatedAt,
    }
}
package service

import (
	"errors"
	"fmt"
	"strconv"

	"rugram-api/internal/cache"
	"rugram-api/internal/dto"
	"rugram-api/internal/models"
	"rugram-api/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	repo     *repository.PostRepository
	cacheSvc *cache.CacheService
}

func NewPostService(repo *repository.PostRepository, cacheSvc *cache.CacheService) *PostService {
	return &PostService{
		repo:     repo,
		cacheSvc: cacheSvc,
	}
}

func (s *PostService) buildListKey(page, limit int) string {
	return s.cacheSvc.BuildKey("posts", "list", "page", strconv.Itoa(page), "limit", strconv.Itoa(limit))
}

func (s *PostService) buildUserListKey(userID string, page, limit int) string {
	return s.cacheSvc.BuildKey("posts", "user", userID, "list", "page", strconv.Itoa(page), "limit", strconv.Itoa(limit))
}

func (s *PostService) buildItemKey(id string) string {
	return s.cacheSvc.BuildKey("posts", "item", id)
}

func (s *PostService) invalidateLists() {
	pattern := s.cacheSvc.BuildKey("posts", "list", "*")
	s.cacheSvc.DelByPattern(pattern)
}

func (s *PostService) invalidateUserLists(userID string) {
	pattern := s.cacheSvc.BuildKey("posts", "user", userID, "list", "*")
	s.cacheSvc.DelByPattern(pattern)
}

func (s *PostService) invalidateItem(id string) {
	key := s.buildItemKey(id)
	s.cacheSvc.Del(key)
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

	s.invalidateLists()
	s.invalidateUserLists(req.UserID)

	return s.toResponse(post), nil
}

func (s *PostService) GetByID(id string) (*dto.PostResponse, error) {
	cacheKey := s.buildItemKey(id)
	var cachedPost dto.PostResponse
	if err := s.cacheSvc.Get(cacheKey, &cachedPost); err == nil && cachedPost.ID != "" {
		return &cachedPost, nil
	}

	post, err := s.repo.FindByIDString(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	response := s.toResponse(post)
	s.cacheSvc.SetWithDefaultTTL(cacheKey, response)
	return response, nil
}

func (s *PostService) Update(id string, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := s.repo.FindByIDString(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	oldUserID := post.UserID

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

	s.invalidateItem(id)
	s.invalidateLists()
	s.invalidateUserLists(oldUserID)
	if oldUserID != post.UserID {
		s.invalidateUserLists(post.UserID)
	}

	return s.toResponse(post), nil
}

func (s *PostService) Delete(id string) error {
	post, err := s.repo.FindByIDString(id)
	if err != nil {
		return fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return errors.New("post not found")
	}

	postID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid post ID format")
	}

	err = s.repo.SoftDelete(postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	s.invalidateItem(id)
	s.invalidateLists()
	s.invalidateUserLists(post.UserID)

	return nil
}

func (s *PostService) GetAll(page, limit int) (*dto.PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	cacheKey := s.buildListKey(page, limit)
	var cachedResult dto.PaginationResponse
	if err := s.cacheSvc.Get(cacheKey, &cachedResult); err == nil && cachedResult.Data != nil {
		return &cachedResult, nil
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

	result := &dto.PaginationResponse{
		Data: responses,
		Meta: dto.MetaData{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}

	s.cacheSvc.SetWithDefaultTTL(cacheKey, result)
	return result, nil
}

func (s *PostService) GetByUserID(userID string, page, limit int) (*dto.PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	cacheKey := s.buildUserListKey(userID, page, limit)
	var cachedResult dto.PaginationResponse
	if err := s.cacheSvc.Get(cacheKey, &cachedResult); err == nil && cachedResult.Data != nil {
		return &cachedResult, nil
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

	result := &dto.PaginationResponse{
		Data: responses,
		Meta: dto.MetaData{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}

	s.cacheSvc.SetWithDefaultTTL(cacheKey, result)
	return result, nil
}

func (s *PostService) toResponse(post *models.Post) *dto.PostResponse {
	return &dto.PostResponse{
		ID:          post.GetID(),
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

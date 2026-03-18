package service

import (
	"RuGramm/internal/domain"
	"RuGramm/internal/dto"
	"RuGramm/internal/repository"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrInvalidInput = errors.New("invalid input")
)

type PostService interface {
	Create(req *dto.CreatePostRequest) (*domain.Post, error)
	GetByID(id string) (*domain.Post, error)
	GetAll(query *dto.PaginationQuery) (*dto.PaginatedResponse, error)
	Update(id string, req *dto.UpdatePostRequest) (*domain.Post, error)
	Delete(id string) error
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) Create(req *dto.CreatePostRequest) (*domain.Post, error) {
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidInput
	}

	post := &domain.Post{
		Title:       req.Title,
		ImageURL:    req.ImageURL,
		Description: req.Description,
	}

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) GetByID(id string) (*domain.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return post, nil
}

func (s *postService) GetAll(query *dto.PaginationQuery) (*dto.PaginatedResponse, error) {
	query.SetDefaults()

	posts, total, err := s.repo.GetAll(query)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Limit
	if int(total)%query.Limit > 0 {
		totalPages++
	}

	response := &dto.PaginatedResponse{
		Data: posts,
		Meta: dto.Meta{
			Total:       total,
			Page:        query.Page,
			Limit:       query.Limit,
			TotalPages:  totalPages,
			HasNext:     query.Page < totalPages,
			HasPrevious: query.Page > 1,
		},
	}

	return response, nil
}

func (s *postService) Update(id string, req *dto.UpdatePostRequest) (*domain.Post, error) {
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidInput
	}

	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.ImageURL != nil {
		post.ImageURL = *req.ImageURL
	}
	if req.Description != nil {
		post.Description = *req.Description
	}

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) Delete(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}

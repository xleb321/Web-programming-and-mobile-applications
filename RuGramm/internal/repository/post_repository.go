package repository

import (
	"RuGramm/internal/domain"
	"RuGramm/internal/dto"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *domain.Post) error
	GetByID(id string) (*domain.Post, error)
	GetAll(query *dto.PaginationQuery) ([]domain.Post, int64, error)
	Update(post *domain.Post) error
	Delete(id string) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id string) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetAll(query *dto.PaginationQuery) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	// Get total count
	if err := r.db.Model(&domain.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Offset(query.GetOffset()).
		Limit(query.Limit).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id string) error {
	return r.db.Delete(&domain.Post{}, "id = ?", id).Error
}

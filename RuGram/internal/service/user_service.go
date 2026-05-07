package service

import (
	"errors"
	"fmt"

	"rugram-api/internal/dto"
	"rugram-api/internal/repository"
	"rugram-api/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(id string) (*dto.UserResponse, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:        user.GetID(),
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) GetUserByEmail(email string) (*dto.UserResponse, error) {
	email = utils.NormalizeEmail(email) // <-- добавить нормализацию

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:        user.GetID(),
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) UpdateUser(id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Email != nil && *req.Email != "" {
		normalized := utils.NormalizeEmail(*req.Email) // <-- нормализовать перед сохранением
		if err := utils.ValidateEmail(normalized); err != nil {
			return nil, err
		}
		user.Email = normalized
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if req.Password != nil && *req.Password != "" {
		if err := utils.ValidatePassword(*req.Password); err != nil {
			return nil, err
		}
		passwordHash, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = passwordHash
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &dto.UserResponse{
		ID:        user.GetID(),
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) DeleteUser(id string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	err = s.userRepo.SoftDelete(userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) GetAllUsers(page, limit int) (*dto.PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, total, err := s.userRepo.FindAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = dto.UserResponse{
			ID:        user.GetID(),
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
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

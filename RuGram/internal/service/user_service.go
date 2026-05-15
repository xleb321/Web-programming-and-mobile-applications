package service

import (
	"errors"
	"fmt"
	"time"

	"rugram-api/internal/dto"
	"rugram-api/internal/models"
	"rugram-api/internal/repository"
	"rugram-api/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
    userRepo *repository.UserRepository
    fileRepo *repository.FileRepository
    minioSvc *MinioService
}

func NewUserService(userRepo *repository.UserRepository, fileRepo *repository.FileRepository, minioSvc *MinioService) *UserService {
    return &UserService{
        userRepo: userRepo,
        fileRepo: fileRepo,
        minioSvc: minioSvc,
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

    return s.toUserResponse(user), nil
}

func (s *UserService) GetUserByEmail(email string) (*dto.UserResponse, error) {
    email = utils.NormalizeEmail(email)

    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }

    if user == nil {
        return nil, errors.New("user not found")
    }

    return s.toUserResponse(user), nil
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
        normalized := utils.NormalizeEmail(*req.Email)
        if err := utils.ValidateEmail(normalized); err != nil {
            return nil, err
        }
        user.Email = normalized
    }

    if req.Phone != nil {
        user.Phone = req.Phone
    }

    if req.DisplayName != nil {
        user.DisplayName = req.DisplayName
    }

    if req.Bio != nil {
        user.Bio = req.Bio
    }

    if req.AvatarFileID != nil && *req.AvatarFileID != "" {
        avatarFileID, err := primitive.ObjectIDFromHex(*req.AvatarFileID)
        if err != nil {
            return nil, errors.New("invalid avatar file ID format")
        }

        file, err := s.fileRepo.FindByID(avatarFileID)
        if err != nil {
            return nil, fmt.Errorf("failed to find avatar file: %w", err)
        }
        if file == nil {
            return nil, errors.New("avatar file not found")
        }

        if file.UserID != userID {
            return nil, errors.New("avatar file does not belong to user")
        }

        user.AvatarFileID = &avatarFileID
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

    return s.toUserResponse(user), nil
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
        responses[i] = *s.toUserResponse(&user)
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

func (s *UserService) GetProfile(userID string) (*dto.ProfileResponse, error) {
    user, err := s.GetUserByID(userID)
    if err != nil {
        return nil, err
    }

    var avatarURL *string
    // Нам нужно получить полную модель пользователя из репозитория, чтобы получить AvatarFileID
    userIDObj, _ := primitive.ObjectIDFromHex(userID)
    fullUser, _ := s.userRepo.FindByID(userIDObj)
    
    if fullUser != nil && fullUser.AvatarFileID != nil {
        file, err := s.fileRepo.FindByID(*fullUser.AvatarFileID)
        if err == nil && file != nil && s.minioSvc != nil {
            url := s.minioSvc.GetFileURL(file.ObjectKey)
            avatarURL = &url
        }
    }

    return &dto.ProfileResponse{
        ID:          user.ID,
        Email:       user.Email,
        Phone:       user.Phone,
        DisplayName: user.DisplayName,
        Bio:         user.Bio,
        AvatarURL:   avatarURL,
        CreatedAt:   user.CreatedAt.Format(time.RFC3339),
        UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
    }, nil
}

func (s *UserService) UpdateProfile(userID string, req *dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
    updateReq := &dto.UpdateUserRequest{
        Email:        req.Email,
        Phone:        req.Phone,
        DisplayName:  req.DisplayName,
        Bio:          req.Bio,
        AvatarFileID: req.AvatarFileID,
    }

    _, err := s.UpdateUser(userID, updateReq)
    if err != nil {
        return nil, err
    }

    return s.GetProfile(userID)
}

func (s *UserService) toUserResponse(user *models.User) *dto.UserResponse {
    resp := &dto.UserResponse{
        ID:          user.GetID(),
        Email:       user.Email,
        Phone:       user.Phone,
        DisplayName: user.DisplayName,
        Bio:         user.Bio,
        CreatedAt:   user.CreatedAt,
        UpdatedAt:   user.UpdatedAt,
    }

    if user.AvatarFileID != nil {
        url := fmt.Sprintf("/api/v1/files/%s", user.AvatarFileID.Hex())
        resp.AvatarURL = &url
    }

    return resp
}
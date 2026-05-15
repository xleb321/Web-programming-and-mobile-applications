package service

import (
	"errors"
	"fmt"
	"time"

	"rugram-api/internal/cache"
	"rugram-api/internal/dto"
	"rugram-api/internal/models"
	"rugram-api/internal/repository"
	"rugram-api/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
    userRepo    *repository.UserRepository
    tokenRepo   *repository.TokenRepository
    cacheSvc    *cache.CacheService
    rabbitMQSvc *RabbitMQService
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, cacheSvc *cache.CacheService, rabbitMQSvc *RabbitMQService) *AuthService {
    return &AuthService{
        userRepo:    userRepo,
        tokenRepo:   tokenRepo,
        cacheSvc:    cacheSvc,
        rabbitMQSvc: rabbitMQSvc,
    }
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
    req.Email = utils.NormalizeEmail(req.Email)

    if err := utils.ValidateEmail(req.Email); err != nil {
        return nil, err
    }
    if err := utils.ValidatePassword(req.Password); err != nil {
        return nil, err
    }

    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }

    passwordHash, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    passwordSalt, err := utils.GenerateSalt()
    if err != nil {
        return nil, fmt.Errorf("failed to generate salt: %w", err)
    }

    var phone *string
    if req.Phone != "" {
        phone = &req.Phone
    }

    user := &models.User{
        Email:        req.Email,
        Phone:        phone,
        PasswordHash: passwordHash,
        PasswordSalt: passwordSalt,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // Publish user.registered event to RabbitMQ for async email sending
    displayName := user.Email
    if user.DisplayName != nil {
        displayName = *user.DisplayName
    }
    if s.rabbitMQSvc != nil && s.rabbitMQSvc.IsConnected() {
        if err := s.rabbitMQSvc.PublishUserRegisteredEvent(user.GetID(), user.Email, displayName); err != nil {
            fmt.Printf("Warning: failed to publish user.registered event: %v\n", err)
        }
    }

    return &dto.UserResponse{
        ID:        user.GetID(),
        Email:     user.Email,
        Phone:     user.Phone,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, string, error) {
    email = utils.NormalizeEmail(email)

    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to find user: %w", err)
    }
    if user == nil {
        return nil, "", "", errors.New("invalid email or password")
    }

    if !utils.VerifyPassword(password, user.PasswordHash) {
        return nil, "", "", errors.New("invalid email or password")
    }

    userIDString := user.GetID()

    accessToken, accessJTI, err := utils.CreateAccessToken(userIDString)
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to create access token: %w", err)
    }

    refreshToken, refreshJTI, err := utils.CreateRefreshToken(userIDString)
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to create refresh token: %w", err)
    }

    accessExpiration := time.Now().Add(15 * time.Minute)
    refreshExpiration := time.Now().Add(7 * 24 * time.Hour)

    accessKey := s.cacheSvc.BuildKey("auth", "user", userIDString, "access", accessJTI)
    refreshKey := s.cacheSvc.BuildKey("auth", "user", userIDString, "refresh", refreshJTI)

    s.cacheSvc.SetString(accessKey, userIDString, time.Until(accessExpiration))
    s.cacheSvc.SetString(refreshKey, userIDString, time.Until(refreshExpiration))

    accessSalt, _ := utils.GenerateSalt()
    refreshSalt, _ := utils.GenerateSalt()

    accessTokenHash := utils.HashToken(accessToken, accessSalt)
    refreshTokenHash := utils.HashToken(refreshToken, refreshSalt)

    accessTokenRecord := &models.UserToken{
        UserID:    user.ID,
        TokenHash: accessTokenHash,
        TokenSalt: accessSalt,
        TokenType: "access",
        ExpiresAt: accessExpiration,
        Revoked:   false,
    }
    if err := s.tokenRepo.Create(accessTokenRecord); err != nil {
        return nil, "", "", fmt.Errorf("failed to save access token: %w", err)
    }

    refreshTokenRecord := &models.UserToken{
        UserID:    user.ID,
        TokenHash: refreshTokenHash,
        TokenSalt: refreshSalt,
        TokenType: "refresh",
        ExpiresAt: refreshExpiration,
        Revoked:   false,
    }
    if err := s.tokenRepo.Create(refreshTokenRecord); err != nil {
        return nil, "", "", fmt.Errorf("failed to save refresh token: %w", err)
    }

    return user, accessToken, refreshToken, nil
}

func (s *AuthService) GetUserFromAccessToken(accessToken string) (*models.User, error) {
    claims, err := utils.DecodeAccessToken(accessToken)
    if err != nil {
        return nil, errors.New("invalid or expired access token")
    }

    if claims.Type != "access" {
        return nil, errors.New("invalid token type")
    }

    accessKey := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "access", claims.JTI)

    exists, err := s.cacheSvc.Exists(accessKey)
    if err != nil {
        // If Redis is unavailable, skip check
    } else if !exists {
        return nil, errors.New("token has been revoked")
    }

    userID, err := primitive.ObjectIDFromHex(claims.Sub)
    if err != nil {
        return nil, errors.New("invalid user ID in token")
    }

    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    return user, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (*models.User, string, string, error) {
    claims, err := utils.DecodeRefreshToken(refreshToken)
    if err != nil {
        return nil, "", "", errors.New("invalid or expired refresh token")
    }

    if claims.Type != "refresh" {
        return nil, "", "", errors.New("invalid token type")
    }

    refreshKey := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "refresh", claims.JTI)
    exists, _ := s.cacheSvc.Exists(refreshKey)
    if !exists {
        return nil, "", "", errors.New("refresh token has been revoked")
    }

    userID, err := primitive.ObjectIDFromHex(claims.Sub)
    if err != nil {
        return nil, "", "", errors.New("invalid user ID in token")
    }

    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to find user: %w", err)
    }
    if user == nil {
        return nil, "", "", errors.New("user not found")
    }

    s.cacheSvc.Del(refreshKey)

    oldAccessPattern := s.cacheSvc.BuildKey("auth", "user", user.GetID(), "access", "*")
    s.cacheSvc.DelByPattern(oldAccessPattern)

    newAccessToken, newAccessJTI, err := utils.CreateAccessToken(user.GetID())
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to create access token: %w", err)
    }

    newRefreshToken, newRefreshJTI, err := utils.CreateRefreshToken(user.GetID())
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to create refresh token: %w", err)
    }

    accessExpiration := time.Now().Add(15 * time.Minute)
    refreshExpiration := time.Now().Add(7 * 24 * time.Hour)

    accessKey := s.cacheSvc.BuildKey("auth", "user", user.GetID(), "access", newAccessJTI)
    s.cacheSvc.SetString(accessKey, user.GetID(), time.Until(accessExpiration))

    newRefreshKey := s.cacheSvc.BuildKey("auth", "user", user.GetID(), "refresh", newRefreshJTI)
    s.cacheSvc.SetString(newRefreshKey, user.GetID(), time.Until(refreshExpiration))

    accessSalt, _ := utils.GenerateSalt()
    refreshSalt, _ := utils.GenerateSalt()

    newAccessTokenHash := utils.HashToken(newAccessToken, accessSalt)
    newRefreshTokenHash := utils.HashToken(newRefreshToken, refreshSalt)

    accessTokenRecord := &models.UserToken{
        UserID:    user.ID,
        TokenHash: newAccessTokenHash,
        TokenSalt: accessSalt,
        TokenType: "access",
        ExpiresAt: accessExpiration,
        Revoked:   false,
    }
    s.tokenRepo.Create(accessTokenRecord)

    refreshTokenRecord := &models.UserToken{
        UserID:    user.ID,
        TokenHash: newRefreshTokenHash,
        TokenSalt: refreshSalt,
        TokenType: "refresh",
        ExpiresAt: refreshExpiration,
        Revoked:   false,
    }
    s.tokenRepo.Create(refreshTokenRecord)

    return user, newAccessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(accessToken string) error {
    claims, err := utils.DecodeAccessToken(accessToken)
    if err != nil {
        return errors.New("invalid access token")
    }

    accessKey := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "access", claims.JTI)
    return s.cacheSvc.Del(accessKey)
}

func (s *AuthService) LogoutAll(accessToken string) error {
    claims, err := utils.DecodeAccessToken(accessToken)
    if err != nil {
        return errors.New("invalid access token")
    }

    accessPattern := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "access", "*")
    refreshPattern := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "refresh", "*")

    s.cacheSvc.DelByPattern(accessPattern)
    s.cacheSvc.DelByPattern(refreshPattern)

    userID, _ := primitive.ObjectIDFromHex(claims.Sub)
    return s.tokenRepo.RevokeAllUserTokens(userID)
}

func (s *AuthService) BuildKey(parts ...string) string {
    return s.cacheSvc.BuildKey(parts...)
}
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

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	tokenRepo   *repository.TokenRepository
	cacheSvc    *cache.CacheService
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, cacheSvc *cache.CacheService) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		cacheSvc:  cacheSvc,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
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

	return &dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, string, error) {
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

	// Создаем токены с JTI
	accessToken, accessJTI, err := utils.CreateAccessToken(user.ID.String())
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, refreshJTI, err := utils.CreateRefreshToken(user.ID.String())
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Сохраняем JTI в Redis для проверки валидности токенов
	accessExpiration := time.Now().Add(15 * time.Minute)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)

	// Ключ для access токена: rugram:auth:user:{userId}:access:{jti}
	accessKey := s.cacheSvc.BuildKey("auth", "user", user.ID.String(), "access", accessJTI)
	s.cacheSvc.SetString(accessKey, user.ID.String(), time.Until(accessExpiration))

	// Ключ для refresh токена: rugram:auth:user:{userId}:refresh:{jti}
	refreshKey := s.cacheSvc.BuildKey("auth", "user", user.ID.String(), "refresh", refreshJTI)
	s.cacheSvc.SetString(refreshKey, user.ID.String(), time.Until(refreshExpiration))

	// Сохраняем токены в БД (хешированные)
	accessSalt, _ := utils.GenerateSalt()
	refreshSalt, _ := utils.GenerateSalt()

	accessTokenHash := utils.HashToken(accessToken, accessSalt)
	refreshTokenHash := utils.HashToken(refreshToken, refreshSalt)

	accessTokenRecord := &models.UserToken{
		UserID:     user.ID,
		TokenHash:  accessTokenHash,
		TokenSalt:  accessSalt,
		TokenType:  "access",
		ExpiresAt:  accessExpiration,
		Revoked:    false,
	}
	if err := s.tokenRepo.Create(accessTokenRecord); err != nil {
		return nil, "", "", fmt.Errorf("failed to save access token: %w", err)
	}

	refreshTokenRecord := &models.UserToken{
		UserID:     user.ID,
		TokenHash:  refreshTokenHash,
		TokenSalt:  refreshSalt,
		TokenType:  "refresh",
		ExpiresAt:  refreshExpiration,
		Revoked:    false,
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

	// Проверяем наличие JTI в Redis (валидность токена)
	accessKey := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "access", claims.JTI)
	exists, err := s.cacheSvc.Exists(accessKey)
	if err != nil {
		// Если Redis недоступен, пропускаем проверку (fail-open)
		// Но логируем предупреждение
	} else if !exists {
		return nil, errors.New("token has been revoked")
	}

	userID, err := uuid.Parse(claims.Sub)
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

	// Проверяем refresh token в Redis
	refreshKey := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "refresh", claims.JTI)
	exists, _ := s.cacheSvc.Exists(refreshKey)
	if !exists {
		return nil, "", "", errors.New("refresh token has been revoked")
	}

	userID, err := uuid.Parse(claims.Sub)
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

	// Удаляем старый refresh токен из Redis
	s.cacheSvc.Del(refreshKey)

	// Удаляем все старые access токены пользователя из Redis
	oldAccessPattern := s.cacheSvc.BuildKey("auth", "user", user.ID.String(), "access", "*")
	s.cacheSvc.DelByPattern(oldAccessPattern)

	// Создаем новые токены
	newAccessToken, newAccessJTI, err := utils.CreateAccessToken(user.ID.String())
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	newRefreshToken, newRefreshJTI, err := utils.CreateRefreshToken(user.ID.String())
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	accessExpiration := time.Now().Add(15 * time.Minute)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)

	// Сохраняем новые JTI в Redis
	accessKey := s.cacheSvc.BuildKey("auth", "user", user.ID.String(), "access", newAccessJTI)
	s.cacheSvc.SetString(accessKey, user.ID.String(), time.Until(accessExpiration))

	newRefreshKey := s.cacheSvc.BuildKey("auth", "user", user.ID.String(), "refresh", newRefreshJTI)
	s.cacheSvc.SetString(newRefreshKey, user.ID.String(), time.Until(refreshExpiration))

	// Сохраняем в БД
	accessSalt, _ := utils.GenerateSalt()
	refreshSalt, _ := utils.GenerateSalt()

	newAccessTokenHash := utils.HashToken(newAccessToken, accessSalt)
	newRefreshTokenHash := utils.HashToken(newRefreshToken, refreshSalt)

	accessTokenRecord := &models.UserToken{
		UserID:     user.ID,
		TokenHash:  newAccessTokenHash,
		TokenSalt:  accessSalt,
		TokenType:  "access",
		ExpiresAt:  accessExpiration,
		Revoked:    false,
	}
	s.tokenRepo.Create(accessTokenRecord)

	refreshTokenRecord := &models.UserToken{
		UserID:     user.ID,
		TokenHash:  newRefreshTokenHash,
		TokenSalt:  refreshSalt,
		TokenType:  "refresh",
		ExpiresAt:  refreshExpiration,
		Revoked:    false,
	}
	s.tokenRepo.Create(refreshTokenRecord)

	return user, newAccessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(accessToken string) error {
	claims, err := utils.DecodeAccessToken(accessToken)
	if err != nil {
		return errors.New("invalid access token")
	}

	// Удаляем access токен из Redis
	accessKey := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "access", claims.JTI)
	return s.cacheSvc.Del(accessKey)
}

func (s *AuthService) LogoutAll(accessToken string) error {
	claims, err := utils.DecodeAccessToken(accessToken)
	if err != nil {
		return errors.New("invalid access token")
	}

	// Удаляем все access и refresh токены пользователя из Redis
	accessPattern := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "access", "*")
	refreshPattern := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "refresh", "*")

	s.cacheSvc.DelByPattern(accessPattern)
	s.cacheSvc.DelByPattern(refreshPattern)

	return s.tokenRepo.RevokeAllUserTokens(uuid.MustParse(claims.Sub))
}

// BuildKey экспортируем метод для использования в других сервисах
func (s *AuthService) BuildKey(parts ...string) string {
	return s.cacheSvc.BuildKey(parts...)
}
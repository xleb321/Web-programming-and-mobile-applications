package service

import (
    "errors"
    "fmt"
    "time"

    "rugram-api/internal/cache"
    "rugram-api/internal/dto"
    "rugram-api/internal/models"
    "rugram-api/internal/repository"
    "rugram-api/internal/utils"          // <-- убедитесь, что путь верный (у вас может быть "rugram-api/pkg/utils")

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
    userRepo  *repository.UserRepository
    tokenRepo *repository.TokenRepository
    cacheSvc  *cache.CacheService
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, cacheSvc *cache.CacheService) *AuthService {
    return &AuthService{
        userRepo:  userRepo,
        tokenRepo: tokenRepo,
        cacheSvc:  cacheSvc,
    }
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
    // Нормализуем email перед валидацией и сохранением
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
        Email:        req.Email,    // уже в нижнем регистре
        Phone:        phone,
        PasswordHash: passwordHash,
        PasswordSalt: passwordSalt,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
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
    // Нормализуем email перед поиском
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

    // --- Остальная часть Login остаётся без изменений ---
    userIDString := user.GetID()
    fmt.Printf("User logged in with ID: %s (hex: %s)\n", user.ID.Hex(), userIDString)

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

    fmt.Printf("Redis keys created:\n")
    fmt.Printf("  Access: %s\n", accessKey)
    fmt.Printf("  Refresh: %s\n", refreshKey)

    s.cacheSvc.SetString(accessKey, userIDString, time.Until(accessExpiration))
    s.cacheSvc.SetString(refreshKey, userIDString, time.Until(refreshExpiration))

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

	fmt.Printf("Token claims - Sub: %s, JTI: %s\n", claims.Sub, claims.JTI)

	// Проверяем наличие JTI в Redis
	accessKey := s.cacheSvc.BuildKey("auth", "user", claims.Sub, "access", claims.JTI)
	fmt.Printf("Checking Redis key: %s\n", accessKey)

	exists, err := s.cacheSvc.Exists(accessKey)
	if err != nil {
		fmt.Printf("Redis error: %v\n", err)
		// Если Redis недоступен, пропускаем проверку
	} else if !exists {
		fmt.Printf("Token not found in Redis - revoked or expired\n")
		return nil, errors.New("token has been revoked")
	}

	// Конвертируем строку в ObjectID
	userID, err := primitive.ObjectIDFromHex(claims.Sub)
	if err != nil {
		fmt.Printf("Failed to parse ObjectID from hex: %s, error: %v\n", claims.Sub, err)
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

	// Удаляем старый refresh токен из Redis
	s.cacheSvc.Del(refreshKey)

	// Удаляем все старые access токены пользователя из Redis
	oldAccessPattern := s.cacheSvc.BuildKey("auth", "user", user.GetID(), "access", "*")
	s.cacheSvc.DelByPattern(oldAccessPattern)

	// Создаем новые токены
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

	// Сохраняем новые JTI в Redis
	accessKey := s.cacheSvc.BuildKey("auth", "user", user.GetID(), "access", newAccessJTI)
	s.cacheSvc.SetString(accessKey, user.GetID(), time.Until(accessExpiration))

	newRefreshKey := s.cacheSvc.BuildKey("auth", "user", user.GetID(), "refresh", newRefreshJTI)
	s.cacheSvc.SetString(newRefreshKey, user.GetID(), time.Until(refreshExpiration))

	// Сохраняем в БД
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

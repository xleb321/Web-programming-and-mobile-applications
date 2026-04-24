package service

import (
    "errors"
    "fmt"
    "time"
    
    "rugram-api/internal/dto"
    "rugram-api/internal/models"
    "rugram-api/internal/repository"
    "rugram-api/internal/utils"
    
    "github.com/google/uuid"
)

type AuthService struct {
    userRepo  *repository.UserRepository
    tokenRepo *repository.TokenRepository
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) *AuthService {
    return &AuthService{
        userRepo:  userRepo,
        tokenRepo: tokenRepo,
    }
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
    // Validate input
    if err := utils.ValidateEmail(req.Email); err != nil {
        return nil, err
    }
    if err := utils.ValidatePassword(req.Password); err != nil {
        return nil, err
    }
    
    // Check if user exists
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    // Hash password (bcrypt includes salt)
    passwordHash, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Generate salt for token hashing
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
    // Find user by email
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to find user: %w", err)
    }
    if user == nil {
        return nil, "", "", errors.New("invalid email or password")
    }
    
    // Verify password
    if !utils.VerifyPassword(password, user.PasswordHash) {
        return nil, "", "", errors.New("invalid email or password")
    }
    
    // Create tokens
    accessToken, err := utils.CreateAccessToken(user.ID.String())
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to create access token: %w", err)
    }
    
    refreshToken, err := utils.CreateRefreshToken(user.ID.String())
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to create refresh token: %w", err)
    }
    
    // Save tokens to database (hashed)
    accessSalt, err := utils.GenerateSalt()
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to generate salt: %w", err)
    }
    
    refreshSalt, err := utils.GenerateSalt()
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to generate salt: %w", err)
    }
    
    accessTokenHash := utils.HashToken(accessToken, accessSalt)
    refreshTokenHash := utils.HashToken(refreshToken, refreshSalt)
    
    accessExpiration := time.Now().Add(15 * time.Minute)
    refreshExpiration := time.Now().Add(7 * 24 * time.Hour)
    
    // Save access token
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
    
    // Save refresh token
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
    // Decode token
    claims, err := utils.DecodeAccessToken(accessToken)
    if err != nil {
        return nil, errors.New("invalid or expired access token")
    }
    
    if claims.Type != "access" {
        return nil, errors.New("invalid token type")
    }
    
    userID, err := uuid.Parse(claims.Sub)
    if err != nil {
        return nil, errors.New("invalid user ID in token")
    }
    
    // Find user
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
    // Decode refresh token
    claims, err := utils.DecodeRefreshToken(refreshToken)
    if err != nil {
        return nil, "", "", errors.New("invalid or expired refresh token")
    }
    
    if claims.Type != "refresh" {
        return nil, "", "", errors.New("invalid token type")
    }
    
    userID, err := uuid.Parse(claims.Sub)
    if err != nil {
        return nil, "", "", errors.New("invalid user ID in token")
    }
    
    // Find user
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to find user: %w", err)
    }
    if user == nil {
        return nil, "", "", errors.New("user not found")
    }
    
    // Find and revoke current refresh token
    refreshTokenHash := utils.HashToken(refreshToken, "")
    currentToken, err := s.tokenRepo.FindValidRefreshToken(userID, refreshTokenHash)
    if err == nil && currentToken != nil {
        s.tokenRepo.RevokeToken(currentToken.ID)
    }
    
    // Revoke old access tokens (optional, but good practice)
    s.tokenRepo.RevokeAllUserTokensByType(userID, "access")
    
    // Create new tokens
    newAccessToken, err := utils.CreateAccessToken(user.ID.String())
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to create access token: %w", err)
    }
    
    newRefreshToken, err := utils.CreateRefreshToken(user.ID.String())
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to create refresh token: %w", err)
    }
    
    // Save new tokens
    accessSalt, _ := utils.GenerateSalt()
    refreshSalt, _ := utils.GenerateSalt()
    
    newAccessTokenHash := utils.HashToken(newAccessToken, accessSalt)
    newRefreshTokenHash := utils.HashToken(newRefreshToken, refreshSalt)
    
    accessExpiration := time.Now().Add(15 * time.Minute)
    refreshExpiration := time.Now().Add(7 * 24 * time.Hour)
    
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
    
    userID, err := uuid.Parse(claims.Sub)
    if err != nil {
        return errors.New("invalid user ID")
    }
    
    // Revoke specific access token
    tokenHash := utils.HashToken(accessToken, "")
    token, err := s.tokenRepo.FindValidAccessToken(userID, tokenHash)
    if err == nil && token != nil {
        return s.tokenRepo.RevokeToken(token.ID)
    }
    
    return nil
}

func (s *AuthService) LogoutAll(accessToken string) error {
    claims, err := utils.DecodeAccessToken(accessToken)
    if err != nil {
        return errors.New("invalid access token")
    }
    
    userID, err := uuid.Parse(claims.Sub)
    if err != nil {
        return errors.New("invalid user ID")
    }
    
    return s.tokenRepo.RevokeAllUserTokens(userID)
}
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"rugram-api/internal/dto"
	"rugram-api/internal/models"
	"rugram-api/internal/repository"
	"rugram-api/internal/utils"
)

type OAuthService struct {
	userRepo   *repository.UserRepository
	tokenRepo  *repository.TokenRepository
	httpClient *http.Client
}

func NewOAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) *OAuthService {
	return &OAuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Yandex OAuth
func (s *OAuthService) GetYandexAuthURL(state string) string {
	clientID := os.Getenv("YANDEX_CLIENT_ID")
	redirectURI := os.Getenv("YANDEX_REDIRECT_URI")

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)

	return "https://oauth.yandex.ru/authorize?" + params.Encode()
}

func (s *OAuthService) HandleYandexCallback(code, state string) (*models.User, string, string, error) {
	tokenData, err := s.exchangeYandexCode(code)
	if err != nil {
		return nil, "", "", err
	}

	userInfo, err := s.getYandexUserInfo(tokenData.AccessToken)
	if err != nil {
		return nil, "", "", err
	}

	user, err := s.findOrCreateYandexUser(userInfo)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, _, err := utils.CreateAccessToken(user.GetID())
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, _, err := utils.CreateRefreshToken(user.GetID())
	if err != nil {
		return nil, "", "", err
	}

	accessSalt, _ := utils.GenerateSalt()
	refreshSalt, _ := utils.GenerateSalt()

	accessTokenHash := utils.HashToken(accessToken, accessSalt)
	refreshTokenHash := utils.HashToken(refreshToken, refreshSalt)

	accessExpiration := time.Now().Add(15 * time.Minute)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)

	accessTokenRecord := &models.UserToken{
		UserID:    user.ID,
		TokenHash: accessTokenHash,
		TokenSalt: accessSalt,
		TokenType: "access",
		ExpiresAt: accessExpiration,
		Revoked:   false,
	}
	s.tokenRepo.Create(accessTokenRecord)

	refreshTokenRecord := &models.UserToken{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		TokenSalt: refreshSalt,
		TokenType: "refresh",
		ExpiresAt: refreshExpiration,
		Revoked:   false,
	}
	s.tokenRepo.Create(refreshTokenRecord)

	return user, accessToken, refreshToken, nil
}

func (s *OAuthService) exchangeYandexCode(code string) (*dto.YandexTokenResponse, error) {
	clientID := os.Getenv("YANDEX_CLIENT_ID")
	clientSecret := os.Getenv("YANDEX_CLIENT_SECRET")
	redirectURI := os.Getenv("YANDEX_REDIRECT_URI")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", "https://oauth.yandex.ru/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yandex token exchange failed: %s", string(body))
	}

	var tokenResp dto.YandexTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (s *OAuthService) getYandexUserInfo(accessToken string) (*dto.YandexUserInfo, error) {
	req, err := http.NewRequest("GET", "https://login.yandex.ru/info?format=json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "OAuth "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo dto.YandexUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *OAuthService) findOrCreateYandexUser(userInfo *dto.YandexUserInfo) (*models.User, error) {
	user, err := s.userRepo.FindByYandexID(userInfo.ID)
	if err == nil && user != nil {
		return user, nil
	}

	email := userInfo.DefaultEmail
	if email == "" && len(userInfo.Emails) > 0 {
		email = userInfo.Emails[0]
	}

	if email != "" {
		user, err = s.userRepo.FindByEmail(email)
		if err == nil && user != nil {
			user.YandexID = &userInfo.ID
			s.userRepo.Update(user)
			return user, nil
		}
	}

	passwordSalt, _ := utils.GenerateSalt()
	randomPassword, _ := utils.GenerateSecureToken(32)
	passwordHash, _ := utils.HashPassword(randomPassword)

	newUser := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		YandexID:     &userInfo.ID,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// VK OAuth
func (s *OAuthService) GetVKAuthURL(state string) string {
	clientID := os.Getenv("VK_CLIENT_ID")
	redirectURI := os.Getenv("VK_REDIRECT_URI")

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("display", "page")
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", "email")
	params.Add("response_type", "code")
	params.Add("v", "5.131")
	params.Add("state", state)

	return "https://oauth.vk.com/authorize?" + params.Encode()
}

func (s *OAuthService) HandleVKCallback(code, state string) (*models.User, string, string, error) {
	tokenData, err := s.exchangeVKCode(code)
	if err != nil {
		return nil, "", "", err
	}

	userInfo, err := s.getVKUserInfo(tokenData.AccessToken, tokenData.UserID)
	if err != nil {
		return nil, "", "", err
	}

	user, err := s.findOrCreateVKUser(userInfo)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, _, err := utils.CreateAccessToken(user.GetID())
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, _, err := utils.CreateRefreshToken(user.GetID())
	if err != nil {
		return nil, "", "", err
	}

	accessSalt, _ := utils.GenerateSalt()
	refreshSalt, _ := utils.GenerateSalt()

	accessTokenHash := utils.HashToken(accessToken, accessSalt)
	refreshTokenHash := utils.HashToken(refreshToken, refreshSalt)

	accessExpiration := time.Now().Add(15 * time.Minute)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)

	accessTokenRecord := &models.UserToken{
		UserID:    user.ID,
		TokenHash: accessTokenHash,
		TokenSalt: accessSalt,
		TokenType: "access",
		ExpiresAt: accessExpiration,
		Revoked:   false,
	}
	s.tokenRepo.Create(accessTokenRecord)

	refreshTokenRecord := &models.UserToken{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		TokenSalt: refreshSalt,
		TokenType: "refresh",
		ExpiresAt: refreshExpiration,
		Revoked:   false,
	}
	s.tokenRepo.Create(refreshTokenRecord)

	return user, accessToken, refreshToken, nil
}

func (s *OAuthService) exchangeVKCode(code string) (*dto.VKTokenResponse, error) {
	clientID := os.Getenv("VK_CLIENT_ID")
	clientSecret := os.Getenv("VK_CLIENT_SECRET")
	redirectURI := os.Getenv("VK_REDIRECT_URI")

	apiURL := fmt.Sprintf(
		"https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
		clientID, clientSecret, redirectURI, code,
	)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vk token exchange failed: %s", string(body))
	}

	var tokenResp dto.VKTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (s *OAuthService) getVKUserInfo(accessToken string, userID int) (*dto.VKUserInfo, error) {
	apiURL := fmt.Sprintf(
		"https://api.vk.com/method/users.get?user_ids=%d&fields=email&access_token=%s&v=5.131",
		userID, accessToken,
	)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vkResponse struct {
		Response []dto.VKUserInfo `json:"response"`
		Error    struct {
			ErrorMsg string `json:"error_msg"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &vkResponse); err != nil {
		return nil, err
	}

	if vkResponse.Error.ErrorMsg != "" {
		return nil, errors.New(vkResponse.Error.ErrorMsg)
	}

	if len(vkResponse.Response) == 0 {
		return nil, errors.New("no user data received")
	}

	return &vkResponse.Response[0], nil
}

func (s *OAuthService) findOrCreateVKUser(userInfo *dto.VKUserInfo) (*models.User, error) {
	vkID := fmt.Sprintf("%d", userInfo.ID)

	user, err := s.userRepo.FindByVkID(vkID)
	if err == nil && user != nil {
		return user, nil
	}

	if userInfo.Email != "" {
		user, err = s.userRepo.FindByEmail(userInfo.Email)
		if err == nil && user != nil {
			user.VkID = &vkID
			s.userRepo.Update(user)
			return user, nil
		}
	}

	passwordSalt, _ := utils.GenerateSalt()
	randomPassword, _ := utils.GenerateSecureToken(32)
	passwordHash, _ := utils.HashPassword(randomPassword)

	email := userInfo.Email
	if email == "" {
		email = fmt.Sprintf("vk_user_%d@temp.local", userInfo.ID)
	}

	newUser := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		VkID:         &vkID,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

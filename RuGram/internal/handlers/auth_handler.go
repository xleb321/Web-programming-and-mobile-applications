package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"rugram-api/internal/dto"
	"rugram-api/internal/middleware"
	"rugram-api/internal/service"
	"rugram-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService  *service.AuthService
	oauthService *service.OAuthService
}

func NewAuthHandler(authService *service.AuthService, oauthService *service.OAuthService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Registration data"
// @Success      201  {object}  utils.StandardResponse{data=dto.UserResponse}
// @Failure      400  {object}  utils.ErrorResponseStruct
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, user)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates user and returns access/refresh tokens in cookies
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200  {object}  utils.StandardResponse{data=dto.LoginResponse}
// @Failure      400  {object}  utils.ErrorResponseStruct
// @Failure      401  {object}  utils.ErrorResponseStruct
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	user, accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)

	response := dto.LoginResponse{
		Message: "login successful",
		User: dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, response)
}

// Whoami godoc
// @Summary      Get current user info
// @Description  Returns authenticated user's profile
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse{data=dto.WhoamiResponse}
// @Failure      401  {object}  utils.ErrorResponseStruct
// @Router       /auth/whoami [get]
func (h *AuthHandler) Whoami(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	currentUser := user.(*middleware.AuthenticatedUser)

	response := dto.WhoamiResponse{
		ID:    currentUser.ID,
		Email: currentUser.Email,
		Phone: currentUser.Phone,
	}

	utils.SuccessResponse(c, http.StatusOK, response)
}

// Refresh godoc
// @Summary      Refresh tokens
// @Description  Uses refresh token from cookie to issue new access/refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.StandardResponse{data=dto.LoginResponse}
// @Failure      401  {object}  utils.ErrorResponseStruct
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "refresh token missing")
		return
	}

	user, accessToken, newRefreshToken, err := h.authService.RefreshTokens(refreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, newRefreshToken)

	response := dto.LoginResponse{
		Message: "tokens refreshed",
		User: dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, response)
}

// Logout godoc
// @Summary      Logout current session
// @Description  Revokes current access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.ErrorResponseStruct
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.authService.Logout(accessToken); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.clearAuthCookies(c)

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "logged out"})
}

// LogoutAll godoc
// @Summary      Logout from all sessions
// @Description  Revokes all access and refresh tokens for the user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.ErrorResponseStruct
// @Router       /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.authService.LogoutAll(accessToken); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.clearAuthCookies(c)

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "logged out from all sessions"})
}

// OAuthYandex godoc
// @Summary      Yandex OAuth login
// @Description  Redirects to Yandex OAuth authorization page
// @Tags         Auth
// @Success      302
// @Router       /auth/oauth/yandex [get]
func (h *AuthHandler) OAuthYandex(c *gin.Context) {
	state := generateOAuthState()
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	authURL := h.oauthService.GetYandexAuthURL(state)
	c.Redirect(http.StatusFound, authURL)
}

// OAuthYandexCallback godoc
// @Summary      Yandex OAuth callback
// @Description  Handles Yandex OAuth callback, creates/authenticates user
// @Tags         Auth
// @Param        code query string true "Authorization code"
// @Param        state query string true "State parameter"
// @Success      302
// @Failure      400  {object}  utils.ErrorResponseStruct
// @Router       /auth/oauth/yandex/callback [get]
func (h *AuthHandler) OAuthYandexCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid state parameter")
		return
	}

	_, accessToken, refreshToken, err := h.oauthService.HandleYandexCallback(code, state)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)
	h.clearOAuthStateCookie(c)

	frontendURL := getFrontendURL()
	c.Redirect(http.StatusFound, frontendURL)
}

// OAuthVK godoc
// @Summary      VK OAuth login
// @Description  Redirects to VK OAuth authorization page
// @Tags         Auth
// @Success      302
// @Router       /auth/oauth/vk [get]
func (h *AuthHandler) OAuthVK(c *gin.Context) {
	state := generateOAuthState()
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	authURL := h.oauthService.GetVKAuthURL(state)
	c.Redirect(http.StatusFound, authURL)
}

// OAuthVKCallback godoc
// @Summary      VK OAuth callback
// @Description  Handles VK OAuth callback, creates/authenticates user
// @Tags         Auth
// @Param        code query string true "Authorization code"
// @Param        state query string true "State parameter"
// @Success      302
// @Failure      400  {object}  utils.ErrorResponseStruct
// @Router       /auth/oauth/vk/callback [get]
func (h *AuthHandler) OAuthVKCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid state parameter")
		return
	}

	_, accessToken, refreshToken, err := h.oauthService.HandleVKCallback(code, state)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)
	h.clearOAuthStateCookie(c)

	frontendURL := getFrontendURL()
	c.Redirect(http.StatusFound, frontendURL)
}

// Helper methods (без изменений)
func (h *AuthHandler) setAccessTokenCookie(c *gin.Context, token string) {
	c.SetCookie("access_token", token, 15*60, "/", "", false, true)
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetCookie("refresh_token", token, 7*24*60*60, "/", "", false, true)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
}

func (h *AuthHandler) clearOAuthStateCookie(c *gin.Context) {
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
}

func generateOAuthState() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getFrontendURL() string {
	return "http://localhost:3000"
}
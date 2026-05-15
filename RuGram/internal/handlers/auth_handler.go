package handlers

import (
    "crypto/rand"
    "encoding/hex"
    "net/http"
    
    "rugram-api/internal/dto"
    "rugram-api/internal/middleware"  // Добавляем импорт middleware
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

// Register handles user registration
// POST /api/v1/auth/register
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

// Login handles user login
// POST /api/v1/auth/login
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
    
    // Set cookies
    h.setAccessTokenCookie(c, accessToken)
    h.setRefreshTokenCookie(c, refreshToken)
    
    response := dto.LoginResponse{
        Message: "login successful",
        User: dto.UserResponse{
            ID:        user.ID.String(),
            Email:     user.Email,
            Phone:     user.Phone,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        },
    }
    
    utils.SuccessResponse(c, http.StatusOK, response)
}

// Whoami returns current user info
// GET /api/v1/auth/whoami
func (h *AuthHandler) Whoami(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }
    
    // Используем middleware.AuthenticatedUser вместо service.AuthenticatedUser
    currentUser := user.(*middleware.AuthenticatedUser)
    
    response := dto.WhoamiResponse{
        ID:    currentUser.ID,
        Email: currentUser.Email,
        Phone: currentUser.Phone,
    }
    
    utils.SuccessResponse(c, http.StatusOK, response)
}

// Refresh handles token refresh
// POST /api/v1/auth/refresh
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
            ID:        user.ID.String(),
            Email:     user.Email,
            Phone:     user.Phone,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        },
    }
    
    utils.SuccessResponse(c, http.StatusOK, response)
}

// Logout handles single session logout
// POST /api/v1/auth/logout
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

// LogoutAll handles logout from all sessions
// POST /api/v1/auth/logout-all
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

// OAuthYandex initiates Yandex OAuth flow
// GET /api/v1/auth/oauth/yandex
func (h *AuthHandler) OAuthYandex(c *gin.Context) {
    state := generateOAuthState()
    // Store state in session or cookie for verification
    c.SetCookie("oauth_state", state, 300, "/", "", false, true)
    
    authURL := h.oauthService.GetYandexAuthURL(state)
    c.Redirect(http.StatusFound, authURL)
}

// OAuthYandexCallback handles Yandex OAuth callback
// GET /api/v1/auth/oauth/yandex/callback
func (h *AuthHandler) OAuthYandexCallback(c *gin.Context) {
    code := c.Query("code")
    state := c.Query("state")
    
    // Verify state
    savedState, err := c.Cookie("oauth_state")
    if err != nil || savedState != state {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid state parameter")
        return
    }
    
    // Убираем неиспользуемую переменную user, заменяем на _
    _, accessToken, refreshToken, err := h.oauthService.HandleYandexCallback(code, state)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    h.setAccessTokenCookie(c, accessToken)
    h.setRefreshTokenCookie(c, refreshToken)
    h.clearOAuthStateCookie(c)
    
    // Redirect to frontend
    frontendURL := getFrontendURL()
    c.Redirect(http.StatusFound, frontendURL)
}

// OAuthVK initiates VK OAuth flow
// GET /api/v1/auth/oauth/vk
func (h *AuthHandler) OAuthVK(c *gin.Context) {
    state := generateOAuthState()
    c.SetCookie("oauth_state", state, 300, "/", "", false, true)
    
    authURL := h.oauthService.GetVKAuthURL(state)
    c.Redirect(http.StatusFound, authURL)
}

// OAuthVKCallback handles VK OAuth callback
// GET /api/v1/auth/oauth/vk/callback
func (h *AuthHandler) OAuthVKCallback(c *gin.Context) {
    code := c.Query("code")
    state := c.Query("state")
    
    savedState, err := c.Cookie("oauth_state")
    if err != nil || savedState != state {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid state parameter")
        return
    }
    
    // Убираем неиспользуемую переменную user, заменяем на _
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

// Helper methods
func (h *AuthHandler) setAccessTokenCookie(c *gin.Context, token string) {
    c.SetCookie(
        "access_token",
        token,
        15*60, // 15 minutes in seconds
        "/",
        "",
        false, // secure in production
        true,  // httpOnly
    )
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
    c.SetCookie(
        "refresh_token",
        token,
        7*24*60*60, // 7 days in seconds
        "/",
        "",
        false,
        true,
    )
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
    // В production должен быть URL фронтенда
    return "http://localhost:3000"
}
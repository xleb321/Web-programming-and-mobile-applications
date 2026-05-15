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
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя с указанными email, паролем и опционально телефоном
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Данные для регистрации"
// @Success      201 {object} utils.StandardResponse{data=dto.UserResponse} "Пользователь успешно создан"
// @Failure      400 {object} utils.ErrorResponseStruct "Неверный формат запроса или email уже существует"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
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
// @Summary      Вход в систему
// @Description  Аутентифицирует пользователя и устанавливает JWT токены в cookies
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Учетные данные пользователя"
// @Success      200 {object} utils.StandardResponse{data=dto.LoginResponse} "Успешный вход"
// @Failure      400 {object} utils.ErrorResponseStruct "Неверный формат запроса"
// @Failure      401 {object} utils.ErrorResponseStruct "Неверный email или пароль"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
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

	// Set cookies
	h.setAccessTokenCookie(c, accessToken)
	h.setRefreshTokenCookie(c, refreshToken)

	response := dto.LoginResponse{
		Message: "login successful",
		User: dto.UserResponse{
			ID:        user.GetID(),
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, response)
}

// Whoami godoc
// @Summary      Получение информации о текущем пользователе
// @Description  Возвращает данные аутентифицированного пользователя
// @Tags         Auth
// @Produce      json
// @Security     CookieAuth
// @Success      200 {object} utils.StandardResponse{data=dto.WhoamiResponse} "Данные пользователя"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
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
// @Summary      Обновление токенов
// @Description  Обновляет access и refresh токены с использованием refresh токена из cookie
// @Tags         Auth
// @Produce      json
// @Success      200 {object} utils.StandardResponse{data=dto.LoginResponse} "Новые токены установлены"
// @Failure      401 {object} utils.ErrorResponseStruct "Refresh токен отсутствует или недействителен"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
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
			ID:        user.GetID(),
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, response)
}

// Logout godoc
// @Summary      Выход из текущей сессии
// @Description  Отзывает текущий access токен и очищает cookies
// @Tags         Auth
// @Security     CookieAuth
// @Success      200 {object} utils.StandardResponse{message=string} "Успешный выход"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
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
// @Summary      Выход из всех сессий
// @Description  Отзывает все токены пользователя и очищает cookies
// @Tags         Auth
// @Security     CookieAuth
// @Success      200 {object} utils.StandardResponse{message=string} "Выход из всех сессий выполнен"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
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
// @Summary      Инициация OAuth2 через Яндекс
// @Description  Перенаправляет пользователя на страницу авторизации Яндекса
// @Tags         OAuth
// @Success      302 "Редирект на страницу авторизации Яндекса"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /auth/oauth/yandex [get]
func (h *AuthHandler) OAuthYandex(c *gin.Context) {
	state := generateOAuthState()
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	authURL := h.oauthService.GetYandexAuthURL(state)
	c.Redirect(http.StatusFound, authURL)
}

// OAuthYandexCallback godoc
// @Summary      Обработка callback от Яндекса
// @Description  Обрабатывает callback от OAuth2 провайдера Яндекса и устанавливает сессию
// @Tags         OAuth
// @Param        code query string true "Код авторизации от Яндекса"
// @Param        state query string true "State параметр для защиты от CSRF"
// @Success      302 "Редирект на фронтенд с установленными cookies"
// @Failure      400 {object} utils.ErrorResponseStruct "Неверный state параметр"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
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
// @Summary      Инициация OAuth2 через ВКонтакте
// @Description  Перенаправляет пользователя на страницу авторизации ВКонтакте
// @Tags         OAuth
// @Success      302 "Редирект на страницу авторизации ВКонтакте"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /auth/oauth/vk [get]
func (h *AuthHandler) OAuthVK(c *gin.Context) {
	state := generateOAuthState()
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	authURL := h.oauthService.GetVKAuthURL(state)
	c.Redirect(http.StatusFound, authURL)
}

// OAuthVKCallback godoc
// @Summary      Обработка callback от ВКонтакте
// @Description  Обрабатывает callback от OAuth2 провайдера ВКонтакте и устанавливает сессию
// @Tags         OAuth
// @Param        code query string true "Код авторизации от ВКонтакте"
// @Param        state query string true "State параметр для защиты от CSRF"
// @Success      302 "Редирект на фронтенд с установленными cookies"
// @Failure      400 {object} utils.ErrorResponseStruct "Неверный state параметр"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
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

// Helper methods
func (h *AuthHandler) setAccessTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		"access_token",
		token,
		15*60,
		"/",
		"",
		false,
		true,
	)
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		"refresh_token",
		token,
		7*24*60*60,
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
	return "http://localhost:3000"
}

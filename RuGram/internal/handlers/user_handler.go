package handlers

import (
	"net/http"
	"strconv"

	"rugram-api/internal/dto"
	"rugram-api/internal/middleware"
	"rugram-api/internal/service"
	"rugram-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetByID godoc
// @Summary      Получение пользователя по ID
// @Description  Возвращает информацию о пользователе (только свои данные)
// @Tags         Users
// @Produce      json
// @Security     CookieAuth
// @Param        id path string true "ID пользователя (ObjectId)"
// @Success      200 {object} utils.StandardResponse{data=dto.UserResponse} "Информация о пользователе"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      403 {object} utils.ErrorResponseStruct "Доступ запрещен"
// @Failure      404 {object} utils.ErrorResponseStruct "Пользователь не найден"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	userVal, exists := c.Get("user")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	currentUser := userVal.(*middleware.AuthenticatedUser)
	if currentUser.ID != id {
		utils.ErrorResponse(c, http.StatusForbidden, "you can only access your own user data")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}

// GetByEmail godoc
// @Summary      Получение пользователя по email
// @Description  Возвращает информацию о пользователе по email
// @Tags         Users
// @Produce      json
// @Security     CookieAuth
// @Param        email path string true "Email пользователя"
// @Success      200 {object} utils.StandardResponse{data=dto.UserResponse} "Информация о пользователе"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      404 {object} utils.ErrorResponseStruct "Пользователь не найден"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /users/email/{email} [get]
func (h *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}

// Update godoc
// @Summary      Обновление пользователя
// @Description  Обновляет данные пользователя (только свои данные)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id path string true "ID пользователя (ObjectId)"
// @Param        request body dto.UpdateUserRequest true "Данные для обновления"
// @Success      200 {object} utils.StandardResponse{data=dto.UserResponse} "Пользователь успешно обновлен"
// @Failure      400 {object} utils.ErrorResponseStruct "Неверный формат запроса"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      403 {object} utils.ErrorResponseStruct "Доступ запрещен"
// @Failure      404 {object} utils.ErrorResponseStruct "Пользователь не найден"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /users/{id} [put]
// @Router       /users/{id} [patch]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	userVal, exists := c.Get("user")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	currentUser := userVal.(*middleware.AuthenticatedUser)
	if currentUser.ID != id {
		utils.ErrorResponse(c, http.StatusForbidden, "you can only update your own user data")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}

// Delete godoc
// @Summary      Удаление пользователя (soft delete)
// @Description  Помечает пользователя как удаленного (только свою учетную запись)
// @Tags         Users
// @Produce      json
// @Security     CookieAuth
// @Param        id path string true "ID пользователя (ObjectId)"
// @Success      204 "Пользователь успешно удален"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      403 {object} utils.ErrorResponseStruct "Доступ запрещен"
// @Failure      404 {object} utils.ErrorResponseStruct "Пользователь не найден"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	userVal, exists := c.Get("user")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	currentUser := userVal.(*middleware.AuthenticatedUser)
	if currentUser.ID != id {
		utils.ErrorResponse(c, http.StatusForbidden, "you can only delete your own user account")
		return
	}

	err := h.userService.DeleteUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.Status(http.StatusNoContent)
}

// GetAll godoc
// @Summary      Получение списка пользователей
// @Description  Возвращает список всех активных пользователей с пагинацией
// @Tags         Users
// @Produce      json
// @Security     CookieAuth
// @Param        page query int false "Номер страницы" default(1) minimum(1)
// @Param        limit query int false "Количество элементов на странице" default(10) minimum(1) maximum(100)
// @Success      200 {object} utils.StandardResponse{data=dto.PaginationResponse} "Список пользователей"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	result, err := h.userService.GetAllUsers(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get users: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

package handlers

import (
	"net/http"
	"strconv"

	"rugram-api/internal/dto"
	"rugram-api/internal/service"
	"rugram-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

// Create godoc
// @Summary      Создание нового поста
// @Description  Создает новый пост с указанными данными
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        request body dto.CreatePostRequest true "Данные для создания поста"
// @Success      201 {object} utils.StandardResponse{data=dto.PostResponse} "Пост успешно создан"
// @Failure      400 {object} utils.ErrorResponseStruct "Неверный формат запроса"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	post, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, post)
}

// GetByID godoc
// @Summary      Получение поста по ID
// @Description  Возвращает информацию о посте по его ID
// @Tags         Posts
// @Produce      json
// @Security     CookieAuth
// @Param        id path string true "ID поста (ObjectId)"
// @Success      200 {object} utils.StandardResponse{data=dto.PostResponse} "Информация о посте"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      404 {object} utils.ErrorResponseStruct "Пост не найден"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	post, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "post not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get post: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, post)
}

// Update godoc
// @Summary      Обновление поста
// @Description  Обновляет данные поста (PUT/PATCH)
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id path string true "ID поста (ObjectId)"
// @Param        request body dto.UpdatePostRequest true "Данные для обновления"
// @Success      200 {object} utils.StandardResponse{data=dto.PostResponse} "Пост успешно обновлен"
// @Failure      400 {object} utils.ErrorResponseStruct "Неверный формат запроса"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      404 {object} utils.ErrorResponseStruct "Пост не найден"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /posts/{id} [put]
// @Router       /posts/{id} [patch]
func (h *PostHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	post, err := h.service.Update(id, &req)
	if err != nil {
		if err.Error() == "post not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update post: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, post)
}

// Delete godoc
// @Summary      Удаление поста (soft delete)
// @Description  Помечает пост как удаленный (soft delete)
// @Tags         Posts
// @Produce      json
// @Security     CookieAuth
// @Param        id path string true "ID поста (ObjectId)"
// @Success      204 "Пост успешно удален"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      404 {object} utils.ErrorResponseStruct "Пост не найден"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(id)
	if err != nil {
		if err.Error() == "post not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete post: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAll godoc
// @Summary      Получение списка постов
// @Description  Возвращает список всех активных постов с пагинацией
// @Tags         Posts
// @Produce      json
// @Security     CookieAuth
// @Param        page query int false "Номер страницы" default(1) minimum(1)
// @Param        limit query int false "Количество элементов на странице" default(10) minimum(1) maximum(100)
// @Success      200 {object} utils.StandardResponse{data=dto.PaginationResponse} "Список постов"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /posts [get]
func (h *PostHandler) GetAll(c *gin.Context) {
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

	result, err := h.service.GetAll(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get posts: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetByUserID godoc
// @Summary      Получение постов пользователя
// @Description  Возвращает все посты указанного пользователя с пагинацией
// @Tags         Posts
// @Produce      json
// @Security     CookieAuth
// @Param        userId path string true "ID пользователя (ObjectId)"
// @Param        page query int false "Номер страницы" default(1) minimum(1)
// @Param        limit query int false "Количество элементов на странице" default(10) minimum(1) maximum(100)
// @Success      200 {object} utils.StandardResponse{data=dto.PaginationResponse} "Список постов пользователя"
// @Failure      401 {object} utils.ErrorResponseStruct "Не аутентифицирован"
// @Failure      500 {object} utils.ErrorResponseStruct "Внутренняя ошибка сервера"
// @Router       /posts/user/{userId} [get]
func (h *PostHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("userId")
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

	result, err := h.service.GetByUserID(userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user posts: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

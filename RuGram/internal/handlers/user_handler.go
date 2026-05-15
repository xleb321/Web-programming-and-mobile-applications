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

// GetByID returns user by ID
// GET /api/v1/users/:id
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

// GetByEmail returns user by email
// GET /api/v1/users/email/:email
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

// Update updates user information
// PUT /api/v1/users/:id
// PATCH /api/v1/users/:id
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

// Delete soft deletes user
// DELETE /api/v1/users/:id
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

// GetAll returns all users with pagination
// GET /api/v1/users
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

// GetProfile godoc
// @Summary Получить профиль текущего пользователя
// @Description Возвращает профиль текущего авторизованного пользователя
// @Tags Profile
// @Security BearerAuth
// @Success 200 {object} utils.StandardResponse{data=dto.ProfileResponse}
// @Failure 401 {object} utils.ErrorResponseStruct
// @Router /api/v1/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
    userVal, exists := c.Get("user")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }
    currentUser := userVal.(*middleware.AuthenticatedUser)

    profile, err := h.userService.GetProfile(currentUser.ID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(c, http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary Обновить профиль текущего пользователя
// @Description Обновляет профиль (displayName, bio, avatar)
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Данные для обновления"
// @Security BearerAuth
// @Success 200 {object} utils.StandardResponse{data=dto.ProfileResponse}
// @Failure 400 {object} utils.ErrorResponseStruct
// @Failure 401 {object} utils.ErrorResponseStruct
// @Router /api/v1/profile [post]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
    userVal, exists := c.Get("user")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }
    currentUser := userVal.(*middleware.AuthenticatedUser)

    var req dto.UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    profile, err := h.userService.UpdateProfile(currentUser.ID, &req)
    if err != nil {
        if err.Error() == "avatar file not found" {
            utils.ErrorResponse(c, http.StatusNotFound, err.Error())
            return
        }
        if err.Error() == "avatar file does not belong to user" {
            utils.ErrorResponse(c, http.StatusForbidden, err.Error())
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(c, http.StatusOK, profile)
}
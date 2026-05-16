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
// @Summary      Get user by ID
// @Description  Returns user details (only own profile)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200  {object}  utils.StandardResponse{data=dto.UserResponse}
// @Failure      401  {object}  utils.ErrorResponseStruct
// @Failure      403  {object}  utils.ErrorResponseStruct
// @Failure      404  {object}  utils.ErrorResponseStruct
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
// @Summary      Get user by email
// @Description  Returns user details by email
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        email path string true "User email"
// @Success      200  {object}  utils.StandardResponse{data=dto.UserResponse}
// @Failure      404  {object}  utils.ErrorResponseStruct
// @Failure      500  {object}  utils.ErrorResponseStruct
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
// @Summary      Update user
// @Description  Updates user profile (email, phone, password)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Param        request body dto.UpdateUserRequest true "Update data"
// @Success      200  {object}  utils.StandardResponse{data=dto.UserResponse}
// @Failure      400  {object}  utils.ErrorResponseStruct
// @Failure      401  {object}  utils.ErrorResponseStruct
// @Failure      403  {object}  utils.ErrorResponseStruct
// @Failure      404  {object}  utils.ErrorResponseStruct
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
// @Summary      Delete user (soft delete)
// @Description  Soft-deletes own user account
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      204  "No Content"
// @Failure      401  {object}  utils.ErrorResponseStruct
// @Failure      403  {object}  utils.ErrorResponseStruct
// @Failure      404  {object}  utils.ErrorResponseStruct
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
// @Summary      List users with pagination
// @Description  Returns a paginated list of all users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page  query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10) maximum(100)
// @Success      200  {object}  utils.StandardResponse{data=dto.PaginationResponse}
// @Failure      500  {object}  utils.ErrorResponseStruct
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
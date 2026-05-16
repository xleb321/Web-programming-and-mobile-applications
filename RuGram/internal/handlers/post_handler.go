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
// @Summary      Create a new post
// @Description  Creates a post with given data
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreatePostRequest true "Post data"
// @Success      201  {object}  utils.StandardResponse{data=dto.PostResponse}
// @Failure      400  {object}  utils.ErrorResponseStruct
// @Failure      500  {object}  utils.ErrorResponseStruct
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
// @Summary      Get post by ID
// @Description  Returns a single post by its ID
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Post ID"
// @Success      200  {object}  utils.StandardResponse{data=dto.PostResponse}
// @Failure      404  {object}  utils.ErrorResponseStruct
// @Failure      500  {object}  utils.ErrorResponseStruct
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
// @Summary      Update a post
// @Description  Updates post fields (full or partial)
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Post ID"
// @Param        request body dto.UpdatePostRequest true "Update data"
// @Success      200  {object}  utils.StandardResponse{data=dto.PostResponse}
// @Failure      400  {object}  utils.ErrorResponseStruct
// @Failure      404  {object}  utils.ErrorResponseStruct
// @Failure      500  {object}  utils.ErrorResponseStruct
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
// @Summary      Delete a post (soft delete)
// @Description  Soft-deletes a post by ID
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Post ID"
// @Success      204  "No Content"
// @Failure      404  {object}  utils.ErrorResponseStruct
// @Failure      500  {object}  utils.ErrorResponseStruct
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
// @Summary      List posts with pagination
// @Description  Returns a paginated list of posts
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page  query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10) maximum(100)
// @Success      200  {object}  utils.StandardResponse{data=dto.PaginationResponse}
// @Failure      500  {object}  utils.ErrorResponseStruct
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
// @Summary      Get posts by user ID
// @Description  Returns paginated posts of a specific user
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userId path string true "User ID"
// @Param        page  query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10) maximum(100)
// @Success      200  {object}  utils.StandardResponse{data=dto.PaginationResponse}
// @Failure      500  {object}  utils.ErrorResponseStruct
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
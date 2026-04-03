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

func (h *PostHandler) GetAll(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    // Validate pagination parameters
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
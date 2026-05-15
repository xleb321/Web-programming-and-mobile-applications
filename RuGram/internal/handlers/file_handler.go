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

type FileHandler struct {
    fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
    return &FileHandler{
        fileService: fileService,
    }
}

// UploadFile godoc
// @Summary Загрузить файл
// @Description Загружает файл в MinIO и сохраняет метаданные
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки (изображение или PDF)"
// @Security BearerAuth
// @Success 201 {object} utils.StandardResponse{data=dto.UploadFileResponse}
// @Failure 400 {object} utils.ErrorResponseStruct
// @Failure 401 {object} utils.ErrorResponseStruct
// @Failure 413 {object} utils.ErrorResponseStruct
// @Router /api/v1/files [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
    userVal, exists := c.Get("user")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }
    currentUser := userVal.(*middleware.AuthenticatedUser)

    file, header, err := c.Request.FormFile("file")
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "failed to get file: "+err.Error())
        return
    }
    defer file.Close()

    response, err := h.fileService.UploadFile(file, header, currentUser.ID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, response)
}

// GetFile godoc
// @Summary Скачать файл
// @Description Скачивает файл по ID (только владелец)
// @Tags Files
// @Produce octet-stream
// @Param id path string true "ID файла"
// @Security BearerAuth
// @Success 200 {file} binary
// @Failure 401 {object} utils.ErrorResponseStruct
// @Failure 403 {object} utils.ErrorResponseStruct
// @Failure 404 {object} utils.ErrorResponseStruct
// @Router /api/v1/files/{id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
    fileID := c.Param("id")

    userVal, exists := c.Get("user")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }
    currentUser := userVal.(*middleware.AuthenticatedUser)

    stream, file, err := h.fileService.GetFileStream(fileID, currentUser.ID)
    if err != nil {
        if err.Error() == "file not found" {
            utils.ErrorResponse(c, http.StatusNotFound, "File not found")
            return
        }
        if err.Error() == "access denied" {
            utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    defer stream.Close()

    c.Header("Content-Type", file.MimeType)
    c.Header("Content-Disposition", "attachment; filename="+file.OriginalName)
    c.Header("Content-Length", string(rune(file.Size)))

    c.DataFromReader(http.StatusOK, file.Size, file.MimeType, stream, nil)
}

// DeleteFile godoc
// @Summary Удалить файл
// @Description Удаляет файл (soft delete + из MinIO)
// @Tags Files
// @Param id path string true "ID файла"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 401 {object} utils.ErrorResponseStruct
// @Failure 403 {object} utils.ErrorResponseStruct
// @Failure 404 {object} utils.ErrorResponseStruct
// @Router /api/v1/files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
    fileID := c.Param("id")

    userVal, exists := c.Get("user")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }
    currentUser := userVal.(*middleware.AuthenticatedUser)

    err := h.fileService.DeleteFile(fileID, currentUser.ID)
    if err != nil {
        if err.Error() == "file not found" {
            utils.ErrorResponse(c, http.StatusNotFound, "File not found")
            return
        }
        if err.Error() == "access denied" {
            utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    c.Status(http.StatusNoContent)
}

// GetMyFiles godoc
// @Summary Получить мои файлы
// @Description Возвращает список файлов пользователя с пагинацией
// @Tags Files
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на странице" default(10)
// @Security BearerAuth
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse}
// @Failure 401 {object} utils.ErrorResponseStruct
// @Router /api/v1/files [get]
func (h *FileHandler) GetMyFiles(c *gin.Context) {
    page := 1
    limit := 10

    if p, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && p > 0 {
        page = p
    }
    if l, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil && l > 0 && l <= 100 {
        limit = l
    }

    userVal, exists := c.Get("user")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }
    currentUser := userVal.(*middleware.AuthenticatedUser)

    files, total, err := h.fileService.GetUserFiles(currentUser.ID, page, limit)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    totalPages := (total + int64(limit) - 1) / int64(limit)

    response := dto.PaginationResponse{
        Data: files,
        Meta: dto.MetaData{
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
        },
    }

    utils.SuccessResponse(c, http.StatusOK, response)
}
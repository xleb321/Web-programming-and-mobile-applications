package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"rugram-api/internal/cache"
	"rugram-api/internal/dto"
	"rugram-api/internal/models"
	"rugram-api/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService struct {
    fileRepo    *repository.FileRepository
    minioSvc    *MinioService
    cacheSvc    *cache.CacheService
    maxFileSize int64
}

func NewFileService(fileRepo *repository.FileRepository, minioSvc *MinioService, cacheSvc *cache.CacheService) (*FileService, error) {
    maxSizeStr := os.Getenv("MAX_FILE_SIZE")
    maxSize := int64(10 * 1024 * 1024) // 10MB default

    if maxSizeStr != "" {
        if size, err := strconv.ParseInt(maxSizeStr, 10, 64); err == nil {
            maxSize = size
        }
    }

    return &FileService{
        fileRepo:    fileRepo,
        minioSvc:    minioSvc,
        cacheSvc:    cacheSvc,
        maxFileSize: maxSize,
    }, nil
}

func (s *FileService) ValidateFile(header *multipart.FileHeader) error {
    if header.Size > s.maxFileSize {
        return fmt.Errorf("file too large: max %d bytes", s.maxFileSize)
    }

    contentType := header.Header.Get("Content-Type")
    allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "application/pdf"}

    allowed := false
    for _, t := range allowedTypes {
        if strings.HasPrefix(contentType, t) {
            allowed = true
            break
        }
    }

    if !allowed {
        return fmt.Errorf("file type not allowed: %s", contentType)
    }

    return nil
}

func (s *FileService) UploadFile(file multipart.File, header *multipart.FileHeader, userIDStr string) (*dto.UploadFileResponse, error) {
    if err := s.ValidateFile(header); err != nil {
        return nil, err
    }

    objectKey, size, mimeType, err := s.minioSvc.UploadFile(file, header, userIDStr)
    if err != nil {
        return nil, fmt.Errorf("failed to upload to MinIO: %w", err)
    }

    userID, err := primitive.ObjectIDFromHex(userIDStr)
    if err != nil {
        return nil, errors.New("invalid user ID format")
    }

    fileModel := &models.File{
        UserID:       userID,
        OriginalName: header.Filename,
        ObjectKey:    objectKey,
        Size:         size,
        MimeType:     mimeType,
        Bucket:       s.minioSvc.Bucket,
    }

    if err := s.fileRepo.Create(fileModel); err != nil {
        s.minioSvc.DeleteFile(objectKey)
        return nil, fmt.Errorf("failed to save file metadata: %w", err)
    }

    s.invalidateUserFileCache(userIDStr)

    return &dto.UploadFileResponse{
        ID:           fileModel.GetID(),
        OriginalName: fileModel.OriginalName,
        Size:         fileModel.Size,
        MimeType:     fileModel.MimeType,
        CreatedAt:    fileModel.CreatedAt.Format(time.RFC3339),
    }, nil
}

func (s *FileService) GetFileStream(fileID, userIDStr string) (io.ReadCloser, *models.File, error) {
    file, err := s.fileRepo.FindByIDString(fileID)
    if err != nil {
        return nil, nil, err
    }
    if file == nil {
        return nil, nil, errors.New("file not found")
    }

    if file.UserID.Hex() != userIDStr {
        return nil, nil, errors.New("access denied")
    }

    stream, err := s.minioSvc.GetFileStream(file.ObjectKey)
    if err != nil {
        return nil, nil, err
    }

    return stream, file, nil
}

func (s *FileService) DeleteFile(fileID, userIDStr string) error {
    file, err := s.fileRepo.FindByIDString(fileID)
    if err != nil {
        return err
    }
    if file == nil {
        return errors.New("file not found")
    }

    if file.UserID.Hex() != userIDStr {
        return errors.New("access denied")
    }

    objectID, err := primitive.ObjectIDFromHex(fileID)
    if err != nil {
        return errors.New("invalid file ID format")
    }

    if err := s.fileRepo.SoftDelete(objectID); err != nil {
        return err
    }

    if err := s.minioSvc.DeleteFile(file.ObjectKey); err != nil {
        fmt.Printf("Warning: failed to delete file from MinIO: %v\n", err)
    }

    s.invalidateUserFileCache(userIDStr)
    s.invalidateFileCache(fileID)

    return nil
}

func (s *FileService) GetUserFiles(userIDStr string, page, limit int) ([]dto.FileResponse, int64, error) {
    userID, err := primitive.ObjectIDFromHex(userIDStr)
    if err != nil {
        return nil, 0, errors.New("invalid user ID format")
    }

    offset := (page - 1) * limit
    files, total, err := s.fileRepo.GetUserFiles(userID, limit, offset)
    if err != nil {
        return nil, 0, err
    }

    responses := make([]dto.FileResponse, len(files))
    for i, file := range files {
        responses[i] = dto.FileResponse{
            ID:           file.GetID(),
            OriginalName: file.OriginalName,
            Size:         file.Size,
            MimeType:     file.MimeType,
            CreatedAt:    file.CreatedAt.Format(time.RFC3339),
        }
    }

    return responses, total, nil
}

func (s *FileService) GetFileMetadata(fileID string) (*models.File, error) {
    cacheKey := s.buildFileCacheKey(fileID)

    var file models.File
    if err := s.cacheSvc.Get(cacheKey, &file); err == nil && !file.ID.IsZero() {
        return &file, nil
    }

    filePtr, err := s.fileRepo.FindByIDString(fileID)
    if err != nil {
        return nil, err
    }
    if filePtr == nil {
        return nil, errors.New("file not found")
    }

    s.cacheSvc.SetWithDefaultTTL(cacheKey, filePtr)

    return filePtr, nil
}

func (s *FileService) buildFileCacheKey(fileID string) string {
    return s.cacheSvc.BuildKey("files", "meta", fileID)
}

func (s *FileService) invalidateFileCache(fileID string) {
    key := s.buildFileCacheKey(fileID)
    s.cacheSvc.Del(key)
}

func (s *FileService) invalidateUserFileCache(userID string) {
    pattern := s.cacheSvc.BuildKey("files", "user", userID, "*")
    s.cacheSvc.DelByPattern(pattern)
}
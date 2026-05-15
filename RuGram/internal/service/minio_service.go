package service

import (
    "context"
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
    client   *minio.Client
    Bucket   string
    endpoint string
    useSSL   bool
}

func NewMinioService() (*MinioService, error) {
    endpoint := os.Getenv("MINIO_ENDPOINT")
    accessKey := os.Getenv("MINIO_ACCESS_KEY")
    secretKey := os.Getenv("MINIO_SECRET_KEY")
    bucket := os.Getenv("MINIO_BUCKET")
    useSSL := os.Getenv("MINIO_USE_SSL") == "true"

    if endpoint == "" {
        endpoint = "localhost:9000"
    }
    if bucket == "" {
        bucket = "rugram-files"
    }
    if accessKey == "" {
        accessKey = "minioadmin"
    }
    if secretKey == "" {
        secretKey = "minioadmin"
    }

    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create MinIO client: %w", err)
    }

    ctx := context.Background()
    exists, err := client.BucketExists(ctx, bucket)
    if err != nil {
        return nil, fmt.Errorf("failed to check bucket existence: %w", err)
    }

    if !exists {
        err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
        if err != nil {
            return nil, fmt.Errorf("failed to create bucket: %w", err)
        }
    }

    return &MinioService{
        client:   client,
        Bucket:   bucket,
        endpoint: endpoint,
        useSSL:   useSSL,
    }, nil
}

func (s *MinioService) UploadFile(file multipart.File, header *multipart.FileHeader, userID string) (string, int64, string, error) {
    ext := ""
    if dotIndex := strings.LastIndex(header.Filename, "."); dotIndex != -1 {
        ext = header.Filename[dotIndex:]
    }
    objectKey := fmt.Sprintf("users/%s/%s%s", userID, uuid.New().String(), ext)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    info, err := s.client.PutObject(ctx, s.Bucket, objectKey, file, header.Size, minio.PutObjectOptions{
        ContentType: header.Header.Get("Content-Type"),
    })
    if err != nil {
        return "", 0, "", fmt.Errorf("failed to upload file: %w", err)
    }

    return objectKey, info.Size, header.Header.Get("Content-Type"), nil
}

func (s *MinioService) GetFileStream(objectKey string) (io.ReadCloser, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    obj, err := s.client.GetObject(ctx, s.Bucket, objectKey, minio.GetObjectOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to get object: %w", err)
    }

    return obj, nil
}

func (s *MinioService) DeleteFile(objectKey string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    err := s.client.RemoveObject(ctx, s.Bucket, objectKey, minio.RemoveObjectOptions{})
    if err != nil {
        return fmt.Errorf("failed to delete object: %w", err)
    }
    return nil
}

func (s *MinioService) FileExists(objectKey string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := s.client.StatObject(ctx, s.Bucket, objectKey, minio.StatObjectOptions{})
    if err != nil {
        errResponse := minio.ToErrorResponse(err)
        if errResponse.Code == "NoSuchKey" {
            return false, nil
        }
        return false, err
    }
    return true, nil
}

func (s *MinioService) GetFileURL(objectKey string) string {
    protocol := "http"
    if s.useSSL {
        protocol = "https"
    }
    return fmt.Sprintf("%s://%s/%s/%s", protocol, s.endpoint, s.Bucket, objectKey)
}
package repository

import (
	"database/sql"
	"time"

	"rugram-api/internal/models"

	"github.com/google/uuid"
)

type PostRepository struct {
    db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
    return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
    query := `
        INSERT INTO posts (id, user_id, title, description, image_url, status, likes_count, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at
    `
    
    now := time.Now()
    post.ID = uuid.New()
    post.CreatedAt = now
    post.UpdatedAt = now
    
    err := r.db.QueryRow(
        query,
        post.ID,
        post.UserID,
        post.Title,
        post.Description,
        post.ImageURL,
        post.Status,
        post.LikesCount,
        post.CreatedAt,
        post.UpdatedAt,
    ).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
    
    return err
}

func (r *PostRepository) FindByID(id uuid.UUID) (*models.Post, error) {
    query := `
        SELECT id, user_id, title, description, image_url, status, likes_count, created_at, updated_at, deleted_at
        FROM posts
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    post := &models.Post{}
    err := r.db.QueryRow(query, id).Scan(
        &post.ID,
        &post.UserID,
        &post.Title,
        &post.Description,
        &post.ImageURL,
        &post.Status,
        &post.LikesCount,
        &post.CreatedAt,
        &post.UpdatedAt,
        &post.DeletedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return post, err
}

func (r *PostRepository) Update(post *models.Post) error {
    query := `
        UPDATE posts
        SET title = $2, description = $3, image_url = $4, status = $5, updated_at = $6
        WHERE id = $1 AND deleted_at IS NULL
        RETURNING updated_at
    `
    
    post.UpdatedAt = time.Now()
    err := r.db.QueryRow(
        query,
        post.ID,
        post.Title,
        post.Description,
        post.ImageURL,
        post.Status,
        post.UpdatedAt,
    ).Scan(&post.UpdatedAt)
    
    return err
}

func (r *PostRepository) SoftDelete(id uuid.UUID) error {
    query := `
        UPDATE posts
        SET deleted_at = $2
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    result, err := r.db.Exec(query, id, time.Now())
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}

func (r *PostRepository) FindAll(limit, offset int) ([]models.Post, int64, error) {
    countQuery := `SELECT COUNT(*) FROM posts WHERE deleted_at IS NULL`
    var total int64
    err := r.db.QueryRow(countQuery).Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    
    query := `
        SELECT id, user_id, title, description, image_url, status, likes_count, created_at, updated_at, deleted_at
        FROM posts
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.Query(query, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    posts := []models.Post{}
    for rows.Next() {
        var post models.Post
        err := rows.Scan(
            &post.ID,
            &post.UserID,
            &post.Title,
            &post.Description,
            &post.ImageURL,
            &post.Status,
            &post.LikesCount,
            &post.CreatedAt,
            &post.UpdatedAt,
            &post.DeletedAt,
        )
        if err != nil {
            return nil, 0, err
        }
        posts = append(posts, post)
    }
    
    return posts, total, nil
}

func (r *PostRepository) FindByUserID(userID string, limit, offset int) ([]models.Post, int64, error) {
    countQuery := `
        SELECT COUNT(*) 
        FROM posts 
        WHERE user_id = $1 AND deleted_at IS NULL
    `
    var total int64
    err := r.db.QueryRow(countQuery, userID).Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    
    query := `
        SELECT id, user_id, title, description, image_url, status, likes_count, created_at, updated_at, deleted_at
        FROM posts
        WHERE user_id = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.Query(query, userID, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    posts := []models.Post{}
    for rows.Next() {
        var post models.Post
        err := rows.Scan(
            &post.ID,
            &post.UserID,
            &post.Title,
            &post.Description,
            &post.ImageURL,
            &post.Status,
            &post.LikesCount,
            &post.CreatedAt,
            &post.UpdatedAt,
            &post.DeletedAt,
        )
        if err != nil {
            return nil, 0, err
        }
        posts = append(posts, post)
    }
    
    return posts, total, nil
}
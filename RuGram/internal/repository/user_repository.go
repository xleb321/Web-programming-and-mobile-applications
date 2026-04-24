package repository

import (
    "database/sql"
    "time"
    
    "rugram-api/internal/models"
    "github.com/google/uuid"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (id, email, phone, password_hash, password_salt, yandex_id, vk_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at
    `
    
    now := time.Now()
    user.ID = uuid.New()
    user.CreatedAt = now
    user.UpdatedAt = now
    
    var phone, yandexID, vkID *string
    if user.Phone != nil {
        phone = user.Phone
    }
    if user.YandexID != nil {
        yandexID = user.YandexID
    }
    if user.VkID != nil {
        vkID = user.VkID
    }
    
    err := r.db.QueryRow(
        query,
        user.ID,
        user.Email,
        phone,
        user.PasswordHash,
        user.PasswordSalt,
        yandexID,
        vkID,
        user.CreatedAt,
        user.UpdatedAt,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    query := `
        SELECT id, email, phone, password_hash, password_salt, yandex_id, vk_id, created_at, updated_at, deleted_at
        FROM users
        WHERE email = $1 AND deleted_at IS NULL
    `
    
    user := &models.User{}
    var phone, yandexID, vkID *string
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Email,
        &phone,
        &user.PasswordHash,
        &user.PasswordSalt,
        &yandexID,
        &vkID,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.DeletedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    user.Phone = phone
    user.YandexID = yandexID
    user.VkID = vkID
    
    return user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
    query := `
        SELECT id, email, phone, password_hash, password_salt, yandex_id, vk_id, created_at, updated_at, deleted_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    user := &models.User{}
    var phone, yandexID, vkID *string
    err := r.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Email,
        &phone,
        &user.PasswordHash,
        &user.PasswordSalt,
        &yandexID,
        &vkID,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.DeletedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    user.Phone = phone
    user.YandexID = yandexID
    user.VkID = vkID
    
    return user, nil
}

func (r *UserRepository) FindByYandexID(yandexID string) (*models.User, error) {
    query := `
        SELECT id, email, phone, password_hash, password_salt, yandex_id, vk_id, created_at, updated_at, deleted_at
        FROM users
        WHERE yandex_id = $1 AND deleted_at IS NULL
    `
    
    user := &models.User{}
    var phone, yandexIDPtr, vkID *string
    err := r.db.QueryRow(query, yandexID).Scan(
        &user.ID,
        &user.Email,
        &phone,
        &user.PasswordHash,
        &user.PasswordSalt,
        &yandexIDPtr,
        &vkID,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.DeletedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    user.Phone = phone
    user.YandexID = yandexIDPtr
    user.VkID = vkID
    
    return user, nil
}

func (r *UserRepository) FindByVkID(vkID string) (*models.User, error) {
    query := `
        SELECT id, email, phone, password_hash, password_salt, yandex_id, vk_id, created_at, updated_at, deleted_at
        FROM users
        WHERE vk_id = $1 AND deleted_at IS NULL
    `
    
    user := &models.User{}
    var phone, yandexID, vkIDPtr *string
    err := r.db.QueryRow(query, vkID).Scan(
        &user.ID,
        &user.Email,
        &phone,
        &user.PasswordHash,
        &user.PasswordSalt,
        &yandexID,
        &vkIDPtr,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.DeletedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    user.Phone = phone
    user.YandexID = yandexID
    user.VkID = vkIDPtr
    
    return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
    query := `
        UPDATE users
        SET email = $2, phone = $3, password_hash = $4, password_salt = $5, 
            yandex_id = $6, vk_id = $7, updated_at = $8
        WHERE id = $1 AND deleted_at IS NULL
        RETURNING updated_at
    `
    
    user.UpdatedAt = time.Now()
    
    var phone, yandexID, vkID *string
    if user.Phone != nil {
        phone = user.Phone
    }
    if user.YandexID != nil {
        yandexID = user.YandexID
    }
    if user.VkID != nil {
        vkID = user.VkID
    }
    
    err := r.db.QueryRow(
        query,
        user.ID,
        user.Email,
        phone,
        user.PasswordHash,
        user.PasswordSalt,
        yandexID,
        vkID,
        user.UpdatedAt,
    ).Scan(&user.UpdatedAt)
    
    return err
}

func (r *UserRepository) SoftDelete(id uuid.UUID) error {
    query := `
        UPDATE users
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

func (r *UserRepository) FindAll(limit, offset int) ([]models.User, int64, error) {
    countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
    var total int64
    err := r.db.QueryRow(countQuery).Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    
    query := `
        SELECT id, email, phone, password_hash, password_salt, yandex_id, vk_id, created_at, updated_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.Query(query, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    users := []models.User{}
    for rows.Next() {
        var user models.User
        var phone, yandexID, vkID *string
        err := rows.Scan(
            &user.ID,
            &user.Email,
            &phone,
            &user.PasswordHash,
            &user.PasswordSalt,
            &yandexID,
            &vkID,
            &user.CreatedAt,
            &user.UpdatedAt,
            &user.DeletedAt,
        )
        if err != nil {
            return nil, 0, err
        }
        user.Phone = phone
        user.YandexID = yandexID
        user.VkID = vkID
        users = append(users, user)
    }
    
    return users, total, nil
}
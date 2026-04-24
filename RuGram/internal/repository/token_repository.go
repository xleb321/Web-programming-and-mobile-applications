package repository

import (
    "database/sql"
    "time"
    
    "rugram-api/internal/models"
    "github.com/google/uuid"
)

type TokenRepository struct {
    db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
    return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *models.UserToken) error {
    query := `
        INSERT INTO user_tokens (id, user_id, token_hash, token_salt, token_type, expires_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at
    `
    
    token.ID = uuid.New()
    token.CreatedAt = time.Now()
    
    err := r.db.QueryRow(
        query,
        token.ID,
        token.UserID,
        token.TokenHash,
        token.TokenSalt,
        token.TokenType,
        token.ExpiresAt,
        token.CreatedAt,
    ).Scan(&token.ID, &token.CreatedAt)
    
    return err
}

func (r *TokenRepository) FindValidAccessToken(userID uuid.UUID, tokenHash string) (*models.UserToken, error) {
    query := `
        SELECT id, user_id, token_hash, token_salt, token_type, expires_at, revoked, created_at
        FROM user_tokens
        WHERE user_id = $1 AND token_type = 'access' AND revoked = false AND expires_at > $2
    `
    
    now := time.Now()
    rows, err := r.db.Query(query, userID, now)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var token models.UserToken
        err := rows.Scan(
            &token.ID,
            &token.UserID,
            &token.TokenHash,
            &token.TokenSalt,
            &token.TokenType,
            &token.ExpiresAt,
            &token.Revoked,
            &token.CreatedAt,
        )
        if err != nil {
            continue
        }
        
        // Возвращаем первый валидный токен (обычно их немного)
        if token.TokenHash == tokenHash {
            return &token, nil
        }
    }
    
    return nil, sql.ErrNoRows
}

func (r *TokenRepository) FindValidRefreshToken(userID uuid.UUID, tokenHash string) (*models.UserToken, error) {
    query := `
        SELECT id, user_id, token_hash, token_salt, token_type, expires_at, revoked, created_at
        FROM user_tokens
        WHERE user_id = $1 AND token_type = 'refresh' AND revoked = false AND expires_at > $2
    `
    
    now := time.Now()
    rows, err := r.db.Query(query, userID, now)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var token models.UserToken
        err := rows.Scan(
            &token.ID,
            &token.UserID,
            &token.TokenHash,
            &token.TokenSalt,
            &token.TokenType,
            &token.ExpiresAt,
            &token.Revoked,
            &token.CreatedAt,
        )
        if err != nil {
            continue
        }
        
        if token.TokenHash == tokenHash {
            return &token, nil
        }
    }
    
    return nil, sql.ErrNoRows
}

func (r *TokenRepository) RevokeToken(tokenID uuid.UUID) error {
    query := `
        UPDATE user_tokens
        SET revoked = true
        WHERE id = $1
    `
    
    _, err := r.db.Exec(query, tokenID)
    return err
}

func (r *TokenRepository) RevokeAllUserTokens(userID uuid.UUID) error {
    query := `
        UPDATE user_tokens
        SET revoked = true
        WHERE user_id = $1 AND revoked = false
    `
    
    _, err := r.db.Exec(query, userID)
    return err
}

func (r *TokenRepository) RevokeAllUserTokensByType(userID uuid.UUID, tokenType string) error {
    query := `
        UPDATE user_tokens
        SET revoked = true
        WHERE user_id = $1 AND token_type = $2 AND revoked = false
    `
    
    _, err := r.db.Exec(query, userID, tokenType)
    return err
}
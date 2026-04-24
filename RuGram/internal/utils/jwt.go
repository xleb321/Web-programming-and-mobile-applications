package utils

import (
    "fmt"
    "os"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    Sub   string `json:"sub"`
    Type  string `json:"type"`
    jwt.RegisteredClaims
}

// CreateAccessToken создает Access JWT токен
func CreateAccessToken(userID string) (string, error) {
    secret := os.Getenv("JWT_ACCESS_SECRET")
    if secret == "" {
        secret = "default_access_secret_change_in_prod"
    }
    
    expirationMinutes := 15 // по умолчанию
    // можно загрузить из env
    
    claims := JWTClaims{
        Sub:  userID,
        Type: "access",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationMinutes) * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// CreateRefreshToken создает Refresh JWT токен
func CreateRefreshToken(userID string) (string, error) {
    secret := os.Getenv("JWT_REFRESH_SECRET")
    if secret == "" {
        secret = "default_refresh_secret_change_in_prod"
    }
    
    expirationDays := 7 // по умолчанию
    
    claims := JWTClaims{
        Sub:  userID,
        Type: "refresh",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationDays) * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// DecodeAccessToken расшифровывает Access токен
func DecodeAccessToken(tokenString string) (*JWTClaims, error) {
    secret := os.Getenv("JWT_ACCESS_SECRET")
    if secret == "" {
        secret = "default_access_secret_change_in_prod"
    }
    
    return decodeToken(tokenString, secret)
}

// DecodeRefreshToken расшифровывает Refresh токен
func DecodeRefreshToken(tokenString string) (*JWTClaims, error) {
    secret := os.Getenv("JWT_REFRESH_SECRET")
    if secret == "" {
        secret = "default_refresh_secret_change_in_prod"
    }
    
    return decodeToken(tokenString, secret)
}

func decodeToken(tokenString, secret string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}
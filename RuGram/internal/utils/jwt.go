package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	Sub   string `json:"sub"`
	Type  string `json:"type"`
	JTI   string `json:"jti"` // JWT ID для идентификации токена
	jwt.RegisteredClaims
}

// CreateAccessToken создает Access JWT токен с JTI
func CreateAccessToken(userID string) (string, string, error) {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = "default_access_secret_change_in_prod"
	}

	jti := uuid.New().String()
	expirationMinutes := 15

	claims := JWTClaims{
		Sub:  userID,
		Type: "access",
		JTI:  jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return tokenString, jti, nil
}

// CreateRefreshToken создает Refresh JWT токен
func CreateRefreshToken(userID string) (string, string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "default_refresh_secret_change_in_prod"
	}

	jti := uuid.New().String()
	expirationDays := 7

	claims := JWTClaims{
		Sub:  userID,
		Type: "refresh",
		JTI:  jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return tokenString, jti, nil
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
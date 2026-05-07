package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings" // <-- добавить импорт

	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует пароль с солью (bcrypt автоматически включает соль)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword проверяет пароль
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSalt генерирует случайную соль
func GenerateSalt() (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// HashToken хеширует токен с солью
func HashToken(token, salt string) string {
	combined := token + salt
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// GenerateSecureToken генерирует безопасный случайный токен
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// NormalizeEmail приводит email к нижнему регистру и обрезает пробелы
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ValidatePassword проверяет сложность пароля
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password too short (minimum 6 characters)")
	}

	// Проверка на кириллицу
	for _, r := range password {
		if r >= 0x0400 && r <= 0x04FF {
			return fmt.Errorf("password should not contain Cyrillic characters")
		}
	}

	return nil
}

// ValidateEmail проверяет формат email
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	// Простая проверка email
	if !contains(email, "@") || !contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexOf(s, substr) != -1))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

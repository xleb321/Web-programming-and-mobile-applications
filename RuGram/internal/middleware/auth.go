package middleware

import (
    "net/http"
    
    "rugram-api/internal/service"
    "rugram-api/pkg/utils"
    
    "github.com/gin-gonic/gin"
)

// AuthenticatedUser - структура для хранения данных аутентифицированного пользователя
type AuthenticatedUser struct {
    ID    string
    Email string
    Phone *string
}

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get access token from cookie
        accessToken, err := c.Cookie("access_token")
        if err != nil {
            utils.ErrorResponse(c, http.StatusUnauthorized, "authentication required")
            c.Abort()
            return
        }
        
        // Validate token and get user
        user, err := authService.GetUserFromAccessToken(accessToken)
        if err != nil {
            utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
            c.Abort()
            return
        }
        
        // Store user info in context
        c.Set("user", &AuthenticatedUser{
            ID:    user.ID.String(),
            Email: user.Email,
            Phone: user.Phone,
        })
        c.Set("userID", user.ID.String())
        
        c.Next()
    }
}

// OptionalAuthMiddleware attempts to authenticate but doesn't require it
func OptionalAuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        accessToken, err := c.Cookie("access_token")
        if err == nil {
            user, err := authService.GetUserFromAccessToken(accessToken)
            if err == nil {
                c.Set("user", &AuthenticatedUser{
                    ID:    user.ID.String(),
                    Email: user.Email,
                    Phone: user.Phone,
                })
                c.Set("userID", user.ID.String())
            }
        }
        c.Next()
    }
}
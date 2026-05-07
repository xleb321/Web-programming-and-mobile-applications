package middleware

import (
	"net/http"

	"rugram-api/internal/service"
	"rugram-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthenticatedUser struct {
	ID    string
	Email string
	Phone *string
}

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "authentication required")
			c.Abort()
			return
		}

		user, err := authService.GetUserFromAccessToken(accessToken)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user", &AuthenticatedUser{
			ID:    user.GetID(),
			Email: user.Email,
			Phone: user.Phone,
		})
		c.Set("userID", user.GetID())

		c.Next()
	}
}

func OptionalAuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err == nil {
			user, err := authService.GetUserFromAccessToken(accessToken)
			if err == nil {
				c.Set("user", &AuthenticatedUser{
					ID:    user.GetID(),
					Email: user.Email,
					Phone: user.Phone,
				})
				c.Set("userID", user.GetID())
			}
		}
		c.Next()
	}
}

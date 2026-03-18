package middleware

import (
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// Логирование ошибок
			for _, err := range c.Errors {
				// Здесь можно добавить логирование
				_ = err
			}

			// Отправляем общий ответ об ошибке, если еще не отправлен
			if !c.Writer.Written() {
				c.JSON(500, gin.H{"error": "Internal server error"})
			}
		}
	}
}

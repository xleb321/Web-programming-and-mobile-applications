package utils

import (
	"github.com/gin-gonic/gin"
)

// StandardResponse represents standard API response
// @Description Стандартный ответ API
type StandardResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponseStruct represents error response
// @Description Ответ с ошибкой
type ErrorResponseStruct struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Invalid request data"`
}

func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponseStruct{
		Success: false,
		Error:   message,
	})
}

func MessageResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Message: message,
	})
}

package utils

import (
	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

type ErrorResponseStruct struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
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
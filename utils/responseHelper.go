package utils

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func ErrorResponses(c *gin.Context, statusCode int, message string, errors []map[string]string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
		"errors":  errors,
	})
}

func ErrorResponse(c *gin.Context, status int, message string, err error) {
	c.JSON(status, gin.H{
		"status":  "error",
		"message": message,
		"error":  err.Error(),
	})
}
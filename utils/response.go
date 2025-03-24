package utils

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

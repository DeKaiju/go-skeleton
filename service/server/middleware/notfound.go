package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"data":    "",
			"message": "api not found",
		})
	}
}

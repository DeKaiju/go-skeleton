package middleware

import "github.com/gin-gonic/gin"

// NotFound handles 404
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    4000,
			"message": "api not found",
			"data":    "",
		})
	}
}

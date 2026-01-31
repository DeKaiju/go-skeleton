package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RecoverAtLast global exception handler
func RecoverAtLast() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(500, gin.H{
					"code":    5000,
					"message": fmt.Sprintf("panic: %v", r),
					"data":    "",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

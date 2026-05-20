package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dekaiju/go-skeleton/pkg/log"
)

func RecoverAtLast() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered: %v", r)
				c.JSON(http.StatusInternalServerError, gin.H{
					"data":    "",
					"message": "internal server error",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

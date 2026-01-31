package response

import (
	"github.com/gin-gonic/gin"
)

// Success response success
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "",
		"data":    data,
		"traceid": c.GetString("traceId"),
	})
}

// Fail response err message
func Fail(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"code":    getCodeByError(err),
		"message": err.Error(),
		"data":    "",
		"traceid": c.GetString("traceId"),
	})
}

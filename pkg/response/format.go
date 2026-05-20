package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "",
	})
}

func Fail(c *gin.Context, code int, err error) {
	message := ""
	if err != nil {
		message = err.Error()
	}

	c.JSON(code, gin.H{
		"data":    "",
		"message": message,
	})
}

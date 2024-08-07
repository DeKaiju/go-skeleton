package route

import (
	"github.com/gin-gonic/gin"

	"github.com/dekaiju/go-skeleton/internal/handler"
	"github.com/dekaiju/go-skeleton/internal/middleware"
)

// Init 初始化路由
func Init() *gin.Engine {
	app := gin.New()
	// 中间件
	app.Use(gin.Logger(), middleware.Cors(), middleware.RecoverAtLast(), middleware.TraceId())
	// 接口不存在
	app.NoRoute(middleware.NotFound())
	// 路由分组
	api := app.Group("/api/")

	// Welcome
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Welcome": "This is go-skeleton, build with Gin and Gorm",
		})
	})

	api.POST("/login", handler.Login)

	return app
}

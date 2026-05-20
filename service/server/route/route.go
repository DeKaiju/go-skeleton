package route

import (
	"github.com/gin-gonic/gin"

	"github.com/dekaiju/go-skeleton/pkg/response"
	"github.com/dekaiju/go-skeleton/service/server/handler"
	"github.com/dekaiju/go-skeleton/service/server/middleware"
)

func Init() *gin.Engine {
	app := gin.New()
	app.Use(gin.Logger(), middleware.Cors(), middleware.RecoverAtLast(), middleware.RequestLogger())
	app.NoRoute(middleware.NotFound())

	app.GET("/", func(c *gin.Context) {
		response.Success(c, gin.H{
			"name":    "go-skeleton",
			"message": "service is running",
		})
	})

	app.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	api := app.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.GET("/nonce", handler.GetNonce)
	authGroup.POST("/login", handler.Login)

	userGroup := api.Group("/user")
	userGroup.Use(middleware.AuthRequired())
	userGroup.GET("/profile", handler.GetProfile)

	return app
}

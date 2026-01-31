package route

import (
	"github.com/gin-gonic/gin"

	"github.com/dekaiju/go-skeleton/internal/handler"
	"github.com/dekaiju/go-skeleton/internal/middleware"
)

// Init initializes routes
func Init() *gin.Engine {
	app := gin.New()
	// Middleware
	app.Use(gin.Logger(), middleware.Cors(), middleware.RecoverAtLast(), middleware.TraceId())
	// Handle not found
	app.NoRoute(middleware.NotFound())
	// Route groups
	api := app.Group("/api/")

	// Welcome endpoint
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Welcome": "This is go-skeleton, build with Gin and Gorm",
		})
	})

	api.POST("/login", handler.Login)

	return app
}

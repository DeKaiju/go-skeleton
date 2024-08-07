package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"github.com/dekaiju/go-skeleton/internal/data"
	"github.com/dekaiju/go-skeleton/pkg/redis"
)

type requestHeader struct {
	Authorization string `header:"Authorization"`
}

// AuthRequired 登录验证
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := redis.Get()
		defer conn.Close()
		// 获取header信息
		header := requestHeader{}
		c.BindHeader(&header)
		if len(header.Authorization) == 0 {
			c.JSON(200, gin.H{
				"code":    4003,
				"message": "auth fail: no Authorization header",
				"data":    "",
			})
			c.Abort()
			return
		}

		// 获取secret
		authSecret := []byte(os.Getenv("JWT_SECRET"))

		// 获取jwt负载,验证
		token, err := jwt.Parse(header.Authorization, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return authSecret, nil
		})
		if err != nil {
			c.JSON(200, gin.H{
				"code":    4003,
				"message": fmt.Sprintf("auth fail: %v", err),
				"data":    "",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			address := claims["address"]
			if address == nil {
				c.JSON(200, gin.H{
					"code":    4003,
					"message": "auth fail: no address in claims",
					"data":    "",
				})
				c.Abort()
				return
			}

			// verify address exist in database
			_, err := data.QueryUser(c, address.(string))
			if err != nil && err == gorm.ErrRecordNotFound {
				c.JSON(200, gin.H{
					"code":    4003,
					"message": "auth fail: invalid address",
					"data":    "",
				})
				c.Abort()
				return
			}

			// verify token in redis
			//tokenString, err := redis.String(conn.Do("GET", address))
			//if err != nil || header.Authorization != tokenString {
			//	c.JSON(200, gin.H{
			//		"code":    4003,
			//		"message": "auth fail: invalid token",
			//		"data":    "",
			//	})
			//	c.Abort()
			//	return
			//}

			// store address for this context
			c.Set("address", address)
		} else {
			c.JSON(200, gin.H{
				"code":    4003,
				"message": "auth fail: can not get claims",
				"data":    "",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

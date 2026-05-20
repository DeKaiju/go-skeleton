package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/dekaiju/go-skeleton/data"
	"github.com/dekaiju/go-skeleton/pkg/auth"
	"github.com/dekaiju/go-skeleton/pkg/mysql"
	"github.com/dekaiju/go-skeleton/pkg/redis"
	"github.com/dekaiju/go-skeleton/pkg/response"
	"github.com/dekaiju/go-skeleton/types"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
		if tokenString == "" {
			response.Fail(c, http.StatusUnauthorized, types.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, types.BearerPrefix))
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, types.ErrInvalidToken)
			c.Abort()
			return
		}

		ctx := mysql.GinContextToContext(c)
		user, err := data.GetUserByID(ctx, claims.UserID)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		if user == nil || !strings.EqualFold(user.Address, claims.Address) {
			response.Fail(c, http.StatusUnauthorized, types.ErrUnauthorized)
			c.Abort()
			return
		}

		conn := redis.Get()
		defer conn.Close()

		cachedToken, err := redigo.String(conn.Do("GET", strings.ToLower(claims.Address)))
		if err == redigo.ErrNil || cachedToken != tokenString {
			response.Fail(c, http.StatusUnauthorized, types.ErrUnauthorized)
			c.Abort()
			return
		}
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.Set(types.ContextUserID, user.ID)
		c.Set(types.ContextAddress, user.Address)
		c.Next()
	}
}

package user

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spruceid/siwe-go"

	"github.com/dekaiju/go-skeleton/internal/types"
	"github.com/dekaiju/go-skeleton/pkg/cache"
	"github.com/dekaiju/go-skeleton/pkg/redis"
)

func Login(c *gin.Context, message, signature string) (*LoginResponse, error) {
	conn := redis.Get()
	defer conn.Close()

	// parse siwe message
	siweMessage, err := siwe.ParseMessage(message)
	if err != nil {
		return nil, err
	}

	address := strings.ToLower(siweMessage.GetAddress().String())
	// memory cache
	nonce, exist := cache.Get(address)
	if !exist {
		return nil, types.ErrInvalidNonce
	}

	ns := nonce.(string)
	_, err = siweMessage.Verify(signature, nil, &ns, nil)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address":   address,
		"timestamp": time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	_, err = conn.Do("SET", address, tokenString)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{JWT: tokenString}, nil
}

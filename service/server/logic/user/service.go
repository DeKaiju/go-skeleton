package user

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spruceid/siwe-go"

	"github.com/dekaiju/go-skeleton/config"
	"github.com/dekaiju/go-skeleton/data"
	"github.com/dekaiju/go-skeleton/pkg/auth"
	"github.com/dekaiju/go-skeleton/pkg/cache"
	"github.com/dekaiju/go-skeleton/pkg/log"
	"github.com/dekaiju/go-skeleton/pkg/mysql"
	"github.com/dekaiju/go-skeleton/pkg/redis"
	"github.com/dekaiju/go-skeleton/types"
)

func GetNonce(_ *gin.Context, req *NonceReq) (*NonceResponse, error) {
	address := strings.ToLower(strings.TrimSpace(req.Address))
	if address == "" {
		return nil, types.ErrInvalidAddress
	}

	nonce := uuid.NewString()
	cache.Set(address, nonce, time.Duration(types.AuthNonceTTLSeconds)*time.Second)

	return &NonceResponse{Nonce: nonce}, nil
}

func Login(c *gin.Context, req *LoginReq) (*LoginResponse, error) {
	siweMessage, err := siwe.ParseMessage(req.Message)
	if err != nil {
		return nil, types.ErrUnauthorized
	}

	address := strings.ToLower(siweMessage.GetAddress().String())
	nonce, exists := cache.Get(address)
	if !exists {
		return nil, types.ErrInvalidNonce
	}

	nonceValue, ok := nonce.(string)
	if !ok {
		return nil, types.ErrInvalidNonce
	}

	if _, err := siweMessage.Verify(req.Signature, nil, &nonceValue, nil); err != nil {
		return nil, types.ErrUnauthorized
	}

	ctx := mysql.GinContextToContext(c)
	user, err := data.CreateUserIfNotExists(ctx, address)
	if err != nil {
		log.Printf("failed to upsert user: %v", err)
		return nil, types.ErrInternalError
	}

	tokenString, err := auth.GenerateToken(user.ID, address)
	if err != nil {
		log.Printf("failed to sign jwt: %v", err)
		return nil, types.ErrInternalError
	}

	conn := redis.Get()
	defer conn.Close()

	if _, err := conn.Do("SETEX", address, config.GetIntEnv("JWT_EXPIRES", 7200), tokenString); err != nil {
		log.Printf("failed to persist session: %v", err)
		return nil, types.ErrInternalError
	}

	cache.Delete(address)

	return &LoginResponse{
		JWT:     tokenString,
		UserID:  user.ID,
		Address: address,
	}, nil
}

func GetProfile(c *gin.Context) (*ProfileResponse, error) {
	userID := c.GetInt64(types.ContextUserID)
	if userID == 0 {
		return nil, types.ErrUnauthorized
	}

	ctx := mysql.GinContextToContext(c)
	user, err := data.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return nil, types.ErrInternalError
	}
	if user == nil {
		return nil, types.ErrUserNotFound
	}

	return &ProfileResponse{
		ID:      user.ID,
		Address: user.Address,
	}, nil
}

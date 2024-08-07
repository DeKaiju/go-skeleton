package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GetToken 获取token
func GetToken(user int) string {
	authSecret := []byte(os.Getenv("JWT_SECRET"))
	jwtExpires, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES"))
	expiresTime := time.Duration(jwtExpires)

	type appClaims struct {
		User int `json:"user"`
		jwt.RegisteredClaims
	}

	claims := appClaims{
		user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresTime * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(authSecret)
	return ss
}

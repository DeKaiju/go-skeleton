package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/dekaiju/go-skeleton/config"
)

type Claims struct {
	UserID  int64  `json:"user_id"`
	Address string `json:"address"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, address string) (string, error) {
	authSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(authSecret) == 0 {
		return "", errors.New("JWT_SECRET is not set")
	}

	expiresAt := time.Now().Add(time.Duration(config.GetIntEnv("JWT_EXPIRES", 7200)) * time.Second)
	claims := Claims{
		UserID:  userID,
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(authSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	authSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(authSecret) == 0 {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return authSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

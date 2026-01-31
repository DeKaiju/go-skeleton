package data

import (
	"github.com/gin-gonic/gin"

	"github.com/dekaiju/go-skeleton/pkg/mysql"
)

func QueryUser(c *gin.Context, address string) (*User, error) {
	user := &User{}
	err := mysql.Instance(c).Where("address = ?", address).First(&user).Error
	return user, err
}

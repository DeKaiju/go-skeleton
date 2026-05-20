package data

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/dekaiju/go-skeleton/pkg/mysql"
)

func GetUserByID(ctx context.Context, id int64) (*User, error) {
	user := new(User)
	err := mysql.Instance(ctx).Where("id = ?", id).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CreateUserIfNotExists(ctx context.Context, address string) (*User, error) {
	user := &User{Address: address}
	if err := mysql.Instance(ctx).Where(User{Address: address}).FirstOrCreate(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

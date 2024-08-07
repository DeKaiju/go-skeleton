package data

import "time"

type User struct {
	ID        int64     `gorm:"column:id"`
	Address   string    `gorm:"column:address"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "user"
}

package data

import "time"

type User struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Address   string    `gorm:"column:address;type:varchar(255);uniqueIndex" json:"address"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

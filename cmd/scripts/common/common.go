package common

import (
	"os"

	"github.com/dekaiju/go-skeleton/config"
	"github.com/dekaiju/go-skeleton/data"
	"github.com/dekaiju/go-skeleton/pkg/log"
	"github.com/dekaiju/go-skeleton/pkg/mysql"
	"github.com/dekaiju/go-skeleton/pkg/redis"
)

func Depend() error {
	config.SetAppRoot(os.Args[0])

	configPath, err := config.LoadEnv(".env")
	if err != nil {
		log.Printf("failed to load config file %s: %v", configPath, err)
		return err
	}

	log.Setup()
	mysql.GetDB()
	if err := mysql.DB.AutoMigrate(&data.User{}); err != nil {
		log.Printf("failed to auto migrate sample models: %v", err)
		return err
	}

	return nil
}

func Release() {
	if mysql.DB != nil {
		sqlDB, err := mysql.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	_ = redis.Close()
}

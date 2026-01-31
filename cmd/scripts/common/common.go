package common

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/dekaiju/go-skeleton/config"
	"github.com/dekaiju/go-skeleton/pkg/mysql"
	"github.com/dekaiju/go-skeleton/pkg/redis"
)

func Depend() {
	// Root directory
	config.SetAppRoot(os.Args[0])
	// Config
	configFile := config.AppRoot + "/.env"
	godotenv.Load(configFile)
	// Connection pools
	mysql.GetDB()
	redis.GetCache()
}

func Release() {
	// Close MySQL
	sqlDB, _ := mysql.DB.DB()
	sqlDB.Close()
	// Close Redis
	redis.Get().Close()
}

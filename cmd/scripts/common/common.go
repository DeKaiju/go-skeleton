package common

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/dekaiju/go-skeleton/config"
	"github.com/dekaiju/go-skeleton/pkg/mysql"
	"github.com/dekaiju/go-skeleton/pkg/redis"
)

func Depend() {
	// 根目录
	config.SetAppRoot(os.Args[0])
	// 配置
	configFile := config.AppRoot + "/.env"
	godotenv.Load(configFile)
	// 连接池
	mysql.GetDB()
	redis.GetCache()
}

func Release() {
	// 关闭 MySQL
	sqlDB, _ := mysql.DB.DB()
	sqlDB.Close()
	// 关闭 Redis
	redis.Get().Close()
}

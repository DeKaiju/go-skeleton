package mysql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/dekaiju/go-skeleton/config"
	"github.com/dekaiju/go-skeleton/pkg/log"
	"github.com/dekaiju/go-skeleton/types"
)

var DB *gorm.DB

func GetDB() {
	DB = connectDbMySQL(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_CHARSET"),
	)

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(config.GetIntEnv("DB_MAX_OPEN_CONNECTIONS", 10))
	sqlDB.SetMaxIdleConns(config.GetIntEnv("DB_MAX_IDLE_CONNECTIONS", config.GetIntEnv("DB_MAX_CONNECTIONS", 10)))
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func connectDbMySQL(host, port, database, user, pass, charset string) *gorm.DB {
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		user,
		pass,
		host,
		port,
		database,
		charset,
	)

	sqlLogLevel := logger.Info
	if !config.GetBoolEnv("DB_LOG", true) {
		sqlLogLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: NewTraceLogger(sqlLogLevel, time.Second),
	})
	if err != nil {
		log.Fatalf("models.InitDbMySQL err: %v", err)
	}

	return db
}

func Instance(ctx context.Context) *gorm.DB {
	if ctx == nil {
		ctx = context.Background()
	}

	return DB.WithContext(ctx)
}

func GinContextToContext(c *gin.Context) context.Context {
	if c == nil {
		return context.Background()
	}

	ctx := c.Request.Context()
	if traceID := c.GetString(types.ContextTraceID); traceID != "" {
		ctx = context.WithValue(ctx, types.ContextTraceID, traceID)
	}

	return ctx
}

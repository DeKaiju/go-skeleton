package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/dekaiju/go-skeleton/config"
	"github.com/dekaiju/go-skeleton/data"
	"github.com/dekaiju/go-skeleton/pkg/log"
	"github.com/dekaiju/go-skeleton/pkg/mysql"
	"github.com/dekaiju/go-skeleton/pkg/redis"
	"github.com/dekaiju/go-skeleton/service/server/route"
)

var (
	configFile string
	port       string
	mode       string
	Server     = &cobra.Command{
		Use:   "server",
		Short: "Start HTTP API server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	Server.PersistentFlags().StringVarP(&configFile, "config", "c", ".env", "Start server with provided configuration file")
	Server.PersistentFlags().StringVarP(&port, "port", "p", "3000", "Tcp port server listening on")
	Server.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "Server mode: dev, test, prod")
}

func setup() error {
	config.SetAppRoot(os.Args[0])

	configPath, err := config.LoadEnv(configFile)
	if err != nil {
		log.Printf("failed to load config file %s: %v", configPath, err)
		return err
	}

	if mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Setup()
	log.Printf("loaded config from %s", configPath)

	mysql.GetDB()
	if err := mysql.DB.AutoMigrate(&data.User{}); err != nil {
		log.Printf("failed to auto migrate sample models: %v", err)
		return err
	}

	return nil
}

func run() error {
	defer func() {
		if mysql.DB == nil {
			return
		}

		sqlDB, err := mysql.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}()

	defer func() {
		_ = redis.Close()
	}()

	handler := route.Init()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	log.Printf("server listening on :%s", port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-quit.Done()
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server exiting")
	return nil
}

package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/dekaiju/go-skeleton/config"
	"github.com/dekaiju/go-skeleton/internal/route"
	"github.com/dekaiju/go-skeleton/pkg/mysql"
	"github.com/dekaiju/go-skeleton/pkg/redis"
)

var (
	configFile string
	port       string
	mode       string
	Server     = &cobra.Command{
		Use: "server",
		PreRun: func(cmd *cobra.Command, args []string) {
			welcome()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	Server.PersistentFlags().StringVarP(&configFile, "config", "c", ".env", "Start server with provided configuration file")
	Server.PersistentFlags().StringVarP(&port, "port", "p", "9551", "Tcp port server listening on")
	Server.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func welcome() {
	usageStr := `starting api server`
	log.Printf("%s\n", usageStr)
}

func setup() {
	config.SetAppRoot(os.Args[0])
	// Default config
	if configFile == ".env" {
		configFile = config.AppRoot + "/.env"
	}
	fmt.Printf("The config file path is %s\n", configFile)
	godotenv.Load(configFile)
	if mode == "prod" {
		// Production mode
		gin.SetMode(gin.ReleaseMode)
	}
	// Connection pools
	mysql.GetDB()
	redis.GetCache()
}

func run() error {
	// Close MySQL
	defer func() {
		sqlDB, _ := mysql.DB.DB()
		sqlDB.Close()
	}()

	// Close Redis
	defer func() {
		redis.Get().Close()
	}()

	// Routes
	handler := route.Init()

	// Port
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	fmt.Printf("The server is running on port %s\n", port)

	// Run service
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return nil
}

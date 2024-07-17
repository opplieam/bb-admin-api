package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/opplieam/bb-admin-api/internal/middleware"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

var build = "dev"

type config struct {
	Addr            string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("service", "bb-admin-api", "build", build)

	if err := run(logger); err != nil {
		logger.Error("start up", "error", err)
	}

}

func run(log *slog.Logger) error {
	readTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_READ_TIMEOUT", "5"))
	writeTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_IDLE_TIMEOUT", "120"))
	shutDownTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_SHUTDOWN_TIMEOUT", "20"))

	cfg := config{
		Addr:            utils.GetEnv("WEB_ADDR", ":3000"),
		WriteTimeout:    time.Duration(writeTimeout) * time.Second,
		ReadTimeout:     time.Duration(readTimeout) * time.Second,
		IdleTimeout:     time.Duration(idleTimeout) * time.Second,
		ShutdownTimeout: time.Duration(shutDownTimeout) * time.Second,
	}

	log.Info("start up", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	r := setupRoutes(log)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      r.Handler(),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("startup", "status", "api router started", "address", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info("shutdown", "status", "shutdown complete", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("could not shutdown gratefuly: %w", err)
		}
	}

	return nil
}

func setupRoutes(log *slog.Logger) *gin.Engine {
	var r *gin.Engine
	if build == "dev" {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(middleware.SLogger(log))
		r.Use(gin.Recovery())
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}

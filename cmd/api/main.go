package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

var build = "dev"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("service", "bb-admin-api", "build", build)

	if err := run(logger); err != nil {
		logger.Error("failed to start server", "error", err)
	}

}

func run(log *slog.Logger) error {

	log.Info("start up", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	r := gin.Default()

	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err := r.Run(":8080")
	if err != nil {
		return err
	}

	return nil
}

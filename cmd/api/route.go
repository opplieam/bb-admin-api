package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/middleware"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

func setupRoutes(log *slog.Logger) *gin.Engine {
	var r *gin.Engine
	if utils.GetEnv("WEB_SERVICE_ENV", "dev") == "dev" {
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

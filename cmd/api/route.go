package main

import (
	"errors"
	"log/slog"
	"math/rand/v2"
	"net/http"

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
	}

	r.Use(gin.Recovery())
	r.Use(middleware.SLogger(log))

	r.GET("/ping", func(c *gin.Context) {
		num := rand.IntN(2)
		if num == 0 {
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}

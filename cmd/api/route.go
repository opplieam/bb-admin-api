package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/middleware"
	"github.com/opplieam/bb-admin-api/internal/utils"
	"github.com/opplieam/bb-admin-api/internal/v1/probe"
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
	r.Use(middleware.SLogger(log, []string{"/v1/liveness"}))

	v1 := r.Group("/v1")
	healthCheck := probe.NewProbe(build)
	v1.GET("/liveness", healthCheck.LivenessHandler)

	return r
}

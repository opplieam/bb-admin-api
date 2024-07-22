package main

import (
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/db/healthcheck"
	"github.com/opplieam/bb-admin-api/internal/middleware"
	"github.com/opplieam/bb-admin-api/internal/utils"
	"github.com/opplieam/bb-admin-api/internal/v1/probe"
)

func setupRoutes(log *slog.Logger, db *sql.DB) *gin.Engine {
	var r *gin.Engine
	if utils.GetEnv("WEB_SERVICE_ENV", "dev") == "dev" {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	r.Use(gin.Recovery())
	r.Use(middleware.SLogger(log, []string{"/v1/liveness", "/v1/readiness"}))

	v1 := r.Group("/v1")

	healthCheckStore := healthcheck.NewStore(db)
	probeH := probe.NewHandler(build, healthCheckStore)
	v1.GET("/liveness", probeH.LivenessHandler)
	v1.GET("/readiness", probeH.ReadinessHandler)

	return r
}

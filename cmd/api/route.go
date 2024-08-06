package main

import (
	"database/sql"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/middleware"
	"github.com/opplieam/bb-admin-api/internal/store"
	"github.com/opplieam/bb-admin-api/internal/utils"
	"github.com/opplieam/bb-admin-api/internal/v1/probe"
	"github.com/opplieam/bb-admin-api/internal/v1/user"
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
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsConfig))
	r.Use(middleware.SLogger(log, []string{"/v1/liveness", "/v1/readiness"}))

	v1 := r.Group("/v1")

	healthCheckStore := store.NewHealthCheckStore(db)
	probeH := probe.NewHandler(build, healthCheckStore)
	v1.GET("/liveness", probeH.LivenessHandler)
	v1.GET("/readiness", probeH.ReadinessHandler)

	userStore := store.NewUserStore(db)
	userH := user.NewHandler(userStore)
	v1.POST("/login", userH.LoginHandler)
	v1.DELETE("/logout", userH.LogoutHandler)

	v1.Use(middleware.AuthorizationMiddleware())
	{
		v1.POST("/user", userH.CreateUser)
		v1.GET("/user", userH.GetAllUsers)
	}

	return r
}

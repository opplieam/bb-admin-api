package main

import (
	"strconv"
	"time"

	"github.com/opplieam/bb-admin-api/internal/utils"
)

type Config struct {
	Web WebConfig
}

type WebConfig struct {
	Addr            string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func NewConfig() *Config {
	writeTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_WRITE_TIMEOUT", "10"))
	readTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_READ_TIMEOUT", "5"))
	idleTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_IDLE_TIMEOUT", "120"))
	shutDownTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_SHUTDOWN_TIMEOUT", "20"))

	return &Config{
		Web: WebConfig{
			Addr:            utils.GetEnv("WEB_ADDR", ":3000"),
			WriteTimeout:    time.Duration(writeTimeout) * time.Second,
			ReadTimeout:     time.Duration(readTimeout) * time.Second,
			IdleTimeout:     time.Duration(idleTimeout) * time.Second,
			ShutdownTimeout: time.Duration(shutDownTimeout) * time.Second,
		},
	}
}

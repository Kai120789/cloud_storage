package app

import (
	"cloud/internal/config"
	"cloud/pkg/logger"
	"fmt"
)

func StartServer() {
	// init config
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// init logger
	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log := zapLog.ZapLogger

	_ = log

	// connect to db

	// init storage

	// init service

	// init handler

	// init router

	// start server
}

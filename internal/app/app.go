package app

import (
	"cloud/internal/config"
	"cloud/internal/service"
	"cloud/internal/storage/postgres"
	"cloud/internal/transport/http/handler"
	"cloud/internal/transport/http/router"
	"cloud/pkg/logger"
	"fmt"
	"net/http"

	"go.uber.org/zap"
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

	// connect to db
	dbConn, err := postgres.Connection(cfg.DBDSN)
	if err != nil {
		log.Fatal("error connect to db", zap.Error(err))
	}
	defer dbConn.Close()

	// init postgres
	db := postgres.New(dbConn, log)

	// init redis

	// init minio

	// init service
	serv := service.New(service.Storager{
		UserStorager: db,
		FileStorager: &db.FileStorage,
	}, log)

	// init handler
	handl := handler.New(handler.Service{
		UserService: &serv.UserService,
		FileService: &serv.FileService,
	}, log, cfg)

	// init router
	router := router.New(handl)

	// start server
	log.Info("starting server", zap.String("address", cfg.ServerAddress))

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}

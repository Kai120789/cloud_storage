package app

import (
	"cloud/internal/config"
	"cloud/internal/service"
	"cloud/internal/storage/postgres"
	"cloud/internal/storage/redis"
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
	db := postgres.NewPostgresStorage(dbConn, log)

	// init redis
	redisAddr := "localhost:6379"
	redisPassword := "" // Задайте пароль, если требуется
	redisDB := 0        // Используем базу данных 0

	// Создаем экземпляр RedisStorage
	redisStorage := redis.NewRedisStorage(redisAddr, redisPassword, redisDB)

	// init minio

	// init service
	serv := service.New(service.Storager{
		UserStorager:  db,
		FileStorager:  &db.FileStorage,
		TokenStorager: redisStorage,
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

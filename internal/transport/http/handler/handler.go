package handler

import (
	"cloud/internal/config"

	"go.uber.org/zap"
)

type Handler struct {
	UserHandler UserHandler
}

type Service struct {
	UserService UserHandlerer
}

func New(s Service, log *zap.Logger, cfg *config.Config) Handler {
	return Handler{
		UserHandler: NewUserHandler(s.UserService, log, cfg),
	}
}

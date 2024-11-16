package service

import (
	"go.uber.org/zap"
)

type Service struct {
	UserService UserService
	FileService FileService
}

type Storager struct {
	FileStorager FileStorager
	UserStorager UserStorager
}

func New(stor Storager, log *zap.Logger) *Service {
	return &TodoService{
		FileService: *NewFileService(stor.FileStorager, log),
		UserService: *NewUserService(stor.UserStorager, log),
	}
}

package service

import (
	"go.uber.org/zap"
)

type Service struct {
	UserService  UserService
	FileService  FileService
	TokenService TokenService
}

type Storager struct {
	FileStorager  FileStorager
	UserStorager  UserStorager
	TokenStorager TokenStorager
}

func New(stor Storager, log *zap.Logger) *Service {
	return &Service{
		FileService:  *NewFileService(stor.FileStorager, log),
		UserService:  *NewUserService(stor.UserStorager, log),
		TokenService: *NewTokenService(stor.TokenStorager, log),
	}
}

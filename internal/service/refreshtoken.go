package service

import (
	"go.uber.org/zap"
)

type TokenService struct {
	storage TokenStorager
	logger  *zap.Logger
}

type TokenStorager interface {
	WriteRefreshToken(userId uint, refreshTokenValue string) error
}

func NewTokenService(s TokenStorager, l *zap.Logger) *TokenService {
	return &TokenService{
		storage: s,
		logger:  l,
	}
}
func (t *TokenService) WriteRefreshToken(userId uint, refreshTokenValue string) error {
	err := t.storage.WriteRefreshToken(userId, refreshTokenValue)
	if err != nil {
		return err
	}

	return nil
}

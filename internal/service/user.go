package service

import (
	"cloud/internal/dto"
	"cloud/internal/models"
	"cloud/internal/utils/hash"

	"go.uber.org/zap"
)

type UserService struct {
	storage UserStorager
	redis   TokenStorager
	logger  *zap.Logger
}

type UserStorager interface {
	RegisterNewUser(body dto.User) (*models.User, error)
	AuthorizateUser(body dto.User) (*uint, *string, error)
}

type TokenStorager interface {
	WriteRefreshToken(userId uint, refreshTokenValue string) error
}

func NewUserService(s UserStorager, l *zap.Logger, r TokenStorager) *UserService {
	return &UserService{
		storage: s,
		logger:  l,
		redis:   r,
	}
}

func (t *UserService) RegisterNewUser(body dto.User) (*models.User, error) {
	passwordHash, err := hash.HashPassword(body.Password)
	if err != nil {
		return nil, err
	}

	body.Password = passwordHash

	token, err := t.storage.RegisterNewUser(body)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *UserService) AuthorizateUser(body dto.User) (*uint, error) {
	id, passwordHash, err := t.storage.AuthorizateUser(body)
	if err != nil {
		return nil, err
	}

	if !hash.CheckPasswordHash(body.Password, *passwordHash) {
		return nil, err
	}

	return id, nil
}

func (t *UserService) WriteRefreshToken(userId uint, refreshTokenValue string) error {
	err := t.redis.WriteRefreshToken(userId, refreshTokenValue)
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"cloud/internal/dto"
	"cloud/internal/models"
	"cloud/internal/utils/hash"

	"go.uber.org/zap"
)

type Service struct {
	storage Storager
	logger  *zap.Logger
}

type Storager interface {
	RegisterNewUser(body dto.User) (*models.UserToken, error)
	AuthorizateUser(body dto.User) (*uint, *string, error)
	WriteRefreshToken(userId uint, refreshTokenValue string) error
	UserLogout(id uint) error
}

func New(s Storager, l *zap.Logger) *Service {
	return &Service{
		storage: s,
		logger:  l,
	}
}

func (t *Service) RegisterNewUser(body dto.User) (*models.UserToken, error) {
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

func (t *Service) AuthorizateUser(body dto.User) (*uint, error) {
	id, passwordHash, err := t.storage.AuthorizateUser(body)
	if err != nil {
		return nil, err
	}

	if !hash.CheckPasswordHash(body.Password, *passwordHash) {
		return nil, err
	}

	return id, nil
}

func (t *Service) UserLogout(id uint) error {
	err := t.storage.UserLogout(uint(id))
	if err != nil {
		return err
	}

	return nil
}

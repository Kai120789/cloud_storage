package service

import (
	"cloud/internal/dto"
	"cloud/internal/models"

	"go.uber.org/zap"
)

type FileService struct {
	storage FileStorager
	logger  *zap.Logger
}

type FileStorager interface {
	UploadFile(dto dto.Object) (*models.Object, error)
	CreateFolder(dto dto.Object) (*models.Object, error)
	DeleteItem(path string) error
	RenameItem(dto dto.Object) (*models.Object, error)
	SearchFiles() error
	ListDirectory() error
}

func NewFileService(s FileStorager, l *zap.Logger) *FileService {
	return &FileService{
		storage: s,
		logger:  l,
	}
}

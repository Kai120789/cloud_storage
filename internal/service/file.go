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
	CreateNewFileOrFold(dto dto.Object) (*models.Object, error)
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

func (s *FileService) UploadFile(dto dto.Object) (*models.Object, error) {
	return nil, nil
}

func (s *FileService) CreateFolder(dto dto.Object) (*models.Object, error) {
	return nil, nil
}

func (s *FileService) DeleteItem(path string) error {
	return nil
}

func (s *FileService) RenameItem(dto dto.Object) (*models.Object, error) {
	return nil, nil
}

func (s *FileService) SearchFiles() error {
	return nil
}

func (s *FileService) ListDirectory() error {
	return nil
}

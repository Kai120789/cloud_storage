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
	file, err := s.storage.CreateNewFileOrFold(dto)
	if err != nil {
		s.logger.Error("upload file error", zap.Error(err))
		return nil, err
	}

	return file, nil
}

func (s *FileService) CreateFolder(dto dto.Object) (*models.Object, error) {
	folder, err := s.storage.CreateNewFileOrFold(dto)
	if err != nil {
		s.logger.Error("create folder error", zap.Error(err))
		return nil, err
	}

	return folder, nil
}

func (s *FileService) DeleteItem(path string) error {
	err := s.storage.DeleteItem(path)
	if err != nil {
		s.logger.Error("delete item error", zap.Error(err))
		return err
	}
	return nil
}

func (s *FileService) RenameItem(dto dto.Object) (*models.Object, error) {
	file, err := s.storage.RenameItem(dto)
	if err != nil {
		s.logger.Error("rename file error", zap.Error(err))
		return nil, err
	}

	return file, nil
}

func (s *FileService) SearchFiles() error {
	return nil
}

func (s *FileService) ListDirectory() error {
	return nil
}

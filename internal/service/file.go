package service

import (
	"cloud/internal/dto"
	"cloud/internal/models"
	"io"

	"go.uber.org/zap"
)

type FileService struct {
	storage FileStorager
	logger  *zap.Logger
}

type FileStorager interface {
	CreateNewFileOrFold(file io.Reader, dto dto.Object) (*models.Object, error)
	DeleteItem(path string) error
	RenameItem(file dto.Object, newName string) (*models.Object, error)
	SearchFiles(query string) ([]models.Object, error)
	ListDirectory(path string) ([]models.Object, error)
}

func NewFileService(s FileStorager, l *zap.Logger) *FileService {
	return &FileService{
		storage: s,
		logger:  l,
	}
}

func (s *FileService) UploadFile(file io.Reader, dto dto.Object) (*models.Object, error) {
	fileRet, err := s.storage.CreateNewFileOrFold(file, dto)
	if err != nil {
		s.logger.Error("upload file error", zap.Error(err))
		return nil, err
	}

	return fileRet, nil
}

func (s *FileService) CreateFolder(dto dto.Object) (*models.Object, error) {
	folder, err := s.storage.CreateNewFileOrFold(nil, dto)
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

func (s *FileService) RenameItem(file dto.Object, newName string) (*models.Object, error) {
	fileRet, err := s.storage.RenameItem(file, newName)
	if err != nil {
		s.logger.Error("rename file error", zap.Error(err))
		return nil, err
	}

	return fileRet, nil
}

func (s *FileService) SearchFiles(query string) ([]models.Object, error) {
	return s.storage.SearchFiles(query)
}

func (s *FileService) ListDirectory(path string) ([]models.Object, error) {
	return s.storage.ListDirectory(path)
}

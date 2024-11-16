package handler

import (
	"cloud/internal/config"
	"cloud/internal/dto"
	"cloud/internal/models"

	"go.uber.org/zap"
)

type FileHandler struct {
	service FileHandlerer
	logger  *zap.Logger
	config  *config.Config
}

type FileHandlerer interface {
	UploadFile(dto dto.Object) (*models.Object, error)
	CreateFolder(dto dto.Object) (*models.Object, error)
	DeleteItem(path string) error
	RenameItem(dto dto.Object) (*models.Object, error)
	SearchFiles() error
	ListDirectory() error
}

func NewFileHandler(s FileHandlerer, l *zap.Logger, c *config.Config) *FileHandler {
	return &FileHandler{
		service: s,
		logger:  l,
		config:  c,
	}
}

func (h *FileHandler) UploadFile(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

func (h *FileHandler) CreateFolder(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

func (h *FileHandler) DeleteItem(path string) error {
	return nil
}

func (h *FileHandler) RenameItem(dto dto.Object) (*models.Object, error) {
	var obj models.Object
	return &obj, nil
}

func (h *FileHandler) SearchFiles() error {
	return nil
}

func (h *FileHandler) ListDirectory() error {
	return nil
}

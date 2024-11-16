package handler

import (
	"cloud/internal/config"
	"cloud/internal/dto"
	"cloud/internal/models"
	"net/http"

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

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {

}

func (h *FileHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
}

func (h *FileHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
}

func (h *FileHandler) RenameItem(w http.ResponseWriter, r *http.Request) {
}

func (h *FileHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {
}

func (h *FileHandler) ListDirectory(w http.ResponseWriter, r *http.Request) {
}

package handler

import (
	"cloud/internal/config"
	"cloud/internal/dto"
	"cloud/internal/models"
	"encoding/json"
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
	var file dto.Object
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if file.Name == "" {
		http.Error(w, "file name cannot be empty", http.StatusBadRequest)
		return
	}

	fileRet, err := h.service.UploadFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fileRet)
}

func (h *FileHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	var folder dto.Object
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if folder.Name == "" {
		http.Error(w, "folder name cannot be empty", http.StatusBadRequest)
		return
	}

	folderRet, err := h.service.UploadFile(folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folderRet)
}

func (h *FileHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
}

func (h *FileHandler) RenameItem(w http.ResponseWriter, r *http.Request) {
}

func (h *FileHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {
}

func (h *FileHandler) ListDirectory(w http.ResponseWriter, r *http.Request) {
}

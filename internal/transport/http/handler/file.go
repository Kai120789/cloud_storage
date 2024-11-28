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
	SearchFiles(query string) ([]models.Object, error)
	ListDirectory(path string) ([]models.Object, error)
}

func NewFileHandler(s FileHandlerer, l *zap.Logger, c *config.Config) *FileHandler {
	return &FileHandler{
		service: s,
		logger:  l,
		config:  c,
	}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	h.logger.Info("File uploaded",
		zap.String("filename", handler.Filename),
		zap.String("path", path),
	)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
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
	var obj dto.Object
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteItem(obj.Path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FileHandler) RenameItem(w http.ResponseWriter, r *http.Request) {
	var obj dto.Object
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if obj.Name == "" {
		http.Error(w, "object name cannot be empty", http.StatusBadRequest)
		return
	}

	objRet, err := h.service.RenameItem(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(objRet)
}

func (h *FileHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {
}

func (h *FileHandler) ListDirectory(w http.ResponseWriter, r *http.Request) {
}

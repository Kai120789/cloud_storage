package handler

import (
	"cloud/internal/config"
	"cloud/internal/dto"
	"cloud/internal/models"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type FileHandler struct {
	service FileHandlerer
	logger  *zap.Logger
	config  *config.Config
}

type FileHandlerer interface {
	UploadFile(file io.Reader, dto dto.Object) (*models.Object, error)
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

	dtoObj := dto.Object{
		Name:    handler.Filename,
		Path:    path + "/" + handler.Filename,
		Content: nil,
	}

	uploadedFile, err := h.service.UploadFile(file, dtoObj)
	if err != nil {
		h.logger.Error("Error uploading file", zap.Error(err))
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	h.logger.Info("File uploaded",
		zap.String("filename", handler.Filename),
		zap.String("path", path),
	)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
	json.NewEncoder(w).Encode(uploadedFile)
}

func (h *FileHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {

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
	query := r.URL.Query().Get("query")

	files, err := h.service.SearchFiles(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (h *FileHandler) ListDirectory(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	objects, err := h.service.ListDirectory(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(objects)
}

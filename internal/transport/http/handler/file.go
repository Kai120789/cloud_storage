package handler

import (
	"cloud/internal/config"
	"cloud/internal/dto"
	"cloud/internal/models"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	RenameItem(file dto.Object, newName string) (*models.Object, error)
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
		Name: handler.Filename,
		Path: path + "/" + handler.Filename,
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
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	dtoObj := dto.Object{
		Name: "",
		Path: path + "/" + "",
	}

	uploadedFile, err := h.service.CreateFolder(dtoObj)
	if err != nil {
		h.logger.Error("Error uploading file", zap.Error(err))
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	h.logger.Info("File uploaded",
		zap.String("filename", ""),
		zap.String("path", path),
	)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
	json.NewEncoder(w).Encode(uploadedFile)
}

func (h *FileHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	name := chi.URLParam(r, "name")

	delPath := path + "/" + name

	if err := h.service.DeleteItem(delPath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FileHandler) RenameItem(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
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
		Name: handler.Filename,
		Path: path,
	}

	uploadedFile, err := h.service.RenameItem(dtoObj, name)
	if err != nil {
		h.logger.Error("Error uploading file", zap.Error(err))
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	h.logger.Info("File uploaded",
		zap.String("filename", name),
		zap.String("path", name),
	)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
	json.NewEncoder(w).Encode(uploadedFile)
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

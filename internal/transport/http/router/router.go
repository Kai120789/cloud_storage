package router

import (
	"cloud/internal/middleware"
	"cloud/internal/transport/http/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	RegisterNewUser(w http.ResponseWriter, r *http.Request)
	AuthorizateUser(w http.ResponseWriter, r *http.Request)
	//UserLogout(w http.ResponseWriter, r *http.Request)
	UploadFile(w http.ResponseWriter, r *http.Request)
	CreateFolder(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	RenameItem(w http.ResponseWriter, r *http.Request)
	SearchFiles(w http.ResponseWriter, r *http.Request)
	ListDirectory(w http.ResponseWriter, r *http.Request)
}

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.UserHandler.RegisterNewUser)
		r.Post("/login", h.UserHandler.AuthorizateUser)
		//r.With(middleware.JWT).Delete("/logout", h.UserHandler.UserLogout)
	})

	r.Route("/api/files", func(r chi.Router) {
		r.With(middleware.JWT).Post("/upload", h.FileHandler.UploadFile)
		r.With(middleware.JWT).Post("/folder", h.FileHandler.CreateFolder)
		r.With(middleware.JWT).Delete("/{name}", h.FileHandler.DeleteItem)
		r.With(middleware.JWT).Put("/{name}", h.FileHandler.RenameItem)
		r.With(middleware.JWT).Get("/", h.FileHandler.ListDirectory)
	})

	r.Route("/api/search", func(r chi.Router) {
		r.With(middleware.JWT).Get("/", h.FileHandler.SearchFiles)
	})

	return r
}

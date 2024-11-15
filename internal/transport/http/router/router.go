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
	UserLogout(w http.ResponseWriter, r *http.Request)
}

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("register", h.RegisterNewUser)
		r.Post("login", h.AuthorizateUser)
		r.With(middleware.JWT).Delete("Logout", h.UserLogout)
	})

	r.Route("api/files", func(r chi.Router) {
		r.With(middleware.JWT).Post("upload", h.UploadFile)
		r.With(middleware.JWT).Post("/folder", h.CreateFolder)
		r.With(middleware.JWT).Delete("/{id}", h.DeleteItem)
		r.With(middleware.JWT).Patch("/{id}", h.RenameItem)
		r.With(middleware.JWT).Get("/", h.ListDirectory)
	})

	r.Route("/api/search", func(r chi.Router) {
		r.With(middleware.JWT).Get("/", h.SearchFiles)
	})

	return r
}

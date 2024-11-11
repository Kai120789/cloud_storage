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

	return r
}

package handler

import (
	"cloud/internal/config"
	"cloud/internal/dto"
	"cloud/internal/models"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Handler struct {
	service Handlerer
	logger  *zap.Logger
	config  *config.Config
}

type Handlerer interface {
	RegisterNewUser(body dto.User) (*models.UserToken, error)
	AuthorizateUser(body dto.User) (*uint, error)
	WriteRefreshToken(userId uint, refreshTokenValue string) error
	UserLogout(id uint) error
}

func New(s Handlerer, l *zap.Logger, c *config.Config) Handler {
	return Handler{
		service: s,
		logger:  l,
		config:  c,
	}
}

// Register new user
func (h *Handler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Login == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	if _, err := h.service.RegisterNewUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login user
func (h *Handler) AuthorizateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Login == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := h.service.AuthorizateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessTokenValue, err := tokens.GenerateJWT(*userID, time.Now().Add(15*time.Minute))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	refreshTokenValue, err := tokens.GenerateJWT(*userID, time.Now().Add(2*time.Hour*24*30))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessTokenCokie := http.Cookie{
		Name:     "access_token",
		Value:    accessTokenValue,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenValue,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, &accessTokenCokie)
	http.SetCookie(w, &refreshTokenCokie)

	err = h.service.WriteRefreshToken(*userID, refreshTokenValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessTokenValue)
}

// Logout user
func (h *Handler) UserLogout(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Login == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := h.service.AuthorizateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.service.UserLogout(*userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiredCookie := time.Now().Add(-1 * time.Hour)

	accessTokenCokie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  expiredCookie,
		HttpOnly: true,
		Secure:   false,
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  expiredCookie,
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, &accessTokenCokie)
	http.SetCookie(w, &refreshTokenCokie)

	w.WriteHeader(http.StatusNoContent)
}

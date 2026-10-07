package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/service"
)

type AuthHandler struct {
	service *service.UserService
}

func NewAuthHandler(service *service.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.service.CreateUser(&user)

	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(request.Email, request.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

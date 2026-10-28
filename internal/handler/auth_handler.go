package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/service"

	"gorm.io/gorm"
)

type AuthHandler struct {
	service *service.UserService
}

func NewAuthHandler(service *service.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" example:"alice@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"alice@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

// Register godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} UserResponse
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Failure 500 {string} string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user := domain.User{
		Email:    req.Email,
		Password: req.Password,
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Email and password"
// @Success 200 {object} LoginResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request LoginRequest

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

	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// GetMe godoc
// @Summary Authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/users/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

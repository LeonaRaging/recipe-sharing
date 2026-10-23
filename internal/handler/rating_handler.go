package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/service"
	"strconv"

	"gorm.io/gorm"
)

type RatingHandler struct {
	service *service.RatingService
}

func NewRatingHandler(service *service.RatingService) *RatingHandler {
	return &RatingHandler{service: service}
}

type RateRecipeRequest struct {
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}

// RateRecipe godoc
// @Summary Rate a recipe
// @Tags ratings
// @Accept json
// @Produce json
// @Param body body RateRecipeRequest true "Rating Score and Comment"
// @Param id path int64 true "Recipe ID"
// @Security BearerAuth
// @Success 201 {object} domain.Rating
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/recipes/{id}/ratings [post]
func (h *RatingHandler) RateRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || recipeID <= 0 {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	var req RateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	rating, err := h.service.AddRating(recipeID, userID, req.Score, req.Comment)
	if err != nil {
		if errors.Is(err, service.ErrInvalidScore) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrRecipeNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to rate recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rating)
}

// GetRecipeRatings godoc
// @Summary Retrieve recipe ratings
// @Tags ratings
// @Produce json
// @Param id path int64 true "Recipe ID"
// @Success 200 {array} domain.Rating
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/recipes/{id}/ratings [get]
func (h *RatingHandler) GetRecipeRatings(w http.ResponseWriter, r *http.Request) {
	recipeID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || recipeID <= 0 {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	ratings, err := h.service.GetRatings(recipeID)
	if err != nil {
		if errors.Is(err, service.ErrRecipeNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve ratings", http.StatusInternalServerError)
		return
	}

	if ratings == nil {
		ratings = []domain.Rating{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ratings)
}

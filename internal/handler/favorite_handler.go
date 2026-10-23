package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"recipe-sharing/internal/service"
)

type FavoriteHandler struct {
	service *service.FavoriteService
}

func NewFavoriteHandler(service *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

func getUserIDFromContext(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value("user_id").(int64)
	return userID, ok
}

// AddFavorite godoc
// @Summary Add recipe to favorites
// @Tags favorites
// @Produce json
// @Param id path int64 true "Recipe ID"
// @Security BearerAuth
// @Success 201 {object} map[string]string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /api/recipes/{id}/favorites [post]
func (h *FavoriteHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recipeIDStr := r.PathValue("id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	if err := h.service.FavoriteRecipe(userID, recipeID); err != nil {
		http.Error(w, "Failed to favorite recipe (maybe already favorited)", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Recipe favorited successfully"}`))
}

// RemoveFavorite godoc
// @Summary Remove recipe from favorites
// @Tags favorites
// @Param id path int64 true "Recipe ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /api/recipes/{id}/favorites [delete]
func (h *FavoriteHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recipeIDStr := r.PathValue("id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	if err := h.service.UnfavoriteRecipe(userID, recipeID); err != nil {
		http.Error(w, "Failed to unfavorite recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetMyFavorites godoc
// @Summary List authenticated user's favorite recipes
// @Tags favorites
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Recipe
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /api/users/me/favorites [get]
func (h *FavoriteHandler) HandleGetMyFavorites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recipes, err := h.service.GetFavorites(userID)
	if err != nil {
		http.Error(w, "Failed to get favorites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recipes); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

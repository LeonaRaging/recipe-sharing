package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/service"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type RecipeHandler struct {
	service *service.RecipeService
}

func NewRecipeHandler(service *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		service: service,
	}
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recipe domain.Recipe

	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateRecipe(&recipe); err != nil {
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		return
	}
}

func (h *RecipeHandler) HandleRecipeByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		w.Header().Set("Allow", "GET, DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/recipes/"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		recipe, err := h.service.GetRecipe(id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Failed to get recipe", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(recipe); err != nil {
			http.Error(w, "Failed to encode recipe", http.StatusInternalServerError)
		}
	case http.MethodDelete:
		if err := h.service.DeleteRecipe(id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Recipe not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *RecipeHandler) HandleRecipes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateRecipe(w, r)
	case http.MethodGet:
		recipes, err := h.service.GetAllRecipes()
		if err != nil {
			http.Error(w, "Failed to get recipes", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recipes)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

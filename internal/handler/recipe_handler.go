package handler

import (
	"encoding/json"
	"net/http"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/service"
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

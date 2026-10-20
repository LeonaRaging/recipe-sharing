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

// @Summary Create recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body domain.Recipe true "Recipe"
// @Security BearerAuth
// @Success 201 {object} domain.Recipe
// @Failure 401 {string} string
// @Failure 400 {string} string
// @Router /api/recipes [post]
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

	recipe.UserID = r.Context().Value("user_id").(int64)

	if err := h.service.CreateRecipe(&recipe); err != nil {
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}

// HandleGetRecipeByID godoc
// @Summary Get recipe by ID
// @Tags recipes
// @Produce json
// @Param id path int64 true "Recipe ID"
// @Success 200 {object} domain.Recipe
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/recipes/{id} [get]
func (h *RecipeHandler) HandleGetRecipeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/recipes/"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

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
}

// HandleDeleteRecipeByID godoc
// @Summary Delete recipe
// @Tags recipes
// @Param id path int64 true "Recipe ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/recipes/{id} [delete]
func (h *RecipeHandler) HandleDeleteRecipeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/recipes/"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	if err := h.service.DeleteRecipe(id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, service.ErrForbidden) {
			http.Error(w, "Only owner can delete this recipe", http.StatusForbidden)
			return
		}

		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetRecipes godoc
// @Summary List recipes
// @Description Returns all recipes, optionally filtered by search text or user ID.
// @Tags recipes
// @Produce json
// @Param search query string false "Search title or description"
// @Param userid query int false "Filter by user ID"
// @Success 200 {array} domain.Recipe
// @Failure 500 {string} string
// @Router /api/recipes [get]
func (h *RecipeHandler) HandleGetRecipes(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	userIDStr := r.URL.Query().Get("userid")
	var userID int64
	if userIDStr != "" {
		var errParse error
		userID, errParse = strconv.ParseInt(userIDStr, 10, 64)
		if errParse != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}
	}
	recipes, err := h.service.GetAllRecipes(search, userID)
	if err != nil {
		http.Error(w, "Failed to get recipes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) HandlePostRecipes(w http.ResponseWriter, r *http.Request) {
	h.CreateRecipe(w, r)
}

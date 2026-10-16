package httpapi

import (
	"net/http"
	"recipe-sharing/internal/auth"
	"recipe-sharing/internal/handler"
)

func NewRouter(
	recipeHandler *handler.RecipeHandler,
	authHandler *handler.AuthHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	mux.HandleFunc("GET /api/recipes", recipeHandler.HandleGetRecipes)
	mux.HandleFunc("POST /api/recipes", auth.RequireAuth(recipeHandler.HandlePostRecipes))
	mux.HandleFunc("GET /api/recipes/{id}", recipeHandler.HandleGetRecipeByID)
	mux.HandleFunc("DELETE /api/recipes/{id}", auth.RequireAuth(recipeHandler.HandleDeleteRecipeByID))

	return mux
}

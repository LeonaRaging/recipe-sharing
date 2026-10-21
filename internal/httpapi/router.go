package httpapi

import (
	"net/http"
	"recipe-sharing/internal/auth"
	"recipe-sharing/internal/handler"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	recipeHandler *handler.RecipeHandler,
	authHandler *handler.AuthHandler,
	favoriteHandler *handler.FavoriteHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	mux.HandleFunc("GET /api/recipes", recipeHandler.HandleGetRecipes)
	mux.HandleFunc("POST /api/recipes", auth.RequireAuth(recipeHandler.HandlePostRecipes))
	mux.HandleFunc("GET /api/recipes/{id}", recipeHandler.HandleGetRecipeByID)
	mux.HandleFunc("DELETE /api/recipes/{id}", auth.RequireAuth(recipeHandler.HandleDeleteRecipeByID))

	mux.HandleFunc("GET /api/users/me/favorites", auth.RequireAuth(favoriteHandler.HandleGetMyFavorites))

	mux.HandleFunc("POST /api/recipes/{id}/favorites", auth.RequireAuth(favoriteHandler.HandleRecipeFavorites))
	mux.HandleFunc("DELETE /api/recipes/{id}/favorites", auth.RequireAuth(favoriteHandler.HandleRecipeFavorites))
	return mux
}

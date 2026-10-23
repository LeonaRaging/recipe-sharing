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
	ratingHandler *handler.RatingHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	mux.HandleFunc("GET /api/recipes", recipeHandler.HandleGetRecipes)
	mux.HandleFunc("POST /api/recipes", auth.RequireAuth(recipeHandler.HandlePostRecipes))
	mux.HandleFunc("GET /api/recipes/{id}", recipeHandler.HandleGetRecipeByID)
	mux.HandleFunc("DELETE /api/recipes/{id}", auth.RequireAuth(recipeHandler.HandleDeleteRecipeByID))

	mux.HandleFunc("GET /api/users/me", auth.RequireAuth(authHandler.GetMe))

	mux.HandleFunc("POST /api/recipes/{id}/ratings", auth.RequireAuth(ratingHandler.RateRecipe))
	mux.HandleFunc("GET /api/recipes/{id}/ratings", ratingHandler.GetRecipeRatings)

	mux.HandleFunc("GET /api/users/me/favorites", auth.RequireAuth(favoriteHandler.HandleGetMyFavorites))

	mux.HandleFunc("POST /api/recipes/{id}/favorites", auth.RequireAuth(favoriteHandler.AddFavorite))
	mux.HandleFunc("DELETE /api/recipes/{id}/favorites", auth.RequireAuth(favoriteHandler.RemoveFavorite))
	return mux
}

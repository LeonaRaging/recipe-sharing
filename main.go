package main

import (
	"fmt"
	"net/http"
	"recipe-sharing/internal/auth"
	"recipe-sharing/internal/handler"
	"recipe-sharing/internal/repository"
	"recipe-sharing/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=recipe_sharing port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		return
	}

	recipeRepository := repository.NewRecipeRepository(db)
	recipeService := service.NewRecipeService(recipeRepository)
	recipeHandler := handler.NewRecipeHandler(recipeService)

	//http.HandleFunc("/api/recipes", recipeHandler.CreateRecipe)
	http.HandleFunc("/api/recipes", auth.RequireAuth(recipeHandler.CreateRecipe))

	authHandler := setupAuth(db)
	http.HandleFunc("/api/auth/register", authHandler.Register)

	http.HandleFunc("/api/auth/login", authHandler.Login)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

func setupAuth(db *gorm.DB) *handler.AuthHandler {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	return handler.NewAuthHandler(userService)
}

package main

import (
	"log"
	"net/http"
	"recipe-sharing/internal/config"
	"recipe-sharing/internal/handler"
	"recipe-sharing/internal/httpapi"
	"recipe-sharing/internal/repository"
	"recipe-sharing/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "recipe-sharing/docs"
)

// @title Recipe Sharing API
// @version 1.0
// @description HTTP API for sharing recipes
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func run() error {
	cfg := config.Load()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})

	if err != nil {
		return err
	}

	recipeRepository := repository.NewRecipeRepository(db)
	recipeService := service.NewRecipeService(recipeRepository)
	recipeHandler := handler.NewRecipeHandler(recipeService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authHandler := handler.NewAuthHandler(userService)

	router := httpapi.NewRouter(recipeHandler, authHandler)

	server := &http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: router,
	}

	log.Println("Server running on :8080")

	return server.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

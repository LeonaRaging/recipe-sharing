package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"recipe-sharing/internal/auth"
	"recipe-sharing/internal/config"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/handler"
	"recipe-sharing/internal/httpapi"
	"recipe-sharing/internal/repository"
	"recipe-sharing/internal/service"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Invalid config: %w", err)
	}

	auth.InitKey(cfg.JWTKey)
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Error,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})

	if err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Recipe{},
		&domain.Rating{},
		&domain.Favorite{},
	); err != nil {
		return fmt.Errorf("run database migrations: %w", err)
	}

	recipeRepository := repository.NewRecipeRepository(db)
	recipeService := service.NewRecipeService(recipeRepository)
	recipeHandler := handler.NewRecipeHandler(recipeService)

	ratingRepo := repository.NewRatingRepository(db)
	ratingService := service.NewRatingService(ratingRepo, recipeRepository)
	ratingHandler := handler.NewRatingHandler(ratingService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authHandler := handler.NewAuthHandler(userService)

	favoriteRepo := repository.NewFavoriteRepository(db)
	favoriteService := service.NewFavoriteService(favoriteRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	router := httpapi.NewRouter(recipeHandler, authHandler, favoriteHandler, ratingHandler)

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

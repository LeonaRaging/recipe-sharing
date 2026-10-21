package service

import (
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/repository"
)

type FavoriteService struct {
	repo *repository.FavoriteRepository
}

func NewFavoriteService(repo *repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) FavoriteRecipe(userID, recipeID int64) error {
	return s.repo.AddFavorite(userID, recipeID)
}

func (s *FavoriteService) UnfavoriteRecipe(userID, recipeID int64) error {
	return s.repo.RemoveFavorite(userID, recipeID)
}

func (s *FavoriteService) GetFavorites(userID int64) ([]domain.Recipe, error) {
	return s.repo.GetUserFavorites(userID)
}

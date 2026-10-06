package service

import (
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/repository"
)

type RecipeService struct {
	repo *repository.RecipeRepository
}

func NewRecipeService(repo *repository.RecipeRepository) *RecipeService {
	return &RecipeService{
		repo: repo,
	}
}

func (s *RecipeService) CreateRecipe(recipe *domain.Recipe) error {
	err := s.repo.Create(recipe)

	if err != nil {
		return err
	}

	return nil
}

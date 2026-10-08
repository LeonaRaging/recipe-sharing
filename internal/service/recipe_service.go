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

func (s *RecipeService) GetRecipe(id int64) (*domain.Recipe, error) {
	return s.repo.GetByID(id)
}

func (s *RecipeService) DeleteRecipe(id int64) error {
	return s.repo.Delete(id)
}

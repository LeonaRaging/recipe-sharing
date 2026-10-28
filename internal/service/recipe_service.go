package service

import (
	"errors"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/repository"
)

var ErrForbidden = errors.New("forbidden")

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

func (s *RecipeService) GetAllRecipes(search string, userID int64, limit, offset int) ([]domain.Recipe, error) {
	return s.repo.GetAll(search, userID, limit, offset)
}

func (s *RecipeService) DeleteRecipe(id int64, userID int64) error {
	recipe, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if recipe.UserID != userID {
		return ErrForbidden
	}

	return s.repo.Delete(id)
}

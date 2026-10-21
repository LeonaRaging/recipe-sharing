package service

import (
	"errors"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/repository"
)

var ErrInvalidScore = errors.New("score must be between 1 and 5")
var ErrRecipeNotFound = errors.New("recipe not found")

type RatingService struct {
	ratingRepo *repository.RatingRepository
	recipeRepo *repository.RecipeRepository
}

func NewRatingService(ratingRepo *repository.RatingRepository, recipeRepo *repository.RecipeRepository) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		recipeRepo: recipeRepo,
	}
}

func (s *RatingService) AddRating(recipeID, userID int64, score int, comment string) (*domain.Rating, error) {
	if score < 1 || score > 5 {
		return nil, ErrInvalidScore
	}

	_, err := s.recipeRepo.GetByID(recipeID)
	if err != nil {
		return nil, ErrRecipeNotFound
	}

	rating := &domain.Rating{
		RecipeID: recipeID,
		UserID:   userID,
		Score:    score,
		Comment:  comment,
	}

	if err := s.ratingRepo.Upsert(rating); err != nil {
		return nil, err
	}

	return rating, nil
}

func (s *RatingService) GetRatings(recipeID int64) ([]domain.Rating, error) {
	_, err := s.recipeRepo.GetByID(recipeID)
	if err != nil {
		return nil, ErrRecipeNotFound
	}

	return s.ratingRepo.GetByRecipeID(recipeID)
}

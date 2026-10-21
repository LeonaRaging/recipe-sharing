package repository

import (
	"recipe-sharing/internal/domain"

	"gorm.io/gorm"
)

type RatingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) *RatingRepository {
	return &RatingRepository{
		db: db,
	}
}

func (r *RatingRepository) Upsert(rating *domain.Rating) error {
	var existing domain.Rating
	err := r.db.Where("recipe_id = ? AND user_id = ?", rating.RecipeID, rating.UserID).First(&existing).Error
	if err == nil {
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"score":      rating.Score,
			"comment":    rating.Comment,
			"created_at": rating.CreatedAt,
		}).Error
	}
	return r.db.Create(rating).Error
}

func (r *RatingRepository) GetByRecipeID(recipeID int64) ([]domain.Rating, error) {
	var ratings []domain.Rating
	err := r.db.
		Where("recipe_id = ?", recipeID).
		Order("created_at desc").
		Find(&ratings).Error
	return ratings, err
}

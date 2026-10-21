package repository

import (
	"recipe-sharing/internal/domain"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) AddFavorite(userID, recipeID int64) error {
	favorite := domain.Favorite{UserID: userID, RecipeID: recipeID}
	return r.db.Create(&favorite).Error
}

func (r *FavoriteRepository) RemoveFavorite(userID, recipeID int64) error {
	return r.db.Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Delete(&domain.Favorite{}).Error
}

func (r *FavoriteRepository) GetUserFavorites(userID int64) ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	err := r.db.Table("recipes").
		Joins("JOIN favorites ON favorites.recipe_id = recipes.id").
		Where("favorites.user_id = ?", userID).
		Find(&recipes).Error
	return recipes, err
}

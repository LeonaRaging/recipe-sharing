package repository

import (
	"recipe-sharing/internal/domain"

	"gorm.io/gorm"
)

type RecipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{
		db: db,
	}
}

func (r *RecipeRepository) Create(recipe *domain.Recipe) error {
	r.db.AutoMigrate(&domain.Recipe{})

	if err := r.db.Create(&recipe).Error; err != nil {
		return err
	}

	return nil
}

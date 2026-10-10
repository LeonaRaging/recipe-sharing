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

func (r *RecipeRepository) GetByID(id int64) (*domain.Recipe, error) {
	var recipe domain.Recipe

	if err := r.db.First(&recipe, id).Error; err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (r *RecipeRepository) GetAll() ([]domain.Recipe, error) {
	var recipes []domain.Recipe

	if err := r.db.Find(&recipes).Error; err != nil {
		return nil, err
	}

	return recipes, nil
}

func (r *RecipeRepository) Delete(id int64) error {
	result := r.db.Delete(&domain.Recipe{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

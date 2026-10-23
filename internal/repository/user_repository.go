package repository

import (
	"recipe-sharing/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	result := r.db.Where("email = ?", email).Limit(1).Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id int64) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

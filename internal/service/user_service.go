package service

import (
	"errors"
	"fmt"
	"recipe-sharing/internal/auth"
	"recipe-sharing/internal/domain"
	"recipe-sharing/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrEmailExists = errors.New("email already exists")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user *domain.User) error {
	_, err := s.repo.FindByEmail(user.Email)

	if err == nil {
		return ErrEmailExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	if err := s.repo.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(email string, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	return auth.GenerateToken(user.ID)
}

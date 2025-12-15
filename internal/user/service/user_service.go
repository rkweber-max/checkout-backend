package service

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/rkweber-max/checkout-backend/internal/user/domain"
	"github.com/rkweber-max/checkout-backend/internal/user/repository"
)

type UserService interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	List() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(user *domain.User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already in use")
	}

	hashPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(hashPassword)

	return s.repo.Create(user)
}

func (s *userService) GetByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetByEmail(email string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	return s.repo.FindByEmail(email)
}

func (s *userService) List() ([]domain.User, error) {
	return s.repo.List()
}

func (s *userService) Update(user *domain.User) error {
	if user.ID == 0 {
		return errors.New("Invalid user id")
	}

	return s.repo.Update(user)
}

func (s *userService) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid user id")
	}

	return s.repo.Delete(id)
}

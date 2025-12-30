package service

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/rkweber-max/checkout-backend/internal/auth"
	"github.com/rkweber-max/checkout-backend/internal/user/domain"
	"github.com/rkweber-max/checkout-backend/internal/user/repository"
	"github.com/rkweber-max/checkout-backend/pkg/config"
)

type UserService interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	List() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	Login(email, password string) (string, error)
}

type userService struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{repo: repo, config: cfg}
}

func (s *userService) Create(user *domain.User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Password = strings.TrimSpace(user.Password)

	if user.Password == "" {
		return errors.New("password cannot be empty")
	}

	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already in use")
	}

	log.Printf("Creating user with email: %s, password length: %d", user.Email, len(user.Password))

	hashPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("Error generating password hash: %v", err)
		return err
	}

	user.Password = string(hashPassword)
	log.Printf("Password hash generated, length: %d, starts with: %s", len(user.Password), user.Password[:min(10, len(user.Password))])

	err = s.repo.Create(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	log.Printf("User created successfully with ID: %d", user.ID)
	return nil
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

func (s *userService) Login(email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	log.Printf("Login attempt - Email: %s, Password length: %d", email, len(password))

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return "", err
	}
	if user == nil {
		log.Printf("User not found for email: %s", email)
		return "", errors.New("invalid credentials")
	}

	if user.Password == "" {
		log.Printf("Password is empty for user ID: %d, Email: %s", user.ID, user.Email)
		return "", errors.New("invalid credentials: password not set for user")
	}

	if len(user.Password) < 10 || !strings.HasPrefix(user.Password, "$2") {
		log.Printf("Invalid password hash format for user ID: %d, Email: %s, Hash: %s", user.ID, user.Email, user.Password[:min(20, len(user.Password))])
		return "", errors.New("invalid credentials: corrupted password hash")
	}

	log.Printf("Comparing password - User ID: %d, Email: %s, Hash length: %d, Hash prefix: %s, Input password length: %d",
		user.ID, user.Email, len(user.Password), user.Password[:min(10, len(user.Password))], len(password))

	testHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	log.Printf("Test hash generated from input password - Length: %d, Prefix: %s", len(testHash), string(testHash)[:min(10, len(testHash))])

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		log.Printf("Password comparison failed for user ID: %d, Error: %v", user.ID, err)
		log.Printf("Stored hash: %s...", user.Password[:min(30, len(user.Password))])
		return "", errors.New("invalid credentials")
	}

	log.Printf("Password comparison successful for user ID: %d", user.ID)

	if s.config.JWTSecret == "" {
		log.Printf("JWT_SECRET not configured")
		return "", errors.New("JWT_SECRET not configured")
	}

	log.Printf("Generating token for user ID: %d", user.ID)
	token, err := auth.GenerateToken(user.ID, s.config.JWTSecret, user.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	log.Printf("Token generated successfully, length: %d", len(token))
	return token, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

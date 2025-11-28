package service

import (
	"context"
	"errors"

	"github.com/rkweber-max/checkout-backend/internal/product"
	"github.com/rkweber-max/checkout-backend/internal/product/repository"
)

type ProductService interface {
	Create(ctx context.Context, p product.Product) (int64, error)
	GetAll(ctx context.Context) ([]product.Product, error)
	GetByID(ctx context.Context, id int64) (*product.Product, error)
	Update(ctx context.Context, p product.Product) error
	Delete(ctx context.Context, id int64) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, p product.Product) (int64, error) {
	if p.Name == "" {
		return 0, errors.New("product name cannot be empty")
	}

	if p.Price < 0 {
		return 0, errors.New("product price cannot be negative")
	}

	return s.repo.Create(ctx, p)
}

func (s *productService) GetAll(ctx context.Context) ([]product.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *productService) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) Update(ctx context.Context, p product.Product) error {
	if p.Name == "" {
		return errors.New("product name cannot be empty")
	}

	if p.Price < 0 {
		return errors.New("product price cannot be negative")
	}

	return s.repo.Update(ctx, p)
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

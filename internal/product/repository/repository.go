package repository

import (
	"context"

	"github.com/rkweber-max/checkout-backend/internal/product"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, p product.Product) (int64, error)
	FindAll(ctx context.Context) ([]product.Product, error)
	FindByID(ctx context.Context, id int64) (*product.Product, error)
	Update(ctx context.Context, p product.Product) error
	Delete(ctx context.Context, id int64) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p product.Product) (int64, error) {
	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return 0, err
	}
	return int64(p.ID), nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]product.Product, error) {
	var products []product.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id int64) (*product.Product, error) {
	var p product.Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) Update(ctx context.Context, p product.Product) error {
	return r.db.WithContext(ctx).Save(&p).Error
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&product.Product{}, id).Error
}

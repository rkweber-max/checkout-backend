package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rkweber-max/checkout-backend/internal/checkout/domain"
	"github.com/rkweber-max/checkout-backend/internal/product/repository"
)

type CheckoutService struct {
	repo repository.ProductRepository
}

func NewCheckoutService(repo repository.ProductRepository) *CheckoutService {
	return &CheckoutService{repo: repo}
}

func (s *CheckoutService) ProcessOrder(order domain.CheckoutRequest) (*domain.Order, error) {
	if len(order.ProductIDs) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	ctx := context.Background()
	var prices []float64

	for _, productID := range order.ProductIDs {
		product, err := s.repo.FindByID(ctx, int64(productID))
		if err != nil {
			return nil, fmt.Errorf("product with ID %d not found: %w", productID, err)
		}
		prices = append(prices, product.Price)
	}

	total := domain.CalculateTotalPrice(prices, order.PaymentType)
	newOrder := &domain.Order{
		Total:       total,
		PaymentType: order.PaymentType,
		Customer:    order.Customer,
	}

	return newOrder, nil
}

package repository

import (
	"context"
	"database/sql"

	"github.com/rkweber-max/checkout-backend/internal/product"
)

type ProductRepository interface {
	Create(ctx context.Context, p product.Product) (int64, error)
	FindAll(ctx context.Context) ([]product.Product, error)
	FindByID(ctx context.Context, id int64) (*product.Product, error)
	Update(ctx context.Context, p product.Product) error
	Delete(ctx context.Context, id int64) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p product.Product) (int64, error) {
	query := "INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id"

	var id int64
	err := r.db.QueryRowContext(ctx, query, p.Name, p.Description, p.Price).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]product.Product, error) {
	query := "SELECT id, name, description, price FROM products"
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		var p product.Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); 
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id int64) (*product.Product, error) {
	query := "SELECT id, name, description, price FROM products WHERE id = $1"

	var p product.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) Update(ctx context.Context, p product.Product) error {
	query := "UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4"

	_, err := r.db.ExecContext(ctx, query, p.Name, p.Description, p.Price, p.ID)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM products WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
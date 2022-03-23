package products

import (
	"context"
)

//go:generate mockgen -destination=../../mocks/mock_product.go -mock_names=Repository=ProductRepository -package=mocks github.com/artback/mvp/pkg/products Repository

type Repository interface {
	Get(ctx context.Context, name string) (*Product, error)
	Update(ctx context.Context, username string, product Update) error
	Insert(ctx context.Context, product Product) error
	Delete(ctx context.Context, username string, name string) error
}

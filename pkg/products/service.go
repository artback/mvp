package products

import "context"

//go:generate mockgen -destination=../../mocks/mock_product_service.go -mock_names=Service=ProductService -package=mocks github.com/artback/mvp/pkg/products Service
type Service interface {
	Get(ctx context.Context, name string) (*Product, error)
	Update(ctx context.Context, product Product) error
	Insert(ctx context.Context, product Product) error
	Delete(ctx context.Context, username string, name string) error
}

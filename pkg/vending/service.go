package vending

import (
	"context"
	"github.com/artback/mvp/pkg/products"
)

//go:generate mockgen -destination=../../mocks/mock_vending_service.go -mock_names=Service=VendingService -package=mocks github.com/artback/mvp/pkg/vending Service
type Service interface {
	IncrementDeposit(ctx context.Context, username string, deposit int) error
	GetAccount(ctx context.Context, username string) (*Response, error)
	BuyProduct(ctx context.Context, username string, product products.Product) error
	SetDeposit(ctx context.Context, username string, deposit int) error
}

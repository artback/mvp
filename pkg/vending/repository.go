package vending

import (
	"context"

	"github.com/artback/mvp/pkg/products"
)

//go:generate mockgen -destination=../../mocks/mock_vending_repository.go -mock_names=Repository=VendingRepsitory -package=mocks github.com/artback/mvp/pkg/vending Repository
type Repository interface {
	IncrementDeposit(ctx context.Context, username string, deposit int) error
	GetAccount(ctx context.Context, username string) (*Account, error)
	BuyProduct(ctx context.Context, username string, product products.Product) error
	SetDeposit(ctx context.Context, username string, deposit int) error
}

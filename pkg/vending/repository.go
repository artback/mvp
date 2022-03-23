package vending

import (
	"context"
	"github.com/artback/mvp/pkg/products"
)

//go:generate mockgen -destination=../../mocks/mock_vending.go -mock_names=Repository=VendingRepsitory -package=mocks github.com/artback/mvp/pkg/vending Repository
type Repository interface {
	IncrementDeposit(ctx context.Context, username string, deposit int) error
	GetAccount(ctx context.Context, username string) (*AccountResponse, error)
	BuyProduct(ctx context.Context, username string, product products.Update) error
	SetDeposit(ctx context.Context, username string, deposit int) error
}

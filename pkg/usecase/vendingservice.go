package usecase

import (
	"context"

	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/vending"
)

type VendingService struct {
	vending.Repository
	coin.Coins
}

func (v VendingService) GetAccount(ctx context.Context, username string) (*vending.Response, error) {
	account, err := v.Repository.GetAccount(ctx, username)
	if err != nil {
		return nil, err
	}

	return &vending.Response{
		Deposit:  change.New(v.Coins, account.Deposit),
		Products: account.Products,
		Spent:    account.Spent,
	}, nil
}

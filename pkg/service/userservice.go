package service

import (
	"context"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/users"
)

type UserService struct {
	coin.Coins
	users.Repository
}

func (u UserService) GetResponse(ctx context.Context, username string) (*users.Response, error) {
	user, err := u.Get(ctx, username)
	if err != nil {
		return nil, err
	}
	return &users.Response{
		Deposit:  change.New(u.Coins, user.Deposit),
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

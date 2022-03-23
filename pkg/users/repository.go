package users

import (
	"context"
)

//go:generate mockgen -destination=../../mocks/mock_users.go -mock_names=Repository=UserRepository -package=mocks github.com/artback/mvp/pkg/users Repository
type Repository interface {
	GetResponse(ctx context.Context, username string) (*Response, error)
	Get(ctx context.Context, username string) (*User, error)
	Insert(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, username string) error
}

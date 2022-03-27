package users

import "context"

//go:generate mockgen -destination=../../mocks/mock_users_service.go -mock_names=Service=UserService -package=mocks github.com/artback/mvp/pkg/users Service
type Service interface {
	GetResponse(ctx context.Context, username string) (*Response, error)
	Get(ctx context.Context, username string) (*User, error)
	Insert(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, username string) error
}

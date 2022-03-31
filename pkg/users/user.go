package users

import (
	"context"
	"github.com/artback/mvp/pkg/change"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     Role   `json:"role"  binding:"required"`
	Deposit  int    `json:"deposit" binding:"required"`
}

type Response struct {
	Username string         `json:"username" binding:"required"`
	Role     Role           `json:"role"  binding:"required"`
	Deposit  change.Deposit `json:"deposit" binding:"required"`
}

type nameKey string

const key = nameKey("username")

func GetUser(ctx context.Context) User {
	return ctx.Value(key).(User)
}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, key, user)
}

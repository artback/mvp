package security

import "context"

type User struct {
	Username string `json:"username" binding:"required"`
	Role     Role   `json:"role"  binding:"required"`
}
type nameKey string

const key = nameKey("username")

func GetUser(ctx context.Context) User {
	return ctx.Value(key).(User)
}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, key, user)
}

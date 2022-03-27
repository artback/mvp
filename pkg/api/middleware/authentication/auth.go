package authentication

import (
	"context"
	"github.com/artback/mvp/pkg/users"
	"net/http"
)

type nameKey string

const key = nameKey("username")

type Auth interface {
	Authenticate(roles ...users.Role) func(next http.Handler) http.Handler
}

func GetUserName(ctx context.Context) string {
	return ctx.Value(key).(string)
}

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, key, username)
}

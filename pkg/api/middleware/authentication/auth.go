package authentication

import (
	"context"
	"github.com/artback/mvp/pkg/users"
	"net/http"
)

type Auth interface {
	Authenticate(roles ...users.Role) func(next http.Handler) http.Handler
}

func FromCtx(ctx context.Context) string {
	return ctx.Value("username").(string)
}

func CtxWithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, "username", username)
}

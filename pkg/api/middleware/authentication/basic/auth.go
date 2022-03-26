package basic

import (
	"context"
	"errors"
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"net/http"
	"time"
)

var (
	WrongPasswordErr      = errors.New("password is wrong")
	MissingBasicHeaderErr = errors.New("missing basic auth header")
	WrongRoleErr          = errors.New("user is of wrong role")
)

type Auth struct {
	users.Service
}

func httpError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	var code int
	switch {
	case errors.Is(err, repository.EmptyErr{}):
		code = http.StatusUnauthorized
	case errors.Is(err, MissingBasicHeaderErr):
		code = http.StatusUnauthorized
	case errors.Is(err, WrongPasswordErr):
		code = http.StatusUnauthorized
	case errors.Is(err, WrongRoleErr):
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), code)
}

func (a Auth) Authenticate(roles ...users.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() {
				httpError(w, err)
			}()
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			u, p, ok := r.BasicAuth()
			if !ok {
				err = MissingBasicHeaderErr
				return
			}
			user, err := a.Get(ctx, u)
			if err != nil {
				return
			}
			if len(roles) > 0 && !user.IsRole(roles...) {
				err = WrongRoleErr
				return
			}
			if user.Password != p {
				err = WrongPasswordErr
				return
			}
			ctx = authentication.CtxWithUsername(r.Context(), user.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

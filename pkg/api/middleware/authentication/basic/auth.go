package basic

import (
	"errors"
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"net/http"
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
	var code int

	switch {
	case errors.Is(err, repository.EmptyError{}):
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
				if err != nil {
					httpError(w, err)
				}
			}()

			u, p, ok := r.BasicAuth()

			if !ok {
				err = MissingBasicHeaderErr
				return
			}

			user, err := a.Get(r.Context(), u)
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

			next.ServeHTTP(w, r.WithContext(authentication.WithUsername(r.Context(), user.Username)))
		}

		return http.HandlerFunc(fn)
	}
}

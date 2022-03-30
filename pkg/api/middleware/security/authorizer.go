package security

import (
	"github.com/artback/mvp/pkg/users"
	"github.com/casbin/casbin/v2"
	"net/http"
)

func Authorize(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user := users.GetUser(r.Context())
			method := r.Method
			path := r.URL.Path
			if ok, err := e.Enforce(user.Role, path, method); ok && err == nil {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, http.StatusText(403), 403)
			}
		}

		return http.HandlerFunc(fn)
	}
}

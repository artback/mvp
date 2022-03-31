package security

import (
	"fmt"
	"github.com/artback/mvp/pkg/users"
	"github.com/casbin/casbin/v2"
	"net/http"
)

func Authorize(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			role := users.GetUser(r.Context()).Role
			method := r.Method
			path := r.URL.Path
			ok, err := e.Enforce(string(role), path, method)
			fmt.Println(ok, role, path, path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			if ok {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			}
		}

		return http.HandlerFunc(fn)
	}
}

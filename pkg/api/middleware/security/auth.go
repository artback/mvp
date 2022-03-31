package security

import (
	"errors"
	"github.com/artback/mvp/pkg/users"
	"log"
	"net/http"
)

var (
	WrongPasswordErr = errors.New("password is wrong")
	MissingHeaderErr = errors.New("missing auth header")
)

// Auth interface is for allowing extension of other authentication protocols
type Auth interface {
	GetUser(r *http.Request) (*users.User, error)
}

func Authenticate(a Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := a.GetUser(r)
			log.Println(err)
			if err != nil {
				user = &users.User{
					Role: users.Anonymous,
				}
			}
			r = r.WithContext(users.WithUser(r.Context(), *user))
			next.ServeHTTP(w, r)
		})
	}
}

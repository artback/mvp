package security

import (
	"errors"
	"log"
	"net/http"
)

var (
	WrongPasswordErr = errors.New("password is wrong")
	MissingHeaderErr = errors.New("missing auth header")
)

// Auth interface is for allowing extension of other authentication protocols
type Auth interface {
	GetUser(r *http.Request) (*User, error)
}

func Authenticate(a Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := a.GetUser(r)
			log.Println(err)
			if err != nil {
				user = &User{
					Role: Anonymous,
				}
			}
			r = r.WithContext(WithUser(r.Context(), *user))
			next.ServeHTTP(w, r)
		})
	}
}

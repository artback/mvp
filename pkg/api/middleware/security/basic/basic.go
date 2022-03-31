package basic

import (
	"github.com/artback/mvp/pkg/api/middleware/security"
	"github.com/artback/mvp/pkg/pass"
	"github.com/artback/mvp/pkg/users"
	"net/http"
)

type Basic struct {
	Service users.Service
}

func (b Basic) GetUser(r *http.Request) (*users.User, error) {
	u, p, ok := r.BasicAuth()
	if !ok {
		return nil, security.MissingHeaderErr
	}
	user, err := b.Service.Get(r.Context(), u)
	if err != nil {
		return nil, err
	}
	if !pass.Compare(user.Password, p) {
		return nil, security.WrongPasswordErr
	}
	return &users.User{Username: user.Username, Password: p, Role: user.Role, Deposit: user.Deposit}, nil
}

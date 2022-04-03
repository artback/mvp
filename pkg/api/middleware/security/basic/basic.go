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

func (b Basic) GetUser(r *http.Request) (*security.User, error) {
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
	return &security.User{Username: user.Username, Role: user.Role}, nil
}

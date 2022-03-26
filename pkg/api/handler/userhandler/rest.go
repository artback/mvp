package userhandler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

var (
	InvalidUserFormErr = errors.New("invalid user form")
)

type restHandler struct {
	service users.Service
}

func httpError(w http.ResponseWriter, err error) {
	var code int
	switch {
	case errors.As(err, &repository.DuplicateErr{}):
		code = http.StatusConflict
	case errors.As(err, &repository.EmptyErr{}):
		code = http.StatusNotFound
	case errors.Is(err, InvalidUserFormErr):
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), code)
}

func (rest restHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := rest.getUser(r)
	if err != nil {
		httpError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&user); err != nil {
		httpError(w, err)
		return
	}
}
func (rest restHandler) getUser(r *http.Request) (*users.Response, error) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	return rest.service.GetResponse(ctx, chi.URLParam(r, "username"))
}

func (rest restHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	err := rest.updateUser(r)
	if err != nil {
		httpError(w, err)
	}
}
func (rest restHandler) updateUser(r *http.Request) error {
	user := users.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return InvalidUserFormErr
	}
	// Overwrite any username input from the request, Only the user can change its own data
	user.Username = authentication.FromCtx(r.Context())
	return rest.service.Update(r.Context(), user)
}
func (rest restHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := rest.deleteUser(r)
	if err != nil {
		httpError(w, err)
	}
}
func (rest restHandler) deleteUser(r *http.Request) error {
	username := authentication.FromCtx(r.Context())
	return rest.service.Delete(r.Context(), username)
}
func (rest restHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := rest.createUser(r)
	if err != nil {
		httpError(w, err)
	}
}
func (rest restHandler) createUser(r *http.Request) error {
	user := users.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return InvalidUserFormErr
	}
	return rest.service.Insert(r.Context(), user)
}

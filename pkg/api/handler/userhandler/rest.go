package userhandler

import (
	"encoding/json"
	"errors"
	"github.com/artback/mvp/pkg/api/middleware/security"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var InvalidUserFormErr = errors.New("invalid user form")

type RestHandler struct {
	users.Service
}

func httpError(w http.ResponseWriter, err error) {
	var code int

	switch {
	case errors.As(err, &repository.DuplicateError{}):
		code = http.StatusConflict
	case errors.As(err, &repository.EmptyError{}):
		code = http.StatusNotFound
	case errors.Is(err, InvalidUserFormErr):
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	http.Error(w, http.StatusText(code), code)
}

func (rest RestHandler) GetUser(w http.ResponseWriter, r *http.Request) {
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

func (rest RestHandler) getUser(r *http.Request) (*users.Response, error) {
	return rest.GetResponse(r.Context(), chi.URLParam(r, "username"))
}

func (rest RestHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if err := rest.updateUser(r); err != nil {
		httpError(w, err)
	}
}

func (rest RestHandler) updateUser(r *http.Request) error {
	user := users.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return InvalidUserFormErr
	}
	// Overwrite any username input from the request, Only the user can change its own data
	user.Username = security.GetUser(r.Context()).Username

	return rest.Update(r.Context(), user)
}

func (rest RestHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := rest.deleteUser(r)
	if err != nil {
		httpError(w, err)
	}
}

func (rest RestHandler) deleteUser(r *http.Request) error {
	username := security.GetUser(r.Context()).Username
	return rest.Delete(r.Context(), username)
}

func (rest RestHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := rest.createUser(r)
	if err != nil {
		httpError(w, err)
	}
}

func (rest RestHandler) createUser(r *http.Request) error {
	user := users.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return InvalidUserFormErr
	}

	return rest.Insert(r.Context(), user)
}

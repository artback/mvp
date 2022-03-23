package userhandler

import (
	"github.com/artback/mvp/pkg/authentication"
	"github.com/artback/mvp/pkg/users"
	"github.com/go-chi/chi/v5"
)

func Routes(auth authentication.Auth, repository users.Repository) chi.Router {
	r := chi.NewRouter()
	controller := restHandler{repository}
	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticate())
		r.Get("/{username}", controller.GetUser)
		r.Put("/", controller.UpdateUser)
		r.Delete("/", controller.DeleteUser)
	})
	r.Group(func(r chi.Router) {
		r.Post("/", controller.CreateUser)
	})
	return r
}

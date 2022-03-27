package userhandler

import (
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/users"
	"github.com/go-chi/chi/v5"
)

func Routes(auth authentication.Auth, service users.Service) chi.Router {
	router := chi.NewRouter()
	controller := restHandler{service}

	router.Group(func(r chi.Router) {
		r.Use(auth.Authenticate())
		r.Get("/{username}", controller.GetUser)
		r.Put("/", controller.UpdateUser)
		r.Delete("/", controller.DeleteUser)
	})

	router.Group(func(r chi.Router) {
		r.Post("/", controller.CreateUser)
	})

	return router
}

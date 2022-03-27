package producthandler

import (
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/users"
	"github.com/go-chi/chi/v5"
)

func Routes(auth authentication.Auth, repository products.Repository) chi.Router {
	controller := restHandler{repository}
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(auth.Authenticate())
		r.Get("/{product_name}", controller.GetProduct)
	})
	router.Group(func(r chi.Router) {
		r.Use(auth.Authenticate(users.Seller))
		r.Post("/", controller.CreateProduct)
		r.Put("/{product_name}", controller.UpdateProduct)
		r.Delete("/{product_name}", controller.DeleteProduct)
	})

	return router
}

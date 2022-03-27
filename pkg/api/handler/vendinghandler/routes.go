package vendinghandler

import (
	"net/http"

	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/users"
	"github.com/artback/mvp/pkg/vending"
	"github.com/go-chi/chi/v5"
)

func Routes(auth authentication.Auth, repository vending.Service) http.Handler {
	rest := restHandler{repository}

	router := chi.NewRouter()
	router.Use(auth.Authenticate(users.Buyer))
	router.Get("/deposit", rest.GetAccount)
	router.Put("/deposit", rest.Deposit)
	router.Post("/buy/{product_name}", rest.BuyProduct)
	router.Delete("/reset", rest.ResetDeposit)

	return router
}

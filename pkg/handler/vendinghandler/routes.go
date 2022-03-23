package vendinghandler

import (
	"github.com/artback/mvp/pkg/authentication"
	"github.com/artback/mvp/pkg/users"
	"github.com/artback/mvp/pkg/vending"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes(auth authentication.Auth, repository vending.Repository) http.Handler {
	rest := restHandler{repository}

	r := chi.NewRouter()
	r.Use(auth.Authenticate(users.Buyer))
	r.Get("/deposit", rest.GetAccount)
	r.Put("/deposit", rest.Deposit)
	r.Post("/buy/{product_name}", rest.BuyProduct)
	r.Delete("/reset", rest.ResetDeposit)
	return r
}

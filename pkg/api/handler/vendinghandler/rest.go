package vendinghandler

import (
	"encoding/json"
	"errors"
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/vending"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type restHandler struct {
	vending.Service
}

func httpError(w http.ResponseWriter, err error) {
	var code int

	switch {
	case errors.Is(err, repository.EmptyError{}):
		code = http.StatusNotFound
	case errors.Is(err, repository.InvalidError{}):
		code = http.StatusNotAcceptable
	default:
		code = http.StatusInternalServerError
	}

	http.Error(w, err.Error(), code)
}

func (re restHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	account, err := re.getAccount(r)
	if err != nil {
		httpError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&account); err != nil {
		httpError(w, err)
	}
}

func (re restHandler) getAccount(r *http.Request) (*vending.Response, error) {
	username := authentication.GetUserName(r.Context())
	return re.Service.GetAccount(r.Context(), username)
}

func (re restHandler) ResetDeposit(w http.ResponseWriter, r *http.Request) {
	err := re.resetDeposit(r)
	if err != nil {
		httpError(w, err)
	}
}

func (re restHandler) resetDeposit(r *http.Request) error {
	username := authentication.GetUserName(r.Context())
	return re.SetDeposit(r.Context(), username, 0)
}

func (re restHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	if err := re.deposit(r); err != nil {
		httpError(w, err)
	}
}

func (re restHandler) deposit(r *http.Request) error {
	deposit := change.Deposit{}
	if err := json.NewDecoder(r.Body).Decode(&deposit); err != nil {
		return err
	}

	username := authentication.GetUserName(r.Context())

	return re.IncrementDeposit(r.Context(), username, deposit.ToAmount())
}

func (re restHandler) BuyProduct(w http.ResponseWriter, r *http.Request) {
	err := re.buyProduct(r)
	if err != nil {
		httpError(w, err)
	}
}

func atoiWithDefault(str string, def int) int {
	amount, err := strconv.Atoi(str)
	if err != nil {
		return def
	}

	return amount
}

func (re restHandler) buyProduct(r *http.Request) error {
	amount := r.URL.Query().Get("amount")
	username := authentication.GetUserName(r.Context())

	return re.Service.BuyProduct(r.Context(), username, products.Product{
		Name:   chi.URLParam(r, "product_name"),
		Amount: atoiWithDefault(amount, 1),
	})
}

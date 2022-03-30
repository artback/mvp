package vendinghandler

import (
	"encoding/json"
	"errors"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"github.com/artback/mvp/pkg/vending"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type RestHandler struct {
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

func (re RestHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	account, err := re.getAccount(r)
	if err != nil {
		httpError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&account); err != nil {
		httpError(w, err)
	}
}

func (re RestHandler) getAccount(r *http.Request) (*vending.Response, error) {
	username := users.GetUser(r.Context()).Username
	return re.Service.GetAccount(r.Context(), username)
}

func (re RestHandler) ResetDeposit(w http.ResponseWriter, r *http.Request) {
	err := re.resetDeposit(r)
	if err != nil {
		httpError(w, err)
	}
}

func (re RestHandler) resetDeposit(r *http.Request) error {
	username := users.GetUser(r.Context()).Username
	return re.SetDeposit(r.Context(), username, 0)
}

func (re RestHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	if err := re.deposit(r); err != nil {
		httpError(w, err)
	}
}

func (re RestHandler) deposit(r *http.Request) error {
	deposit := change.Deposit{}
	if err := json.NewDecoder(r.Body).Decode(&deposit); err != nil {
		return err
	}

	username := users.GetUser(r.Context()).Username
	return re.IncrementDeposit(r.Context(), username, deposit.ToAmount())
}

func (re RestHandler) BuyProduct(w http.ResponseWriter, r *http.Request) {
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

func (re RestHandler) buyProduct(r *http.Request) error {
	amount := r.URL.Query().Get("amount")
	username := users.GetUser(r.Context()).Username

	return re.Service.BuyProduct(r.Context(), username, products.Product{
		Name:   chi.URLParam(r, "product_name"),
		Amount: atoiWithDefault(amount, 1),
	})
}

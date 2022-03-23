package vendinghandler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/artback/mvp/pkg/authentication"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/vending"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

type restHandler struct {
	Repository vending.Repository
}

func httpError(w http.ResponseWriter, err error) {
	var code int
	switch {
	case errors.Is(err, repository.EmptyErr{}):
		code = http.StatusNotFound
	case errors.Is(err, repository.InvalidErr{}):
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
func (re restHandler) getAccount(r *http.Request) (*vending.AccountResponse, error) {
	username := authentication.FromCtx(r.Context())
	return re.Repository.GetAccount(r.Context(), username)
}

func (re restHandler) ResetDeposit(w http.ResponseWriter, r *http.Request) {
	err := re.resetDeposit(r)
	if err != nil {
		httpError(w, err)
	}
}
func (re restHandler) resetDeposit(r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	username := authentication.FromCtx(ctx)
	return re.Repository.SetDeposit(ctx, username, 0)
}

func (re restHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	err := re.deposit(r)
	if err != nil {
		httpError(w, err)
	}
}
func (re restHandler) deposit(r *http.Request) error {

	deposit := change.Deposit{}
	if err := json.NewDecoder(r.Body).Decode(&deposit); err != nil {
		log.Println(err)
		return err
	}
	username := authentication.FromCtx(r.Context())
	return re.Repository.IncrementDeposit(r.Context(), username, deposit.ToAmount())
}

func (re restHandler) BuyProduct(w http.ResponseWriter, r *http.Request) {
	err := re.buyProduct(r)
	if err != nil {
		httpError(w, err)
	}
}
func atoiWithDefault(string string, def int) int {
	amount, err := strconv.Atoi(string)
	if err != nil {
		return def
	}
	return amount
}
func (re restHandler) buyProduct(r *http.Request) error {
	amount := r.URL.Query().Get("amount")
	username := authentication.FromCtx(r.Context())
	return re.Repository.BuyProduct(r.Context(), username, products.Update{
		Name:   chi.URLParam(r, "product_name"),
		Amount: atoiWithDefault(amount, 1),
	})
}

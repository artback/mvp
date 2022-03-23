package producthandler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/artback/mvp/pkg/authentication"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

var (
	JsonErr = errors.New("error parsing json body")
)

type restHandler struct {
	Repository products.Repository
}

func httpError(w http.ResponseWriter, err error) {
	var code int
	switch {
	case errors.Is(err, repository.EmptyErr{}):
		code = http.StatusNotFound
	case errors.Is(err, repository.DuplicateErr{}):
		code = http.StatusConflict
	case errors.Is(err, JsonErr):
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), code)
}

func (rest restHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	err := rest.createProduct(r)
	if err != nil {
		httpError(w, err)
	}
}
func (rest restHandler) createProduct(r *http.Request) error {
	product := products.Product{}
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return err
	}
	product.SellerId = authentication.FromCtx(r.Context())
	return rest.Repository.Insert(r.Context(), product)
}

func (rest restHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	p, err := rest.getProduct(r)
	if err != nil {
		httpError(w, err)
	}
	if err := json.NewEncoder(w).Encode(&p); err != nil {
		httpError(w, err)
	}
}
func (rest restHandler) getProduct(r *http.Request) (*products.Product, error) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	return rest.Repository.Get(ctx, chi.URLParam(r, "product_name"))
}

func (rest restHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	err := rest.updateProduct(r)
	if err == nil {
		return
	}
	httpError(w, err)
}
func (rest restHandler) updateProduct(r *http.Request) error {
	req := products.Update{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return JsonErr
	}
	req.Name = chi.URLParam(r, "product_name")
	username := authentication.FromCtx(r.Context())
	return rest.Repository.Update(r.Context(), username, req)
}

func (rest restHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	err := rest.deleteProduct(r)
	if err == nil {
		return
	}
	httpError(w, err)
}
func (rest restHandler) deleteProduct(r *http.Request) error {
	username := authentication.FromCtx(r.Context())
	return rest.Repository.Delete(r.Context(), username, chi.URLParam(r, "product_name"))
}

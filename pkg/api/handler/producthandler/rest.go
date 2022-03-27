package producthandler

import (
	"encoding/json"
	"errors"
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var JsonErr = errors.New("error parsing json body")

type restHandler struct {
	products.Service
}

func httpError(w http.ResponseWriter, err error) {
	var code int

	switch {
	case errors.Is(err, repository.EmptyError{}):
		code = http.StatusNotFound
	case errors.Is(err, repository.DuplicateError{}):
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

	product.SellerID = authentication.GetUserName(r.Context())

	return rest.Insert(r.Context(), product)
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
	return rest.Get(r.Context(), chi.URLParam(r, "product_name"))
}

func (rest restHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	err := rest.updateProduct(r)
	if err == nil {
		return
	}

	httpError(w, err)
}

func (rest restHandler) updateProduct(r *http.Request) error {
	req := products.Product{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return JsonErr
	}

	req.Name = chi.URLParam(r, "product_name")
	req.SellerID = authentication.GetUserName(r.Context())

	return rest.Update(r.Context(), req)
}

func (rest restHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	err := rest.deleteProduct(r)
	if err == nil {
		return
	}

	httpError(w, err)
}

func (rest restHandler) deleteProduct(r *http.Request) error {
	username := authentication.GetUserName(r.Context())
	return rest.Delete(r.Context(), username, chi.URLParam(r, "product_name"))
}

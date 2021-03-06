package producthandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/artback/mvp/pkg/api/middleware/security"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/golang/mock/gomock"
)

type ServiceResponse struct {
	times   int
	Product *products.Product
	err     error
}

func TestController_CreateProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
		ServiceResponse
		insert products.Product
		body   []byte
		want   int
	}{
		{
			name:            "successful create",
			body:            []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:          products.Product{Name: "product1", SellerID: "mike"},
			username:        "mike",
			ServiceResponse: ServiceResponse{times: 1},
			want:            http.StatusOK,
		},
		{
			name:            "unsuccessful create, json decode",
			body:            []byte(`{"name: "product1"}`),
			username:        "mike",
			ServiceResponse: ServiceResponse{times: 0},
			want:            http.StatusInternalServerError,
		},
		{
			name:            "unsuccessful create, service error",
			body:            []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:          products.Product{Name: "product1", SellerID: "mike"},
			want:            http.StatusInternalServerError,
			username:        "mike",
			ServiceResponse: ServiceResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:            "unsuccessful create, insert error",
			body:            []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:          products.Product{Name: "product1", SellerID: "mike"},
			want:            http.StatusInternalServerError,
			username:        "mike",
			ServiceResponse: ServiceResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:            "unsuccessful create, Duplicate error",
			body:            []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:          products.Product{Name: "product1", SellerID: "mike"},
			want:            http.StatusConflict,
			username:        "mike",
			ServiceResponse: ServiceResponse{err: repository.DuplicateError{}, times: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			service := mocks.NewProductService(mockCtrl)
			service.EXPECT().Insert(gomock.Any(), tt.insert).Return(tt.ServiceResponse.err).Times(tt.ServiceResponse.times)
			co := RestHandler{Service: service}
			w := httptest.NewRecorder()

			ctx := security.WithUser(context.Background(), security.User{Username: tt.username})
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", bytes.NewReader(tt.body))
			co.CreateProduct(w, req)
			if status := w.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func TestController_GetProduct(t *testing.T) {
	t.Parallel()

	type want struct {
		code int
		body products.Product
	}

	tests := []struct {
		name string
		want
		ServiceResponse
	}{
		{
			name:            "successful get",
			ServiceResponse: ServiceResponse{Product: &products.Product{Name: "product1"}, times: 1},
			want: want{
				code: http.StatusOK,
				body: products.Product{Name: "product1"},
			},
		},
		{
			name:            "unsuccessful get,error service",
			ServiceResponse: ServiceResponse{err: errors.New("something happened"), times: 1},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:            "unsuccessful get,error empty response",
			ServiceResponse: ServiceResponse{err: repository.EmptyError{}, times: 1},
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewProductService(mockCtrl)
			service.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tt.ServiceResponse.Product, tt.ServiceResponse.err).Times(tt.ServiceResponse.times)
			co := RestHandler{Service: service}
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "", nil)
			co.GetProduct(recorder, req)
			if status := recorder.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
			var p products.Product
			_ = json.NewDecoder(recorder.Body).Decode(&p)
			if !reflect.DeepEqual(p, tt.want.body) {
				t.Errorf("handler returned wrong body: got %v want %v", p, tt.want.body)
			}
		})
	}
}

func TestController_UpdateProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      []byte
		urlParams string
		want      int
		Service   ServiceResponse
		update    products.Product
		username  string
	}{
		{
			name:   "successful update",
			body:   []byte(`{"price": 5}`),
			update: products.Product{Price: 5, SellerID: "mike"},
			want:   http.StatusOK,
			Service: ServiceResponse{
				times: 1,
			},
			username: "mike",
		},
		{
			name:   "unsuccessful update, no products error",
			body:   []byte(`{"price": 5}`),
			update: products.Product{Price: 5, SellerID: "mike"},
			want:   http.StatusNotFound,
			Service: ServiceResponse{
				err:   repository.EmptyError{},
				times: 1,
			},
			username: "mike",
		},
		{
			name:     "unsuccessful update, json body is malformed",
			body:     []byte(`{name: "product1"}`),
			Service:  ServiceResponse{times: 0},
			username: "mike",
			want:     http.StatusBadRequest,
		},
		{
			name:     "unsuccessful update, insert error",
			body:     []byte(`{"price": 5}`),
			update:   products.Product{Price: 5, SellerID: "mike"},
			want:     http.StatusInternalServerError,
			username: "mike",
			Service: ServiceResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			service := mocks.NewProductService(mockCtrl)
			service.EXPECT().Update(gomock.Any(), tt.update).Return(tt.Service.err).Times(tt.Service.times)
			co := RestHandler{Service: service}
			w := httptest.NewRecorder()
			ctx := security.WithUser(context.Background(), security.User{Username: tt.username})
			req, _ := http.NewRequestWithContext(ctx, http.MethodPut, "/", bytes.NewReader(tt.body))
			co.UpdateProduct(w, req)
			if status := w.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func TestController_DeleteProduct(t *testing.T) {
	t.Parallel()

	type want struct {
		code int
	}

	tests := []struct {
		name     string
		username string
		want
		Service ServiceResponse
	}{
		{
			name:     "successful delete",
			Service:  ServiceResponse{times: 1},
			username: "mike",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:     "unsuccessful get,error service",
			username: "mike",
			Service: ServiceResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:     "unsuccessful get,error empty response service",
			username: "mike",
			Service: ServiceResponse{
				err:   repository.EmptyError{},
				times: 1,
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewProductService(mockCtrl)
			service.EXPECT().Delete(gomock.Any(), tt.username, gomock.Any()).Return(tt.Service.err).Times(tt.Service.times)

			co := RestHandler{Service: service}
			w := httptest.NewRecorder()

			ctx := security.WithUser(context.Background(), security.User{Username: tt.username})
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			co.DeleteProduct(w, req)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
		})
	}
}

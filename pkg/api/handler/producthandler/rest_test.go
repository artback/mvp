package producthandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type ServiceResponse struct {
	times   int
	Product *products.Product
	err     error
}

func TestController_CreateProduct(t *testing.T) {
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
			insert:          products.Product{Name: "product1", SellerId: "mike"},
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
			insert:          products.Product{Name: "product1", SellerId: "mike"},
			want:            http.StatusInternalServerError,
			username:        "mike",
			ServiceResponse: ServiceResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:            "unsuccessful create, insert error",
			body:            []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:          products.Product{Name: "product1", SellerId: "mike"},
			want:            http.StatusInternalServerError,
			username:        "mike",
			ServiceResponse: ServiceResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:            "unsuccessful create, Duplicate error",
			body:            []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:          products.Product{Name: "product1", SellerId: "mike"},
			want:            http.StatusConflict,
			username:        "mike",
			ServiceResponse: ServiceResponse{err: repository.DuplicateErr{}, times: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewProductService(mockCtrl)
			rep.EXPECT().Insert(gomock.Any(), tt.insert).Return(tt.ServiceResponse.err).Times(tt.ServiceResponse.times)
			co := restHandler{Service: rep}
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.username), http.MethodGet, "/", bytes.NewReader(tt.body))
			co.CreateProduct(w, req)
			if status := w.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func TestController_GetProduct(t *testing.T) {
	type want struct {
		code int
		body products.Product
	}
	type Param struct {
		name string
	}
	tests := []struct {
		name string
		want
		Param
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
			ServiceResponse: ServiceResponse{err: repository.EmptyErr{}, times: 1},
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	for _, tt := range tests {
		rep := mocks.NewProductService(mockCtrl)
		rep.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tt.ServiceResponse.Product, tt.ServiceResponse.err).Times(tt.ServiceResponse.times)
		co := restHandler{Service: rep}
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.Param), nil)
			co.GetProduct(w, req)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
			p := products.Product{}
			_ = json.NewDecoder(w.Body).Decode(&p)
			if !reflect.DeepEqual(p, tt.want.body) {
				t.Errorf("handler returned wrong body: got %v want %v", p, tt.want.body)
			}
		})
	}
}

func TestController_UpdateProduct(t *testing.T) {
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
			update: products.Product{Price: 5, SellerId: "mike"},
			want:   http.StatusOK,
			Service: ServiceResponse{
				times: 1,
			},
			username: "mike",
		},
		{
			name:   "unsuccessful update, no products error",
			body:   []byte(`{"price": 5}`),
			update: products.Product{Price: 5, SellerId: "mike"},
			want:   http.StatusNotFound,
			Service: ServiceResponse{
				err:   repository.EmptyErr{},
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
			update:   products.Product{Price: 5, SellerId: "mike"},
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewProductService(mockCtrl)
			rep.EXPECT().Update(gomock.Any(), tt.update).Return(tt.Service.err).Times(tt.Service.times)
			co := restHandler{Service: rep}
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.username), http.MethodPut, "/", bytes.NewReader(tt.body))
			co.UpdateProduct(w, req)
			if status := w.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func TestController_DeleteProduct(t *testing.T) {
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
				err:   repository.EmptyErr{},
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
		rep := mocks.NewProductService(mockCtrl)
		rep.EXPECT().Delete(gomock.Any(), tt.username, gomock.Any()).Return(tt.Service.err).Times(tt.Service.times)
		co := restHandler{Service: rep}
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.username), http.MethodGet, "/", nil)
			co.DeleteProduct(w, req)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
		})
	}
}

package producthandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/authentication"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type RepositoryResponse struct {
	times   int
	Product *products.Product
	err     error
}

func TestController_CreateProduct(t *testing.T) {
	tests := []struct {
		name     string
		username string
		RepositoryResponse
		insert products.Product
		body   []byte
		want   int
	}{
		{
			name:               "successful create",
			body:               []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:             products.Product{Name: "product1", SellerId: "mike"},
			username:           "mike",
			RepositoryResponse: RepositoryResponse{times: 1},
			want:               http.StatusOK,
		},
		{
			name:               "unsuccessful create, json decode",
			body:               []byte(`{"name: "product1"}`),
			username:           "mike",
			RepositoryResponse: RepositoryResponse{times: 0},
			want:               http.StatusInternalServerError,
		},
		{
			name:               "unsuccessful create, repository error",
			body:               []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:             products.Product{Name: "product1", SellerId: "mike"},
			want:               http.StatusInternalServerError,
			username:           "mike",
			RepositoryResponse: RepositoryResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:               "unsuccessful create, insert error",
			body:               []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:             products.Product{Name: "product1", SellerId: "mike"},
			want:               http.StatusInternalServerError,
			username:           "mike",
			RepositoryResponse: RepositoryResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:               "unsuccessful create, Duplicate error",
			body:               []byte(`{"name": "product1","seller_id": "mike"}`),
			insert:             products.Product{Name: "product1", SellerId: "mike"},
			want:               http.StatusConflict,
			username:           "mike",
			RepositoryResponse: RepositoryResponse{err: repository.DuplicateErr{}, times: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewProductRepository(mockCtrl)
			rep.EXPECT().Insert(gomock.Any(), tt.insert).Return(tt.RepositoryResponse.err).Times(tt.RepositoryResponse.times)
			co := restHandler{Repository: rep}
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
		RepositoryResponse
	}{
		{
			name:               "successful get",
			RepositoryResponse: RepositoryResponse{Product: &products.Product{Name: "product1"}, times: 1},
			want: want{
				code: http.StatusOK,
				body: products.Product{Name: "product1"},
			},
		},
		{
			name:               "unsuccessful get,error repository",
			RepositoryResponse: RepositoryResponse{err: errors.New("something happened"), times: 1},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:               "unsuccessful get,error empty response",
			RepositoryResponse: RepositoryResponse{err: repository.EmptyErr{}, times: 1},
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	for _, tt := range tests {
		rep := mocks.NewProductRepository(mockCtrl)
		rep.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tt.RepositoryResponse.Product, tt.RepositoryResponse.err).Times(tt.RepositoryResponse.times)
		co := restHandler{Repository: rep}
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
		name       string
		body       []byte
		urlParams  string
		want       int
		Repository RepositoryResponse
		update     products.Update
		username   string
	}{
		{
			name:   "successful update",
			body:   []byte(`{"price": 5}`),
			update: products.Update{Price: 5},
			want:   http.StatusOK,
			Repository: RepositoryResponse{
				times: 1,
			},
			username: "mike",
		},
		{
			name:   "unsuccessful update, no products",
			body:   []byte(`{"price": 5}`),
			update: products.Update{Price: 5},
			want:   http.StatusNotFound,
			Repository: RepositoryResponse{
				err:   repository.EmptyErr{},
				times: 1,
			},
			username: "mike",
		},
		{
			name:       "unsuccessful update, json updater",
			body:       []byte(`{"name: "product1"}`),
			Repository: RepositoryResponse{times: 0},
			username:   "mike",
			want:       http.StatusBadRequest,
		},
		{
			name:     "unsuccessful update, insert error",
			body:     []byte(`{"price": 5}`),
			update:   products.Update{Price: 5},
			want:     http.StatusInternalServerError,
			username: "mike",
			Repository: RepositoryResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewProductRepository(mockCtrl)
			rep.EXPECT().Update(gomock.Any(), tt.username, tt.update).Return(tt.Repository.err).Times(tt.Repository.times)
			co := restHandler{Repository: rep}
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
		Repository RepositoryResponse
	}{
		{
			name:       "successful delete",
			Repository: RepositoryResponse{times: 1},
			username:   "mike",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:     "unsuccessful get,error repository",
			username: "mike",
			Repository: RepositoryResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:     "unsuccessful get,error empty response repository",
			username: "mike",
			Repository: RepositoryResponse{
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
		rep := mocks.NewProductRepository(mockCtrl)
		rep.EXPECT().Delete(gomock.Any(), tt.username, gomock.Any()).Return(tt.Repository.err).Times(tt.Repository.times)
		co := restHandler{Repository: rep}
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

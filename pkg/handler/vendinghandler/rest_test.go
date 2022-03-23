package vendinghandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/authentication"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/vending"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type RepositoryResponse struct {
	*vending.AccountResponse
	err   error
	times int
}

func TestController_BuyProduct(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name       string
		BuyProduct RepositoryResponse
		Username   string
		want       want
	}{
		{
			name: "successful",
			BuyProduct: RepositoryResponse{
				times: 1,
			},
			Username: "mike",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "unsuccessful,Empty error repository",
			BuyProduct: RepositoryResponse{
				err:   repository.EmptyErr{},
				times: 1,
			},
			Username: "mike",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name: "unsuccessful,Invalid error repository",
			BuyProduct: RepositoryResponse{
				err:   repository.InvalidErr{},
				times: 1,
			},
			Username: "mike",
			want: want{
				code: http.StatusNotAcceptable,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewVendingRepsitory(mockCtrl)
			rep.EXPECT().BuyProduct(gomock.Any(), tt.Username, gomock.Any()).Return(tt.BuyProduct.err).Times(tt.BuyProduct.times)
			co := restHandler{Repository: rep}
			r, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.Username), http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			co.BuyProduct(w, r)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
		})
	}
}

func TestController_Deposit(t *testing.T) {
	type want struct {
		code    int
		deposit int
	}
	tests := []struct {
		name string
		RepositoryResponse
		Username string
		body     []byte
		want     want
	}{
		{
			name: "successful",
			body: []byte(`{"5": 2,"10": 0, "20": 5, "50": 0, "100": 1}`),
			RepositoryResponse: RepositoryResponse{
				times: 1,
			},
			want: want{
				code:    http.StatusOK,
				deposit: 210,
			},
		},
		{
			name: "unsuccessful, error marshal",
			body: []byte(`{"5": 2,10: 0, "20": 5, "50": 0, "100": 1}`),
			RepositoryResponse: RepositoryResponse{
				times: 0,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "unsuccessful, error repository",
			body: []byte(`{"5": 2,"10": 0, "20": 5, "50": 0, "100": 1}`),
			RepositoryResponse: RepositoryResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
			want: want{
				code:    http.StatusInternalServerError,
				deposit: 210,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			rep := mocks.NewVendingRepsitory(mockCtrl)
			rep.EXPECT().IncrementDeposit(gomock.Any(), tt.Username, tt.want.deposit).Return(tt.err).Times(tt.times)
			co := restHandler{Repository: rep}
			r, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.Username), http.MethodGet, "/", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			co.Deposit(w, r)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
		})
	}
}

func TestController_ResetDeposit(t *testing.T) {
	type want struct {
		code    int
		deposit int
	}
	tests := []struct {
		name string
		RepositoryResponse
		Username string
		body     []byte
		want     want
	}{
		{
			name: "successful",
			RepositoryResponse: RepositoryResponse{
				times: 1,
			},
			want: want{
				code:    http.StatusOK,
				deposit: 0,
			},
		},
		{
			name: "unsuccessful,error repository",
			RepositoryResponse: RepositoryResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
			want: want{
				code:    http.StatusInternalServerError,
				deposit: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			rep := mocks.NewVendingRepsitory(mockCtrl)
			rep.EXPECT().SetDeposit(gomock.Any(), tt.Username, tt.want.deposit).Return(tt.err).Times(tt.times)
			co := restHandler{Repository: rep}
			r, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.Username), http.MethodGet, "/", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			co.ResetDeposit(w, r)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
		})
	}
}

func TestController_GetAccount(t *testing.T) {
	type want struct {
		code int
		body vending.AccountResponse
	}
	tests := []struct {
		name string
		RepositoryResponse
		Username string
		want     want
	}{
		{
			name: "successful",
			RepositoryResponse: RepositoryResponse{
				times: 1,
				AccountResponse: &vending.AccountResponse{
					Deposit:  change.New(coin.Coins{100}, 400),
					Products: []products.Update{{Name: "cheesecake", Amount: 1}},
					Spent:    100,
				},
			},
			Username: "mike",
			want: want{
				code: http.StatusOK,
				body: vending.AccountResponse{
					Deposit:  change.New(coin.Coins{50, 100}, 400),
					Products: []products.Update{{Name: "cheesecake", Amount: 1}},
					Spent:    100,
				},
			},
		},
		{
			name: "unsuccessful,error repository",
			RepositoryResponse: RepositoryResponse{
				times: 1,
				err:   errors.New("something happened"),
			},
			Username: "mike",
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "unsuccessful,empty response",
			RepositoryResponse: RepositoryResponse{
				times: 1,
				err:   repository.EmptyErr{},
			},
			Username: "mike",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewVendingRepsitory(mockCtrl)
			rep.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(tt.AccountResponse, tt.err).Times(tt.times)
			co := restHandler{Repository: rep}
			r, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.Username), http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			co.GetAccount(w, r)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
			a := vending.AccountResponse{}
			_ = json.NewDecoder(w.Body).Decode(&a)
			if !reflect.DeepEqual(a, tt.want.body) {
				t.Errorf("handler returned wrong body: got %v want %v",
					a, tt.want.body)
			}
		})
	}
}

func Test_toAmount(t *testing.T) {
	type args struct {
		query    string
		defaults int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "default amount",
			args: args{
				query:    "",
				defaults: 1,
			},
			want: 1,
		},
		{
			name: "non default amount",
			args: args{
				query:    "10",
				defaults: 1,
			},
			want: 10,
		},
		{
			name: "non number amount,return default",
			args: args{
				query:    "hello",
				defaults: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := atoiWithDefault(tt.args.query, tt.args.defaults); got != tt.want {
				t.Errorf("toAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

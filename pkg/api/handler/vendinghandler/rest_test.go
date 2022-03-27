package vendinghandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/api/middleware/authentication"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/vending"
	"github.com/golang/mock/gomock"
)

type ServiceResponse struct {
	*vending.Response
	err   error
	times int
}

func TestController_BuyProduct(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name       string
		BuyProduct ServiceResponse
		Username   string
		want       want
	}{
		{
			name: "successful request",
			BuyProduct: ServiceResponse{
				times: 1,
			},
			Username: "mike",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "unsuccessful,Empty error repository",
			BuyProduct: ServiceResponse{
				err:   repository.EmptyError{},
				times: 1,
			},
			Username: "mike",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name: "unsuccessful,Invalid error repository",
			BuyProduct: ServiceResponse{
				err:   repository.InvalidError{},
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
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := mocks.NewVendingService(mockCtrl)
			s.EXPECT().BuyProduct(gomock.Any(), tt.Username, gomock.Any()).Return(tt.BuyProduct.err).Times(tt.BuyProduct.times)
			co := restHandler{Service: s}
			r, _ := http.NewRequestWithContext(authentication.WithUsername(context.Background(), tt.Username), http.MethodGet, "/", nil)
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
	t.Parallel()

	type want struct {
		code    int
		deposit int
	}

	tests := []struct {
		name string
		ServiceResponse
		Username string
		body     []byte
		want     want
	}{
		{
			name: "successful",
			body: []byte(`{"5": 2,"10": 0, "20": 5, "50": 0, "100": 1}`),
			ServiceResponse: ServiceResponse{
				times: 1,
			},
			want: want{
				code:    http.StatusOK,
				deposit: 210,
			},
		},
		{
			name: "unsuccessful, error json marshal",
			body: []byte(`{"5": 2,10: 0, "20": 5, "50": 0, "100": 1}`),
			ServiceResponse: ServiceResponse{
				times: 0,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "unsuccessful, error repository",
			body: []byte(`{"5": 2,"10": 0, "20": 5, "50": 0, "100": 1}`),
			ServiceResponse: ServiceResponse{
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
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			s := mocks.NewVendingService(mockCtrl)
			s.EXPECT().IncrementDeposit(gomock.Any(), tt.Username, tt.want.deposit).Return(tt.err).Times(tt.times)
			co := restHandler{Service: s}
			r, _ := http.NewRequestWithContext(authentication.WithUsername(context.Background(), tt.Username), http.MethodGet, "/", bytes.NewReader(tt.body))
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
	t.Parallel()

	type want struct {
		code    int
		deposit int
	}

	tests := []struct {
		name string
		ServiceResponse
		Username string
		body     []byte
		want     want
	}{
		{
			name: "successful",
			ServiceResponse: ServiceResponse{
				times: 1,
			},
			want: want{
				code:    http.StatusOK,
				deposit: 0,
			},
		},
		{
			name: "unsuccessful,error repository",
			ServiceResponse: ServiceResponse{
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
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			s := mocks.NewVendingService(mockCtrl)
			s.EXPECT().SetDeposit(gomock.Any(), tt.Username, tt.want.deposit).Return(tt.err).Times(tt.times)
			co := restHandler{Service: s}
			r, _ := http.NewRequestWithContext(authentication.WithUsername(context.Background(), tt.Username), http.MethodGet, "/", bytes.NewReader(tt.body))
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
	t.Parallel()

	type want struct {
		code int
		body vending.Response
	}

	tests := []struct {
		name string
		ServiceResponse
		Username string
		want     want
	}{
		{
			name: "successful",
			ServiceResponse: ServiceResponse{
				times: 1,
				Response: &vending.Response{
					Deposit:  change.New(coin.Coins{100}, 400),
					Products: []products.Product{{Name: "cheesecake", Amount: 1}},
					Spent:    100,
				},
			},
			Username: "mike",
			want: want{
				code: http.StatusOK,
				body: vending.Response{
					Deposit:  change.New(coin.Coins{50, 100}, 400),
					Products: []products.Product{{Name: "cheesecake", Amount: 1}},
					Spent:    100,
				},
			},
		},
		{
			name: "unsuccessful,error repository",
			ServiceResponse: ServiceResponse{
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
			ServiceResponse: ServiceResponse{
				times: 1,
				err:   repository.EmptyError{},
			},
			Username: "mike",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := mocks.NewVendingService(mockCtrl)
			s.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(tt.Response, tt.err).Times(tt.times)
			co := restHandler{Service: s}
			r, _ := http.NewRequestWithContext(authentication.WithUsername(context.Background(), tt.Username), http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			co.GetAccount(w, r)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
			a := vending.Response{}
			_ = json.NewDecoder(w.Body).Decode(&a)
			if !reflect.DeepEqual(a, tt.want.body) {
				t.Errorf("handler returned wrong body: got %v want %v",
					a, tt.want.body)
			}
		})
	}
}

func Test_toAmount(t *testing.T) {
	t.Parallel()

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
			t.Parallel()
			if got := atoiWithDefault(tt.args.query, tt.args.defaults); got != tt.want {
				t.Errorf("toAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

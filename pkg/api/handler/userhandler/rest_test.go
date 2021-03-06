package userhandler

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
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"github.com/golang/mock/gomock"
)

type serviceResponse struct {
	times    int
	response users.Response
	err      error
}

func TestController_GetUser(t *testing.T) {
	t.Parallel()

	type want struct {
		code int
		body users.Response
	}

	tests := []struct {
		name string
		want
		Service serviceResponse
	}{
		{
			name: "successful get",
			Service: serviceResponse{
				response: users.Response{Username: "user_1"},
				times:    1,
			},
			want: want{
				code: http.StatusOK,
				body: users.Response{Username: "user_1"},
			},
		},
		{
			name: "unsuccessful get,error repository",
			Service: serviceResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "unsuccessful get,error empty response repository",
			Service: serviceResponse{
				err:   repository.EmptyError{},
				times: 1,
			},
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
			service := mocks.NewUserService(mockCtrl)
			service.EXPECT().GetResponse(gomock.Any(), gomock.Any()).Return(&tt.Service.response, tt.Service.err).Times(tt.Service.times)
			co := RestHandler{service}
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			recorder := httptest.NewRecorder()
			co.GetUser(recorder, req)
			if status := recorder.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
			var body users.Response
			_ = json.Unmarshal(recorder.Body.Bytes(), &body)
			if !reflect.DeepEqual(body, tt.want.body) {
				t.Errorf("handler returned wrong body: got %v want %v",
					body, tt.want.body)
			}
		})
	}
}

func TestController_UpdateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		body     []byte
		want     int
		username string
		user     users.User
		Service  serviceResponse
	}{
		{
			name:     "successful update",
			body:     []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			want:     200,
			username: "mike",
			user:     users.User{Username: "mike", Password: "password", Role: security.Buyer, Deposit: 100},
			Service:  serviceResponse{times: 1},
		},
		{
			name:    "unsuccessful update, json decoder",
			body:    []byte(`{"username: "user1"}`),
			Service: serviceResponse{times: 0},
			want:    http.StatusBadRequest,
		},
		{
			name:     "unsuccessful update, insert error",
			body:     []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			username: "mike",
			user:     users.User{Username: "mike", Password: "password", Role: security.Buyer, Deposit: 100},
			want:     http.StatusInternalServerError,
			Service:  serviceResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:     "unsuccessful update ,error empty response repository",
			body:     []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			username: "mike",
			user:     users.User{Username: "mike", Password: "password", Role: security.Buyer, Deposit: 100},
			Service:  serviceResponse{err: repository.EmptyError{}, times: 1},
			want:     http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			service := mocks.NewUserService(mockCtrl)
			co := RestHandler{service}
			service.EXPECT().Update(gomock.Any(), tt.user).Return(tt.Service.err).Times(tt.Service.times)
			w := httptest.NewRecorder()
			ctx := security.WithUser(context.Background(), security.User{Username: tt.username})
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", bytes.NewReader(tt.body))
			co.UpdateUser(w, req)
			if status := w.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func TestController_DeleteUser(t *testing.T) {
	t.Parallel()

	type want struct {
		code int
	}

	tests := []struct {
		name string
		want
		username string
		Service  serviceResponse
	}{
		{
			name:     "successful delete",
			Service:  serviceResponse{times: 1, err: nil},
			username: "mike",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:     "unsuccessful get,error repository",
			Service:  serviceResponse{err: errors.New("something happened"), times: 1},
			username: "mike",
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:     "unsuccessful get,error empty response repository",
			Service:  serviceResponse{err: repository.EmptyError{}, times: 1},
			username: "mike",
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
			service := mocks.NewUserService(mockCtrl)
			co := RestHandler{service}
			service.EXPECT().Delete(gomock.Any(), tt.username).Return(tt.Service.err).Times(tt.Service.times)
			w := httptest.NewRecorder()
			ctx := security.WithUser(context.Background(), security.User{Username: tt.username})
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			co.DeleteUser(w, req)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
		})
	}
}

func TestController_CreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    []byte
		want    int
		Service serviceResponse
	}{
		{
			name:    "successful decode",
			body:    []byte(`{"username": "user1","password":"password","role": "buyer"}`),
			want:    http.StatusOK,
			Service: serviceResponse{times: 1},
		},
		{
			name:    "unsuccessful decode, json decoder",
			body:    []byte(`{"username: "user1","role": "buyer","deposit": 100}`),
			Service: serviceResponse{times: 0},
			want:    http.StatusBadRequest,
		},
		{
			name:    "unsuccessful decode, insert error",
			body:    []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			want:    http.StatusInternalServerError,
			Service: serviceResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:    "unsuccessful decode, duplicate error",
			body:    []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			want:    http.StatusConflict,
			Service: serviceResponse{err: repository.DuplicateError{}, times: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			service := mocks.NewUserService(mockCtrl)
			co := RestHandler{service}
			service.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(tt.Service.err).Times(tt.Service.times)
			req, _ := http.NewRequest(http.MethodGet, "/", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			co.CreateUser(w, req)
			if status := w.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

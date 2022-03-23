package userhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/authentication"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type RepositoryResponse struct {
	times    int
	response users.Response
	err      error
}

func TestController_GetUser(t *testing.T) {
	type want struct {
		code int
		body users.Response
	}
	tests := []struct {
		name string
		want
		Repository RepositoryResponse
	}{
		{
			name: "successful get",
			Repository: RepositoryResponse{
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
			Repository: RepositoryResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "unsuccessful get,error empty response repository",
			Repository: RepositoryResponse{
				err:   repository.EmptyErr{},
				times: 1,
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewUserRepository(mockCtrl)
			rep.EXPECT().GetResponse(gomock.Any(), gomock.Any()).Return(&tt.Repository.response, tt.Repository.err).Times(tt.Repository.times)
			co := restHandler{Repository: rep}
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			co.GetUser(w, req)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
			var body users.Response
			_ = json.Unmarshal(w.Body.Bytes(), &body)
			if !reflect.DeepEqual(body, tt.want.body) {
				t.Errorf("handler returned wrong body: got %v want %v",
					body, tt.want.body)
			}
		})
	}
}

func TestController_UpdateUser(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		want       int
		username   string
		user       users.User
		Repository RepositoryResponse
	}{
		{
			name:       "successful update",
			body:       []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			want:       200,
			username:   "mike",
			user:       users.User{Username: "mike", Password: "password", Role: users.Buyer},
			Repository: RepositoryResponse{times: 1},
		},
		{
			name:       "unsuccessful update, json decoder",
			body:       []byte(`{"username: "user1"}`),
			Repository: RepositoryResponse{times: 0},
			want:       http.StatusBadRequest,
		},
		{
			name:       "unsuccessful update, insert error",
			body:       []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			username:   "mike",
			user:       users.User{Username: "mike", Password: "password", Role: users.Buyer},
			want:       http.StatusInternalServerError,
			Repository: RepositoryResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:       "unsuccessful update ,error empty response repository",
			body:       []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			username:   "mike",
			user:       users.User{Username: "mike", Password: "password", Role: users.Buyer},
			Repository: RepositoryResponse{err: repository.EmptyErr{}, times: 1},
			want:       http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewUserRepository(mockCtrl)
			co := restHandler{Repository: rep}
			rep.EXPECT().Update(gomock.Any(), tt.user).Return(tt.Repository.err).Times(tt.Repository.times)
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.username), http.MethodGet, "/", bytes.NewReader(tt.body))
			co.UpdateUser(w, req)
			if status := w.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func TestController_DeleteUser(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		want
		username   string
		Repository RepositoryResponse
	}{
		{
			name:       "successful delete",
			Repository: RepositoryResponse{times: 1, err: nil},
			username:   "mike",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:       "unsuccessful get,error repository",
			Repository: RepositoryResponse{err: errors.New("something happened"), times: 1},
			username:   "mike",
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:       "unsuccessful get,error empty response repository",
			Repository: RepositoryResponse{err: repository.EmptyErr{}, times: 1},
			username:   "mike",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewUserRepository(mockCtrl)
			co := restHandler{Repository: rep}
			rep.EXPECT().Delete(gomock.Any(), tt.username).Return(tt.Repository.err).Times(tt.Repository.times)
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(authentication.CtxWithUsername(context.Background(), tt.username), http.MethodGet, "/", nil)
			co.DeleteUser(w, req)
			if status := w.Code; status != tt.want.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.code)
			}
		})
	}
}

func TestController_CreateUser(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		want       int
		Repository RepositoryResponse
	}{
		{
			name:       "successful decode",
			body:       []byte(`{"username": "user1","password":"password","role": "buyer"}`),
			want:       http.StatusOK,
			Repository: RepositoryResponse{times: 1},
		},
		{
			name:       "unsuccessful decode, json decoder",
			body:       []byte(`{"username: "user1","role": "buyer","deposit": 100}`),
			Repository: RepositoryResponse{times: 0},
			want:       http.StatusBadRequest,
		},
		{
			name:       "unsuccessful decode, insert error",
			body:       []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			want:       http.StatusInternalServerError,
			Repository: RepositoryResponse{err: errors.New("something happened"), times: 1},
		},
		{
			name:       "unsuccessful decode, duplicate error",
			body:       []byte(`{"username": "user1","password":"password","role": "buyer","deposit": 100}`),
			want:       http.StatusConflict,
			Repository: RepositoryResponse{err: repository.DuplicateErr{}, times: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewUserRepository(mockCtrl)
			co := restHandler{Repository: rep}
			rep.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(tt.Repository.err).Times(tt.Repository.times)
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

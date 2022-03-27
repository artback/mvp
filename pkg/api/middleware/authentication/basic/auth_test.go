package basic

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"github.com/golang/mock/gomock"
)

type ServiceResponse struct {
	times int
	user  *users.User
	err   error
}

func emptySuccessResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuth_Authorize(t *testing.T) {
	t.Parallel()

	type auth struct {
		password string
		username string
	}

	tests := []struct {
		name string
		*auth
		roles []users.Role
		ServiceResponse
		want int
	}{
		{
			name: "successful authorization",
			auth: &auth{
				username: "mike",
				password: "password",
			},
			ServiceResponse: ServiceResponse{
				user:  &users.User{Username: "mike", Password: "password", Role: users.Seller},
				times: 1,
			},
			want: http.StatusOK,
		},
		{
			name: "successful authorization,with roles",
			auth: &auth{
				username: "mike",
				password: "password",
			},
			roles: []users.Role{users.Seller},
			ServiceResponse: ServiceResponse{
				user:  &users.User{Username: "mike", Password: "password", Role: users.Seller},
				times: 1,
			},
			want: http.StatusOK,
		},
		{
			name: "successful authorization",
			auth: &auth{
				username: "mike",
				password: "password",
			},
			ServiceResponse: ServiceResponse{
				user:  &users.User{Username: "mike", Password: "password"},
				times: 1,
			},
			want: http.StatusOK,
		},
		{
			name: "unsuccessful authorization,with roles",
			auth: &auth{
				username: "mike",
				password: "password",
			},
			roles: []users.Role{users.Seller},
			ServiceResponse: ServiceResponse{
				user:  &users.User{Username: "mike", Password: "password", Role: users.Buyer},
				times: 1,
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "unsuccessful authorization wrong password",
			auth: &auth{
				username: "mike",
				password: "password",
			},
			ServiceResponse: ServiceResponse{
				user:  &users.User{Username: "mike", Password: "pass"},
				times: 1,
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "unsuccessful authorization error service",
			ServiceResponse: ServiceResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
			auth: &auth{
				username: "mike",
				password: "password",
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "unsuccessful authorization empty response",
			ServiceResponse: ServiceResponse{
				err:   repository.EmptyError{},
				times: 1,
			},
			auth: &auth{
				username: "mike",
				password: "password",
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "unsuccessful authorization missing header",
			ServiceResponse: ServiceResponse{
				times: 0,
			},
			want: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := mocks.NewUserService(mockCtrl)
			s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tt.ServiceResponse.user, tt.ServiceResponse.err).Times(tt.times)
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			if tt.auth != nil {
				req.SetBasicAuth(tt.username, tt.password)
			}
			co := Auth{s}
			co.Authenticate(tt.roles...)(http.HandlerFunc(emptySuccessResponse)).ServeHTTP(w, req)
			if status := w.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

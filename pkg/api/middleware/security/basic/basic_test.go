package basic

import (
	"errors"
	"github.com/artback/mvp/pkg/api/middleware/security"
	"github.com/artback/mvp/pkg/pass"
	"github.com/artback/mvp/pkg/repository"
	"net/http"
	"reflect"
	"testing"

	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/users"
	"github.com/golang/mock/gomock"
)

type ServiceResponse struct {
	times int
	user  *users.User
	err   error
}

func TestBasic_GetUser(t *testing.T) {
	t.Parallel()
	type fields struct {
		Service users.Service
	}
	type args struct {
		User *users.User
	}

	tests := []struct {
		name            string
		fields          fields
		ServiceResponse ServiceResponse
		args            args
		want            *security.User
		wantErr         bool
	}{
		{
			name: "successful authorization",
			args: args{
				User: &users.User{Username: "mike", Password: "password", Role: security.Seller},
			},
			ServiceResponse: ServiceResponse{
				user:  &users.User{Username: "mike", Password: "password", Role: security.Seller},
				times: 1,
			},
			want: &security.User{Username: "mike", Role: security.Seller},
		},
		{
			name: "unsuccessful authorization wrong password",
			args: args{
				User: &users.User{Username: "mike", Password: "password", Role: security.Seller},
			},
			ServiceResponse: ServiceResponse{
				user:  &users.User{Username: "mike", Password: "pass"},
				times: 1,
			},
			wantErr: true,
		},
		{
			name: "unsuccessful authorization error service",
			ServiceResponse: ServiceResponse{
				err:   errors.New("something happened"),
				times: 1,
			},
			args: args{
				User: &users.User{Username: "mike", Password: "password", Role: security.Seller},
			},
			wantErr: true,
		},
		{
			name: "unsuccessful authorization empty response",
			ServiceResponse: ServiceResponse{
				err:   repository.EmptyError{},
				times: 1,
			},
			args: args{
				User: &users.User{Username: "mike", Password: "password", Role: security.Seller},
			},
			wantErr: true,
		},
		{
			name: "unsuccessful authorization missing header",
			ServiceResponse: ServiceResponse{
				times: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := mocks.NewUserService(mockCtrl)
			if tt.ServiceResponse.user != nil {
				tt.ServiceResponse.user.Password, _ = pass.HashAndSalt(tt.ServiceResponse.user.Password)
			}
			s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tt.ServiceResponse.user, tt.ServiceResponse.err).Times(tt.ServiceResponse.times)
			b := Basic{
				Service: s,
			}
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			if tt.args.User != nil {
				req.SetBasicAuth(tt.args.User.Username, tt.args.User.Password)
			}
			got, err := b.GetUser(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
	"github.com/golang/mock/gomock"
)

type userResponse struct {
	*users.User
	err   error
	times int
}

func TestUserService_GetResponse(t *testing.T) {
	t.Parallel()

	type fields struct {
		Coins coin.Coins
	}

	type args struct {
		ctx      context.Context
		username string
	}

	tests := []struct {
		name       string
		fields     fields
		Repository userResponse
		args       args
		want       *users.Response
		wantErr    bool
	}{
		{
			name:   "successful get",
			fields: fields{Coins: coin.Coins{5, 10, 20, 100}},
			Repository: userResponse{
				User:  &users.User{Username: "user_1"},
				times: 1,
			},
			want: &users.Response{Username: "user_1", Deposit: map[coin.Coin]int{}},
		},
		{
			name:   "successful get with deposit",
			fields: fields{Coins: coin.Coins{5, 10, 20, 100}},
			Repository: userResponse{
				User:  &users.User{Username: "user_1", Deposit: 100},
				times: 1,
			},
			want: &users.Response{Username: "user_1", Deposit: map[coin.Coin]int{100: 1}},
		},
		{
			name:   "error repository get with deposit",
			fields: fields{Coins: coin.Coins{5, 10, 20, 100}},
			Repository: userResponse{
				err:   repository.EmptyError{},
				times: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rep := mocks.NewUserRepository(mockCtrl)
			rep.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tt.Repository.User, tt.Repository.err).Times(tt.Repository.times)
			u := UserService{
				Coins:      tt.fields.Coins,
				Repository: rep,
			}

			got, err := u.GetResponse(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResponse() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

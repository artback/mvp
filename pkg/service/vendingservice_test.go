package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/artback/mvp/mocks"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/vending"
	"github.com/golang/mock/gomock"
)

func TestVendingService_GetAccount(t *testing.T) {
	t.Parallel()

	type fields struct {
		Coins coin.Coins
	}

	type args struct {
		username string
	}

	type mockArg struct {
		times   int
		account vending.Account
		err     error
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		mockArg mockArg
		want    *vending.Response
		wantErr bool
	}{
		{
			name: "successful",
			fields: fields{
				Coins: coin.Coins{5, 10, 20, 50, 100},
			},
			mockArg: mockArg{
				times: 1,
				account: vending.Account{
					Deposit:  400,
					Products: []products.Product{{Name: "cheesecake", Amount: 1}},
					Spent:    100,
				},
			},
			args: args{username: "mike"},
			want: &vending.Response{
				Deposit:  change.New(coin.Coins{5, 10, 20, 50, 100}, 400),
				Products: []products.Product{{Name: "cheesecake", Amount: 1}},
				Spent:    100,
			},
		},
		{
			name: "unsuccessful,error repository",
			mockArg: mockArg{
				times: 1,
				err:   errors.New("something happened"),
			},
			args:    args{username: "mike"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			r := mocks.NewVendingRepsitory(mockCtrl)
			r.EXPECT().GetAccount(gomock.Any(), tt.args.username).Return(&tt.mockArg.account, tt.mockArg.err).Times(tt.mockArg.times)
			v := VendingService{
				Repository: r,
				Coins:      tt.fields.Coins,
			}
			got, err := v.GetAccount(context.Background(), tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccount() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

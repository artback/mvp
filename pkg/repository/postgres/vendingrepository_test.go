package postgres_test

import (
	"context"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository/postgres"
	"github.com/artback/mvp/pkg/vending"
	"log"
	"math"
	"reflect"
	"testing"
)

func vendingReposity(fn ...func(r vending.Repository)) vending.Repository {
	productRepo := postgres.ProductRepository{DB: db}
	if err := productRepo.Insert(context.Background(), defaultProduct); err != nil {
		if err := productRepo.Update(context.Background(), defaultProduct); err != nil {
			log.Fatal("vendingRepository: ", err)
		}
	}
	repo := postgres.VendingRepository{DB: db}
	_, err := repo.Exec("DELETE FROM transactions")
	if err != nil {
		log.Fatal("vendingRepository ", err)
	}
	for _, f := range fn {
		f(repo)
	}
	return repo
}
func TestVendingRepository_BuyProduct(t *testing.T) {
	type args struct {
		username string
		product  products.Product
	}
	tests := []struct {
		name    string
		args    args
		setup   func(r vending.Repository)
		wantErr bool
	}{
		{
			name: "buy existing product with existing user without deposit",
			args: args{username: defaultBuyer.Username, product: defaultProduct},
			setup: func(r vending.Repository) {
				err := r.SetDeposit(context.Background(), defaultBuyer.Username, 0)
				if err != nil {
					t.Error(err)
				}
			},
			wantErr: true,
		},
		{
			name: "buy existing product with existing user with deposit",
			args: args{username: defaultBuyer.Username, product: products.Product{
				Name: defaultProduct.Name, Amount: 5,
			}},
			setup: func(r vending.Repository) {
				err := r.SetDeposit(context.Background(), defaultBuyer.Username, 100)
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "buy non existing product with existing user with deposit",
			args: args{username: defaultBuyer.Username, product: products.Product{Name: "durex"}},
			setup: func(r vending.Repository) {
				err := r.SetDeposit(context.Background(), defaultBuyer.Username, 100)
				if err != nil {
					t.Error(err)
				}
			},
			wantErr: true,
		},
		{
			name: "buy existing product with non existing user ",
			args: args{username: "non existing", product: defaultProduct},
			setup: func(r vending.Repository) {
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if err := vendingReposity(tt.setup).BuyProduct(ctx, tt.args.username, tt.args.product); (err != nil) != tt.wantErr {
				t.Errorf("BuyProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVendingRepository_GetAccount(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(r vending.Repository)
		want    *vending.Account
		wantErr bool
	}{
		{
			name: "get existing account",
			args: args{
				username: defaultBuyer.Username,
			},
			setup: func(r vending.Repository) {
				err := r.SetDeposit(context.Background(), defaultBuyer.Username, 100)
				if err != nil {
					t.Error(err)
				}
			},
			want: &vending.Account{
				Deposit:  100,
				Products: []products.Product{},
				Spent:    0,
			},
		},
		{
			name: "get existing account with transactions",
			args: args{
				username: defaultBuyer.Username,
			},
			setup: func(r vending.Repository) {
				ctx := context.Background()
				if err := r.SetDeposit(ctx, defaultBuyer.Username, 100); err != nil {
					t.Error(err)
				}
				if err := r.BuyProduct(ctx, defaultBuyer.Username, products.Product{
					Name: defaultProduct.Name, Amount: 5,
				}); err != nil {
					t.Error(err)
				}
				if err := r.BuyProduct(ctx, defaultBuyer.Username, products.Product{
					Name: defaultProduct.Name, Amount: 5,
				}); err != nil {
					t.Error(err)
				}
			},
			want: &vending.Account{
				Deposit:  50,
				Products: []products.Product{{Name: defaultProduct.Name, Price: defaultProduct.Price, Amount: 10}},
				Spent:    50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vendingReposity(tt.setup).GetAccount(context.Background(), tt.args.username)
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

func TestVendingRepository_IncrementDeposit(t *testing.T) {
	type args struct {
		username string
		deposit  int
	}
	tests := []struct {
		name    string
		args    args
		setup   func(r vending.Repository)
		want    int
		wantErr bool
	}{
		{
			name: "increment existing user",
			args: args{
				username: defaultBuyer.Username,
				deposit:  100,
			},
			setup: func(r vending.Repository) {
				err := r.SetDeposit(context.Background(), defaultBuyer.Username, 0)
				if err != nil {
					t.Error(err)
				}
			},
			want: 100,
		},
		{
			name: "decrement existing user",
			setup: func(r vending.Repository) {
				err := r.SetDeposit(context.Background(), defaultBuyer.Username, 0)
				if err != nil {
					t.Error(err)
				}
			},
			args: args{
				username: defaultBuyer.Username,
				deposit:  -100,
			},
			want: -100,
		},
		{
			name:  "increment non existing user",
			setup: func(r vending.Repository) {},
			args: args{
				username: "non existing",
				deposit:  100,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := vendingReposity(tt.setup)
			if err := repo.IncrementDeposit(context.Background(), tt.args.username, tt.args.deposit); (err != nil) != tt.wantErr {
				t.Errorf("IncrementDeposit() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := repo.GetAccount(context.Background(), tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.Deposit, tt.want) {
				t.Errorf("GetAccount().Deposit got = %v, want %v", got.Deposit, tt.want)
			}
		})
	}
}

func TestVendingRepository_SetDeposit(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
		deposit  int
	}
	tests := []struct {
		name    string
		args    args
		setup   func(r vending.Repository)
		want    int
		wantErr bool
	}{
		{
			name: "set for existing user",
			args: args{
				username: defaultBuyer.Username,
				deposit:  500,
			},
			want: 500,
			setup: func(r vending.Repository) {
				err := r.SetDeposit(context.Background(), defaultBuyer.Username, 0)
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "set for existing user,out of boundary",
			args: args{
				username: defaultBuyer.Username,
				deposit:  math.MaxInt,
			},
			wantErr: true,
			setup: func(r vending.Repository) {
				err := r.SetDeposit(context.Background(), defaultBuyer.Username, 0)
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:  "set for non existing user",
			setup: func(r vending.Repository) {},
			args: args{
				username: "non existing",
				deposit:  100,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := vendingReposity(tt.setup)
			if err := repo.SetDeposit(context.Background(), tt.args.username, tt.args.deposit); (err != nil) != tt.wantErr {
				t.Errorf("SetDeposit() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				got, err := repo.GetAccount(context.Background(), tt.args.username)
				if err != nil {
					t.Errorf("GetAccount() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Deposit, tt.want) {
					t.Errorf("GetAccount().Deposit got = %v, want %v", got.Deposit, tt.want)
				}
			}
		})
	}
}

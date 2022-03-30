//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository/postgres"
	"reflect"
	"testing"
)

func productReposity(fn ...func(r products.Repository)) products.Repository {
	repo := postgres.ProductRepository{DB: db}
	if err := repo.Insert(context.Background(), defaultProduct); err != nil {
		repo.Update(context.Background(), defaultProduct)
	}

	for _, f := range fn {
		f(repo)
	}
	return repo
}
func TestProductRepository_Delete(t *testing.T) {
	type args struct {
		username string
		name     string
	}

	tests := []struct {
		name    string
		setup   func(r products.Repository)
		args    args
		wantErr bool
	}{
		{
			name: "delete existing product owned by user",
			args: args{name: "toy", username: "defaultSeller"},
			setup: func(r products.Repository) {
				r.Insert(context.Background(), products.Product{Name: "toy", SellerID: defaultSeller.Username, Price: 100, Amount: 1})
			},
		},
		{
			name: "delete existing product not owned by user",
			args: args{name: "toy", username: "defaultBuyer"},
			setup: func(r products.Repository) {
				r.Insert(context.Background(), products.Product{Name: "toy", SellerID: defaultSeller.Username, Price: 100, Amount: 1})
			},
			wantErr: true,
		},
		{
			name: "delete non existing product",
			args: args{name: "nonExisting", username: "defaultSeller"},
			setup: func(r products.Repository) {
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := productReposity(tt.setup).Delete(context.Background(), tt.args.username, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductRepository_Get(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name    string
		setup   func(r products.Repository)
		args    args
		want    *products.Product
		wantErr bool
	}{
		{
			name: "get existing product",
			args: args{name: defaultProduct.Name},
			setup: func(r products.Repository) {
				r.Insert(context.Background(), defaultProduct)
			},
			want: &defaultProduct,
		},
		{
			name:    "get non existing product",
			args:    args{name: "fanta"},
			setup:   func(r products.Repository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := productReposity(tt.setup).Get(context.Background(), tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductRepository_Insert(t *testing.T) {
	type args struct {
		product products.Product
	}

	tests := []struct {
		name    string
		setup   func(r products.Repository)
		args    args
		wantErr bool
	}{
		{
			name: "insert non existing product",
			setup: func(r products.Repository) {
			},
			args: args{
				product: products.Product{Name: "heineken", SellerID: defaultSeller.Username, Price: 5, Amount: 1},
			},
		},
		{
			name: "insert existing product",
			setup: func(r products.Repository) {
				r.Insert(context.Background(), defaultProduct)
			},
			args: args{
				product: defaultProduct,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := productReposity(tt.setup).Insert(context.Background(), tt.args.product); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductRepository_Update(t *testing.T) {
	type args struct {
		username string
		product  products.Product
	}

	tests := []struct {
		name    string
		setup   func(r products.Repository)
		args    args
		wantErr bool
	}{
		{
			name: "update existing product",
			setup: func(r products.Repository) {
				r.Insert(context.Background(), defaultProduct)
			},
			args: args{
				username: defaultSeller.Username,
				product:  defaultProduct,
			},
		},
		{
			name: "update existing product not owned by user",
			setup: func(r products.Repository) {
				r.Insert(context.Background(), defaultProduct)
			},
			args: args{
				username: defaultBuyer.Username,
				product:  products.Product{Name: defaultProduct.Name, SellerID: "non default", Price: 10, Amount: 2},
			},
			wantErr: true,
		},
		{
			name:  "update non existing product",
			setup: func(r products.Repository) {},
			args: args{
				username: defaultSeller.Username,
				product:  products.Product{Name: "Eluxadolin", SellerID: defaultSeller.Username, Price: 5, Amount: 1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := productReposity(tt.setup).Update(context.Background(), tt.args.product); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

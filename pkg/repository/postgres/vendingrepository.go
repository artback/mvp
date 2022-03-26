package postgres

import (
	"context"
	"database/sql"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/vending"
)

type VendingRepository struct {
	*sql.DB
}

func (v VendingRepository) GetAccount(ctx context.Context, username string) (*vending.Account, error) {
	response, err := v.getAccount(ctx, username)
	return response, DomainError(err)
}
func (v VendingRepository) getAccount(ctx context.Context, username string) (*vending.Account, error) {
	tx, err := v.BeginTx(ctx, nil)
	rows, err := tx.QueryContext(ctx,
		`SELECT price,SUM(amount),product_name
    FROM transactions inner join users on transactions.username = users.username 
    where users.username = $1 group by product_name,price;`,
		username)
	if err != nil {
		return nil, err
	}
	var (
		total int
		// to prevent empty slice to be null in json
		productRequest = make([]products.Product, 0)
	)
	for rows.Next() {
		var product products.Product
		if err := rows.Scan(&product.Price, &product.Amount, &product.Name); err != nil {
			return nil, err
		}
		total += product.Price * product.Amount
		productRequest = append(productRequest, product)
	}

	var deposit int
	err = tx.QueryRowContext(ctx,
		`SELECT deposit FROM users WHERE username = $1`, username).Scan(&deposit)
	if err != nil {
		return nil, err
	}

	return &vending.Account{
		Deposit:  deposit,
		Products: productRequest,
		Spent:    total,
	}, nil
}

func (v VendingRepository) IncrementDeposit(ctx context.Context, username string, deposit int) error {
	return DomainError(v.incrementDeposit(ctx, username, deposit))
}
func (v VendingRepository) incrementDeposit(ctx context.Context, username string, deposit int) error {
	result, err := v.ExecContext(ctx, `UPDATE users set deposit=deposit+$1 where username = $2`, deposit, username)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return repository.EmptyErr{}
	}
	return err
}

func (v VendingRepository) BuyProduct(ctx context.Context, username string, product products.Product) error {
	_, err := v.ExecContext(ctx, `INSERT INTO transactions(product_name, username, amount) VALUES ($1,$2,$3)`, product.Name, username, product.Amount)
	return DomainError(err)
}

func (v VendingRepository) SetDeposit(ctx context.Context, username string, deposit int) error {
	return DomainError(v.setDeposit(ctx, username, deposit))
}
func (v VendingRepository) setDeposit(ctx context.Context, username string, deposit int) error {
	result, err := v.ExecContext(ctx, `UPDATE users set deposit=$1 where username = $2`, deposit, username)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return repository.EmptyErr{}
	}
	return err
}

package postgres

import (
	"context"
	"database/sql"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/vending"
)

type VendingRepository struct {
	coin.Coins
	*sql.DB
}

func (v VendingRepository) GetAccount(ctx context.Context, username string) (*vending.AccountResponse, error) {
	tx, err := v.BeginTx(ctx, nil)
	defer func() {
		err = DomainError(err)
	}()
	rows, err := tx.QueryContext(ctx,
		`SELECT COALESCE(SUM(amount * price), 0) as spent,COALESCE(SUM(amount),0),product_name
    FROM transactions inner join users on transactions.username = users.username 
    where users.username = $1 group by product_name;`,
		username)
	if err != nil {
		return nil, err
	}
	var (
		total          int
		deposit        int
		productRequest []products.Update
	)
	// to prevent empty slice to be null in json
	productRequest = make([]products.Update, 0)
	for rows.Next() {
		var (
			spent   int
			product products.Update
		)
		err := rows.Scan(&spent, &product.Amount, &product.Name)
		if err != nil {
			return nil, DomainError(err)
		}
		total += spent
		productRequest = append(productRequest, product)
	}
	err = tx.QueryRowContext(ctx,
		`SELECT deposit FROM users WHERE username = $1`, username).Scan(&deposit)
	if err != nil {
		return nil, err
	}
	return &vending.AccountResponse{
		Deposit:  change.New(v.Coins, deposit),
		Products: productRequest,
		Spent:    total,
	}, nil
}

func (v VendingRepository) IncrementDeposit(ctx context.Context, username string, deposit int) error {
	_, err := v.ExecContext(ctx, `UPDATE users set deposit = deposit + $1 where username = $2`, deposit, username)
	return DomainError(err)
}

func (v VendingRepository) BuyProduct(ctx context.Context, username string, product products.Update) error {
	_, err := v.ExecContext(ctx, `INSERT INTO transactions(product_name, username, amount) VALUES ($1,$2, $3)`, product.Name, username, product.Amount)
	return DomainError(err)
}

func (v VendingRepository) SetDeposit(ctx context.Context, username string, deposit int) error {
	_, err := v.ExecContext(ctx, `UPDATE users set deposit = $1 where username = $2`, deposit, username)
	return DomainError(err)
}

package postgres

import (
	"context"
	"database/sql"

	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository"
)

type ProductRepository struct {
	*sql.DB
}

func (p ProductRepository) Insert(ctx context.Context, product products.Product) error {
	return DomainError(p.insert(ctx, product))
}

func (p ProductRepository) insert(ctx context.Context, product products.Product) error {
	tx, err := p.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}

		_ = tx.Commit()
	}()

	if _, err = tx.ExecContext(ctx,
		`INSERT INTO products(name,seller_id) VALUES ($1,$2)`, product.Name, product.SellerID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx,
		`INSERT INTO inventory(product_name,amount,price) VALUES ($1,$2,$3)`, product.Name, product.Amount, product.Price); err != nil {
		return err
	}

	return nil
}

func (p ProductRepository) Update(ctx context.Context, product products.Product) error {
	return DomainError(p.update(ctx, product))
}

func (p ProductRepository) update(ctx context.Context, product products.Product) error {
	result, err := p.ExecContext(ctx,
		`UPDATE inventory as i SET amount = $1 ,price = $2  FROM products as p,users as u where i.product_name = p.name and p.name = $3 and u.username = $4`,
		product.Amount, product.Price, product.Name, product.SellerID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if affected == 0 {
		return repository.EmptyError{}
	}

	return err
}

func (p ProductRepository) Get(ctx context.Context, name string) (*products.Product, error) {
	pr, err := p.get(ctx, name)

	return pr, DomainError(err)
}

func (p ProductRepository) get(ctx context.Context, name string) (*products.Product, error) {
	product := products.Product{}
	if err := p.QueryRowContext(ctx, `
		SELECT name,seller_id,price,amount from products 
		    INNER JOIN inventory i on products.name = i.product_name where name= $1`, name).Scan(&product.Name, &product.SellerID, &product.Price, &product.Amount); err != nil {
		return nil, err
	}

	return &product, nil
}

func (p ProductRepository) Delete(ctx context.Context, username string, name string) error {
	return DomainError(p.delete(ctx, username, name))
}

func (p ProductRepository) delete(ctx context.Context, username string, name string) error {
	result, err := p.ExecContext(ctx, `DELETE FROM products p USING users u  where p.seller_id = u.username and u.username = $1 and p.name =$2 and u.role = 'seller'`, username, name)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if affected == 0 {
		return repository.EmptyError{}
	}

	return err
}

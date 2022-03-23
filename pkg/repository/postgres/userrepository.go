package postgres

import (
	"context"
	"database/sql"
	"github.com/artback/mvp/pkg/change"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
)

type UserRepository struct {
	coin.Coins
	*sql.DB
}

func (u UserRepository) Get(ctx context.Context, username string) (*users.User, error) {
	var (
		password string
		role     users.Role
	)
	err := u.QueryRowContext(ctx,
		`SELECT password, role FROM users where username = $1`,
		username).Scan(&password, &role)
	if err != nil {
		return nil, err
	}
	return &users.User{
		Username: username,
		Password: password,
		Role:     role,
	}, nil
}

func (u UserRepository) GetResponse(ctx context.Context, username string) (*users.Response, error) {
	var err error
	defer func() {
		err = DomainError(err)
	}()
	var (
		role    users.Role
		deposit int
	)
	err = u.QueryRowContext(ctx,
		`SELECT deposit,role FROM users WHERE username = $1`, username).Scan(&deposit, &role)
	if err != nil {
		return nil, err
	}
	return &users.Response{
		Deposit:  change.New(u.Coins, deposit),
		Username: username,
		Role:     role,
	}, nil
}
func (u UserRepository) Insert(ctx context.Context, user users.User) error {
	var err error
	defer func() {
		err = DomainError(err)
	}()
	_, err = u.ExecContext(ctx,
		`INSERT INTO users(username,password,role) VALUES ($1,$2,$3)`,
		user.Username, user.Password, user.Role,
	)
	return err
}

func (u UserRepository) Update(ctx context.Context, user users.User) error {
	var err error
	defer func() {
		err = DomainError(err)
	}()
	result, err := u.ExecContext(ctx,
		`update users set password = COALESCE(NULLIF($1,""),password), role = COALESCE(NULLIF($2,""),role) where username = $3`,
		user.Password, user.Role, user.Username,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return repository.EmptyErr{}
	}
	return err
}
func (u UserRepository) Delete(ctx context.Context, username string) error {
	var err error
	defer func() {
		err = DomainError(err)
	}()
	result, err := u.ExecContext(ctx,
		`delete FROM users where username = $1`,
		username,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return repository.EmptyErr{}
	}
	return err
}

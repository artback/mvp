package postgres

import (
	"context"
	"database/sql"

	"github.com/artback/mvp/pkg/repository"
	"github.com/artback/mvp/pkg/users"
)

type UserRepository struct {
	*sql.DB
}

func (u UserRepository) Get(ctx context.Context, username string) (*users.User, error) {
	user, err := u.get(ctx, username)

	return user, DomainError(err)
}

func (u UserRepository) get(ctx context.Context, username string) (*users.User, error) {
	var (
		password string
		role     users.Role
		deposit  int
	)

	if err := u.QueryRowContext(ctx, `SELECT password, role,deposit FROM users where username = $1`, username).Scan(&password, &role, &deposit); err != nil {
		return nil, err
	}

	return &users.User{
		Username: username,
		Password: password,
		Role:     role,
		Deposit:  deposit,
	}, nil
}

func (u UserRepository) Insert(ctx context.Context, user users.User) error {
	return DomainError(u.insert(ctx, user))
}

func (u UserRepository) insert(ctx context.Context, user users.User) error {
	_, err := u.ExecContext(ctx,
		`INSERT INTO users(username,password,role) VALUES ($1,$2,$3)`,
		user.Username, user.Password, user.Role,
	)

	return err
}

func (u UserRepository) Update(ctx context.Context, user users.User) error {
	return DomainError(u.update(ctx, user))
}

func (u UserRepository) update(ctx context.Context, user users.User) error {
	result, err := u.ExecContext(ctx,
		`update users set password = COALESCE(NULLIF($1,''),password), role = COALESCE(NULLIF($2,''),role) where username = $3`,
		user.Password, user.Role, user.Username,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if affected == 0 {
		return repository.EmptyError{}
	}

	return err
}

func (u UserRepository) Delete(ctx context.Context, username string) error {
	return DomainError(u.delete(ctx, username))
}

func (u UserRepository) delete(ctx context.Context, username string) error {
	result, err := u.ExecContext(ctx,
		`delete FROM users where username = $1`,
		username,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if affected == 0 {
		return repository.EmptyError{}
	}

	return err
}

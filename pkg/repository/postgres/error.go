package postgres

import (
	"database/sql"
	"errors"
	"github.com/artback/mvp/pkg/repository"
	"github.com/lib/pq"
)

// DomainError translates lower sql errors to domain errors
func DomainError(err error) error {
	if err == nil {
		return nil
	}
	pqErr, ok := err.(*pq.Error)
	if errors.Is(err, sql.ErrNoRows) {
		return repository.EmptyErr{}
	}
	if !ok {
		return err
	}
	switch pqErr.Constraint {
	case "users_pkey", "products_pkey":
		return repository.DuplicateErr{Err: pqErr}
	case "fk_product_name":
		return repository.EmptyErr{}
	default:
		return repository.InvalidErr{Title: pqErr.Message}
	}
}

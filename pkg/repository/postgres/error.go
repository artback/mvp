package postgres

import (
	"database/sql"
	"errors"

	"github.com/artback/mvp/pkg/repository"
	"github.com/lib/pq"
)

// DomainError translates lower sql errors to domain errors.
func DomainError(err error) error {
	pqErr, ok := err.(*pq.Error)

	if errors.Is(err, sql.ErrNoRows) {
		return repository.EmptyError{}
	}

	if !ok {
		return err
	}

	switch pqErr.Code {
	case "23505", "23503":
		return repository.DuplicateError{Err: pqErr, Constraint: pqErr.Constraint}
	default:
		return repository.InvalidError{Title: pqErr.Message}
	}
}

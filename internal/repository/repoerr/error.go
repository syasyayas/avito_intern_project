package repoerr

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = errors.New("not found")

var ErrAlreadyExists = errors.New("already exists")

var ErrInvalidExpiresAt = errors.New("invalid expires_at value")

func PgErrorWrapper(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			err = ErrAlreadyExists
		case "42703":
			err = ErrNotFound
		default:
		}
	}
	return err
}

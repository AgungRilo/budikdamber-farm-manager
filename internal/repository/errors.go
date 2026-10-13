package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDuplicate = errors.New("data duplikat")
	ErrInvalid   = errors.New("data melanggar aturan")
)

// translatePgError mengubah kode error PostgreSQL jadi error yang dipahami service/handler.
func translatePgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrDuplicate
		case "23514": // check_violation
			return ErrInvalid
		}
	}
	return err
}

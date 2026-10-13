package repository

import (
	"errors"
	"fmt"

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

var (
	ErrPondUnavailable = errors.New("wadah tidak ditemukan atau tidak aktif")
	ErrSpeciesNotFound = errors.New("jenis ikan tidak ditemukan")
)

// CapacityError dikembalikan saat tebar melebihi kapasitas wadah (R01).
type CapacityError struct {
	MaxStock  int   `json:"max_stock"`
	Used      int64 `json:"used"`
	Remaining int64 `json:"remaining"`
}

func (e *CapacityError) Error() string {
	return fmt.Sprintf("kapasitas wadah tidak cukup: maksimal %d ekor, terisi %d, sisa %d",
		e.MaxStock, e.Used, e.Remaining)
}

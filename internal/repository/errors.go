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

var (
	ErrCycleNotFound  = errors.New("siklus tidak ditemukan")
	ErrCycleClosed    = errors.New("siklus sudah ditutup, data tidak bisa diubah")
	ErrFeedNotFound   = errors.New("pakan tidak ditemukan")
	ErrLogBeforeStart = errors.New("tanggal log tidak boleh sebelum tanggal tebar")
)

var (
	ErrCycleNotClosed = errors.New("siklus belum ditutup")
	ErrEndBeforeStart = errors.New("tanggal selesai tidak boleh sebelum tanggal tebar")
)

// AliveExceededError: jumlah mati/jual melebihi ikan hidup (R02).
type AliveExceededError struct {
	Requested int   `json:"requested"`
	Alive     int64 `json:"alive"`
}

func (e *AliveExceededError) Error() string {
	return fmt.Sprintf("jumlah melebihi ikan hidup: diminta %d, sisa hidup %d ekor", e.Requested, e.Alive)
}

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

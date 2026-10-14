package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// lockActiveCycle mengunci baris siklus (FOR UPDATE) dan memastikan statusnya active.
// Dipakai di dalam transaksi sebelum menambah/menghapus log, supaya tidak bentrok
// dengan proses lain (misalnya tutup siklus) yang berjalan bersamaan.
func lockActiveCycle(ctx context.Context, tx pgx.Tx, cycleID int64) (time.Time, error) {
	var status string
	var start time.Time
	err := tx.QueryRow(ctx,
		`SELECT status::text, start_date FROM cycles WHERE id = $1 FOR UPDATE`, cycleID).
		Scan(&status, &start)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrCycleNotFound
	}
	if err != nil {
		return time.Time{}, err
	}
	if status != "active" {
		return time.Time{}, ErrCycleClosed
	}
	return start, nil
}

func cycleExists(ctx context.Context, db *pgxpool.Pool, cycleID int64) error {
	var ok bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM cycles WHERE id = $1)`, cycleID).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return ErrCycleNotFound
	}
	return nil
}

// nilIfEmpty: string kosong → NULL, supaya filter tanggal opsional.
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// closeCycle menutup siklus dan mencatat siapa yang menutup (R05).
// Wajib dipanggil di dalam transaksi yang sudah mengunci baris siklus.
func closeCycle(ctx context.Context, tx pgx.Tx, cycleID int64, endDate string, closedBy int64) error {
	_, err := tx.Exec(ctx,
		`UPDATE cycles SET status = 'closed', end_date = $2::date, closed_by = $3 WHERE id = $1`,
		cycleID, endDate, closedBy)
	return err
}

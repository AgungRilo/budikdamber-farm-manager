package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

const mortalityLogSelect = `
	SELECT ml.id, ml.cycle_id, to_char(ml.log_date, 'YYYY-MM-DD'), ml.qty, ml.cause, ml.notes,
	       ml.created_by, u.name, ml.created_at
	FROM mortality_logs ml
	JOIN users u ON u.id = ml.created_by`

type MortalityLogRepository struct {
	db *pgxpool.Pool
}

func NewMortalityLogRepository(db *pgxpool.Pool) *MortalityLogRepository {
	return &MortalityLogRepository{db: db}
}

func scanMortalityLog(row pgx.Row) (model.MortalityLog, error) {
	var l model.MortalityLog
	err := row.Scan(&l.ID, &l.CycleID, &l.LogDate, &l.Qty, &l.Cause, &l.Notes,
		&l.CreatedBy, &l.CreatedByName, &l.CreatedAt)
	return l, err
}

func (r *MortalityLogRepository) List(ctx context.Context, f model.LogFilter) ([]model.MortalityLog, int64, error) {
	if err := cycleExists(ctx, r.db, f.CycleID); err != nil {
		return nil, 0, err
	}
	const where = ` WHERE ml.cycle_id = $1
		AND ($2::date IS NULL OR ml.log_date >= $2::date)
		AND ($3::date IS NULL OR ml.log_date <= $3::date)`
	from, to := nilIfEmpty(f.DateFrom), nilIfEmpty(f.DateTo)

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM mortality_logs ml`+where, f.CycleID, from, to).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, mortalityLogSelect+where+` ORDER BY ml.log_date DESC, ml.id DESC LIMIT $4 OFFSET $5`,
		f.CycleID, from, to, f.Limit, (f.Page-1)*f.Limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]model.MortalityLog, 0)
	for rows.Next() {
		l, err := scanMortalityLog(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, l)
	}
	return list, total, rows.Err()
}

func (r *MortalityLogRepository) FindByID(ctx context.Context, id int64) (*model.MortalityLog, error) {
	l, err := scanMortalityLog(r.db.QueryRow(ctx, mortalityLogSelect+` WHERE ml.id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *MortalityLogRepository) Create(ctx context.Context, cycleID int64, in model.MortalityLogInput, createdBy int64) (*model.MortalityCreateResult, error) {
	// Mulai transaksi: insert log dan tutup siklus harus berhasil bersama atau gagal bersama.
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// Kalau fungsi keluar sebelum Commit (ada error), semua perubahan dibatalkan.
	// Setelah Commit, Rollback ini tidak berpengaruh apa-apa.
	defer tx.Rollback(ctx)

	// R02: kunci baris siklus dulu. Request lain untuk siklus yang sama
	// harus menunggu sampai transaksi ini selesai, jadi hitungan "sisa hidup" aman.
	// Sekalian memastikan siklus ada dan masih active, lalu mengambil tanggal tebar.
	start, err := lockActiveCycle(ctx, tx, cycleID)
	if err != nil {
		return nil, err
	}
	// Format tanggal sudah divalidasi di layer handler, jadi error parse diabaikan.
	logDate, _ := time.Parse("2006-01-02", in.LogDate)
	// Ikan tidak mungkin mati sebelum ditebar.
	if logDate.Before(start) {
		return nil, ErrLogBeforeStart
	}

	// Ambil jumlah ikan yang masih hidup dari view (tebar - mati - terjual).
	// Dibaca setelah lock, jadi angkanya tidak bisa berubah oleh request lain.
	var alive int64
	if err := tx.QueryRow(ctx,
		`SELECT alive_qty FROM v_cycle_status WHERE cycle_id = $1`, cycleID).Scan(&alive); err != nil {
		return nil, err
	}
	// Jumlah mati tidak boleh melebihi sisa hidup, supaya stok tidak minus.
	if int64(in.Qty) > alive {
		return nil, &AliveExceededError{Requested: in.Qty, Alive: alive}
	}

	// Simpan log kematian; nilai dikirim lewat placeholder ($1..$6), bukan disambung ke SQL.
	var id int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO mortality_logs (cycle_id, log_date, qty, cause, notes, created_by)
		VALUES ($1, $2::date, $3, $4, $5, $6)
		RETURNING id`,
		cycleID, in.LogDate, in.Qty, in.Cause, in.Notes, createdBy).Scan(&id); err != nil {
		// Ubah error PostgreSQL (misalnya pelanggaran constraint) jadi error yang dipahami service.
		return nil, translatePgError(err)
	}

	// R05: semua ikan habis → siklus otomatis ditutup.
	// Sisa hidup setelah log ini, dikembalikan juga ke client di response.
	remaining := alive - int64(in.Qty)
	// Penanda untuk response: apakah request ini ikut menutup siklus.
	closed := false
	if remaining == 0 {
		// if _, err := tx.Exec(ctx,
		// 	`UPDATE cycles SET status = 'closed', end_date = $2::date WHERE id = $1`,
		// 	cycleID, in.LogDate); err != nil {
		// 	return nil, err
		// }
		// closed = true
		if remaining == 0 {
			// Tutup siklus di transaksi yang sama, dengan end_date = tanggal log
			// dan closed_by = user yang mencatat.
			if err := closeCycle(ctx, tx, cycleID, in.LogDate, createdBy); err != nil {
				return nil, err
			}
			closed = true
		}
	}

	// Simpan permanen semua perubahan dan lepas lock siklus.
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	// Baca ulang log lengkap (termasuk nama pencatat dari JOIN users) untuk response.
	l, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Kembalikan log beserta sisa hidup dan status apakah siklus ikut ditutup.
	return &model.MortalityCreateResult{MortalityLog: *l, AliveQtyAfter: remaining, CycleClosed: closed}, nil
}

func (r *MortalityLogRepository) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var cycleID int64
	err = tx.QueryRow(ctx, `SELECT cycle_id FROM mortality_logs WHERE id = $1`, id).Scan(&cycleID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if _, err := lockActiveCycle(ctx, tx, cycleID); err != nil { // R04
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM mortality_logs WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

const feedLogSelect = `
	SELECT fl.id, fl.cycle_id, fl.feed_id, f.name, to_char(fl.log_date, 'YYYY-MM-DD'),
	       fl.qty_gram, fl.price_per_kg, fl.cost, fl.notes, fl.created_by, u.name, fl.created_at
	FROM feed_logs fl
	JOIN feeds f ON f.id = fl.feed_id
	JOIN users u ON u.id = fl.created_by`

type FeedLogRepository struct {
	db *pgxpool.Pool
}

func NewFeedLogRepository(db *pgxpool.Pool) *FeedLogRepository {
	return &FeedLogRepository{db: db}
}

func scanFeedLog(row pgx.Row) (model.FeedLog, error) {
	var l model.FeedLog
	err := row.Scan(&l.ID, &l.CycleID, &l.FeedID, &l.FeedName, &l.LogDate,
		&l.QtyGram, &l.PricePerKg, &l.Cost, &l.Notes, &l.CreatedBy, &l.CreatedByName, &l.CreatedAt)
	return l, err
}

func (r *FeedLogRepository) List(ctx context.Context, f model.LogFilter) ([]model.FeedLog, int64, error) {
	if err := cycleExists(ctx, r.db, f.CycleID); err != nil {
		return nil, 0, err
	}
	const where = ` WHERE fl.cycle_id = $1
		AND ($2::date IS NULL OR fl.log_date >= $2::date)
		AND ($3::date IS NULL OR fl.log_date <= $3::date)`
	from, to := nilIfEmpty(f.DateFrom), nilIfEmpty(f.DateTo)

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM feed_logs fl`+where, f.CycleID, from, to).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, feedLogSelect+where+` ORDER BY fl.log_date DESC, fl.id DESC LIMIT $4 OFFSET $5`,
		f.CycleID, from, to, f.Limit, (f.Page-1)*f.Limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]model.FeedLog, 0)
	for rows.Next() {
		l, err := scanFeedLog(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, l)
	}
	return list, total, rows.Err()
}

func (r *FeedLogRepository) FindByID(ctx context.Context, id int64) (*model.FeedLog, error) {
	l, err := scanFeedLog(r.db.QueryRow(ctx, feedLogSelect+` WHERE fl.id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *FeedLogRepository) Create(ctx context.Context, cycleID int64, in model.FeedLogInput, createdBy int64) (*model.FeedLog, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	start, err := lockActiveCycle(ctx, tx, cycleID)
	if err != nil {
		return nil, err
	}
	logDate, _ := time.Parse("2006-01-02", in.LogDate) // format sudah divalidasi binding
	if logDate.Before(start) {
		return nil, ErrLogBeforeStart
	}

	// R03: harga pakan DISALIN saat ini juga, biaya dihitung di DB.
	// Kalau pakan tidak ada / sudah dihapus, SELECT tidak menghasilkan baris → ErrNoRows.
	var id int64
	err = tx.QueryRow(ctx, `
		INSERT INTO feed_logs (cycle_id, feed_id, log_date, qty_gram, price_per_kg, cost, notes, created_by)
		SELECT $1::bigint, f.id, $3::date, $4::int, f.price_per_kg,
		       ROUND($4::int * f.price_per_kg / 1000.0)::bigint, $5::text, $6::bigint
		FROM feeds f
		WHERE f.id = $2 AND f.deleted_at IS NULL
		RETURNING id`,
		cycleID, in.FeedID, in.LogDate, in.QtyGram, in.Notes, createdBy).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFeedNotFound
	}
	if err != nil {
		return nil, translatePgError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *FeedLogRepository) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var cycleID int64
	err = tx.QueryRow(ctx, `SELECT cycle_id FROM feed_logs WHERE id = $1`, id).Scan(&cycleID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if _, err := lockActiveCycle(ctx, tx, cycleID); err != nil { // R04
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM feed_logs WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

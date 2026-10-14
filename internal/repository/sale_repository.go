package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

const saleSelect = `
	SELECT s.id, s.cycle_id, to_char(s.sale_date, 'YYYY-MM-DD'), s.sale_type::text, s.channel::text,
	       s.qty, s.weight_kg::float8, s.price_per_fish, s.total, s.buyer_name, s.is_paid, s.notes,
	       s.created_by, u.name, s.created_at, s.updated_at
	FROM sales s
	JOIN users u ON u.id = s.created_by`

type SaleRepository struct {
	db *pgxpool.Pool
}

func NewSaleRepository(db *pgxpool.Pool) *SaleRepository {
	return &SaleRepository{db: db}
}

func scanSale(row pgx.Row) (model.Sale, error) {
	var s model.Sale
	err := row.Scan(&s.ID, &s.CycleID, &s.SaleDate, &s.SaleType, &s.Channel,
		&s.Qty, &s.WeightKg, &s.PricePerFish, &s.Total, &s.BuyerName, &s.IsPaid, &s.Notes,
		&s.CreatedBy, &s.CreatedByName, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *SaleRepository) List(ctx context.Context, f model.SaleFilter) ([]model.Sale, int64, error) {
	if err := cycleExists(ctx, r.db, f.CycleID); err != nil {
		return nil, 0, err
	}
	const where = ` WHERE s.cycle_id = $1
		AND ($2::date IS NULL OR s.sale_date >= $2::date)
		AND ($3::date IS NULL OR s.sale_date <= $3::date)
		AND ($4::text IS NULL OR s.sale_type::text = $4::text)
		AND ($5::bool IS NULL OR s.is_paid = $5::bool)`
	args := []any{f.CycleID, nilIfEmpty(f.DateFrom), nilIfEmpty(f.DateTo), nilIfEmpty(f.SaleType), f.IsPaid}

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM sales s`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, (f.Page-1)*f.Limit)
	rows, err := r.db.Query(ctx, saleSelect+where+` ORDER BY s.sale_date DESC, s.id DESC LIMIT $6 OFFSET $7`, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]model.Sale, 0)
	for rows.Next() {
		s, err := scanSale(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}
	return list, total, rows.Err()
}

func (r *SaleRepository) FindByID(ctx context.Context, id int64) (*model.Sale, error) {
	s, err := scanSale(r.db.QueryRow(ctx, saleSelect+` WHERE s.id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SaleRepository) Create(ctx context.Context, cycleID int64, in model.SaleInput, createdBy int64) (*model.SaleCreateResult, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	start, err := lockActiveCycle(ctx, tx, cycleID) // R02: kunci siklus
	if err != nil {
		return nil, err
	}
	saleDate, _ := time.Parse("2006-01-02", in.SaleDate)
	if saleDate.Before(start) {
		return nil, ErrLogBeforeStart
	}

	var alive int64
	if err := tx.QueryRow(ctx, `SELECT alive_qty FROM v_cycle_status WHERE cycle_id = $1`, cycleID).Scan(&alive); err != nil {
		return nil, err
	}
	if int64(in.Qty) > alive {
		return nil, &AliveExceededError{Requested: in.Qty, Alive: alive}
	}

	// Harga: own_use selalu 0 (R07); sale tanpa harga → default dari jenis ikan.
	var price int64
	switch {
	case in.SaleType == "own_use":
		price = 0
	case in.PricePerFish != nil:
		price = *in.PricePerFish
	default:
		if err := tx.QueryRow(ctx, `
			SELECT fs.default_price_per_fish
			FROM cycles c JOIN fish_species fs ON fs.id = c.species_id
			WHERE c.id = $1`, cycleID).Scan(&price); err != nil {
			return nil, err
		}
	}
	total := int64(in.Qty) * price // R06: total dihitung server

	var id int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO sales (cycle_id, sale_date, sale_type, channel, qty, weight_kg,
		                   price_per_fish, total, buyer_name, is_paid, notes, created_by)
		VALUES ($1, $2::date, $3::sale_type, $4::sale_channel, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`,
		cycleID, in.SaleDate, in.SaleType, in.Channel, in.Qty, in.WeightKg,
		price, total, in.BuyerName, *in.IsPaid, in.Notes, createdBy).Scan(&id); err != nil {
		return nil, translatePgError(err)
	}

	remaining := alive - int64(in.Qty)
	closed := false
	if remaining == 0 { // R05: ikan habis terjual → siklus selesai
		if err := closeCycle(ctx, tx, cycleID, in.SaleDate, createdBy); err != nil {
			return nil, err
		}
		closed = true
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	s, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.SaleCreateResult{Sale: *s, AliveQtyAfter: remaining, CycleClosed: closed}, nil
}

// SetPaid sengaja TIDAK mengecek status siklus: pelunasan boleh setelah panen.
func (r *SaleRepository) SetPaid(ctx context.Context, id int64, isPaid bool) (*model.Sale, error) {
	tag, err := r.db.Exec(ctx, `UPDATE sales SET is_paid = $2 WHERE id = $1`, id, isPaid)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *SaleRepository) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var cycleID int64
	err = tx.QueryRow(ctx, `SELECT cycle_id FROM sales WHERE id = $1`, id).Scan(&cycleID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if _, err := lockActiveCycle(ctx, tx, cycleID); err != nil { // R04
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM sales WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

const cycleSelect = `
SELECT c.id, c.code, c.pond_id, p.name, c.species_id, f.name,
       to_char(c.start_date, 'YYYY-MM-DD'), to_char(c.end_date, 'YYYY-MM-DD'),
       (COALESCE(c.end_date, CURRENT_DATE) - c.start_date),
       s.initial_qty, s.dead_qty, s.sold_qty, s.alive_qty,
       c.seed_size_cm::float8, c.seed_price_per_fish, c.status::text, c.notes,
       c.created_by, c.created_at, c.updated_at
FROM cycles c
JOIN ponds p          ON p.id = c.pond_id
JOIN fish_species f   ON f.id = c.species_id
JOIN v_cycle_status s ON s.cycle_id = c.id`

type CycleRepository struct {
	db *pgxpool.Pool
}

func NewCycleRepository(db *pgxpool.Pool) *CycleRepository {
	return &CycleRepository{db: db}
}

func scanCycle(row pgx.Row) (model.Cycle, error) {
	var c model.Cycle
	err := row.Scan(&c.ID, &c.Code, &c.PondID, &c.PondName, &c.SpeciesID, &c.SpeciesName,
		&c.StartDate, &c.EndDate, &c.AgeDays,
		&c.InitialQty, &c.DeadQty, &c.SoldQty, &c.AliveQty,
		&c.SeedSizeCm, &c.SeedPricePerFish, &c.Status, &c.Notes,
		&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CycleRepository) List(ctx context.Context, f model.CycleFilter) ([]model.Cycle, int64, error) {
	where := []string{"TRUE"}
	args := []any{}
	add := func(cond string, v any) {
		args = append(args, v)
		where = append(where, fmt.Sprintf(cond, len(args)))
	}
	if f.Status != "" {
		add("c.status::text = $%d", f.Status)
	}
	if f.PondID > 0 {
		add("c.pond_id = $%d", f.PondID)
	}
	if f.SpeciesID > 0 {
		add("c.species_id = $%d", f.SpeciesID)
	}
	whereSQL := strings.Join(where, " AND ")

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM cycles c WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, (f.Page-1)*f.Limit)
	query := fmt.Sprintf("%s WHERE %s ORDER BY c.start_date DESC, c.id DESC LIMIT $%d OFFSET $%d",
		cycleSelect, whereSQL, len(args)-1, len(args))
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]model.Cycle, 0)
	for rows.Next() {
		c, err := scanCycle(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, c)
	}
	return list, total, rows.Err()
}

func (r *CycleRepository) FindByID(ctx context.Context, id int64) (*model.Cycle, error) {
	c, err := scanCycle(r.db.QueryRow(ctx, cycleSelect+" WHERE c.id = $1", id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Create membuat siklus dalam satu transaksi, termasuk cek kapasitas wadah (R01).
func (r *CycleRepository) Create(ctx context.Context, in model.CycleInput, createdBy int64) (*model.Cycle, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // otomatis batal kalau belum Commit

	// Kunci baris wadah: dua request bersamaan ke wadah yang sama harus antre,
	// supaya keduanya tidak sama-sama lolos cek kapasitas.
	var maxStock int
	var isActive bool
	err = tx.QueryRow(ctx,
		`SELECT max_stock, is_active FROM ponds WHERE id = $1 AND deleted_at IS NULL FOR UPDATE`,
		in.PondID).Scan(&maxStock, &isActive)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && !isActive) {
		return nil, ErrPondUnavailable
	}
	if err != nil {
		return nil, err
	}

	var speciesExists bool
	if err := tx.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM fish_species WHERE id = $1 AND deleted_at IS NULL)`,
		in.SpeciesID).Scan(&speciesExists); err != nil {
		return nil, err
	}
	if !speciesExists {
		return nil, ErrSpeciesNotFound
	}

	var used int64
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(SUM(s.alive_qty), 0)
		FROM cycles c JOIN v_cycle_status s ON s.cycle_id = c.id
		WHERE c.pond_id = $1 AND c.status = 'active'`, in.PondID).Scan(&used); err != nil {
		return nil, err
	}
	if used+int64(in.InitialQty) > int64(maxStock) {
		return nil, &CapacityError{MaxStock: maxStock, Used: used, Remaining: int64(maxStock) - used}
	}

	// Ambil id lebih dulu supaya kode siklus bisa memuat id-nya.
	var id int64
	if err := tx.QueryRow(ctx, `SELECT nextval('cycles_id_seq')`).Scan(&id); err != nil {
		return nil, err
	}

	// "2026-10-13" -> "261013"
	code := fmt.Sprintf("SK-%s-%d", strings.ReplaceAll(in.StartDate, "-", "")[2:], id)
	_, err = tx.Exec(ctx, `
	INSERT INTO cycles (id, code, pond_id, species_id, start_date, initial_qty,
	                    seed_size_cm, seed_price_per_fish, notes, created_by)
	VALUES ($1, $2, $3, $4, $5::date, $6, $7, $8, $9, $10)`,
		id, code, in.PondID, in.SpeciesID, in.StartDate, in.InitialQty,
		in.SeedSizeCm, in.SeedPricePerFish, in.Notes, createdBy)
	if err != nil {
		return nil, translatePgError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

// Close: tutup siklus manual (owner).
func (r *CycleRepository) Close(ctx context.Context, id int64, endDate string, closedBy int64) (*model.Cycle, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	start, err := lockActiveCycle(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	end, _ := time.Parse("2006-01-02", endDate)
	if end.Before(start) {
		return nil, ErrEndBeforeStart
	}
	if err := closeCycle(ctx, tx, id, endDate, closedBy); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

// Reopen: buka kembali siklus yang sudah ditutup, dengan cek ulang kapasitas wadah.
func (r *CycleRepository) Reopen(ctx context.Context, id int64) (*model.Cycle, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var status string
	var pondID int64
	err = tx.QueryRow(ctx, `SELECT status::text, pond_id FROM cycles WHERE id = $1 FOR UPDATE`, id).Scan(&status, &pondID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCycleNotFound
	}
	if err != nil {
		return nil, err
	}
	if status != "closed" {
		return nil, ErrCycleNotClosed
	}

	// Kunci wadah (urutan sama seperti Create) lalu cek kapasitas.
	var maxStock int
	var isActive bool
	err = tx.QueryRow(ctx,
		`SELECT max_stock, is_active FROM ponds WHERE id = $1 AND deleted_at IS NULL FOR UPDATE`, pondID).
		Scan(&maxStock, &isActive)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && !isActive) {
		return nil, ErrPondUnavailable
	}
	if err != nil {
		return nil, err
	}

	var used, alive int64
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(SUM(s.alive_qty), 0)
		FROM cycles c JOIN v_cycle_status s ON s.cycle_id = c.id
		WHERE c.pond_id = $1 AND c.status = 'active'`, pondID).Scan(&used); err != nil {
		return nil, err
	}
	if err := tx.QueryRow(ctx, `SELECT alive_qty FROM v_cycle_status WHERE cycle_id = $1`, id).Scan(&alive); err != nil {
		return nil, err
	}
	if used+alive > int64(maxStock) {
		return nil, &CapacityError{MaxStock: maxStock, Used: used, Remaining: int64(maxStock) - used}
	}

	if _, err := tx.Exec(ctx,
		`UPDATE cycles SET status = 'active', end_date = NULL, closed_by = NULL WHERE id = $1`, id); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

const pondColumns = `id, name, type::text, volume_liter, max_stock, is_active, notes, created_at, updated_at`

type PondRepository struct {
	db *pgxpool.Pool
}

func NewPondRepository(db *pgxpool.Pool) *PondRepository {
	return &PondRepository{db: db}
}

func scanPond(row pgx.Row) (model.Pond, error) {
	var p model.Pond
	err := row.Scan(&p.ID, &p.Name, &p.Type, &p.VolumeLiter, &p.MaxStock,
		&p.IsActive, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PondRepository) List(ctx context.Context, f model.PondFilter) ([]model.Pond, int64, error) {
	where := []string{"deleted_at IS NULL"}
	args := []any{}
	add := func(cond string, v any) {
		args = append(args, v)
		where = append(where, fmt.Sprintf(cond, len(args)))
	}
	if f.Search != "" {
		add("name ILIKE $%d", "%"+f.Search+"%")
	}
	if f.Type != "" {
		add("type::text = $%d", f.Type)
	}
	if f.IsActive != nil {
		add("is_active = $%d", *f.IsActive)
	}
	whereSQL := strings.Join(where, " AND ")

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM ponds WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, (f.Page-1)*f.Limit)
	query := fmt.Sprintf("SELECT %s FROM ponds WHERE %s ORDER BY id DESC LIMIT $%d OFFSET $%d",
		pondColumns, whereSQL, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	ponds := make([]model.Pond, 0) // supaya JSON-nya [] bukan null
	for rows.Next() {
		p, err := scanPond(rows)
		if err != nil {
			return nil, 0, err
		}
		ponds = append(ponds, p)
	}
	return ponds, total, rows.Err()
}

func (r *PondRepository) FindByID(ctx context.Context, id int64) (*model.Pond, error) {
	p, err := scanPond(r.db.QueryRow(ctx,
		`SELECT `+pondColumns+` FROM ponds WHERE id = $1 AND deleted_at IS NULL`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PondRepository) Create(ctx context.Context, in model.PondInput) (*model.Pond, error) {
	p, err := scanPond(r.db.QueryRow(ctx, `
		INSERT INTO ponds (name, type, volume_liter, max_stock, is_active, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+pondColumns,
		in.Name, in.Type, in.VolumeLiter, in.MaxStock, *in.IsActive, in.Notes))
	if err != nil {
		return nil, translatePgError(err)
	}
	return &p, nil
}

func (r *PondRepository) Update(ctx context.Context, id int64, in model.PondInput) (*model.Pond, error) {
	p, err := scanPond(r.db.QueryRow(ctx, `
		UPDATE ponds
		SET name = $1, type = $2, volume_liter = $3, max_stock = $4, is_active = $5, notes = $6
		WHERE id = $7 AND deleted_at IS NULL
		RETURNING `+pondColumns,
		in.Name, in.Type, in.VolumeLiter, in.MaxStock, *in.IsActive, in.Notes, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, translatePgError(err)
	}
	return &p, nil
}

func (r *PondRepository) SoftDelete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE ponds SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

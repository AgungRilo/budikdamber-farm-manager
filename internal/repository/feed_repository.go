package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

const feedColumns = `id, name, brand, protein_pct::float8, price_per_kg, created_at, updated_at`

type FeedRepository struct {
	db *pgxpool.Pool
}

func NewFeedRepository(db *pgxpool.Pool) *FeedRepository {
	return &FeedRepository{db: db}
}

func scanFeed(row pgx.Row) (model.Feed, error) {
	var f model.Feed
	err := row.Scan(&f.ID, &f.Name, &f.Brand, &f.ProteinPct, &f.PricePerKg, &f.CreatedAt, &f.UpdatedAt)
	return f, err
}

func (r *FeedRepository) List(ctx context.Context, f model.FeedFilter) ([]model.Feed, int64, error) {
	const where = `deleted_at IS NULL AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR brand ILIKE '%' || $1 || '%')`

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM feeds WHERE `+where, f.Search).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+feedColumns+` FROM feeds WHERE `+where+` ORDER BY name LIMIT $2 OFFSET $3`,
		f.Search, f.Limit, (f.Page-1)*f.Limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]model.Feed, 0)
	for rows.Next() {
		fd, err := scanFeed(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, fd)
	}
	return list, total, rows.Err()
}

func (r *FeedRepository) FindByID(ctx context.Context, id int64) (*model.Feed, error) {
	f, err := scanFeed(r.db.QueryRow(ctx,
		`SELECT `+feedColumns+` FROM feeds WHERE id = $1 AND deleted_at IS NULL`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FeedRepository) Create(ctx context.Context, in model.FeedInput) (*model.Feed, error) {
	f, err := scanFeed(r.db.QueryRow(ctx, `
		INSERT INTO feeds (name, brand, protein_pct, price_per_kg)
		VALUES ($1, $2, $3, $4)
		RETURNING `+feedColumns,
		in.Name, in.Brand, in.ProteinPct, in.PricePerKg))
	if err != nil {
		return nil, translatePgError(err)
	}
	return &f, nil
}

func (r *FeedRepository) Update(ctx context.Context, id int64, in model.FeedInput) (*model.Feed, error) {
	f, err := scanFeed(r.db.QueryRow(ctx, `
		UPDATE feeds
		SET name = $1, brand = $2, protein_pct = $3, price_per_kg = $4
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING `+feedColumns,
		in.Name, in.Brand, in.ProteinPct, in.PricePerKg, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, translatePgError(err)
	}
	return &f, nil
}

func (r *FeedRepository) SoftDelete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE feeds SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

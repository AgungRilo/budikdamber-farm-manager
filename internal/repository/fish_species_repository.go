package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

const speciesColumns = `id, name, target_days, default_price_per_fish, default_price_per_kg, created_at, updated_at`

type FishSpeciesRepository struct {
	db *pgxpool.Pool
}

func NewFishSpeciesRepository(db *pgxpool.Pool) *FishSpeciesRepository {
	return &FishSpeciesRepository{db: db}
}

func scanSpecies(row pgx.Row) (model.FishSpecies, error) {
	var s model.FishSpecies
	err := row.Scan(&s.ID, &s.Name, &s.TargetDays, &s.DefaultPricePerFish,
		&s.DefaultPricePerKg, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *FishSpeciesRepository) List(ctx context.Context, f model.FishSpeciesFilter) ([]model.FishSpecies, int64, error) {
	const where = `deleted_at IS NULL AND ($1 = '' OR name ILIKE '%' || $1 || '%')`

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM fish_species WHERE `+where, f.Search).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+speciesColumns+` FROM fish_species WHERE `+where+` ORDER BY name LIMIT $2 OFFSET $3`,
		f.Search, f.Limit, (f.Page-1)*f.Limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]model.FishSpecies, 0)
	for rows.Next() {
		s, err := scanSpecies(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}
	return list, total, rows.Err()
}

func (r *FishSpeciesRepository) FindByID(ctx context.Context, id int64) (*model.FishSpecies, error) {
	s, err := scanSpecies(r.db.QueryRow(ctx,
		`SELECT `+speciesColumns+` FROM fish_species WHERE id = $1 AND deleted_at IS NULL`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *FishSpeciesRepository) Create(ctx context.Context, in model.FishSpeciesInput) (*model.FishSpecies, error) {
	s, err := scanSpecies(r.db.QueryRow(ctx, `
		INSERT INTO fish_species (name, target_days, default_price_per_fish, default_price_per_kg)
		VALUES ($1, $2, $3, $4)
		RETURNING `+speciesColumns,
		in.Name, in.TargetDays, in.DefaultPricePerFish, in.DefaultPricePerKg))
	if err != nil {
		return nil, translatePgError(err)
	}
	return &s, nil
}

func (r *FishSpeciesRepository) Update(ctx context.Context, id int64, in model.FishSpeciesInput) (*model.FishSpecies, error) {
	s, err := scanSpecies(r.db.QueryRow(ctx, `
		UPDATE fish_species
		SET name = $1, target_days = $2, default_price_per_fish = $3, default_price_per_kg = $4
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING `+speciesColumns,
		in.Name, in.TargetDays, in.DefaultPricePerFish, in.DefaultPricePerKg, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, translatePgError(err)
	}
	return &s, nil
}

func (r *FishSpeciesRepository) SoftDelete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE fish_species SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

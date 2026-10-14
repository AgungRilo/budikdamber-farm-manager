package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
)

const expenseSelect = `
	SELECT e.id, e.cycle_id, c.code, to_char(e.expense_date, 'YYYY-MM-DD'), e.category::text,
	       e.description, e.amount, e.created_by, u.name, e.created_at, e.updated_at
	FROM expenses e
	LEFT JOIN cycles c ON c.id = e.cycle_id
	JOIN users u ON u.id = e.created_by`

type ExpenseRepository struct {
	db *pgxpool.Pool
}

func NewExpenseRepository(db *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func scanExpense(row pgx.Row) (model.Expense, error) {
	var e model.Expense
	err := row.Scan(&e.ID, &e.CycleID, &e.CycleCode, &e.ExpenseDate, &e.Category,
		&e.Description, &e.Amount, &e.CreatedBy, &e.CreatedByName, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func (r *ExpenseRepository) List(ctx context.Context, f model.ExpenseFilter) ([]model.Expense, int64, error) {
	const where = ` WHERE ($1::bigint IS NULL OR e.cycle_id = $1::bigint)
		AND (NOT $2::bool OR e.cycle_id IS NULL)
		AND ($3::text IS NULL OR e.category::text = $3::text)
		AND ($4::date IS NULL OR e.expense_date >= $4::date)
		AND ($5::date IS NULL OR e.expense_date <= $5::date)`
	args := []any{f.CycleID, f.GeneralOnly, nilIfEmpty(f.Category), nilIfEmpty(f.DateFrom), nilIfEmpty(f.DateTo)}

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM expenses e`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, (f.Page-1)*f.Limit)
	rows, err := r.db.Query(ctx, expenseSelect+where+` ORDER BY e.expense_date DESC, e.id DESC LIMIT $6 OFFSET $7`, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]model.Expense, 0)
	for rows.Next() {
		e, err := scanExpense(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, e)
	}
	return list, total, rows.Err()
}

func (r *ExpenseRepository) FindByID(ctx context.Context, id int64) (*model.Expense, error) {
	e, err := scanExpense(r.db.QueryRow(ctx, expenseSelect+` WHERE e.id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// checkCycleForExpense: kalau pengeluaran terikat siklus, siklus harus active (R04)
// dan tanggal pengeluaran tidak boleh sebelum tanggal tebar.
func checkCycleForExpense(ctx context.Context, tx pgx.Tx, cycleID *int64, date string) error {
	if cycleID == nil {
		return nil
	}
	start, err := lockActiveCycle(ctx, tx, *cycleID)
	if err != nil {
		return err
	}
	d, _ := time.Parse("2006-01-02", date)
	if d.Before(start) {
		return ErrLogBeforeStart
	}
	return nil
}

func (r *ExpenseRepository) Create(ctx context.Context, in model.ExpenseInput, createdBy int64) (*model.Expense, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := checkCycleForExpense(ctx, tx, in.CycleID, in.ExpenseDate); err != nil {
		return nil, err
	}
	var id int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO expenses (cycle_id, expense_date, category, description, amount, created_by)
		VALUES ($1, $2::date, $3::expense_category, $4, $5, $6)
		RETURNING id`,
		in.CycleID, in.ExpenseDate, in.Category, in.Description, in.Amount, createdBy).Scan(&id); err != nil {
		return nil, translatePgError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *ExpenseRepository) Update(ctx context.Context, id int64, in model.ExpenseInput) (*model.Expense, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var oldCycleID *int64
	err = tx.QueryRow(ctx, `SELECT cycle_id FROM expenses WHERE id = $1 FOR UPDATE`, id).Scan(&oldCycleID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	// Siklus lama harus masih active (data di siklus tertutup tidak boleh diubah)...
	if oldCycleID != nil {
		if _, err := lockActiveCycle(ctx, tx, *oldCycleID); err != nil {
			return nil, err
		}
	}
	// ...dan siklus tujuan juga harus active.
	if err := checkCycleForExpense(ctx, tx, in.CycleID, in.ExpenseDate); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `
		UPDATE expenses
		SET cycle_id = $1, expense_date = $2::date, category = $3::expense_category, description = $4, amount = $5
		WHERE id = $6`,
		in.CycleID, in.ExpenseDate, in.Category, in.Description, in.Amount, id); err != nil {
		return nil, translatePgError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *ExpenseRepository) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var cycleID *int64
	err = tx.QueryRow(ctx, `SELECT cycle_id FROM expenses WHERE id = $1 FOR UPDATE`, id).Scan(&cycleID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if cycleID != nil {
		if _, err := lockActiveCycle(ctx, tx, *cycleID); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx, `DELETE FROM expenses WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

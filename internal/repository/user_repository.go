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

var ErrNotFound = errors.New("data tidak ditemukan")

const userColumns = `id, name, email, password_hash, role::text, is_active, created_at, updated_at`

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.findOne(ctx, `SELECT `+userColumns+` FROM users WHERE LOWER(email) = LOWER($1)`, email)
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.findOne(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
}

func (r *UserRepository) findOne(ctx context.Context, query string, arg any) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx, query, arg).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func scanUser(row pgx.Row) (model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *UserRepository) List(ctx context.Context, f model.UserFilter) ([]model.User, int64, error) {
	where := []string{"TRUE"}
	args := []any{}
	add := func(cond string, v any) {
		args = append(args, v)
		where = append(where, fmt.Sprintf(cond, len(args)))
	}
	if f.Search != "" {
		args = append(args, "%"+f.Search+"%")
		n := len(args)
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d)", n, n))
	}
	if f.Role != "" {
		add("role::text = $%d", f.Role)
	}
	if f.IsActive != nil {
		add("is_active = $%d", *f.IsActive)
	}
	whereSQL := strings.Join(where, " AND ")

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, (f.Page-1)*f.Limit)
	query := fmt.Sprintf("SELECT %s FROM users WHERE %s ORDER BY id LIMIT $%d OFFSET $%d",
		userColumns, whereSQL, len(args)-1, len(args))
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *UserRepository) Create(ctx context.Context, name, email, passwordHash, role string) (*model.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING `+userColumns,
		name, email, passwordHash, role))
	if err != nil {
		return nil, translatePgError(err)
	}
	return &u, nil
}

// Update mengubah profil; password_hash hanya diganti kalau newHash tidak nil.
func (r *UserRepository) Update(ctx context.Context, id int64, name, email, role string, newHash *string) (*model.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx, `
		UPDATE users
		SET name = $1, email = $2, role = $3, password_hash = COALESCE($4, password_hash)
		WHERE id = $5
		RETURNING `+userColumns,
		name, email, role, newHash, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, translatePgError(err)
	}
	return &u, nil
}

func (r *UserRepository) SetActive(ctx context.Context, id int64, active bool) (*model.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx,
		`UPDATE users SET is_active = $1 WHERE id = $2 RETURNING `+userColumns, active, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Register membuat akun worker nonaktif (menunggu persetujuan owner).
func (r *UserRepository) Register(ctx context.Context, name, email, passwordHash string) (*model.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, 'worker', false)
		RETURNING `+userColumns,
		name, email, passwordHash))
	if err != nil {
		return nil, translatePgError(err)
	}
	return &u, nil
}

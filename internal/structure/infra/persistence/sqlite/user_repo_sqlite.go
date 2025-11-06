// internal/repos/user_repo.go
package repos

import (
	"context"
	"database/sql"
	"day2/internal/structure/domain"
	"errors"
	"time"
)

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	List(ctx context.Context, limit, offset int) ([]domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id int) error
}

// internal/repos/user_repo_sqlite.go (скелет)
type sqliteUserRepo struct{ db *sql.DB }

func NewSQLiteUserRepo(db *sql.DB) UserRepo { return &sqliteUserRepo{db: db} }

// Реализации методов: контекст с таймаутом, SQL, scan/err -> ошибки домена

// перед реализациями
const repoTimeout = 3 * time.Second

func (r *sqliteUserRepo) Create(ctx context.Context, u *domain.User) error {
	// проверка уникальности email
	if _, err := r.GetByEmail(ctx, u.Email); err == nil {
		return domain.ErrEmailExists
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	ctx2, cancel := context.WithTimeout(ctx, repoTimeout)
	defer cancel()

	const q = `INSERT INTO users(id,name,email,age,is_active,created_at) VALUES(?,?,?,?,?,?)`
	_, err := r.db.ExecContext(ctx2, q, u.ID, u.Name, u.Email, u.Age, u.IsActive, u.CreatedAt)
	return err
}

func (r *sqliteUserRepo) GetByID(ctx context.Context, id int) (domain.User, error) {
	ctx2, cancel := context.WithTimeout(ctx, repoTimeout)
	defer cancel()

	const q = `SELECT id,name,email,age,is_active,created_at FROM users WHERE id=? LIMIT 1`
	row := r.db.QueryRowContext(ctx2, q, id)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.IsActive, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *sqliteUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	ctx2, cancel := context.WithTimeout(ctx, repoTimeout)
	defer cancel()

	const q = `SELECT id,name,email,age,is_active,created_at FROM users WHERE email=? LIMIT 1`
	row := r.db.QueryRowContext(ctx2, q, email)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.IsActive, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *sqliteUserRepo) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	ctx2, cancel := context.WithTimeout(ctx, repoTimeout)
	defer cancel()

	const q = `SELECT id,name,email,age,is_active,created_at FROM users ORDER BY id LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx2, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *sqliteUserRepo) Delete(ctx context.Context, id int) error {
	// убедись, что существует
	if _, err := r.GetByID(ctx, id); err != nil {
		return err
	}
	ctx2, cancel := context.WithTimeout(ctx, repoTimeout)
	defer cancel()

	const q = `DELETE FROM users WHERE id=?`
	_, err := r.db.ExecContext(ctx2, q, id)
	return err
}

func (r *sqliteUserRepo) Update(ctx context.Context, u *domain.User) error {
	// убедись, что пользователь существует
	if _, err := r.GetByID(ctx, u.ID); err != nil {
		return err
	}
	// проверка уникальности email при изменении
	if other, err := r.GetByEmail(ctx, u.Email); err == nil && other.ID != u.ID {
		return domain.ErrEmailExists
	} else if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	ctx2, cancel := context.WithTimeout(ctx, repoTimeout)
	defer cancel()
	const q = `UPDATE users SET name=?, email=?, age=?, is_active=?, created_at=? WHERE id=?`
	_, err := r.db.ExecContext(ctx2, q, u.Name, u.Email, u.Age, u.IsActive, u.CreatedAt, u.ID)
	return err
}

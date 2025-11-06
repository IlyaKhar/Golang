package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Migration struct {
	ID   string // "0001_create_users"
	Up   string // SQL для применения
	Down string // SQL для отката
}

const migTimeout = 5 * time.Second

// список миграций — дополняй по мере нужды
var migrations = []Migration{
	{
		ID: "0001_create_users",
		Up: `
CREATE TABLE IF NOT EXISTS users (
  id         INTEGER PRIMARY KEY,
  name       TEXT NOT NULL,
  email      TEXT NOT NULL,
  age        INTEGER NOT NULL,
  is_active  BOOLEAN NOT NULL,
  created_at TIMESTAMP NOT NULL
);`,
		Down: `DROP TABLE IF EXISTS users;`,
	},
	{
		ID:   "0002_add_index_email",
		Up:   `CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		Down: `DROP INDEX IF EXISTS idx_users_email;`,
	},
}

func ensureMigrationsTable(db *sql.DB) error {
	const q = `
CREATE TABLE IF NOT EXISTS schema_migrations (
  id TEXT PRIMARY KEY,
  applied_at TIMESTAMP NOT NULL
);`
	ctx, cancel := context.WithTimeout(context.Background(), migTimeout)
	defer cancel()
	_, err := db.ExecContext(ctx, q)
	return err
}

func isApplied(db *sql.DB, id string) (bool, error) {
	const q = `SELECT 1 FROM schema_migrations WHERE id = ? LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), migTimeout)
	defer cancel()
	var one int
	err := db.QueryRowContext(ctx, q, id).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func markApplied(db *sql.DB, id string) error {
	const q = `INSERT INTO schema_migrations(id, applied_at) VALUES(?, ?)`
	ctx, cancel := context.WithTimeout(context.Background(), migTimeout)
	defer cancel()
	_, err := db.ExecContext(ctx, q, id, time.Now())
	return err
}

// ApplyMigrations — прогоняет все не применённые миграции по порядку.
func ApplyMigrations(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("ensure migrations table: %w", err)
	}
	for _, m := range migrations {
		applied, err := isApplied(db, m.ID)
		if err != nil {
			return fmt.Errorf("check applied %s: %w", m.ID, err)
		}
		if applied {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), migTimeout)
		_, execErr := db.ExecContext(ctx, m.Up)
		cancel()
		if execErr != nil {
			return fmt.Errorf("apply %s: %w", m.ID, execErr)
		}
		if err := markApplied(db, m.ID); err != nil {
			return fmt.Errorf("mark applied %s: %w", m.ID, err)
		}
	}
	return nil
}

// helper: найти миграцию по ID
func findMigration(id string) *Migration {
	for i := range migrations {
		if migrations[i].ID == id {
			return &migrations[i]
		}
	}
	return nil
}

// LastAppliedID возвращает последний применённый ID миграции.
func LastAppliedID(db *sql.DB) (string, error) {
	if err := ensureMigrationsTable(db); err != nil {
		return "", err
	}
	const q = `SELECT id FROM schema_migrations ORDER BY applied_at DESC LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), migTimeout)
	defer cancel()
	var id string
	err := db.QueryRowContext(ctx, q).Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return id, err
}

// RollbackLast откатывает последнюю применённую миграцию (Down) и удаляет запись из schema_migrations.
func RollbackLast(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("ensure migrations table: %w", err)
	}
	id, err := LastAppliedID(db)
	if err != nil {
		return fmt.Errorf("get last applied: %w", err)
	}
	if id == "" { // нечего откатывать
		return nil
	}
	m := findMigration(id)
	if m == nil {
		return fmt.Errorf("migration %s not found in code", id)
	}
	if m.Down == "" {
		return fmt.Errorf("migration %s has no Down SQL", id)
	}
	// выполняем Down
	ctx, cancel := context.WithTimeout(context.Background(), migTimeout)
	_, execErr := db.ExecContext(ctx, m.Down)
	cancel()
	if execErr != nil {
		return fmt.Errorf("down %s: %w", id, execErr)
	}
	// удаляем запись из schema_migrations
	const del = `DELETE FROM schema_migrations WHERE id = ?`
	ctx2, cancel2 := context.WithTimeout(context.Background(), migTimeout)
	defer cancel2()
	if _, err := db.ExecContext(ctx2, del, id); err != nil {
		return fmt.Errorf("remove migration record %s: %w", id, err)
	}
	return nil
}

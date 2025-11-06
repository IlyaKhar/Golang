package db

import (
    "context"
    "database/sql"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"    // driver: pgx
    _ "modernc.org/sqlite"                 // driver: sqlite

    "day2/data"
)

const defaultTimeout = 3 * time.Second

// OpenSQLite открывает соединение с SQLite.
// Пример DSN: file:./app.db?cache=shared&mode=rwc
func OpenSQLite(dsn string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(1)
    return db, db.Ping()
}

// OpenPostgres открывает соединение с Postgres (через pgx stdlib).
// Пример DSN: postgres://user:pass@localhost:5432/mydb?sslmode=disable
func OpenPostgres(dsn string) (*sql.DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(5)
    db.SetMaxIdleConns(5)
    return db, db.Ping()
}

// Migrate создаёт таблицу users, если её нет.
// Поля синхронизированы с data.User (упрощённо, без уникальных индексов и т.п.).
func Migrate(db *sql.DB) error {
    const createUsers = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    age INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL
);`
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()
    _, err := db.ExecContext(ctx, createUsers)
    return err
}

// Ниже — заготовки CRUD. Заполни SQL и сканирование данных сам.

// CreateUserDB добавляет пользователя в БД.
func CreateUserDB(db *sql.DB, u data.User) error {
    // TODO: INSERT INTO users (id, name, email, age, is_active, created_at) VALUES (?,?,?,?,?,?)
    query := `INSERT INTO users (id, name, email, age, is_active, created_at) VALUES (?,?,?,?,?,?)`
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()
    _, err := db.ExecContext(ctx, query, u.ID, u.Name, u.Email, u.Age, u.IsActive, u.CreatedAt)
    return err
    // Для Postgres плейсхолдеры $1..$6. Можно написать две ветки, но для учебы оставь одну.
}

// GetUserDB возвращает пользователя по id.
func GetUserDB(db *sql.DB, id int) (data.User, error) {
    // TODO: SELECT id, name, email, age, is_active, created_at FROM users WHERE id = ? LIMIT 1
    // Используй ctx, row := db.QueryRowContext(ctx, query, id) и row.Scan(...)
    query := `SELECT id, name, email, age, is_active, created_at FROM users WHERE id = ? LIMIT 1`
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()
    row := db.QueryRowContext(ctx, query, id)
    var u data.User
    err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.IsActive, &u.CreatedAt)
    return u, err
}

// ListUsersDB возвращает всех пользователей (ограничение по желанию).
func ListUsersDB(db *sql.DB) ([]data.User, error) {
    // TODO: SELECT id, name, email, age, is_active, created_at FROM users ORDER BY id
    query := `SELECT id, name, email, age, is_active, created_at FROM users ORDER BY id`
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []data.User
    for rows.Next() {
        var u data.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.IsActive, &u.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    
    // Вот эта проверка после цикла
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return users, nil
}

// UpdateUserDB обновляет данные пользователя по id.
func UpdateUserDB(db *sql.DB, u data.User) error {
    // TODO: UPDATE users SET name=?, email=?, age=?, is_active=?, created_at=? WHERE id=?
    query := `UPDATE users SET name=?, email=?, age=?, is_active=?, created_at=? WHERE id=?`
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()
    _, err := db.ExecContext(ctx, query, u.Name, u.Email, u.Age, u.IsActive, u.CreatedAt, u.ID)
    return err
}

// DeleteUserDB удаляет пользователя по id.
func DeleteUserDB(db *sql.DB, id int) error {
    // TODO: DELETE FROM users WHERE id=?
    query := `DELETE FROM users WHERE id=?`
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()
    _, err := db.ExecContext(ctx, query, id)
    return err
}

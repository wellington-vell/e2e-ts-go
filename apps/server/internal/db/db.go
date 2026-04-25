package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"server/internal"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
)

var (
	DB    *sql.DB
	Query *Queries
)

func InitQueries() {
	Query = New(DB)
}

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Queries struct {
	db      DBTX
	timeout time.Duration
}

func New(db *sql.DB) *Queries {
	return &Queries{
		db:      db,
		timeout: 5 * time.Second,
	}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:      tx,
		timeout: q.timeout,
	}
}

func (q *Queries) WithTimeout(t time.Duration) *Queries {
	return &Queries{
		db:      q.db,
		timeout: t,
	}
}

func (q *Queries) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, q.timeout)
	defer cancel()
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return q.db.QueryContext(ctx, query, args...)
}

func (q *Queries) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	ctx, cancel := context.WithTimeout(ctx, q.timeout)
	defer cancel()
	return q.db.QueryRowContext(ctx, query, args...)
}

func InitDB() error {
	dbURL := internal.Env("DATABASE_URL")

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func runMigrations() error {
	provider, err := goose.NewProvider(
		database.DialectPostgres,
		DB,
		os.DirFS("./internal/db/migrations"),
	)
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}

	if _, err := provider.Up(context.Background()); err != nil {
		return fmt.Errorf("goose migration failed: %w", err)
	}

	return nil
}

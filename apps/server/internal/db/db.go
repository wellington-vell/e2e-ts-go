package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"server/internal/lib"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

var DB *bun.DB

func InitDB() error {
	dbUrl := lib.Env.DatabaseURL

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))

	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(5 * time.Minute)

	if err := sqldb.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = bun.NewDB(sqldb, pgdialect.New())

	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func runMigrations() error {
	migrator := migrate.NewMigrator(DB, Migrations)

	if err := migrator.Init(context.Background()); err != nil {
		return fmt.Errorf("failed to init migrator: %w", err)
	}

	if _, err := migrator.Migrate(context.Background()); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

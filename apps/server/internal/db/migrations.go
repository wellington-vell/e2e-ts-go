package db

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

var Migrations = migrate.NewMigrations()

func init() {
	if err := Migrations.Discover(migrationFiles); err != nil {
		panic(err)
	}
}

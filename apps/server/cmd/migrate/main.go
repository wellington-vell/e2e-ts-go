package main

import (
	"context"
	"fmt"
	"os"

	"server/internal/auth"
	"server/internal/db"
	"server/internal/lib"

	"github.com/uptrace/bun/migrate"
)

func main() {
	lib.LoadEnv()

	_, err := auth.NewAuthula()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize auth: %v\n", err)
		os.Exit(1)
	}

	if err := db.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.DB.Close()

	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up|down|status|reset]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	ctx := context.Background()

	migrator := migrate.NewMigrator(db.DB, db.Migrations)
	if err := migrator.Init(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init migrator: %v\n", err)
		os.Exit(1)
	}

	switch cmd {
	case "up":
		if err := migrateUp(ctx, migrator); err != nil {
			fmt.Fprintf(os.Stderr, "Migration up failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err := migrateDown(ctx, migrator); err != nil {
			fmt.Fprintf(os.Stderr, "Migration down failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations rolled back successfully")
	case "status":
		if err := migrateStatus(ctx, migrator); err != nil {
			fmt.Fprintf(os.Stderr, "Migration status failed: %v\n", err)
			os.Exit(1)
		}
	case "reset":
		if err := migrator.Reset(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Migration reset failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations reset successfully")
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Usage: migrate [up|down|status|reset]")
		os.Exit(1)
	}
}

func migrateUp(ctx context.Context, migrator *migrate.Migrator) error {
	_, err := migrator.Migrate(ctx)
	return err
}

func migrateDown(ctx context.Context, migrator *migrate.Migrator) error {
	_, err := migrator.Rollback(ctx)
	return err
}

func migrateStatus(ctx context.Context, migrator *migrate.Migrator) error {
	migrations, err := migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Migrations: %d total\n", len(migrations))
	for _, m := range migrations {
		status := "pending"
		if m.IsApplied() {
			status = "applied"
		}
		fmt.Printf("  [%s] %s\n", status, m.Name)
	}
	return nil
}

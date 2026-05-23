package main

import (
	"context"
	"fmt"
	"os"

	"server/internal/auth"
	"server/internal/db"
	"server/internal/lib"
)

func main() {
	lib.LoadEnv()

	if lib.Env.NodeEnv != lib.NodeEnvDevelopment {
		fmt.Fprintln(os.Stderr, "Error: seed command can only be run in development environment")
		os.Exit(1)
	}

	authInstance, err := auth.NewAuthula()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize auth: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := authInstance.ClosePlugins(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close auth plugins: %v\n", err)
		}
	}()

	if err := db.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.DB.Close()

	ctx := context.Background()
	if err := db.SeedDB(ctx, db.DB, authInstance); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to seed database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database seeded successfully")
}

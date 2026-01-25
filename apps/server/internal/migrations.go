package internal

import (
	"fmt"
)

func MigrateDB() error {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			text VARCHAR(250) NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`

	if _, err := DB.Exec(query); err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}

	return nil
}

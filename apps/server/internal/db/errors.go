package db

import (
	"errors"
	"strings"
)

var (
	ErrNotFound            = errors.New("record not found")
	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrTimeout             = errors.New("query timeout")
	ErrCancellation        = errors.New("context cancelled")
)

func Error(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	if strings.Contains(errStr, "context deadline exceeded") ||
		strings.Contains(errStr, "timeout") {
		return ErrTimeout
	}

	if strings.Contains(errStr, "context canceled") {
		return ErrCancellation
	}

	if strings.Contains(errStr, "pq: duplicate key") ||
		strings.Contains(errStr, "unique constraint") {
		return ErrUniqueViolation
	}

	if strings.Contains(errStr, "pq: foreign key") ||
		strings.Contains(errStr, "foreign key constraint") {
		return ErrForeignKeyViolation
	}

	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}

	return err
}

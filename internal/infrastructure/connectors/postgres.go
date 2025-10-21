package connectors

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/orderservice?sslmode=disable"

	pgxconf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	db, err := pgxpool.NewWithConfig(ctx, pgxconf)
	if err != nil {
		return nil, fmt.Errorf("new pgx pool: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}

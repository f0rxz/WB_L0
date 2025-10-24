package connectors

import (
	"context"
	"fmt"
	"orderservice/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := cfg.PostgresDSN

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

package connectors

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/orderservice?sslmode=disable"

	pgxconf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Error while parsing config: %v", err)
		return nil, err
	}
	db, err := pgxpool.NewWithConfig(ctx, pgxconf)
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Cant connect to db: %v", err)
		return nil, err
	}

	return db, nil
}

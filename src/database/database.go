package database

import (
	"context"
	"fmt"

	zerologadapter "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type Database struct {
	pool *pgxpool.Pool
}

type configOpt func(*pgxpool.Config)

func New(ctx context.Context, connString string, opts ...configOpt) (*Database, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("database#New - error parsing database connection string: %w", err)
	}

	for _, opt := range opts {
		opt(cfg)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("database#New - error getting database pool with config: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database#New - error pinging database: %w", err)
	}

	return &Database{
		pool: pool,
	}, nil
}

func WithLogger(log zerolog.Logger) configOpt {
	return func(cfg *pgxpool.Config) {
		cfg.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   zerologadapter.NewLogger(log),
			LogLevel: tracelog.LogLevelInfo,
		}
	}
}

func (db *Database) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("database#Begin - error starting transaction: %w", err)
	}
	return tx, nil
}

func (db *Database) Close() {
	db.pool.Close()
}

type TxnFunc func(ctx context.Context, tx pgx.Tx) error

func (db *Database) WithTxn(ctx context.Context, fn TxnFunc) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("database#WithTxn - error starting transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			fmt.Printf("database#WithTxn - error rolling back transaction: %v\n", err)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return fmt.Errorf("database#WithTxn - error executing transaction function: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("database#WithTxn - error committing transaction: %w", err)
	}

	return nil
}

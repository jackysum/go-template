package database_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackysum/go-template/src/database"
	testhelper "github.com/jackysum/go-template/src/test/helper"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ctx := t.Context()
	cfg := testhelper.Config(t)

	tests := map[string]struct {
		ctx        context.Context
		connString string
		wantErr    bool
	}{
		"success":                         {ctx: ctx, connString: cfg.DatabaseConnString, wantErr: false},
		"error parsing connection string": {ctx: ctx, connString: "invalid-connection-string", wantErr: true},
		"error pinging database": {
			ctx:        ctx,
			connString: "postgres://invalid:invalid@localhost:5432/invalid",
			wantErr:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, gotErr := database.New(tc.ctx, tc.connString, database.WithLogger(zerolog.New(nil)))

			if tc.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
		})
	}
}

func TestDatabase_Begin(t *testing.T) {
	ctx := t.Context()
	cfg := testhelper.Config(t)
	db, err := database.New(ctx, cfg.DatabaseConnString, database.WithLogger(zerolog.New(nil)))
	require.NoError(t, err)

	tests := map[string]struct {
		ctx     context.Context
		wantErr bool
	}{
		"success": {
			ctx:     ctx,
			wantErr: false,
		},
		"error starting transaction": {
			ctx:     testhelper.CancelledContext(t),
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tx, gotErr := db.Begin(tc.ctx)

			if tc.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
			require.NotNil(t, tx)

			tx.Rollback(context.Background())
		})
	}
}

func TestDatabase_WithTxn(t *testing.T) {
	ctx := t.Context()
	cfg := testhelper.Config(t)
	db, err := database.New(ctx, cfg.DatabaseConnString, database.WithLogger(zerolog.New(nil)))
	require.NoError(t, err)

	tests := map[string]struct {
		ctx     context.Context
		fn      database.TxnFunc
		wantErr bool
	}{
		"success": {
			ctx: ctx,
			fn: func(ctx context.Context, tx pgx.Tx) error {
				return nil
			},
			wantErr: false,
		},
		"error starting transaction": {
			ctx: testhelper.CancelledContext(t),
			fn: func(ctx context.Context, tx pgx.Tx) error {
				return nil
			},
			wantErr: true,
		},
		"error in transaction function": {
			ctx: ctx,
			fn: func(ctx context.Context, tx pgx.Tx) error {
				return fmt.Errorf("transaction error")
			},
			wantErr: true,
		},
		"error committing transaction": {
			ctx: ctx,
			fn: func(ctx context.Context, tx pgx.Tx) error {
				tx.Exec(ctx, `BAD SQL`)
				return nil
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := db.WithTxn(tc.ctx, tc.fn)

			if tc.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
		})
	}
}

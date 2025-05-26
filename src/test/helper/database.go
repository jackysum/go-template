package testhelper

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackysum/go-template/src/database"
	"github.com/stretchr/testify/require"
)

type mockDatabase struct {
	tx pgx.Tx
}

func (md *mockDatabase) WithTxn(ctx context.Context, fn database.TxnFunc) {
	fn(ctx, md.tx)
}

func Database(t *testing.T) *mockDatabase {
	t.Helper()

	ctx := t.Context()
	cfg := Config(t)

	db, err := database.New(ctx, cfg.DatabaseConnString)
	require.NoError(t, err)
	t.Cleanup(db.Close)

	tx, err := db.Begin(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := tx.Rollback(ctx)
		require.NoError(t, err)
	})

	return &mockDatabase{
		tx: tx,
	}
}

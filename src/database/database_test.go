package database_test

import (
	"context"
	"testing"

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

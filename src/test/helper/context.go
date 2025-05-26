package testhelper

import (
	"context"
	"testing"
)

func CancelledContext(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	return ctx
}

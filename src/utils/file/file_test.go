package file_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jackysum/go-template/src/utils/file"
	"github.com/stretchr/testify/require"
)

func TestAbsolutePath(t *testing.T) {
	path := "path/to/file"

	_, b, _, _ := runtime.Caller(0)
	want := filepath.Join(filepath.Dir(b), "../../..", path)

	require.Equal(t, want, file.AbsolutePath(path))
}

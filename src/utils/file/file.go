package file

import (
	"path/filepath"
	"runtime"
)

func AbsolutePath(path string) string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../../..", path)
}

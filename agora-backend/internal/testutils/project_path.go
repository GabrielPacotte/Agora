package testutils

import (
	"path/filepath"
	"runtime"
)

func ProjectRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

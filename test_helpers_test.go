package goexmars

import (
	"os"
	"path/filepath"
	"testing"
)

func configureTestLibraryPath(tb testing.TB) {
	tb.Helper()

	name := exmarsLibraryName()

	wd, err := os.Getwd()
	if err != nil {
		tb.Fatalf("getwd: %v", err)
	}
	src := filepath.Join(wd, "lib", name)
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			tb.Skipf("shared library not found at %s", src)
		}
		tb.Fatalf("stat source library: %v", err)
	}

	tb.Setenv(exmarsLibraryPathEnv, src)
}

package goexmars

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/ebitengine/purego"
)

const exmarsLibraryPathEnv = "GOEXMARS_LIB_PATH"

var (
	loadOnce sync.Once
	loadErr  error

	fight1Warrior  func(string, int32, int32, int32, int32, int32, int32, int32, int32, *int32, *int32, *int32)
	fight2Warriors func(string, string, int32, int32, int32, int32, int32, int32, int32, int32, *int32, *int32, *int32)
)

func loadLibrary() error {
	loadOnce.Do(func() {
		libPath, err := exmarsLibraryPath()
		if err != nil {
			loadErr = err
			return
		}

		handle, err := purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			loadErr = fmt.Errorf("open exmars library %q: %w", libPath, err)
			return
		}

		purego.RegisterLibFunc(&fight1Warrior, handle, "Fight1Warrior")
		purego.RegisterLibFunc(&fight2Warriors, handle, "Fight2Warriors")
	})

	return loadErr
}

func exmarsLibraryPath() (string, error) {
	if p := os.Getenv(exmarsLibraryPathEnv); p != "" {
		return p, nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable path: %w", err)
	}

	libDir := filepath.Join(filepath.Dir(exePath), "lib")
	return filepath.Join(libDir, exmarsLibraryName()), nil
}

func exmarsLibraryName() string {
	switch runtime.GOOS {
	case "darwin":
		return "libexmars.dylib"
	case "linux":
		return "libexmars.so"
	case "windows":
		return "exmars.dll"
	default:
		return "libexmars.so"
	}
}

func requireLibrary() {
	if err := loadLibrary(); err != nil {
		panic(err)
	}
}

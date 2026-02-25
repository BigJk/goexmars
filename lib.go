package goexmars

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

const exmarsLibraryPathEnv = "GOEXMARS_LIB_PATH"

var (
	loadOnce sync.Once
	loadErr  error

	fight1 func(string, unsafe.Pointer, unsafe.Pointer, int32, *int32, unsafe.Pointer, int32, *int32)
	fight2 func(string, string, unsafe.Pointer, unsafe.Pointer, int32, *int32, unsafe.Pointer, int32, *int32)
	fight3 func(string, string, string, unsafe.Pointer, unsafe.Pointer, int32, *int32, unsafe.Pointer, int32, *int32)
	fight4 func(string, string, string, string, unsafe.Pointer, unsafe.Pointer, int32, *int32, unsafe.Pointer, int32, *int32)
	fight5 func(string, string, string, string, string, unsafe.Pointer, unsafe.Pointer, int32, *int32, unsafe.Pointer, int32, *int32)
	fight6 func(string, string, string, string, string, string, unsafe.Pointer, unsafe.Pointer, int32, *int32, unsafe.Pointer, int32, *int32)
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

		purego.RegisterLibFunc(&fight1, handle, "fight_1")
		purego.RegisterLibFunc(&fight2, handle, "fight_2")
		purego.RegisterLibFunc(&fight3, handle, "fight_3")
		purego.RegisterLibFunc(&fight4, handle, "fight_4")
		purego.RegisterLibFunc(&fight5, handle, "fight_5")
		purego.RegisterLibFunc(&fight6, handle, "fight_6")
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

# goexmars

exMars binding for go. This uses the [purego](https://github.com/ebitengine/purego) package to call the shared library, without any cgo.

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/BigJk/goexmars"
)

func main() {
	imp := `
        ;redcode-94
        ;name Imp
        MOV 0, 1
        END
`

	dwarf := `
        ;redcode-94
        ;name Dwarf
        ADD #4, 3
        MOV 2, @2
        JMP -2, 0
        DAT #0, #0
        END
`

	cfg := goexmars.DefaultConfig.SetRounds(50)

	result, err := goexmars.Fight([]string{imp, dwarf}, cfg)
	if err != nil {
		log.Printf("fight failed: %v", err)
	}

	fmt.Printf("wins=%v ties=%d\n", result.Wins, result.Ties)
	if result.Diagnostics != "" {
		fmt.Printf("diagnostics:\n%s\n", result.Diagnostics)
	}
}
```

### Shared Library

You need the shared library to run the code. You can build it yourself or download it from the releases page. It needs to be placed in the `./lib` directory relative to the executable, or define the `GOEXMARS_LIB_PATH` environment variable to point to the shared library.

## Build shared library

- The C sources are under `./exmars`
- Build the shared library before building Go code:

```sh
cd exmars
./build.sh
```

This produces `libexmars.dylib` on macOS or `libexmars.so` on Linux.

# Credits

- [exMars](http://corewar.co.uk/ankerl/exmars.htm)

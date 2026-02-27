#!/usr/bin/env sh
set -eu

CC=${CC:-cc}
OUT_DIR=${OUT_DIR:-.}

mkdir -p "$OUT_DIR"

case "$(uname -s)" in
Darwin)
  OUT_LIB="$OUT_DIR/../lib/libexmars.dylib"
  "$CC" -O2 -fPIC -dynamiclib \
    -o "$OUT_LIB" \
    pmars.c sim.c pspace.c
  ;;
Linux)
  OUT_LIB="$OUT_DIR/../lib/libexmars.so"
  "$CC" -O2 -fPIC -shared \
    -o "$OUT_LIB" \
    pmars.c sim.c pspace.c
  ;;
MINGW*|MSYS*|CYGWIN*)
  OUT_LIB="$OUT_DIR/../lib/exmars.dll"
  "$CC" -O2 -shared -Wl,--export-all-symbols \
    -o "$OUT_LIB" \
    pmars.c sim.c pspace.c
  ;;
*)
  echo "unsupported platform: $(uname -s)" >&2
  exit 1
  ;;
esac

echo "built $OUT_LIB"

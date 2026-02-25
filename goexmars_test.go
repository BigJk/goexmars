package goexmars

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFight2WarriorsRoundsSum(t *testing.T) {
	configureTestLibraryPath(t)

	const imp = `
;redcode-94
;name Imp
MOV 0, 1
END
`

	const dwarf = `
;redcode-94
;name Dwarf
ADD #4, 3
MOV 2, @2
JMP -2, 0
DAT #0, #0
END
`

	const rounds = 10
	cfg := FightConfig{
		CoreSize:      8000,
		Cycles:        80000,
		MaxProcess:    8000,
		Rounds:        rounds,
		MaxWarriorLen: 100,
		MinSep:        100,
		FixPos:        0,
	}

	win1, win2, equal := Fight2Warriors(
		imp,
		dwarf,
		cfg,
	)

	if got := win1 + win2 + equal; got != rounds {
		t.Fatalf("unexpected total rounds: win1=%d win2=%d equal=%d total=%d want=%d", win1, win2, equal, got, rounds)
	}
}

func TestFight1WarriorMalformedReturnsMinusOne(t *testing.T) {
	configureTestLibraryPath(t)

	const malformed = `
;redcode-94
;name Broken
MOV.Z 0, 1
END
`

	cfg := FightConfig{
		CoreSize:      8000,
		Cycles:        80000,
		MaxProcess:    8000,
		Rounds:        10,
		MaxWarriorLen: 100,
		MinSep:        100,
		FixPos:        0,
	}

	win1, win2, equal := Fight1Warrior(malformed, cfg)
	if win1 != -1 || win2 != -1 || equal != -1 {
		t.Fatalf("expected malformed warrior to return -1/-1/-1, got %d/%d/%d", win1, win2, equal)
	}
}

func BenchmarkFight2WarriorsTwoImps(b *testing.B) {
	configureTestLibraryPath(b)

	const imp = `
;redcode-94
;name Imp
MOV 0, 1
END
`

	const rounds = 1
	cfg := FightConfig{
		CoreSize:      8000,
		Cycles:        80000,
		MaxProcess:    8000,
		Rounds:        rounds,
		MaxWarriorLen: 100,
		MinSep:        100,
		FixPos:        0,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		win1, win2, equal := Fight2Warriors(
			imp,
			imp,
			cfg,
		)

		if got := win1 + win2 + equal; got != rounds {
			b.Fatalf("unexpected total rounds: win1=%d win2=%d equal=%d total=%d want=%d", win1, win2, equal, got, rounds)
		}
	}
}

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

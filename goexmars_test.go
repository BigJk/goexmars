package goexmars

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFightTwoWarriorsRoundsSum(t *testing.T) {
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

	result, err := Fight([]string{imp, dwarf}, cfg)
	if err != nil {
		t.Fatalf("Fight returned unexpected error: %v", err)
	}
	if len(result.Wins) != 2 {
		t.Fatalf("expected 2 wins entries, got %d", len(result.Wins))
	}
	if got := result.Wins[0] + result.Wins[1] + result.Ties; got != rounds {
		t.Fatalf("unexpected total rounds: wins=%v ties=%d total=%d want=%d", result.Wins, result.Ties, got, rounds)
	}
}

func TestFightMalformedWarriorReturnsSentinelResult(t *testing.T) {
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

	result, err := Fight([]string{malformed}, cfg)
	if err == nil {
		t.Fatalf("expected malformed warrior to return error")
	}
	if len(result.Wins) != 1 || result.Wins[0] != -1 || result.Ties != -1 {
		t.Fatalf("expected malformed warrior sentinel result, got wins=%v ties=%d", result.Wins, result.Ties)
	}
}

func TestFightReturnsDiagnosticsOnMalformedWarrior(t *testing.T) {
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

	result, err := Fight([]string{malformed}, cfg)
	if err == nil {
		t.Fatalf("expected malformed warrior to return error")
	}
	if result.Diagnostics == "" {
		t.Fatalf("expected diagnostics for malformed warrior, got empty string")
	}
	if !strings.Contains(result.Diagnostics, "Missing 'modifier'") {
		t.Fatalf("expected diagnostics to contain parse error, got:\n%s", result.Diagnostics)
	}
	if err.Error() != result.Diagnostics {
		t.Fatalf("expected error to equal diagnostics")
	}
}

func TestFightThreeWarriorsRoundsSum(t *testing.T) {
	configureTestLibraryPath(t)

	const imp = `
;redcode-94
;name Imp
MOV 0, 1
END
`

	const rounds = 5
	cfg := FightConfig{
		CoreSize:      8000,
		Cycles:        80000,
		MaxProcess:    8000,
		Rounds:        rounds,
		MaxWarriorLen: 100,
		MinSep:        100,
		FixPos:        0,
	}

	result, err := Fight([]string{imp, imp, imp}, cfg)
	if err != nil {
		t.Fatalf("Fight returned unexpected error: %v", err)
	}
	if len(result.Wins) != 3 {
		t.Fatalf("expected 3 wins entries, got %d", len(result.Wins))
	}
	total := result.Ties
	for _, w := range result.Wins {
		total += w
	}
	if total != rounds {
		t.Fatalf("unexpected total rounds: wins=%v ties=%d total=%d want=%d diagnostics=%q", result.Wins, result.Ties, total, rounds, result.Diagnostics)
	}
}

func BenchmarkFightTwoImps(b *testing.B) {
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
		result, err := Fight([]string{imp, imp}, cfg)
		if err != nil {
			b.Fatalf("Fight returned unexpected error: %v", err)
		}
		if got := result.Wins[0] + result.Wins[1] + result.Ties; got != rounds {
			b.Fatalf("unexpected total rounds: wins=%v ties=%d total=%d want=%d", result.Wins, result.Ties, got, rounds)
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

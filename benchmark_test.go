package goexmars

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBenchmarkScore(t *testing.T) {
	configureTestLibraryPath(t)

	const imp = `
;redcode-94
;name Imp
MOV 0, 1
END
`

	cfg := DefaultConfig
	cfg.Rounds = 3

	parsed, err := AssembleParsed(imp, cfg)
	if err != nil {
		t.Fatalf("AssembleParsed failed: %v", err)
	}

	bm := Benchmark{
		Warriors: []ParsedWarrior{parsed, parsed},
		Config:   cfg,
	}

	score, err := bm.Score(parsed)
	if err != nil {
		t.Fatalf("Score returned error: %v", err)
	}
	if got, want := score.Rounds(), cfg.Rounds*len(bm.Warriors); got != want {
		t.Fatalf("unexpected total benchmark rounds: got %d want %d", got, want)
	}
	if score.Points() < 0 || score.Points() > 3*score.Rounds() {
		t.Fatalf("unexpected points: %d", score.Points())
	}
	if score.Performance() < 0 || score.Performance() > 1 {
		t.Fatalf("unexpected performance: %v", score.Performance())
	}
}

func TestBenchmarkScoreString(t *testing.T) {
	configureTestLibraryPath(t)

	const imp = `
;redcode-94
;name Imp
MOV 0, 1
END
`

	parsed, err := AssembleParsed(imp, DefaultConfig)
	if err != nil {
		t.Fatalf("AssembleParsed failed: %v", err)
	}

	bm := Benchmark{
		Warriors: []ParsedWarrior{parsed},
		Config:   DefaultConfig,
	}

	score, err := bm.ScoreString(imp)
	if err != nil {
		t.Fatalf("ScoreString returned error: %v", err)
	}
	if got, want := score.Rounds(), DefaultConfig.Rounds; got != want {
		t.Fatalf("unexpected rounds: got %d want %d", got, want)
	}
	if score.Points() < 0 || score.Points() > 3*score.Rounds() {
		t.Fatalf("unexpected points: %d", score.Points())
	}
}

func TestBenchmarkFromFolder(t *testing.T) {
	configureTestLibraryPath(t)

	dir := t.TempDir()

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

	if err := os.WriteFile(filepath.Join(dir, "a.red"), []byte(imp), 0o644); err != nil {
		t.Fatalf("write a.red: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.red"), []byte(dwarf), 0o644); err != nil {
		t.Fatalf("write b.red: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ignore.txt"), []byte("noop"), 0o644); err != nil {
		t.Fatalf("write ignore.txt: %v", err)
	}

	bm, err := BenchmarkFromFolder(dir, DefaultConfig)
	if err != nil {
		t.Fatalf("BenchmarkFromFolder returned error: %v", err)
	}
	if len(bm.Warriors) != 2 {
		t.Fatalf("expected 2 warriors, got %d", len(bm.Warriors))
	}
	if bm.Config.CoreSize != DefaultConfig.CoreSize {
		t.Fatalf("expected benchmark config to be set")
	}
}

func TestBenchmarkScorePointsHelpers(t *testing.T) {
	score := BenchmarkScore{Wins: 2, Losses: 1, Ties: 3}
	if score.Points() != 9 {
		t.Fatalf("unexpected points: %d", score.Points())
	}
	if score.Rounds() != 6 {
		t.Fatalf("unexpected rounds: %d", score.Rounds())
	}
	if score.Performance() != 0.5 {
		t.Fatalf("unexpected performance: %v", score.Performance())
	}
}

package goexmars

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Benchmark is a set of parsed warriors used for fight-based scoring.
type Benchmark struct {
	// Warriors contains the parsed benchmark warriors.
	Warriors []ParsedWarrior
	// Config is used for benchmark fights and for ScoreString assembly/parsing.
	Config FightConfig
}

// BenchmarkScore is the aggregated result of fighting against a benchmark set.
//
// Points are ICWS-style: 3 per win, 1 per tie, 0 per loss.
type BenchmarkScore struct {
	Wins   int
	Losses int
	Ties   int
}

// Points returns the ICWS-style point total (3*Wins + Ties).
func (s BenchmarkScore) Points() int {
	return (3 * s.Wins) + s.Ties
}

// Rounds returns the total number of benchmark rounds represented by the score.
func (s BenchmarkScore) Rounds() int {
	return s.Wins + s.Losses + s.Ties
}

// Performance returns the normalized ICWS-style performance in [0,1].
func (s BenchmarkScore) Performance() float64 {
	rounds := s.Rounds()
	if rounds == 0 {
		return 0
	}
	return float64(s.Points()) / float64(3*rounds)
}

// Score fights warrior against all benchmark warriors and aggregates wins/losses/ties.
func (b Benchmark) Score(warrior ParsedWarrior) (BenchmarkScore, error) {
	cfg := b.Config
	if cfg == (FightConfig{}) {
		cfg = DefaultConfig
	}
	if err := cfg.Validate(); err != nil {
		return BenchmarkScore{}, err
	}

	var total BenchmarkScore
	if len(b.Warriors) == 0 {
		return total, nil
	}

	candidate := warrior.String()
	for i, opponent := range b.Warriors {
		result, err := Fight([]string{candidate, opponent.String()}, cfg)
		if err != nil {
			return BenchmarkScore{}, fmt.Errorf("fight vs benchmark warrior %d: %w", i, err)
		}
		if len(result.Wins) != 2 {
			return BenchmarkScore{}, fmt.Errorf("unexpected result size for benchmark warrior %d: %d", i, len(result.Wins))
		}
		total.Wins += result.Wins[0]
		total.Losses += result.Wins[1]
		total.Ties += result.Ties
	}
	return total, nil
}

// ScoreString assembles and parses warrior, then fights it against the benchmark set.
func (b Benchmark) ScoreString(warrior string) (BenchmarkScore, error) {
	cfg := b.Config
	if cfg == (FightConfig{}) {
		cfg = DefaultConfig
	}
	parsed, err := AssembleParsed(warrior, cfg)
	if err != nil {
		return BenchmarkScore{}, err
	}
	return b.Score(parsed)
}

// BenchmarkFromFolder loads all .red files in folder into a benchmark.
//
// Files are loaded in lexicographic filename order for deterministic results.
func BenchmarkFromFolder(folder string, cfg FightConfig) (Benchmark, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return Benchmark{}, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".red") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	warriors := make([]ParsedWarrior, 0, len(names))
	for _, name := range names {
		path := filepath.Join(folder, name)
		src, err := os.ReadFile(path)
		if err != nil {
			return Benchmark{}, fmt.Errorf("read %s: %w", path, err)
		}
		parsed, err := AssembleParsed(string(src), cfg)
		if err != nil {
			return Benchmark{}, fmt.Errorf("assemble %s: %w", path, err)
		}
		warriors = append(warriors, parsed)
	}

	return Benchmark{
		Warriors: warriors,
		Config:   cfg,
	}, nil
}

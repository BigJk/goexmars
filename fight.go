package goexmars

import (
	"errors"
	"fmt"
	"sort"
	"unsafe"
)

const diagnosticsBufferSize = 16 * 1024

type cFightCfg struct {
	CoreSize      int32
	Cycles        int32
	MaxProcess    int32
	Rounds        int32
	MaxWarriorLen int32
	MinSep        int32
	PSpaceSize    int32
	FixPos        int32
}

func toCFightCfg(cfg FightConfig) cFightCfg {
	return cFightCfg{
		CoreSize:      int32(cfg.CoreSize),
		Cycles:        int32(cfg.Cycles),
		MaxProcess:    int32(cfg.MaxProcess),
		Rounds:        int32(cfg.Rounds),
		MaxWarriorLen: int32(cfg.MaxWarriorLen),
		MinSep:        int32(cfg.MinSep),
		PSpaceSize:    int32(cfg.PSpaceSize),
		FixPos:        int32(cfg.FixPos),
	}
}

func diagnosticsString(buf []byte, diagLen int32) string {
	if diagLen > int32(len(buf)-1) {
		diagLen = int32(len(buf) - 1)
	}
	if diagLen < 0 {
		diagLen = 0
	}
	return string(buf[:diagLen])
}

// FightResult contains the outcome of a fight.
type FightResult struct {
	// Wins contains the sole-win count for each warrior in the input order.
	Wins []int
	// Ties contains rounds without a sole winner.
	Ties int
	// Diagnostics contains exmars warnings/errors captured during assembly/fight setup.
	Diagnostics string
}

// FightNamedResult is a FightResult with name-based lookup helpers.
type FightNamedResult struct {
	FightResult
	index map[string]int
}

// Get returns the sole-win count for name and the shared tie count.
//
// If name is unknown, Get returns 0 and the tie count.
func (r FightNamedResult) Get(name string) (int, int) {
	if i, ok := r.index[name]; ok && i >= 0 && i < len(r.Wins) {
		return r.Wins[i], r.Ties
	}
	return 0, r.Ties
}

// Failed reports whether the result encodes a sentinel failure from the C layer.
func (r FightResult) Failed() bool {
	if r.Ties < 0 {
		return true
	}
	for _, w := range r.Wins {
		if w < 0 {
			return true
		}
	}
	return false
}

// Validate performs a quick validity check for a single warrior.
//
// It runs a single-round self-fight and returns a non-nil error when exmars
// reports an assembly/setup failure. The returned error message is the captured
// diagnostics string when available.
func Validate(warrior string, cfg FightConfig) error {
	cfg.Rounds = 1
	result, err := Fight([]string{warrior}, cfg)
	if err != nil {
		return err
	}
	if result.Failed() {
		if result.Diagnostics != "" {
			return errors.New(result.Diagnostics)
		}
		return errors.New("warrior validation failed")
	}
	return nil
}

// Assemble assembles a single warrior and returns its normalized instruction listing.
//
// The returned source is a disassembly of the assembled instructions (not the
// original source with labels/macros/comments). If exmars reports an assembly
// failure, Assemble returns an error containing diagnostics when available.
func Assemble(warrior string, cfg FightConfig) (string, error) {
	requireLibrary()

	cfg.Rounds = 1
	cfgC := toCFightCfg(cfg)
	outBuf := make([]byte, diagnosticsBufferSize)
	diagBuf := make([]byte, diagnosticsBufferSize)
	var outLen int32
	var diagLen int32

	rc := assemble1(
		warrior,
		unsafe.Pointer(&cfgC),
		unsafe.Pointer(&outBuf[0]), int32(len(outBuf)), &outLen,
		unsafe.Pointer(&diagBuf[0]), int32(len(diagBuf)), &diagLen,
	)

	diag := diagnosticsString(diagBuf, diagLen)
	if rc != 0 {
		if diag != "" {
			return "", errors.New(diag)
		}
		return "", errors.New("assembly failed")
	}

	if outLen > int32(len(outBuf)-1) {
		outLen = int32(len(outBuf) - 1)
	}
	if outLen < 0 {
		outLen = 0
	}
	return string(outBuf[:outLen]), nil
}

// FightNamed runs a fight for 1 to 6 named warriors.
//
// Warrior map keys are used as stable identifiers for result lookup. Internally
// names are sorted to make evaluation order deterministic. Use result.Get(name)
// to map a warrior name back to its sole-win count and the shared tie count.
func FightNamed(warriors map[string]string, cfg FightConfig) (FightNamedResult, error) {
	if len(warriors) < 1 || len(warriors) > 6 {
		return FightNamedResult{}, fmt.Errorf("FightNamed supports 1 to 6 warriors, got %d", len(warriors))
	}

	names := make([]string, 0, len(warriors))
	for name := range warriors {
		names = append(names, name)
	}
	sort.Strings(names)

	sources := make([]string, len(names))
	index := make(map[string]int, len(names))
	for i, name := range names {
		sources[i] = warriors[name]
		index[name] = i
	}

	result, err := Fight(sources, cfg)
	return FightNamedResult{FightResult: result, index: index}, err
}

// Fight runs a fight for 1 to 6 warriors and returns the fight result.
//
// On parser/setup failure, the returned FightResult contains sentinel values
// (negative wins/ties) and error is set to the diagnostics string when available.
func Fight(warriors []string, cfg FightConfig) (FightResult, error) {
	requireLibrary()

	if len(warriors) < 1 || len(warriors) > 6 {
		return FightResult{}, fmt.Errorf("Fight supports 1 to 6 warriors, got %d", len(warriors))
	}

	cfgC := toCFightCfg(cfg)
	wins32 := make([]int32, len(warriors))
	var ties32 int32
	var diagLen int32
	diagBuf := make([]byte, diagnosticsBufferSize)

	switch len(warriors) {
	case 1:
		fight1(
			warriors[0],
			unsafe.Pointer(&cfgC),
			unsafe.Pointer(&wins32[0]), int32(len(wins32)),
			&ties32,
			unsafe.Pointer(&diagBuf[0]), int32(len(diagBuf)), &diagLen,
		)
	case 2:
		fight2(
			warriors[0], warriors[1],
			unsafe.Pointer(&cfgC),
			unsafe.Pointer(&wins32[0]), int32(len(wins32)),
			&ties32,
			unsafe.Pointer(&diagBuf[0]), int32(len(diagBuf)), &diagLen,
		)
	case 3:
		fight3(
			warriors[0], warriors[1], warriors[2],
			unsafe.Pointer(&cfgC),
			unsafe.Pointer(&wins32[0]), int32(len(wins32)),
			&ties32,
			unsafe.Pointer(&diagBuf[0]), int32(len(diagBuf)), &diagLen,
		)
	case 4:
		fight4(
			warriors[0], warriors[1], warriors[2], warriors[3],
			unsafe.Pointer(&cfgC),
			unsafe.Pointer(&wins32[0]), int32(len(wins32)),
			&ties32,
			unsafe.Pointer(&diagBuf[0]), int32(len(diagBuf)), &diagLen,
		)
	case 5:
		fight5(
			warriors[0], warriors[1], warriors[2], warriors[3], warriors[4],
			unsafe.Pointer(&cfgC),
			unsafe.Pointer(&wins32[0]), int32(len(wins32)),
			&ties32,
			unsafe.Pointer(&diagBuf[0]), int32(len(diagBuf)), &diagLen,
		)
	case 6:
		fight6(
			warriors[0], warriors[1], warriors[2], warriors[3], warriors[4], warriors[5],
			unsafe.Pointer(&cfgC),
			unsafe.Pointer(&wins32[0]), int32(len(wins32)),
			&ties32,
			unsafe.Pointer(&diagBuf[0]), int32(len(diagBuf)), &diagLen,
		)
	}

	result := FightResult{
		Wins:        make([]int, len(wins32)),
		Ties:        int(ties32),
		Diagnostics: diagnosticsString(diagBuf, diagLen),
	}
	for i, v := range wins32 {
		result.Wins[i] = int(v)
	}

	if result.Failed() {
		if result.Diagnostics != "" {
			return result, errors.New(result.Diagnostics)
		}
		return result, errors.New("fight failed")
	}
	return result, nil
}

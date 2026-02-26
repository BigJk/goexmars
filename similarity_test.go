package goexmars

import "testing"

func TestDistanceAndSimilarityIdentical(t *testing.T) {
	a := ParsedWarrior{
		Commands: []Command{
			{OpCode: OpCodeMOV, Modifier: ModifierI, AddressingModeA: AddressingDirect, A: 0, AddressingModeB: AddressingDirect, B: 1},
			{OpCode: OpCodeJMP, Modifier: ModifierB, AddressingModeA: AddressingDirect, A: -1, AddressingModeB: AddressingDirect, B: 0},
		},
	}
	b := a

	if d := Distance(a, b); d != 0 {
		t.Fatalf("expected zero distance, got %v", d)
	}
	if s := Similarity(a, b); s != 1 {
		t.Fatalf("expected similarity 1, got %v", s)
	}
}

func TestDistanceWeightedSubstitution(t *testing.T) {
	a := ParsedWarrior{Commands: []Command{{OpCode: OpCodeMOV, Modifier: ModifierI, AddressingModeA: AddressingDirect, A: 0, AddressingModeB: AddressingDirect, B: 1}}}
	b := ParsedWarrior{Commands: []Command{{OpCode: OpCodeMOV, Modifier: ModifierI, AddressingModeA: AddressingDirect, A: 0, AddressingModeB: AddressingDirect, B: 2}}}

	d := Distance(a, b)
	if d <= 0 || d >= 1 {
		t.Fatalf("expected weighted substitution distance between 0 and 1, got %v", d)
	}
	s := Similarity(a, b)
	if s <= 0 || s >= 1 {
		t.Fatalf("expected similarity between 0 and 1, got %v", s)
	}
}

func TestDistanceInsertDelete(t *testing.T) {
	a := ParsedWarrior{Commands: []Command{{OpCode: OpCodeDAT, Modifier: ModifierF, AddressingModeA: AddressingImmediate, A: 0, AddressingModeB: AddressingImmediate, B: 0}}}
	b := ParsedWarrior{}

	d := Distance(a, b)
	if d != 1 {
		t.Fatalf("expected single delete distance 1, got %v", d)
	}
	s := Similarity(a, b)
	if s != 0 {
		t.Fatalf("expected similarity 0, got %v", s)
	}
}

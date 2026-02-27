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

func TestEditScriptIdenticalUsesKeeps(t *testing.T) {
	a := ParsedWarrior{
		Commands: []Command{
			{OpCode: OpCodeMOV, Modifier: ModifierI, AddressingModeA: AddressingDirect, A: 0, AddressingModeB: AddressingDirect, B: 1},
			{OpCode: OpCodeDAT, Modifier: ModifierF, AddressingModeA: AddressingImmediate, A: 0, AddressingModeB: AddressingImmediate, B: 0},
		},
	}
	ops := EditScript(a, a)
	if len(ops) != len(a.Commands) {
		t.Fatalf("expected %d ops, got %d", len(a.Commands), len(ops))
	}
	for i, op := range ops {
		if op.Kind != EditKeep {
			t.Fatalf("op %d: expected keep, got %s", i, op.Kind.String())
		}
		if op.Cost != 0 {
			t.Fatalf("op %d: expected zero cost, got %v", i, op.Cost)
		}
	}
}

func TestEditScriptWeightedSubstitution(t *testing.T) {
	a := ParsedWarrior{Commands: []Command{{OpCode: OpCodeMOV, Modifier: ModifierI, AddressingModeA: AddressingDirect, A: 0, AddressingModeB: AddressingDirect, B: 1}}}
	b := ParsedWarrior{Commands: []Command{{OpCode: OpCodeMOV, Modifier: ModifierI, AddressingModeA: AddressingDirect, A: 0, AddressingModeB: AddressingDirect, B: 2}}}

	ops := EditScript(a, b)
	if len(ops) != 1 {
		t.Fatalf("expected single op, got %d", len(ops))
	}
	if ops[0].Kind != EditSubstitute {
		t.Fatalf("expected substitute op, got %s", ops[0].Kind.String())
	}
	if ops[0].Cost <= 0 || ops[0].Cost >= 1 {
		t.Fatalf("expected weighted substitution cost between 0 and 1, got %v", ops[0].Cost)
	}
}

func TestEditScriptInsertDelete(t *testing.T) {
	a := ParsedWarrior{}
	b := ParsedWarrior{Commands: []Command{{OpCode: OpCodeDAT, Modifier: ModifierF, AddressingModeA: AddressingImmediate, A: 0, AddressingModeB: AddressingImmediate, B: 0}}}

	ops := EditScript(a, b)
	if len(ops) != 1 {
		t.Fatalf("expected single op, got %d", len(ops))
	}
	if ops[0].Kind != EditInsert {
		t.Fatalf("expected insert op, got %s", ops[0].Kind.String())
	}
	if ops[0].Cost != 1 {
		t.Fatalf("expected insert cost 1, got %v", ops[0].Cost)
	}
}

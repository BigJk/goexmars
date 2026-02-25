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

func TestValidateValidWarrior(t *testing.T) {
	configureTestLibraryPath(t)

	const imp = `
;redcode-94
;name Imp
MOV 0, 1
END
`

	if err := Validate(imp, DefaultConfig); err != nil {
		t.Fatalf("expected valid warrior, got error: %v", err)
	}
}

func TestValidateMalformedWarrior(t *testing.T) {
	configureTestLibraryPath(t)

	const malformed = `
;redcode-94
;name Broken
MOV.Z 0, 1
END
`

	err := Validate(malformed, DefaultConfig)
	if err == nil {
		t.Fatalf("expected malformed warrior to fail validation")
	}
	if !strings.Contains(err.Error(), "Missing 'modifier'") {
		t.Fatalf("expected validation error to contain parse diagnostics, got: %v", err)
	}
}

func TestAssembleValidWarrior(t *testing.T) {
	configureTestLibraryPath(t)

	const imp = `
;redcode-94x
;name Four Winds
;author John Metcalf
;strategy oneshot

     step  equ 694
     diff  equ 5

sc:  sub   inc,            ptr
ptr: sne.b >step*4+diff,   >step*4
     djn.f sc,             <ptr

inc: spl   #-step,         <-step-1
     mov   @bptr,          >ptr
     mov   @bptr,          >ptr
bptr:djn.f -2,             {clr

     dat   -5,             8
clr: spl   #-101,          16

     end   ptr

`

	out, err := Assemble(imp, DefaultConfig)
	if err != nil {
		t.Fatalf("expected Assemble to succeed, got error: %v", err)
	}
	if out == "" {
		t.Fatalf("expected assembled output, got empty string")
	}
	if !strings.Contains(out, "MOV.") {
		t.Fatalf("expected normalized instruction output, got:\n%s", out)
	}
}

func TestAssembleMalformedWarrior(t *testing.T) {
	configureTestLibraryPath(t)

	const malformed = `
;redcode-94
;name Broken
MOV.Z 0, 1
END
`

	out, err := Assemble(malformed, DefaultConfig)
	if err == nil {
		t.Fatalf("expected Assemble to fail for malformed warrior")
	}
	if out != "" {
		t.Fatalf("expected no assembled output on failure, got:\n%s", out)
	}
	if !strings.Contains(err.Error(), "Missing 'modifier'") {
		t.Fatalf("expected assembly error to contain diagnostics, got: %v", err)
	}
}

func TestAssembleParsedIncludesMetadataAndCommands(t *testing.T) {
	configureTestLibraryPath(t)

	const warrior = `
;redcode-94
;name Example
;author Test Author
step DAT #0, #0
start MOV step, >step
JMP 0
END 0
`

	parsed, err := AssembleParsed(warrior, DefaultConfig)
	if err != nil {
		t.Fatalf("expected AssembleParsed to succeed, got error: %v", err)
	}
	if parsed.Name != "Example" {
		t.Fatalf("unexpected name: %q", parsed.Name)
	}
	if parsed.Author != "Test Author" {
		t.Fatalf("unexpected author: %q", parsed.Author)
	}
	if parsed.End != 0 {
		t.Fatalf("unexpected END target: %d", parsed.End)
	}
	if len(parsed.Commands) == 0 {
		t.Fatalf("expected parsed commands")
	}
}

func TestParseAssembledCommands(t *testing.T) {
	assembled := "MOV.I $     0, $     1\nDAT.F #     0, #     0\n"
	cmds, err := ParseAssembledCommands(assembled)
	if err != nil {
		t.Fatalf("ParseAssembledCommands returned error: %v", err)
	}
	if len(cmds) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(cmds))
	}
	if cmds[0].OpCode != OpCodeMOV || cmds[0].Modifier != ModifierI {
		t.Fatalf("unexpected first command opcode/modifier: %+v", cmds[0])
	}
	if cmds[1].OpCode != OpCodeDAT || cmds[1].Modifier != ModifierF {
		t.Fatalf("unexpected second command opcode/modifier: %+v", cmds[1])
	}
}

func TestCommandString(t *testing.T) {
	cmd := Command{
		OpCode:          OpCodeMOV,
		Modifier:        ModifierI,
		AddressingModeA: AddressingDirect,
		A:               0,
		AddressingModeB: AddressingDirect,
		B:               1,
	}
	if got, want := cmd.String(), "MOV.I $0, $1"; got != want {
		t.Fatalf("unexpected command string: got %q want %q", got, want)
	}
}

func TestParsedWarriorString(t *testing.T) {
	pw := ParsedWarrior{
		Name:   "Example",
		Author: "Tester",
		End:    0,
		Commands: []Command{
			{
				OpCode:          OpCodeDAT,
				Modifier:        ModifierF,
				AddressingModeA: AddressingImmediate,
				A:               0,
				AddressingModeB: AddressingImmediate,
				B:               0,
			},
			{
				OpCode:          OpCodeJMP,
				Modifier:        ModifierB,
				AddressingModeA: AddressingDirect,
				A:               0,
				AddressingModeB: AddressingDirect,
				B:               0,
			},
		},
	}
	out := pw.String()
	if !strings.Contains(out, ";name Example") || !strings.Contains(out, ";author Tester") {
		t.Fatalf("expected metadata in parsed warrior string, got:\n%s", out)
	}
	if !strings.Contains(out, "DAT.F #0, #0") || !strings.Contains(out, "JMP.B $0, $0") {
		t.Fatalf("expected commands in parsed warrior string, got:\n%s", out)
	}
	if !strings.Contains(out, "END 0") {
		t.Fatalf("expected END line in parsed warrior string, got:\n%s", out)
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

func TestFightNamedGet(t *testing.T) {
	configureTestLibraryPath(t)

	const imp = `
;redcode-94
;name Imp
MOV 0, 1
END
`

	cfg := FightConfig{
		CoreSize:      8000,
		Cycles:        80000,
		MaxProcess:    8000,
		Rounds:        3,
		MaxWarriorLen: 100,
		MinSep:        100,
		FixPos:        0,
	}

	result, err := FightNamed(map[string]string{
		"alpha": imp,
		"beta":  imp,
		"gamma": imp,
	}, cfg)
	if err != nil {
		t.Fatalf("FightNamed returned unexpected error: %v", err)
	}

	wa, ta := result.Get("alpha")
	wb, tb := result.Get("beta")
	wg, tg := result.Get("gamma")
	if ta != result.Ties || tb != result.Ties || tg != result.Ties {
		t.Fatalf("expected Get ties to match result ties, got %d/%d/%d want %d", ta, tb, tg, result.Ties)
	}
	if wa+wb+wg+result.Ties != cfg.Rounds {
		t.Fatalf("unexpected total rounds from named results: %d", wa+wb+wg+result.Ties)
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

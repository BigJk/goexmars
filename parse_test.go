package goexmars

import (
	"strings"
	"testing"
)

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

func TestParsedWarriorFormatOptions(t *testing.T) {
	pw := ParsedWarrior{
		Name:   "Example",
		Author: "Tester",
		End:    3,
		Commands: []Command{
			{OpCode: OpCodeDAT, Modifier: ModifierF, AddressingModeA: AddressingImmediate, A: 0, AddressingModeB: AddressingImmediate, B: 0},
		},
	}

	out := pw.Format(RedcodeFormatOptions{
		IncludeName:   false,
		IncludeAuthor: false,
		IncludeEnd:    false,
	})
	if strings.Contains(out, ";name") || strings.Contains(out, ";author") || strings.Contains(out, "END ") {
		t.Fatalf("unexpected metadata/end in formatted output: %q", out)
	}
	if !strings.Contains(out, "DAT.F #0, #0") {
		t.Fatalf("expected command in formatted output, got %q", out)
	}
}

func TestAssembleNormalizedOptions(t *testing.T) {
	configureTestLibraryPath(t)

	const warrior = `
;redcode-94
;name Example
;author Test Author
MOV 0, 1
END 0
`
	out, err := AssembleNormalized(warrior, DefaultConfig, RedcodeFormatOptions{
		IncludeName:   false,
		IncludeAuthor: false,
		IncludeEnd:    false,
	})
	if err != nil {
		t.Fatalf("AssembleNormalized returned error: %v", err)
	}
	if strings.Contains(out, ";name") || strings.Contains(out, ";author") || strings.Contains(out, "END ") {
		t.Fatalf("expected output without metadata/end, got:\n%s", out)
	}
}

func TestParsedWarriorFingerprint(t *testing.T) {
	a := ParsedWarrior{
		Name:   "One",
		Author: "A",
		End:    0,
		Commands: []Command{
			{OpCode: OpCodeMOV, Modifier: ModifierI, AddressingModeA: AddressingDirect, A: 0, AddressingModeB: AddressingDirect, B: 1},
		},
	}
	b := a
	b.Name = "Two"
	b.Author = "B"

	if a.Fingerprint() != b.Fingerprint() {
		t.Fatalf("expected default fingerprint to ignore metadata")
	}
	if a.FingerprintWithOptions(FingerprintOptions{IncludeName: true, IncludeAuthor: true, IncludeEnd: true}) ==
		b.FingerprintWithOptions(FingerprintOptions{IncludeName: true, IncludeAuthor: true, IncludeEnd: true}) {
		t.Fatalf("expected metadata-inclusive fingerprint to differ")
	}
}

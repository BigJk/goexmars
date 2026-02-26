package goexmars

import (
	"fmt"
	"strconv"
	"strings"
)

// OpCode is a normalized Redcode opcode.
type OpCode byte

// Supported opcodes.
const (
	OpCodeDAT OpCode = iota
	OpCodeMOV
	OpCodeADD
	OpCodeSUB
	OpCodeMUL
	OpCodeDIV
	OpCodeMOD
	OpCodeJMP
	OpCodeJMZ
	OpCodeJMN
	OpCodeDJN
	OpCodeSPL
	OpCodeCMP
	OpCodeSEQ
	OpCodeSNE
	OpCodeSLT
	OpCodeLDP
	OpCodeSTP
	OpCodeNOP
)

// Modifier is a Redcode instruction modifier.
type Modifier byte

// Supported modifiers.
const (
	ModifierF Modifier = iota
	ModifierA
	ModifierB
	ModifierAB
	ModifierBA
	ModifierX
	ModifierI
)

// AddressingMode is a Redcode operand addressing mode.
type AddressingMode byte

// Supported addressing modes.
const (
	AddressingImmediate AddressingMode = iota
	AddressingDirect
	AddressingAIndirect
	AddressingBIndirect
	AddressingAIndirectPre
	AddressingBIndirectPre
	AddressingAIndirectPost
	AddressingBIndirectPost
)

// Command is a single normalized Redcode instruction.
type Command struct {
	OpCode          OpCode
	Modifier        Modifier
	AddressingModeA AddressingMode
	A               int
	AddressingModeB AddressingMode
	B               int
}

// String renders the command as normalized Redcode.
func (c Command) String() string {
	return fmt.Sprintf(
		"%s.%s %s%d, %s%d",
		c.OpCode.String(),
		c.Modifier.String(),
		c.AddressingModeA.String(),
		c.A,
		c.AddressingModeB.String(),
		c.B,
	)
}

// ParsedWarrior is a structured representation of an assembled warrior.
type ParsedWarrior struct {
	Name      string
	Author    string
	End       int
	Commands  []Command
	Assembled string
}

// RedcodeFormatOptions controls how a ParsedWarrior is rendered back to Redcode text.
type RedcodeFormatOptions struct {
	IncludeName   bool
	IncludeAuthor bool
	IncludeEnd    bool
}

// DefaultRedcodeFormatOptions returns the default Redcode rendering options.
func DefaultRedcodeFormatOptions() RedcodeFormatOptions {
	return RedcodeFormatOptions{
		IncludeName:   true,
		IncludeAuthor: true,
		IncludeEnd:    true,
	}
}

// String renders the parsed warrior back to Redcode text.
func (w ParsedWarrior) String() string {
	return w.Format(DefaultRedcodeFormatOptions())
}

// Format renders the parsed warrior to Redcode text using the provided options.
func (w ParsedWarrior) Format(opts RedcodeFormatOptions) string {
	var b strings.Builder
	if opts.IncludeName && w.Name != "" {
		b.WriteString(";name ")
		b.WriteString(w.Name)
		b.WriteByte('\n')
	}
	if opts.IncludeAuthor && w.Author != "" {
		b.WriteString(";author ")
		b.WriteString(w.Author)
		b.WriteByte('\n')
	}
	for _, cmd := range w.Commands {
		b.WriteString(cmd.String())
		b.WriteByte('\n')
	}
	if opts.IncludeEnd {
		b.WriteString("END ")
		b.WriteString(strconv.Itoa(w.End))
	}
	return b.String()
}

// AssembleParsed assembles warrior and parses the normalized result into a Go struct.
//
// Name and Author are parsed from the original source if present.
// End and Commands are parsed from the normalized assembled Redcode returned by Assemble.
func AssembleParsed(warrior string, cfg FightConfig) (ParsedWarrior, error) {
	assembled, err := Assemble(warrior, cfg)
	if err != nil {
		return ParsedWarrior{}, err
	}
	cmds, err := ParseAssembledCommands(assembled)
	if err != nil {
		return ParsedWarrior{}, err
	}
	end, err := parseAssembledEnd(assembled)
	if err != nil {
		return ParsedWarrior{}, err
	}
	name, author := parseWarriorMetadata(warrior)
	return ParsedWarrior{
		Name:      name,
		Author:    author,
		End:       end,
		Commands:  cmds,
		Assembled: assembled,
	}, nil
}

// ParseAssembledCommands parses normalized assembled Redcode instructions into commands.
func ParseAssembledCommands(assembled string) ([]Command, error) {
	lines := strings.Split(assembled, "\n")
	out := make([]Command, 0, len(lines))
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(strings.ToUpper(line), "END") {
			continue
		}
		cmd, err := parseAssembledLine(line)
		if err != nil {
			return nil, fmt.Errorf("parse assembled line %d: %w", i+1, err)
		}
		out = append(out, cmd)
	}
	return out, nil
}

func parseAssembledLine(line string) (Command, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return Command{}, fmt.Errorf("unexpected instruction format: %q", line)
	}

	opmod := strings.SplitN(fields[0], ".", 2)
	if len(opmod) != 2 {
		return Command{}, fmt.Errorf("missing opcode modifier: %q", fields[0])
	}
	op, ok := parseOpCode(opmod[0])
	if !ok {
		return Command{}, fmt.Errorf("unknown opcode: %q", opmod[0])
	}
	mod, ok := parseModifier(opmod[1])
	if !ok {
		return Command{}, fmt.Errorf("unknown modifier: %q", opmod[1])
	}

	modeA, ok := parseAddressingMode(fields[1])
	if !ok {
		return Command{}, fmt.Errorf("unknown A mode: %q", fields[1])
	}
	a, err := strconv.Atoi(strings.TrimSuffix(fields[2], ","))
	if err != nil {
		return Command{}, fmt.Errorf("invalid A operand: %w", err)
	}

	modeB, ok := parseAddressingMode(fields[3])
	if !ok {
		return Command{}, fmt.Errorf("unknown B mode: %q", fields[3])
	}
	b, err := strconv.Atoi(fields[4])
	if err != nil {
		return Command{}, fmt.Errorf("invalid B operand: %w", err)
	}

	return Command{
		OpCode:          op,
		Modifier:        mod,
		AddressingModeA: modeA,
		A:               a,
		AddressingModeB: modeB,
		B:               b,
	}, nil
}

func parseOpCode(s string) (OpCode, bool) {
	switch strings.ToUpper(s) {
	case "DAT":
		return OpCodeDAT, true
	case "MOV":
		return OpCodeMOV, true
	case "ADD":
		return OpCodeADD, true
	case "SUB":
		return OpCodeSUB, true
	case "MUL":
		return OpCodeMUL, true
	case "DIV":
		return OpCodeDIV, true
	case "MOD":
		return OpCodeMOD, true
	case "JMP":
		return OpCodeJMP, true
	case "JMZ":
		return OpCodeJMZ, true
	case "JMN":
		return OpCodeJMN, true
	case "DJN":
		return OpCodeDJN, true
	case "SPL":
		return OpCodeSPL, true
	case "CMP":
		return OpCodeCMP, true
	case "SEQ":
		return OpCodeSEQ, true
	case "SNE":
		return OpCodeSNE, true
	case "SLT":
		return OpCodeSLT, true
	case "LDP":
		return OpCodeLDP, true
	case "STP":
		return OpCodeSTP, true
	case "NOP":
		return OpCodeNOP, true
	default:
		return 0, false
	}
}

// String returns the Redcode mnemonic for the opcode.
func (o OpCode) String() string {
	switch o {
	case OpCodeDAT:
		return "DAT"
	case OpCodeMOV:
		return "MOV"
	case OpCodeADD:
		return "ADD"
	case OpCodeSUB:
		return "SUB"
	case OpCodeMUL:
		return "MUL"
	case OpCodeDIV:
		return "DIV"
	case OpCodeMOD:
		return "MOD"
	case OpCodeJMP:
		return "JMP"
	case OpCodeJMZ:
		return "JMZ"
	case OpCodeJMN:
		return "JMN"
	case OpCodeDJN:
		return "DJN"
	case OpCodeSPL:
		return "SPL"
	case OpCodeCMP:
		return "CMP"
	case OpCodeSEQ:
		return "SEQ"
	case OpCodeSNE:
		return "SNE"
	case OpCodeSLT:
		return "SLT"
	case OpCodeLDP:
		return "LDP"
	case OpCodeSTP:
		return "STP"
	case OpCodeNOP:
		return "NOP"
	default:
		return "DAT"
	}
}

func parseModifier(s string) (Modifier, bool) {
	switch strings.ToUpper(s) {
	case "F":
		return ModifierF, true
	case "A":
		return ModifierA, true
	case "B":
		return ModifierB, true
	case "AB":
		return ModifierAB, true
	case "BA":
		return ModifierBA, true
	case "X":
		return ModifierX, true
	case "I":
		return ModifierI, true
	default:
		return 0, false
	}
}

// String returns the Redcode mnemonic for the modifier.
func (m Modifier) String() string {
	switch m {
	case ModifierF:
		return "F"
	case ModifierA:
		return "A"
	case ModifierB:
		return "B"
	case ModifierAB:
		return "AB"
	case ModifierBA:
		return "BA"
	case ModifierX:
		return "X"
	case ModifierI:
		return "I"
	default:
		return "F"
	}
}

func parseAddressingMode(s string) (AddressingMode, bool) {
	if len(s) != 1 {
		return 0, false
	}
	switch s[0] {
	case '#':
		return AddressingImmediate, true
	case '$':
		return AddressingDirect, true
	case '*':
		return AddressingAIndirect, true
	case '@':
		return AddressingBIndirect, true
	case '{':
		return AddressingAIndirectPre, true
	case '<':
		return AddressingBIndirectPre, true
	case '}':
		return AddressingAIndirectPost, true
	case '>':
		return AddressingBIndirectPost, true
	default:
		return 0, false
	}
}

// String returns the Redcode operand prefix for the addressing mode.
func (m AddressingMode) String() string {
	switch m {
	case AddressingImmediate:
		return "#"
	case AddressingDirect:
		return "$"
	case AddressingAIndirect:
		return "*"
	case AddressingBIndirect:
		return "@"
	case AddressingAIndirectPre:
		return "{"
	case AddressingBIndirectPre:
		return "<"
	case AddressingAIndirectPost:
		return "}"
	case AddressingBIndirectPost:
		return ">"
	default:
		return "$"
	}
}

func parseAssembledEnd(assembled string) (int, error) {
	for _, raw := range strings.Split(assembled, "\n") {
		line := strings.TrimSpace(raw)
		if !strings.HasPrefix(strings.ToUpper(line), "END") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return 0, fmt.Errorf("assembled END line missing value")
		}
		v, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, fmt.Errorf("parse assembled END value %q: %w", parts[1], err)
		}
		return v, nil
	}
	return 0, fmt.Errorf("assembled END line not found")
}

func parseWarriorMetadata(src string) (name string, author string) {
	for _, raw := range strings.Split(src, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, ";name") {
			name = strings.TrimSpace(line[len(";name"):])
			continue
		}
		if strings.HasPrefix(lower, ";author") {
			author = strings.TrimSpace(line[len(";author"):])
			continue
		}
		if strings.HasPrefix(strings.ToUpper(line), "END") {
			continue
		}
	}
	return name, author
}

package goexmars

import "math"

// CommandDistanceWeights controls weighted edit-distance costs for commands.
type CommandDistanceWeights struct {
	Insert          float64
	Delete          float64
	Opcode          float64
	Modifier        float64
	AddressingModeA float64
	AddressingModeB float64
	OperandA        float64
	OperandB        float64
	CoreSize        int
}

// DefaultCommandDistanceWeights returns the default weighted edit-distance costs.
func DefaultCommandDistanceWeights() CommandDistanceWeights {
	return CommandDistanceWeights{
		Insert:          1.0,
		Delete:          1.0,
		Opcode:          1.0,
		Modifier:        0.5,
		AddressingModeA: 0.5,
		AddressingModeB: 0.5,
		OperandA:        1.0,
		OperandB:        1.0,
		CoreSize:        DefaultConfig.CoreSize,
	}
}

// CommandSubstitutionCost returns the weighted substitution cost between two commands.
func CommandSubstitutionCost(a, b Command, w CommandDistanceWeights) float64 {
	if a == b {
		return 0
	}
	var cost float64
	if a.OpCode != b.OpCode {
		cost += w.Opcode
	}
	if a.Modifier != b.Modifier {
		cost += w.Modifier
	}
	if a.AddressingModeA != b.AddressingModeA {
		cost += w.AddressingModeA
	}
	if a.AddressingModeB != b.AddressingModeB {
		cost += w.AddressingModeB
	}
	cost += w.OperandA * operandDistanceCost(a.A, b.A, w.CoreSize)
	cost += w.OperandB * operandDistanceCost(a.B, b.B, w.CoreSize)
	return cost
}

// CommandDistance computes a weighted Levenshtein distance between two command sequences.
func CommandDistance(a, b []Command) float64 {
	return CommandDistanceWithWeights(a, b, DefaultCommandDistanceWeights())
}

// CommandDistanceWithWeights computes a weighted Levenshtein distance between two command sequences.
func CommandDistanceWithWeights(a, b []Command, w CommandDistanceWeights) float64 {
	if len(a) == 0 {
		return float64(len(b)) * w.Insert
	}
	if len(b) == 0 {
		return float64(len(a)) * w.Delete
	}

	prev := make([]float64, len(b)+1)
	curr := make([]float64, len(b)+1)
	for j := 1; j <= len(b); j++ {
		prev[j] = prev[j-1] + w.Insert
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = prev[0] + w.Delete
		for j := 1; j <= len(b); j++ {
			del := prev[j] + w.Delete
			ins := curr[j-1] + w.Insert
			sub := prev[j-1] + CommandSubstitutionCost(a[i-1], b[j-1], w)
			curr[j] = minf(del, minf(ins, sub))
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

// CommandDistanceWithCoreSize computes a weighted Levenshtein distance between two command sequences.
func CommandDistanceWithCoreSize(a, b []Command, coreSize int) float64 {
	def := DefaultCommandDistanceWeights()
	def.CoreSize = coreSize
	return CommandDistanceWithWeights(a, b, def)
}

// Distance computes a weighted distance between two parsed warriors.
func Distance(a, b ParsedWarrior) float64 {
	return CommandDistance(a.Commands, b.Commands)
}

// DistanceWithWeights computes a weighted distance between two parsed warriors.
func DistanceWithWeights(a, b ParsedWarrior, w CommandDistanceWeights) float64 {
	return CommandDistanceWithWeights(a.Commands, b.Commands, w)
}

// DistanceWithCoreSize computes a weighted distance between two parsed warriors with a custom core size.
func DistanceWithCoreSize(a, b ParsedWarrior, coreSize int) float64 {
	return CommandDistanceWithCoreSize(a.Commands, b.Commands, coreSize)
}

// Similarity computes a normalized similarity score in [0,1].
//
// 1 means identical command sequences. 0 means maximally different under the
// normalization denominator (worst-case insert/delete for sequence lengths).
func Similarity(a, b ParsedWarrior) float64 {
	return SimilarityWithWeights(a, b, DefaultCommandDistanceWeights())
}

// SimilarityWithWeights computes a normalized similarity score in [0,1].
func SimilarityWithWeights(a, b ParsedWarrior, w CommandDistanceWeights) float64 {
	dist := CommandDistanceWithWeights(a.Commands, b.Commands, w)
	maxCost := math.Max(float64(len(a.Commands))*w.Delete, float64(len(b.Commands))*w.Insert)
	if maxCost == 0 {
		return 1
	}
	sim := 1 - (dist / maxCost)
	if sim < 0 {
		return 0
	}
	if sim > 1 {
		return 1
	}
	return sim
}

// SimilarityWithCoreSize computes a normalized similarity score in [0,1].
func SimilarityWithCoreSize(a, b ParsedWarrior, coreSize int) float64 {
	def := DefaultCommandDistanceWeights()
	def.CoreSize = coreSize
	return SimilarityWithWeights(a, b, def)
}

func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func operandDistanceCost(a, b, coreSize int) float64 {
	if a == b {
		return 0
	}
	if coreSize <= 0 {
		return minf(1, math.Abs(float64(a-b)))
	}

	delta := absInt(a - b)
	if delta > coreSize {
		delta %= coreSize
	}
	wrapped := delta
	if alt := coreSize - delta; alt < wrapped {
		wrapped = alt
	}
	if wrapped < 0 {
		wrapped = 0
	}
	return float64(wrapped) / float64(coreSize)
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

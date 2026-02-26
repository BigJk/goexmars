package goexmars

import "math"

// EditKind identifies the type of command edit operation.
type EditKind byte

// Supported edit operation kinds.
const (
	EditKeep EditKind = iota
	EditSubstitute
	EditInsert
	EditDelete
)

// String returns the symbolic name of the edit kind.
func (k EditKind) String() string {
	switch k {
	case EditKeep:
		return "keep"
	case EditSubstitute:
		return "substitute"
	case EditInsert:
		return "insert"
	case EditDelete:
		return "delete"
	default:
		return "unknown"
	}
}

// EditOp describes a single weighted edit operation between command sequences.
type EditOp struct {
	Kind   EditKind
	AIndex int
	BIndex int
	From   Command
	To     Command
	Cost   float64
}

// EditScript computes a weighted edit script between two parsed warriors.
func EditScript(a, b ParsedWarrior) []EditOp {
	return EditScriptWithWeights(a, b, DefaultCommandDistanceWeights())
}

// EditScriptWithWeights computes a weighted edit script between two parsed warriors.
func EditScriptWithWeights(a, b ParsedWarrior, w CommandDistanceWeights) []EditOp {
	return CommandEditScriptWithWeights(a.Commands, b.Commands, w)
}

// EditScriptWithCoreSize computes a weighted edit script between two parsed warriors.
func EditScriptWithCoreSize(a, b ParsedWarrior, coreSize int) []EditOp {
	w := DefaultCommandDistanceWeights()
	w.CoreSize = coreSize
	return CommandEditScriptWithWeights(a.Commands, b.Commands, w)
}

// CommandEditScript computes a weighted edit script between two command sequences.
func CommandEditScript(a, b []Command) []EditOp {
	return CommandEditScriptWithWeights(a, b, DefaultCommandDistanceWeights())
}

// CommandEditScriptWithWeights computes a weighted edit script between two command sequences.
func CommandEditScriptWithWeights(a, b []Command, w CommandDistanceWeights) []EditOp {
	n, m := len(a), len(b)
	dp := make([][]float64, n+1)
	trace := make([][]byte, n+1)
	for i := range dp {
		dp[i] = make([]float64, m+1)
		trace[i] = make([]byte, m+1)
	}

	for i := 1; i <= n; i++ {
		dp[i][0] = dp[i-1][0] + w.Delete
		trace[i][0] = byte(EditDelete)
	}
	for j := 1; j <= m; j++ {
		dp[0][j] = dp[0][j-1] + w.Insert
		trace[0][j] = byte(EditInsert)
	}

	const eps = 1e-9
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			subCost := CommandSubstitutionCost(a[i-1], b[j-1], w)
			best := dp[i-1][j-1] + subCost
			bestKind := byte(EditSubstitute)

			del := dp[i-1][j] + w.Delete
			if del+eps < best {
				best = del
				bestKind = byte(EditDelete)
			}

			ins := dp[i][j-1] + w.Insert
			if ins+eps < best {
				best = ins
				bestKind = byte(EditInsert)
			}

			dp[i][j] = best
			trace[i][j] = bestKind
		}
	}

	ops := make([]EditOp, 0, maxInt(n, m))
	i, j := n, m
	for i > 0 || j > 0 {
		switch EditKind(trace[i][j]) {
		case EditInsert:
			j--
			ops = append(ops, EditOp{
				Kind:   EditInsert,
				AIndex: -1,
				BIndex: j,
				To:     b[j],
				Cost:   w.Insert,
			})
		case EditDelete:
			i--
			ops = append(ops, EditOp{
				Kind:   EditDelete,
				AIndex: i,
				BIndex: -1,
				From:   a[i],
				Cost:   w.Delete,
			})
		default:
			i--
			j--
			cost := CommandSubstitutionCost(a[i], b[j], w)
			kind := EditSubstitute
			if nearlyZero(cost) {
				kind = EditKeep
			}
			ops = append(ops, EditOp{
				Kind:   kind,
				AIndex: i,
				BIndex: j,
				From:   a[i],
				To:     b[j],
				Cost:   cost,
			})
		}
	}

	reverseEditOps(ops)
	return ops
}

func reverseEditOps(v []EditOp) {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
}

func nearlyZero(v float64) bool {
	return math.Abs(v) < 1e-9
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

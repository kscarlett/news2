package news2

import "fmt"

// RiskLevel is the clinical risk category that NEWS2 maps an aggregate score
// to. It determines the urgency of the clinical response and the minimum
// observation frequency (see the package documentation and README for the
// full thresholds-and-triggers table).
type RiskLevel int

const (
	// Low risk: aggregate score 0-4 with no single parameter scoring 3.
	// Ward-based response; routine monitoring.
	Low RiskLevel = iota
	// LowMedium risk: aggregate score below 5 but with a score of 3 in any
	// single parameter. Urgent ward-based review.
	LowMedium
	// Medium risk: aggregate score 5-6. Key threshold for an urgent
	// response by a clinician competent in assessing acute illness.
	Medium
	// High risk: aggregate score 7 or more. Emergency response, usually
	// including critical care assessment; continuous monitoring.
	High
)

var riskLevelName = map[RiskLevel]string{
	Low:       "Low",
	LowMedium: "Low-Medium",
	Medium:    "Medium",
	High:      "High",
}

// String returns the string representation of the RiskLevel.
// Returns "Unknown" if the level is not recognized.
func (r RiskLevel) String() string {
	if name, ok := riskLevelName[r]; ok {
		return name
	}
	return "Unknown"
}

// Result is the full outcome of a NEWS2 calculation: the aggregate score,
// the subscore each parameter contributed, and the escalation indicators
// derived from them.
type Result struct {
	// Total is the aggregate NEWS2 score (0-20).
	Total int

	// Per-parameter subscores, as they would appear on a NEWS2 chart.
	RespRate       int
	Saturations    int
	SupplementalO2 int
	SystolicBP     int
	Pulse          int
	Consciousness  int
	Temperature    int

	// RedFlag reports whether any single physiological parameter scored 3.
	// On its own (with a total below 5) this triggers an urgent ward-based
	// review regardless of the aggregate score.
	RedFlag bool

	// Risk is the clinical risk category for the aggregate score and
	// red-flag state.
	Risk RiskLevel
}

// String returns a human-readable summary of the result.
func (r Result) String() string {
	return fmt.Sprintf("NEWS2: %d (risk: %s, red flag: %t)", r.Total, r.Risk, r.RedFlag)
}

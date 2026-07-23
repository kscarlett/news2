package news2

import (
	"fmt"
)

// ConsciousnessLevel represents the patient's level of consciousness on the
// ACVPU scale used by NEWS2.
type ConsciousnessLevel int

const (
	// Alert means the patient is fully awake and responsive. Patients with
	// chronic, baseline confusion (e.g. long-standing dementia) are recorded
	// as Alert unless their confusion is new or worse than baseline.
	Alert ConsciousnessLevel = iota
	// Confused means new-onset confusion or delirium, or confusion that is
	// worse than the patient's known baseline.
	Confused
	// Voice means the patient responds only to a verbal stimulus.
	Voice
	// Pain means the patient responds only to a painful stimulus.
	Pain
	// Unresponsive means the patient does not respond to any stimulus.
	Unresponsive
)

var consciousnessLevelName = map[ConsciousnessLevel]string{
	Alert:        "Alert",
	Confused:     "Confused",
	Voice:        "Voice",
	Pain:         "Pain",
	Unresponsive: "Unresponsive",
}

// String returns the string representation of the ConsciousnessLevel.
// Returns "Unknown" if the level is not recognized.
func (c ConsciousnessLevel) String() string {
	if name, ok := consciousnessLevelName[c]; ok {
		return name
	}
	return "Unknown"
}

// SpO2Scale selects which NEWS2 oxygen saturation scoring scale applies to
// the patient. The zero value is Scale1, which is correct for the vast
// majority of patients.
type SpO2Scale int

const (
	// Scale1 is the standard oxygen saturation scale, used for all patients
	// unless Scale2 has been explicitly prescribed.
	Scale1 SpO2Scale = iota
	// Scale2 is used only for patients with confirmed hypercapnic
	// respiratory failure (usually due to COPD) with a prescribed oxygen
	// saturation target of 88-92%. The decision to use Scale2 must be made
	// by a competent clinical decision maker and recorded in the patient's
	// notes.
	Scale2
)

// String returns the string representation of the SpO2Scale.
// Returns "Unknown" if the scale is not recognized.
func (s SpO2Scale) String() string {
	switch s {
	case Scale1:
		return "Scale 1"
	case Scale2:
		return "Scale 2"
	}
	return "Unknown"
}

// VitalSigns holds one complete set of patient observations for NEWS2
// scoring. All fields must be populated; Calculate and CalculateScore
// validate them and return an error for values that cannot be meaningfully
// scored, which also guards against accidentally scoring an incomplete
// struct. See Validate for what is checked.
type VitalSigns struct {
	// RespRate is the respiration rate in breaths per minute.
	RespRate int
	// OxygenSat is the peripheral oxygen saturation (SpO2) as a percentage.
	OxygenSat int
	// SystolicBP is the systolic blood pressure in mmHg.
	SystolicBP int
	// Pulse is the heart rate in beats per minute.
	Pulse int
	// Temp is the body temperature in degrees Celsius.
	Temp float64
	// ConsciousnessLevel is the patient's consciousness on the ACVPU scale.
	ConsciousnessLevel ConsciousnessLevel
	// OnOxygen reports whether the patient is receiving supplemental oxygen.
	OnOxygen bool
	// SpO2Scale selects the saturation scoring scale. The zero value is
	// Scale1, the correct default for most patients; only set Scale2 when
	// it has been clinically prescribed.
	SpO2Scale SpO2Scale
}

// String returns a string representation of the VitalSigns struct.
// It formats the vital signs in a human-readable way.
func (v VitalSigns) String() string {
	return fmt.Sprintf(
		"RespRate: %d, OxygenSat: %d%% (%s), OnOxygen: %t, SystolicBP: %d mmHg, Pulse: %d BPM, ConsciousnessLevel: %s, Temp: %.1f°C",
		v.RespRate, v.OxygenSat, v.SpO2Scale, v.OnOxygen, v.SystolicBP, v.Pulse, v.ConsciousnessLevel, v.Temp,
	)
}

// Calculate computes the full NEWS2 result from a set of vital signs,
// including the per-parameter subscores, the aggregate score, the red-flag
// indicator (a score of 3 in any single physiological parameter) and the
// resulting clinical risk category.
//
// It returns an error without scoring if the vital signs fail validation;
// see VitalSigns.Validate.
func Calculate(v VitalSigns) (Result, error) {
	if err := v.Validate(); err != nil {
		return Result{}, err
	}

	r := Result{
		RespRate:      calculateRespiratoryRateScore(v.RespRate),
		Saturations:   calculateSaturationsScore(v.OxygenSat, v.OnOxygen, v.SpO2Scale),
		SystolicBP:    calculateSystolicBPScore(v.SystolicBP),
		Pulse:         calculatePulseScore(v.Pulse),
		Consciousness: calculateConsciousnessScore(v.ConsciousnessLevel),
		Temperature:   calculateTemperatureScore(v.Temp),
	}
	if v.OnOxygen {
		r.SupplementalO2 = 2
	}

	r.Total = r.RespRate + r.Saturations + r.SupplementalO2 + r.SystolicBP +
		r.Pulse + r.Consciousness + r.Temperature
	r.RedFlag = r.RespRate == 3 || r.Saturations == 3 || r.SystolicBP == 3 ||
		r.Pulse == 3 || r.Consciousness == 3 || r.Temperature == 3
	r.Risk = riskLevel(r.Total, r.RedFlag)

	return r, nil
}

// CalculateScore computes the aggregate NEWS2 score from vital signs. It is
// a convenience wrapper around Calculate for callers that only need the
// numeric score; prefer Calculate when the clinical risk category or the
// per-parameter breakdown is needed, since escalation decisions cannot be
// derived from the aggregate score alone.
func CalculateScore(v VitalSigns) (int, error) {
	r, err := Calculate(v)
	if err != nil {
		return 0, err
	}
	return r.Total, nil
}

func riskLevel(total int, redFlag bool) RiskLevel {
	switch {
	case total >= 7:
		return High
	case total >= 5:
		return Medium
	case redFlag:
		return LowMedium
	default:
		return Low
	}
}

func calculateRespiratoryRateScore(respRate int) int {
	if respRate <= 8 {
		return 3
	} else if respRate <= 11 {
		return 1
	} else if respRate <= 20 {
		return 0
	} else if respRate <= 24 {
		return 2
	} else {
		return 3
	}
}

func calculateSaturationsScore(oxygenSat int, onOxygen bool, scale SpO2Scale) int {
	if scale == Scale2 {
		if onOxygen {
			if oxygenSat >= 97 {
				return 3
			} else if oxygenSat >= 95 {
				return 2
			} else if oxygenSat >= 93 {
				return 1
			} else if oxygenSat >= 88 {
				return 0
			} else if oxygenSat >= 86 {
				return 1
			} else if oxygenSat >= 84 {
				return 2
			} else {
				return 3
			}
		}
		if oxygenSat >= 88 {
			return 0
		} else if oxygenSat >= 86 {
			return 1
		} else if oxygenSat >= 84 {
			return 2
		} else {
			return 3
		}
	}
	if oxygenSat >= 96 {
		return 0
	} else if oxygenSat >= 94 {
		return 1
	} else if oxygenSat >= 92 {
		return 2
	} else {
		return 3
	}
}

func calculateSystolicBPScore(systolicBP int) int {
	if systolicBP >= 220 {
		return 3
	} else if systolicBP >= 111 {
		return 0
	} else if systolicBP >= 101 {
		return 1
	} else if systolicBP >= 91 {
		return 2
	} else {
		return 3
	}
}

func calculatePulseScore(pulse int) int {
	if pulse <= 40 {
		return 3
	} else if pulse <= 50 {
		return 1
	} else if pulse <= 90 {
		return 0
	} else if pulse <= 110 {
		return 1
	} else if pulse <= 130 {
		return 2
	} else {
		return 3
	}
}

func calculateConsciousnessScore(level ConsciousnessLevel) int {
	if level == Alert {
		return 0
	}
	return 3
}

func calculateTemperatureScore(temp float64) int {
	if temp <= 35.0 {
		return 3
	} else if temp <= 36.0 {
		return 1
	} else if temp <= 38.0 {
		return 0
	} else if temp <= 39.0 {
		return 1
	} else {
		return 2
	}
}

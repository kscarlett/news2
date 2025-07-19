package news2

import (
	"fmt"
)

// VitalSigns holds the input data for NEWS2 scoring.
type VitalSigns struct {
	RespRate           int
	OxygenSat          int
	SystolicBP         int
	Pulse              int
	Temp               float64
	ConsciousnessLevel ConsciousnessLevel
	OnOxygen           bool
}

// ConsciousnessLevel represents the level of consciousness using ACVPU.
type ConsciousnessLevel int

const (
	Alert ConsciousnessLevel = iota
	Confused
	Voice
	Pain
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

// String returns a string representation of the VitalSigns struct.
// It formats the vital signs in a human-readable way.
func (v VitalSigns) String() string {
	return fmt.Sprintf(
		"RespRate: %d, OxygenSat: %d%%, OnOxygen: %t, SystolicBP: %d mmHg, Pulse: %d BPM, ConsciousnessLevel: %s, Temp: %.1f°C",
		v.RespRate, v.OxygenSat, v.OnOxygen, v.SystolicBP, v.Pulse, v.ConsciousnessLevel.String(), v.Temp,
	)
}

// CalculateScore computes the NEWS2 score from vital signs.
func CalculateScore(v VitalSigns, scale1 bool) int {
	score := 0
	score += calculateRespiratoryRateScore(v.RespRate)
	score += calculateSaturationsScore(v.OxygenSat, v.OnOxygen, scale1)
	if v.OnOxygen {
		score += 2
	}
	score += calculateSystolicBPScore(v.SystolicBP)
	score += calculatePulseScore(v.Pulse)
	if v.ConsciousnessLevel != Alert {
		score += 3
	}
	score += calculateTemperatureScore(v.Temp)
	return score
}

func calculateRespiratoryRateScore(respRate int) int {
	if respRate <= 8 {
		return 3
	} else if respRate <= 11 {
		return 1
	} else if respRate <= 20 {
		return 0
	} else if respRate <= 24 {
		return 1
	} else {
		return 3
	}
}

func calculateSaturationsScore(oxygenSat int, onOxygen, scale1 bool) int {
	if scale1 {
		if oxygenSat >= 96 {
			return 0
		} else if oxygenSat >= 94 {
			return 1
		} else if oxygenSat >= 92 {
			return 2
		} else {
			return 3
		}
	} else {
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
		if oxygenSat >= 93 {
			return 0
		} else if oxygenSat >= 86 {
			return 1
		} else if oxygenSat >= 84 {
			return 2
		} else {
			return 3
		}
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

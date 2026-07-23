package news2

import "testing"

// These tests pin every scoring band to the RCP NEWS2 chart (Royal College of
// Physicians, 2017) by testing both edges of each band.

func TestRespiratoryRateBoundaries(t *testing.T) {
	tests := []struct {
		respRate int
		expected int
	}{
		{4, 3},
		{8, 3},
		{9, 1},
		{11, 1},
		{12, 0},
		{20, 0},
		{21, 2},
		{24, 2},
		{25, 3},
		{40, 3},
	}
	for _, tt := range tests {
		if got := calculateRespiratoryRateScore(tt.respRate); got != tt.expected {
			t.Errorf("calculateRespiratoryRateScore(%d) = %d, want %d", tt.respRate, got, tt.expected)
		}
	}
}

func TestSaturationsScale1Boundaries(t *testing.T) {
	tests := []struct {
		oxygenSat int
		expected  int
	}{
		{85, 3},
		{91, 3},
		{92, 2},
		{93, 2},
		{94, 1},
		{95, 1},
		{96, 0},
		{100, 0},
	}
	for _, tt := range tests {
		// Scale 1 saturation scoring is independent of supplemental oxygen.
		for _, onOxygen := range []bool{false, true} {
			if got := calculateSaturationsScore(tt.oxygenSat, onOxygen, Scale1); got != tt.expected {
				t.Errorf("calculateSaturationsScore(%d, %t, scale1) = %d, want %d", tt.oxygenSat, onOxygen, got, tt.expected)
			}
		}
	}
}

func TestSaturationsScale2OnAirBoundaries(t *testing.T) {
	tests := []struct {
		oxygenSat int
		expected  int
	}{
		{80, 3},
		{83, 3},
		{84, 2},
		{85, 2},
		{86, 1},
		{87, 1},
		{88, 0},
		{92, 0},
		{93, 0},
		{100, 0},
	}
	for _, tt := range tests {
		if got := calculateSaturationsScore(tt.oxygenSat, false, Scale2); got != tt.expected {
			t.Errorf("calculateSaturationsScore(%d, on air, scale2) = %d, want %d", tt.oxygenSat, got, tt.expected)
		}
	}
}

func TestSaturationsScale2OnOxygenBoundaries(t *testing.T) {
	tests := []struct {
		oxygenSat int
		expected  int
	}{
		{80, 3},
		{83, 3},
		{84, 2},
		{85, 2},
		{86, 1},
		{87, 1},
		{88, 0},
		{92, 0},
		{93, 1},
		{94, 1},
		{95, 2},
		{96, 2},
		{97, 3},
		{100, 3},
	}
	for _, tt := range tests {
		if got := calculateSaturationsScore(tt.oxygenSat, true, Scale2); got != tt.expected {
			t.Errorf("calculateSaturationsScore(%d, on oxygen, scale2) = %d, want %d", tt.oxygenSat, got, tt.expected)
		}
	}
}

func TestSystolicBPBoundaries(t *testing.T) {
	tests := []struct {
		systolicBP int
		expected   int
	}{
		{70, 3},
		{90, 3},
		{91, 2},
		{100, 2},
		{101, 1},
		{110, 1},
		{111, 0},
		{219, 0},
		{220, 3},
		{250, 3},
	}
	for _, tt := range tests {
		if got := calculateSystolicBPScore(tt.systolicBP); got != tt.expected {
			t.Errorf("calculateSystolicBPScore(%d) = %d, want %d", tt.systolicBP, got, tt.expected)
		}
	}
}

func TestPulseBoundaries(t *testing.T) {
	tests := []struct {
		pulse    int
		expected int
	}{
		{30, 3},
		{40, 3},
		{41, 1},
		{50, 1},
		{51, 0},
		{90, 0},
		{91, 1},
		{110, 1},
		{111, 2},
		{130, 2},
		{131, 3},
		{180, 3},
	}
	for _, tt := range tests {
		if got := calculatePulseScore(tt.pulse); got != tt.expected {
			t.Errorf("calculatePulseScore(%d) = %d, want %d", tt.pulse, got, tt.expected)
		}
	}
}

func TestConsciousnessScore(t *testing.T) {
	tests := []struct {
		level    ConsciousnessLevel
		expected int
	}{
		{Alert, 0},
		{Confused, 3},
		{Voice, 3},
		{Pain, 3},
		{Unresponsive, 3},
	}
	for _, tt := range tests {
		if got := calculateConsciousnessScore(tt.level); got != tt.expected {
			t.Errorf("calculateConsciousnessScore(%s) = %d, want %d", tt.level, got, tt.expected)
		}
	}
}

func TestTemperatureBoundaries(t *testing.T) {
	tests := []struct {
		temp     float64
		expected int
	}{
		{34.0, 3},
		{35.0, 3},
		{35.1, 1},
		{36.0, 1},
		{36.1, 0},
		{38.0, 0},
		{38.1, 1},
		{39.0, 1},
		{39.1, 2},
		{41.0, 2},
	}
	for _, tt := range tests {
		if got := calculateTemperatureScore(tt.temp); got != tt.expected {
			t.Errorf("calculateTemperatureScore(%.1f) = %d, want %d", tt.temp, got, tt.expected)
		}
	}
}

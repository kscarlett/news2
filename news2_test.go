package news2

import (
	"testing"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		vitals   VitalSigns
		scale1   bool
		expected int
	}{
		{
			name: "All normal, scale1, no oxygen",
			vitals: VitalSigns{
				RespRate:           16,
				OxygenSat:          98,
				SystolicBP:         120,
				Pulse:              80,
				Temp:               37.0,
				ConsciousnessLevel: Alert,
				OnOxygen:           false,
			},
			scale1:   true,
			expected: 0,
		},
		{
			name: "Low sats, on oxygen, scale1",
			vitals: VitalSigns{
				RespRate:           18,
				OxygenSat:          90,
				SystolicBP:         100,
				Pulse:              100,
				Temp:               36.5,
				ConsciousnessLevel: Voice,
				OnOxygen:           true,
			},
			scale1:   true,
			expected: 3 + 2 + 2 + 1 + 3, // sats + oxygen + BP + pulse + consciousness
		},
		{
			name: "Critical values",
			vitals: VitalSigns{
				RespRate:           7,
				OxygenSat:          80,
				SystolicBP:         85,
				Pulse:              140,
				Temp:               34.5,
				ConsciousnessLevel: Unresponsive,
				OnOxygen:           true,
			},
			scale1:   true,
			expected: 3 + 3 + 2 + 3 + 3 + 3 + 3, // resp + sats + oxygen + BP + pulse + consciousness + temp
		},
		{
			name: "All normal, scale2, high sats oxygen",
			vitals: VitalSigns{
				RespRate:           16,
				OxygenSat:          100,
				SystolicBP:         120,
				Pulse:              80,
				Temp:               37.0,
				ConsciousnessLevel: Alert,
				OnOxygen:           true,
			},
			scale1:   false,
			expected: 3 + 2, // sats + oxygen
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateScore(tt.vitals, tt.scale1)
			if got != tt.expected {
				t.Errorf("CalculateScore() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestConsciousnessLevelString(t *testing.T) {
	tests := []struct {
		level    ConsciousnessLevel
		expected string
	}{
		{Alert, "Alert"},
		{Confused, "Confused"},
		{Voice, "Voice"},
		{Pain, "Pain"},
		{Unresponsive, "Unresponsive"},
		{ConsciousnessLevel(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.level.String()
		if got != tt.expected {
			t.Errorf("ConsciousnessLevel.String() = %s, want %s", got, tt.expected)
		}
	}
}

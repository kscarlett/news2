package news2

import (
	"errors"
	"testing"
)

// normalVitals returns a fully-populated set of observations that scores 0.
func normalVitals() VitalSigns {
	return VitalSigns{
		RespRate:           16,
		OxygenSat:          98,
		SystolicBP:         120,
		Pulse:              80,
		Temp:               37.0,
		ConsciousnessLevel: Alert,
		OnOxygen:           false,
		SpO2Scale:          Scale1,
	}
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		vitals   VitalSigns
		expected int
	}{
		{
			name:     "All normal, scale1, no oxygen",
			vitals:   normalVitals(),
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
				SpO2Scale:          Scale1,
			},
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
				SpO2Scale:          Scale1,
			},
			expected: 3 + 3 + 2 + 3 + 3 + 3 + 3, // resp + sats + oxygen + BP + pulse + consciousness + temp
		},
		{
			name: "Raised resp rate scores 2",
			vitals: func() VitalSigns {
				v := normalVitals()
				v.RespRate = 22
				return v
			}(),
			expected: 2,
		},
		{
			name: "Scale2, high sats on oxygen",
			vitals: func() VitalSigns {
				v := normalVitals()
				v.OxygenSat = 100
				v.OnOxygen = true
				v.SpO2Scale = Scale2
				return v
			}(),
			expected: 3 + 2, // sats + oxygen
		},
		{
			name: "Scale2, target sats on air score 0",
			vitals: func() VitalSigns {
				v := normalVitals()
				v.OxygenSat = 90
				v.SpO2Scale = Scale2
				return v
			}(),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateScore(tt.vitals)
			if err != nil {
				t.Fatalf("CalculateScore() unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("CalculateScore() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestCalculateSubscores(t *testing.T) {
	vitals := VitalSigns{
		RespRate:           22,       // 2
		OxygenSat:          94,       // 1 (scale 1)
		SystolicBP:         105,      // 1
		Pulse:              45,       // 1
		Temp:               35.5,     // 1
		ConsciousnessLevel: Confused, // 3
		OnOxygen:           true,     // 2
		SpO2Scale:          Scale1,
	}
	r, err := Calculate(vitals)
	if err != nil {
		t.Fatalf("Calculate() unexpected error: %v", err)
	}
	want := Result{
		Total:          11,
		RespRate:       2,
		Saturations:    1,
		SupplementalO2: 2,
		SystolicBP:     1,
		Pulse:          1,
		Consciousness:  3,
		Temperature:    1,
		RedFlag:        true,
		Risk:           High,
	}
	if r != want {
		t.Errorf("Calculate() = %+v, want %+v", r, want)
	}
}

func TestCalculateRiskLevels(t *testing.T) {
	tests := []struct {
		name        string
		modify      func(*VitalSigns)
		wantTotal   int
		wantRedFlag bool
		wantRisk    RiskLevel
	}{
		{
			name:     "score 0 is low risk",
			modify:   func(v *VitalSigns) {},
			wantRisk: Low,
		},
		{
			name: "score 1-4 without red flag is low risk",
			modify: func(v *VitalSigns) {
				v.Pulse = 95    // 1
				v.Temp = 38.5   // 1
				v.RespRate = 10 // 1
			},
			wantTotal: 3,
			wantRisk:  Low,
		},
		{
			name: "single parameter scoring 3 is low-medium risk even with low total",
			modify: func(v *VitalSigns) {
				v.Pulse = 35 // 3
			},
			wantTotal:   3,
			wantRedFlag: true,
			wantRisk:    LowMedium,
		},
		{
			name: "score 5-6 is medium risk",
			modify: func(v *VitalSigns) {
				v.RespRate = 22   // 2
				v.SystolicBP = 95 // 2
				v.Pulse = 95      // 1
			},
			wantTotal: 5,
			wantRisk:  Medium,
		},
		{
			name: "score 7+ is high risk",
			modify: func(v *VitalSigns) {
				v.RespRate = 26   // 3
				v.SystolicBP = 95 // 2
				v.Pulse = 115     // 2
			},
			wantTotal:   7,
			wantRedFlag: true,
			wantRisk:    High,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vitals := normalVitals()
			tt.modify(&vitals)
			r, err := Calculate(vitals)
			if err != nil {
				t.Fatalf("Calculate() unexpected error: %v", err)
			}
			if r.Total != tt.wantTotal {
				t.Errorf("Total = %d, want %d", r.Total, tt.wantTotal)
			}
			if r.RedFlag != tt.wantRedFlag {
				t.Errorf("RedFlag = %t, want %t", r.RedFlag, tt.wantRedFlag)
			}
			if r.Risk != tt.wantRisk {
				t.Errorf("Risk = %s, want %s", r.Risk, tt.wantRisk)
			}
		})
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*VitalSigns)
	}{
		{"zero value struct", func(v *VitalSigns) { *v = VitalSigns{} }},
		{"resp rate unset", func(v *VitalSigns) { v.RespRate = 0 }},
		{"resp rate negative", func(v *VitalSigns) { v.RespRate = -1 }},
		{"resp rate beyond sane bound", func(v *VitalSigns) { v.RespRate = 101 }},
		{"oxygen sat unset", func(v *VitalSigns) { v.OxygenSat = 0 }},
		{"oxygen sat above 100", func(v *VitalSigns) { v.OxygenSat = 101 }},
		{"systolic BP unset", func(v *VitalSigns) { v.SystolicBP = 0 }},
		{"systolic BP beyond sane bound", func(v *VitalSigns) { v.SystolicBP = 301 }},
		{"pulse unset", func(v *VitalSigns) { v.Pulse = 0 }},
		{"pulse beyond sane bound", func(v *VitalSigns) { v.Pulse = 301 }},
		{"temperature unset", func(v *VitalSigns) { v.Temp = 0 }},
		{"temperature negative", func(v *VitalSigns) { v.Temp = -1.5 }},
		{"temperature beyond sane bound", func(v *VitalSigns) { v.Temp = 60.1 }},
		{"unknown consciousness level", func(v *VitalSigns) { v.ConsciousnessLevel = ConsciousnessLevel(99) }},
		{"negative consciousness level", func(v *VitalSigns) { v.ConsciousnessLevel = ConsciousnessLevel(-1) }},
		{"unknown SpO2 scale", func(v *VitalSigns) { v.SpO2Scale = SpO2Scale(3) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vitals := normalVitals()
			tt.modify(&vitals)

			if _, err := Calculate(vitals); !errors.Is(err, ErrInvalidVitalSigns) {
				t.Errorf("Calculate() error = %v, want ErrInvalidVitalSigns", err)
			}
			if _, err := CalculateScore(vitals); !errors.Is(err, ErrInvalidVitalSigns) {
				t.Errorf("CalculateScore() error = %v, want ErrInvalidVitalSigns", err)
			}
		})
	}
}

func TestValidVitalsPassValidation(t *testing.T) {
	if err := normalVitals().Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

// Validation is deliberately minimal: clinically implausible but
// mechanically scoreable values are accepted, since range enforcement is
// the calling software's responsibility.
func TestValidationAcceptsImplausibleValues(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*VitalSigns)
	}{
		{"very high resp rate", func(v *VitalSigns) { v.RespRate = 90 }},
		{"very high systolic BP", func(v *VitalSigns) { v.SystolicBP = 290 }},
		{"very high pulse", func(v *VitalSigns) { v.Pulse = 290 }},
		{"very high temperature", func(v *VitalSigns) { v.Temp = 45.0 }},
		{"very low oxygen sat", func(v *VitalSigns) { v.OxygenSat = 1 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vitals := normalVitals()
			tt.modify(&vitals)
			if err := vitals.Validate(); err != nil {
				t.Errorf("Validate() = %v, want nil", err)
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

func TestRiskLevelString(t *testing.T) {
	tests := []struct {
		level    RiskLevel
		expected string
	}{
		{Low, "Low"},
		{LowMedium, "Low-Medium"},
		{Medium, "Medium"},
		{High, "High"},
		{RiskLevel(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.level.String()
		if got != tt.expected {
			t.Errorf("RiskLevel.String() = %s, want %s", got, tt.expected)
		}
	}
}

func TestSpO2ScaleString(t *testing.T) {
	tests := []struct {
		scale    SpO2Scale
		expected string
	}{
		{Scale1, "Scale 1"},
		{Scale2, "Scale 2"},
		{SpO2Scale(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.scale.String()
		if got != tt.expected {
			t.Errorf("SpO2Scale.String() = %s, want %s", got, tt.expected)
		}
	}
}

func FuzzCalculate(f *testing.F) {
	f.Add(16, 98, 120, 80, 37.0, 0, false, 0)
	f.Add(7, 80, 85, 140, 34.5, 4, true, 0)
	f.Add(22, 90, 100, 100, 38.5, 1, true, 1)

	f.Fuzz(func(t *testing.T, respRate, oxygenSat, systolicBP, pulse int, temp float64, level int, onOxygen bool, scale int) {
		vitals := VitalSigns{
			RespRate:           respRate,
			OxygenSat:          oxygenSat,
			SystolicBP:         systolicBP,
			Pulse:              pulse,
			Temp:               temp,
			ConsciousnessLevel: ConsciousnessLevel(level),
			OnOxygen:           onOxygen,
			SpO2Scale:          SpO2Scale(scale),
		}
		r, err := Calculate(vitals)
		if err != nil {
			if !errors.Is(err, ErrInvalidVitalSigns) {
				t.Errorf("Calculate() unexpected error type: %v", err)
			}
			return
		}
		if r.Total < 0 || r.Total > 20 {
			t.Errorf("Total = %d, outside possible NEWS2 range [0, 20]", r.Total)
		}
		sum := r.RespRate + r.Saturations + r.SupplementalO2 + r.SystolicBP +
			r.Pulse + r.Consciousness + r.Temperature
		if sum != r.Total {
			t.Errorf("Total = %d, but subscores sum to %d", r.Total, sum)
		}
		if r.SupplementalO2 != 0 && r.SupplementalO2 != 2 {
			t.Errorf("SupplementalO2 = %d, want 0 or 2", r.SupplementalO2)
		}
	})
}

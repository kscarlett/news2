package news2

import (
	"strings"
	"testing"
)

func TestResponseForEachRiskLevel(t *testing.T) {
	tests := []struct {
		risk          RiskLevel
		wantFrequency string
	}{
		{Low, "minimum 12 hourly"},
		{LowMedium, "minimum 1 hourly"},
		{Medium, "minimum 1 hourly"},
		{High, "continuous monitoring of vital signs"},
	}
	for _, tt := range tests {
		t.Run(tt.risk.String(), func(t *testing.T) {
			got := ResponseFor(tt.risk)
			if got.MonitoringFrequency != tt.wantFrequency {
				t.Errorf("MonitoringFrequency = %q, want %q", got.MonitoringFrequency, tt.wantFrequency)
			}
			if got.ClinicalResponse == "" {
				t.Errorf("ClinicalResponse is empty for %s", tt.risk)
			}
		})
	}
}

func TestResponseForUnknownRisk(t *testing.T) {
	if got := ResponseFor(RiskLevel(99)); got != (Response{}) {
		t.Errorf("ResponseFor(unknown) = %+v, want zero Response", got)
	}
}

func TestResultResponseMatchesRisk(t *testing.T) {
	vitals := normalVitals()
	vitals.RespRate = 26   // 3
	vitals.SystolicBP = 95 // 2
	vitals.Pulse = 115     // 2, total 7 -> High
	r, err := Calculate(vitals)
	if err != nil {
		t.Fatalf("Calculate() unexpected error: %v", err)
	}
	if r.Risk != High {
		t.Fatalf("Risk = %s, want High", r.Risk)
	}
	if r.Response() != ResponseFor(High) {
		t.Errorf("Result.Response() = %+v, want %+v", r.Response(), ResponseFor(High))
	}
}

func TestResponseString(t *testing.T) {
	got := ResponseFor(Low).String()
	if !strings.Contains(got, "minimum 12 hourly") {
		t.Errorf("Response.String() = %q, missing monitoring frequency", got)
	}
}

package news2

import (
	"errors"
	"fmt"
)

// ErrInvalidVitalSigns is wrapped by all validation errors returned from
// Calculate, CalculateScore and VitalSigns.Validate, so callers can detect
// them with errors.Is.
var ErrInvalidVitalSigns = errors.New("invalid vital signs")

// Validate checks that the VitalSigns can be meaningfully scored. It is
// deliberately minimal: it rejects only values that make scoring nonsense —
// non-positive measurements (which also catches unset fields, whose zero
// values would otherwise silently score as a critically ill patient), an
// oxygen saturation above 100%, and unknown enum values. Enforcing
// clinically plausible ranges is left to the calling software.
//
// It returns nil if the vital signs can be scored, or an error wrapping
// ErrInvalidVitalSigns describing the first invalid field found.
func (v VitalSigns) Validate() error {
	if v.RespRate < 1 {
		return fmt.Errorf("%w: respiration rate must be positive, got %d", ErrInvalidVitalSigns, v.RespRate)
	}
	if v.OxygenSat < 1 || v.OxygenSat > 100 {
		return fmt.Errorf("%w: oxygen saturation must be a percentage in [1, 100], got %d", ErrInvalidVitalSigns, v.OxygenSat)
	}
	if v.SystolicBP < 1 {
		return fmt.Errorf("%w: systolic blood pressure must be positive, got %d", ErrInvalidVitalSigns, v.SystolicBP)
	}
	if v.Pulse < 1 {
		return fmt.Errorf("%w: pulse must be positive, got %d", ErrInvalidVitalSigns, v.Pulse)
	}
	if v.Temp <= 0 {
		return fmt.Errorf("%w: temperature must be positive, got %.1f°C", ErrInvalidVitalSigns, v.Temp)
	}
	if v.ConsciousnessLevel < Alert || v.ConsciousnessLevel > Unresponsive {
		return fmt.Errorf("%w: unknown consciousness level %d", ErrInvalidVitalSigns, int(v.ConsciousnessLevel))
	}
	if v.SpO2Scale != Scale1 && v.SpO2Scale != Scale2 {
		return fmt.Errorf("%w: unknown SpO2 scale %d", ErrInvalidVitalSigns, int(v.SpO2Scale))
	}
	return nil
}

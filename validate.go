package news2

import (
	"errors"
	"fmt"
)

// ErrInvalidVitalSigns is wrapped by all validation errors returned from
// Calculate, CalculateScore and VitalSigns.Validate, so callers can detect
// them with errors.Is.
var ErrInvalidVitalSigns = errors.New("invalid vital signs")

// Sanity bounds. These are NOT clinical ranges — enforcing plausible ranges
// is the calling software's job. They are deliberately extreme outer limits
// that no real reading could reach, so their only purpose is to reject
// impossible values such as unset fields, garbage data entry and unit
// mix-ups, while accepting anything that could conceivably be a real
// observation.
const (
	maxRespRate   = 100  // breaths/min
	maxSystolicBP = 300  // mmHg
	maxPulse      = 300  // BPM
	maxTemp       = 60.0 // °C
)

// Validate checks that the VitalSigns can be meaningfully scored. It is
// deliberately minimal: it rejects only values that cannot be a real
// observation — non-positive measurements (which also catches unset fields,
// whose zero values would otherwise silently score as a critically ill
// patient), an oxygen saturation above 100%, values beyond the extreme
// sanity bounds above, and unknown enum values. Enforcing clinically
// plausible ranges is left to the calling software.
//
// It returns nil if the vital signs can be scored, or an error wrapping
// ErrInvalidVitalSigns describing the first invalid field found.
func (v VitalSigns) Validate() error {
	if v.RespRate < 1 || v.RespRate > maxRespRate {
		return fmt.Errorf("%w: respiration rate %d breaths/min outside sane range [1, %d]", ErrInvalidVitalSigns, v.RespRate, maxRespRate)
	}
	if v.OxygenSat < 1 || v.OxygenSat > 100 {
		return fmt.Errorf("%w: oxygen saturation must be a percentage in [1, 100], got %d", ErrInvalidVitalSigns, v.OxygenSat)
	}
	if v.SystolicBP < 1 || v.SystolicBP > maxSystolicBP {
		return fmt.Errorf("%w: systolic blood pressure %d mmHg outside sane range [1, %d]", ErrInvalidVitalSigns, v.SystolicBP, maxSystolicBP)
	}
	if v.Pulse < 1 || v.Pulse > maxPulse {
		return fmt.Errorf("%w: pulse %d BPM outside sane range [1, %d]", ErrInvalidVitalSigns, v.Pulse, maxPulse)
	}
	if v.Temp <= 0 || v.Temp > maxTemp {
		return fmt.Errorf("%w: temperature %.1f°C outside sane range (0, %.1f]", ErrInvalidVitalSigns, v.Temp, maxTemp)
	}
	if v.ConsciousnessLevel < Alert || v.ConsciousnessLevel > Unresponsive {
		return fmt.Errorf("%w: unknown consciousness level %d", ErrInvalidVitalSigns, int(v.ConsciousnessLevel))
	}
	if v.SpO2Scale != Scale1 && v.SpO2Scale != Scale2 {
		return fmt.Errorf("%w: unknown SpO2 scale %d", ErrInvalidVitalSigns, int(v.SpO2Scale))
	}
	return nil
}

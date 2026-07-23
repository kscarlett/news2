package news2

import (
	"errors"
	"fmt"
)

// ErrInvalidVitalSigns is wrapped by all validation errors returned from
// Calculate, CalculateScore and VitalSigns.Validate, so callers can detect
// them with errors.Is.
var ErrInvalidVitalSigns = errors.New("invalid vital signs")

// Validation bounds. These are deliberately generous — they are plausibility
// limits on what a real observation could read, not normal ranges. Their
// main purpose is to reject unset fields (a zero value would otherwise
// silently score as a critically ill patient) and unit mix-ups.
const (
	minRespRate   = 1
	maxRespRate   = 100
	minOxygenSat  = 1
	maxOxygenSat  = 100
	minSystolicBP = 10
	maxSystolicBP = 350
	minPulse      = 1
	maxPulse      = 350
	minTemp       = 25.0
	maxTemp       = 45.0
)

// Validate checks that every field of the VitalSigns is within a
// physiologically plausible range and that enum fields hold known values.
// It returns nil if the vital signs can be scored, or an error wrapping
// ErrInvalidVitalSigns describing the first invalid field found.
func (v VitalSigns) Validate() error {
	if v.RespRate < minRespRate || v.RespRate > maxRespRate {
		return fmt.Errorf("%w: respiration rate %d breaths/min outside plausible range [%d, %d]",
			ErrInvalidVitalSigns, v.RespRate, minRespRate, maxRespRate)
	}
	if v.OxygenSat < minOxygenSat || v.OxygenSat > maxOxygenSat {
		return fmt.Errorf("%w: oxygen saturation %d%% outside plausible range [%d, %d]",
			ErrInvalidVitalSigns, v.OxygenSat, minOxygenSat, maxOxygenSat)
	}
	if v.SystolicBP < minSystolicBP || v.SystolicBP > maxSystolicBP {
		return fmt.Errorf("%w: systolic blood pressure %d mmHg outside plausible range [%d, %d]",
			ErrInvalidVitalSigns, v.SystolicBP, minSystolicBP, maxSystolicBP)
	}
	if v.Pulse < minPulse || v.Pulse > maxPulse {
		return fmt.Errorf("%w: pulse %d BPM outside plausible range [%d, %d]",
			ErrInvalidVitalSigns, v.Pulse, minPulse, maxPulse)
	}
	if v.Temp < minTemp || v.Temp > maxTemp {
		return fmt.Errorf("%w: temperature %.1f°C outside plausible range [%.1f, %.1f] (is it in Celsius?)",
			ErrInvalidVitalSigns, v.Temp, minTemp, maxTemp)
	}
	if v.ConsciousnessLevel < Alert || v.ConsciousnessLevel > Unresponsive {
		return fmt.Errorf("%w: unknown consciousness level %d",
			ErrInvalidVitalSigns, int(v.ConsciousnessLevel))
	}
	if v.SpO2Scale != Scale1 && v.SpO2Scale != Scale2 {
		return fmt.Errorf("%w: unknown SpO2 scale %d",
			ErrInvalidVitalSigns, int(v.SpO2Scale))
	}
	return nil
}

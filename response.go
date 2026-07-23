package news2

import "fmt"

// Response describes the clinical response recommended by the Royal College
// of Physicians for a NEWS2 result: the minimum frequency of monitoring and
// the escalation of care.
//
// The text is a plain-language summary of the RCP "Clinical response to the
// NEWS2 trigger thresholds" chart (Chart 4), keyed on the clinical risk
// category. It is provided for information only. Do not rely on it to direct
// patient care: it is not a medical device, may be summarised, abbreviated
// or out of date relative to current RCP guidance and local escalation
// policy, and the definitive source is the RCP publication itself. Local
// protocols always take precedence.
type Response struct {
	// MonitoringFrequency is the recommended minimum frequency of
	// observations.
	MonitoringFrequency string
	// ClinicalResponse summarises the recommended escalation of care.
	ClinicalResponse string
}

// String returns a human-readable summary of the response.
func (r Response) String() string {
	return fmt.Sprintf("Monitor %s. %s", r.MonitoringFrequency, r.ClinicalResponse)
}

// Response returns the RCP-recommended clinical response for the result's
// clinical risk category. See ResponseFor and the Response type for the
// important limitations on this information.
func (r Result) Response() Response {
	return ResponseFor(r.Risk)
}

// ResponseFor returns the RCP-recommended clinical response for a given
// clinical risk category, as summarised from the RCP "Clinical response to
// the NEWS2 trigger thresholds" chart (Chart 4).
//
// This information is provided for reference only and is not a substitute
// for the RCP guidance or local escalation policy. See the Response type for
// the full disclaimer. An unrecognised RiskLevel returns a zero Response.
func ResponseFor(risk RiskLevel) Response {
	switch risk {
	case Low:
		return Response{
			MonitoringFrequency: "minimum 12 hourly",
			ClinicalResponse:    "Continue routine NEWS2 monitoring.",
		}
	case LowMedium:
		return Response{
			MonitoringFrequency: "minimum 1 hourly",
			ClinicalResponse:    "Registered nurse to inform the medical team caring for the patient, who will review and decide whether escalation of care is necessary.",
		}
	case Medium:
		return Response{
			MonitoringFrequency: "minimum 1 hourly",
			ClinicalResponse:    "Registered nurse to immediately inform the medical team and request an urgent assessment by a clinician or team with competencies in the care of acutely ill patients. Provide clinical care in an environment with monitoring facilities.",
		}
	case High:
		return Response{
			MonitoringFrequency: "continuous monitoring of vital signs",
			ClinicalResponse:    "Registered nurse to immediately inform the medical team (at least specialist registrar level) and arrange emergency assessment by a team with critical care competencies, including advanced airway management. Consider transfer to a level 2 or 3 (higher-dependency or intensive care) facility.",
		}
	}
	return Response{}
}

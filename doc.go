// Package news2 calculates the National Early Warning Score 2 (NEWS2), the
// Royal College of Physicians' standardised score for detecting acute
// clinical deterioration in adult patients.
//
// The score aggregates six physiological parameters — respiration rate,
// oxygen saturation, systolic blood pressure, pulse, level of consciousness
// (ACVPU) and temperature — plus whether the patient is on supplemental
// oxygen. Calculate returns the full result including per-parameter
// subscores and the clinical risk category that drives escalation;
// CalculateScore returns just the aggregate score.
//
// Scoring follows the RCP publication "National Early Warning Score (NEWS)
// 2: Standardising the assessment of acute-illness severity in the NHS"
// (2017). NEWS2 applies to adults aged 16 or over and is not validated for
// use in pregnancy. SpO2 Scale 2 must only be selected for patients with
// confirmed hypercapnic respiratory failure, on the direction of a
// competent clinical decision maker.
//
// Disclaimer: this package is not a medical device and is not designed,
// certified or warranted for use in clinical care or in any system where
// reliability or safe operation is required. It was written to support
// creation and management of simulation cases. No guarantee is provided of
// its stability, reliability or accuracy.
package news2

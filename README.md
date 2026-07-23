# NEWS2

[![CI](https://github.com/kscarlett/news2/actions/workflows/ci.yml/badge.svg)](https://github.com/kscarlett/news2/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kscarlett/news2.svg)](https://pkg.go.dev/github.com/kscarlett/news2)

A Go package for calculating the National Early Warning Score 2 (NEWS2) from patient vital signs, including the per-parameter breakdown and the clinical risk category used to drive escalation.

> **NOTE:** This package is not a medical device and was not designed for use in any system where reliability or safe operation is needed. It is used as part of my workflow to create and manage simulation cases. I cannot provide any guarantee of the stability, reliability or accuracy of this code.

## Install

```sh
go get github.com/kscarlett/news2
```

## Usage

```go
import "github.com/kscarlett/news2"

vitals := news2.VitalSigns{
    RespRate:           22,          // breaths per minute
    OxygenSat:          93,          // SpO2 %
    SystolicBP:         98,          // mmHg
    Pulse:              112,         // beats per minute
    Temp:               38.4,        // °C
    ConsciousnessLevel: news2.Alert, // ACVPU
    OnOxygen:           false,       // receiving supplemental oxygen
    // SpO2Scale defaults to news2.Scale1, which is correct for most
    // patients. Only set news2.Scale2 for patients with confirmed
    // hypercapnic respiratory failure (target sats 88–92%), as
    // prescribed by a competent clinical decision maker.
}

result, err := news2.Calculate(vitals)
if err != nil {
    // Inputs are validated: physiologically implausible values (including
    // the zero value of VitalSigns) return an error wrapping
    // news2.ErrInvalidVitalSigns instead of a misleading score.
    log.Fatal(err)
}

fmt.Println(result.Total)   // 9
fmt.Println(result.Risk)    // High
fmt.Println(result.RedFlag) // false — no single parameter scored 3
fmt.Println(result.Pulse)   // 2 — each parameter's subscore is available
```

If you only need the aggregate number:

```go
score, err := news2.CalculateScore(vitals)
```

## Interpreting the score

NEWS2 is a trigger system: the aggregate score (and a "red score" of 3 in any single parameter) maps to a clinical risk category, a minimum observation frequency, and an escalation response. `Result.Risk` and `Result.RedFlag` encode this mapping:

| Trigger | `Result.Risk` | Response (per RCP guidance) |
| --- | --- | --- |
| Score 0 | Low | Routine monitoring (minimum 12-hourly) |
| Score 1–4 | Low | 4–6 hourly observations; registered nurse decides on escalation |
| Score of 3 in any single parameter | Low-Medium | Urgent ward-based review; minimum hourly observations |
| Score 5–6 | Medium | Urgent review by clinician competent in acute illness; minimum hourly observations |
| Score ≥ 7 | High | Emergency response, usually with critical care involvement; continuous monitoring |

A score of 5 or more is also a common threshold to screen for sepsis.

## Scope

- NEWS2 applies to adults aged 16 and over. It is not validated for children or for use in pregnancy.
- `Confused` on the ACVPU scale means **new-onset** confusion (or worse than baseline). Patients with chronic baseline confusion are recorded as `Alert`.
- SpO2 `Scale2` is only for patients with confirmed hypercapnic respiratory failure; the decision to use it is clinical and should be recorded in the patient's notes.

## Reference

This package implements the NEWS2 scoring system as described by the Royal College of Physicians (2017; the December 2022 update left the scoring unchanged).
**Source:** Royal College of Physicians. _National Early Warning Score (NEWS) 2: Standardising the assessment of acute-illness severity in the NHS._ London: RCP, 2017.

## License

MIT License

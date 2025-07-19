# NEWS2

A Go package for calculating National Early Warning Score 2 (NEWS2) based on patient vital signs.

## Features

- Calculate NEWS2 scores from respiratory rate, oxygen saturation, systolic blood pressure, pulse, temperature, consciousness level, and oxygen supplementation.

> NOTE: This was not designed for use in any systems where any reliability or safe operation is needed. This package is used as part of my workflow to create and manage simulation cases. I cannot provide any guarantee of the stability, reliability or accuracy of this code.

## Usage

```go
import "github.com/kscarlett/news2"

vitals := news2.VitalSigns{
    RespRate:           18,
    OxygenSat:          95,
    SystolicBP:         120,
    Pulse:              80,
    Temp:               36.8,
    ConsciousnessLevel: news2.Alert,
    OnOxygen:           false,
}

score := news2.CalculateScore(vitals, true) // true for Scale 1, false for Scale 2
```

## Reference

This package implements the NEWS2 scoring system as described by the Royal College of Physicians (2017).  
**Source:** Royal College of Physicians. _National Early Warning Score (NEWS) 2: Standardising the assessment of acute-illness severity in the NHS._ London: RCP, 2017.

## License

MIT License
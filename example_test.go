package news2_test

import (
	"errors"
	"fmt"

	"github.com/kscarlett/news2"
)

func ExampleCalculate() {
	vitals := news2.VitalSigns{
		RespRate:           22,
		OxygenSat:          93,
		SystolicBP:         98,
		Pulse:              112,
		Temp:               38.4,
		ConsciousnessLevel: news2.Alert,
		OnOxygen:           false,
		// SpO2Scale defaults to news2.Scale1.
	}

	result, err := news2.Calculate(vitals)
	if err != nil {
		fmt.Println("invalid observations:", err)
		return
	}

	fmt.Println("Total:", result.Total)
	fmt.Println("Risk:", result.Risk)
	fmt.Println("Red flag:", result.RedFlag)
	// Output:
	// Total: 9
	// Risk: High
	// Red flag: false
}

func ExampleCalculateScore() {
	vitals := news2.VitalSigns{
		RespRate:           18,
		OxygenSat:          95,
		SystolicBP:         120,
		Pulse:              80,
		Temp:               36.8,
		ConsciousnessLevel: news2.Alert,
		OnOxygen:           false,
	}

	score, err := news2.CalculateScore(vitals)
	if err != nil {
		fmt.Println("invalid observations:", err)
		return
	}

	fmt.Println("NEWS2 score:", score)
	// Output:
	// NEWS2 score: 1
}

func ExampleVitalSigns_Validate() {
	// The zero value is rejected rather than silently scored.
	var vitals news2.VitalSigns

	_, err := news2.Calculate(vitals)
	fmt.Println(errors.Is(err, news2.ErrInvalidVitalSigns))
	// Output:
	// true
}

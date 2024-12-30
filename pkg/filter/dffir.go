package filter

import (
	"encoding/json"
	"fmt"
)

// dffir is a direct form FIR filter structure
type dffir struct {
	Arithmetic          string    `json:"Arithmetic"`
	Numerator           []float64 `json:"Numerator"`
	FilterStructure     string    `json:"FilterStructure"`
	States              []int     `json:"States"`
	NumSamplesProcessed int       `json:"NumSamplesProcessed"`
	PersistentMemory    bool      `json:"PersistentMemory"`
	RateChangeFactor    []int     `json:"RateChangeFactor"`
}

func (f *dffir) Filter(signal []float64) ([]float64, error) {
	return filterUsingDirectFormFIR(signal, f.Numerator)
}

func (f *dffir) unmarshal(js []byte) error {
	if err := json.Unmarshal(js, f); err != nil {
		return err
	}

	if f.FilterStructure != "Direct-Form FIR" {
		return fmt.Errorf("invalid filter structure: %s", f.FilterStructure)
	}

	return nil
}

func filterUsingDirectFormFIR(signal []float64, numerator []float64) ([]float64, error) {
	if signal == nil || len(signal) == 0 {
		return nil, fmt.Errorf("signal cannot be nil or empty")
	}
	if numerator == nil || len(numerator) == 0 {
		return nil, fmt.Errorf("numerator coefficients cannot be nil or empty")
	}

	filterLength := len(numerator)
	output := make([]float64, len(signal))

	for i := 0; i < len(signal); i++ {
		var sum float64
		for j := 0; j < filterLength; j++ {
			signalIndex := i - j
			if signalIndex >= 0 {
				sum += signal[signalIndex] * numerator[j]
			}
		}
		output[i] = sum
	}

	return output, nil
}

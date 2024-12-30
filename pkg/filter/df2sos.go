package filter

import (
	"encoding/json"
	"fmt"
)

// df2sos is a direct form II SOS filter structure
type df2sos struct {
	OptimizeScaleValues bool        `json:"OptimizeScaleValues"`
	Arithmetic          string      `json:"Arithmetic"`
	SosMatrix           [][]float64 `json:"sosMatrix"`
	ScaleValues         []float64   `json:"ScaleValues"`
	FilterStructure     string      `json:"FilterStructure"`
	States              [][]int     `json:"States"`
	NumSamplesProcessed int         `json:"NumSamplesProcessed"`
	PersistentMemory    bool        `json:"PersistentMemory"`
	RateChangeFactor    []int       `json:"RateChangeFactor"`
}

func (f *df2sos) unmarshal(js []byte) error {
	if err := json.Unmarshal(js, f); err != nil {
		return err
	}

	if f.FilterStructure != "Direct-Form II, Second-Order Sections" {
		return fmt.Errorf("invalid filter structure: %s", f.FilterStructure)
	}

	return nil
}

func (f *df2sos) Filter(signal []float64) ([]float64, error) {
	return filterUsingDirectFormIISOS(signal, f.SosMatrix, f.ScaleValues)
}

// filterUsingDirectFormIISOS applies a Second Order Sections (SOS) filter to a signal using the Direct Form II Transposed structure.
//
// signal: The input signal.
// sosMatrix: The SOS matrix. Each row represents a second-order section. The columns are [b0, b1, b2, a0, a1, a2].
// scaleValues: Scaling factors for each section. These may not always be necessary.  If unused, provide a nil or empty slice.
//
// Returns: The filtered signal. Returns an error if the signal is empty, the SOS matrix is invalid,
// or if scaleValues is provided and its length doesn't match the number of sections.
func filterUsingDirectFormIISOS(signal []float64, sosMatrix [][]float64, scaleValues []float64) ([]float64, error) {
	if len(signal) == 0 {
		return nil, fmt.Errorf("signal cannot be empty")
	}
	if len(sosMatrix) == 0 {
		return nil, fmt.Errorf("the SOS matrix cannot be empty")
	}
	if len(scaleValues) != 0 && len(scaleValues) != len(sosMatrix)+1 {
		return nil, fmt.Errorf("length of scaleValues must match number of SOS sections + 1 or be 0")
	}

	numSections := len(sosMatrix)
	output := make([]float64, len(signal))

	// Section states (Direct Form II Transposed needs two states per section)
	w := make([][]float64, numSections)
	for i := range w {
		w[i] = make([]float64, 2) // Initialize states to 0
	}

	for n := 0; n < len(signal); n++ {
		x := signal[n]
		for i := 0; i < numSections; i++ {
			section := sosMatrix[i]
			if len(section) != 6 {
				return nil, fmt.Errorf("invalid SOS matrix: each row must have 6 coefficients [b0, b1, b2, a0, a1, a2]")
			}

			b0, b1, b2 := section[0], section[1], section[2]
			a0, a1, a2 := section[3], section[4], section[5]

			scale := 1.0
			if len(scaleValues) > 0 {
				scale = scaleValues[i]
			}

			// Direct Form II Transposed
			w0 := x - a1/a0*w[i][0] - a2/a0*w[i][1]
			y := (b0/a0)*w0 + (b1/a0)*w[i][0] + (b2/a0)*w[i][1]

			w[i][1] = w[i][0]
			w[i][0] = w0

			x = y * scale
		}
		output[n] = x * scaleValues[numSections] // Apply the final overall gain
	}

	return output, nil
}

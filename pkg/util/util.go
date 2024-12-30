package util

import (
	"encoding/json"
	"math"
)

const (
	tolerance = 1e-10
)

// GetSignalFromJSON fetches a signal out of its JSON representation
func GetSignalFromJSON(s []byte) ([]float64, error) {
	var out []float64
	if err := json.Unmarshal(s, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IsAlmostEqual compares two signals for equality within a tolerance value
func IsAlmostEqual(x1, x2 []float64) bool {
	if len(x1) != len(x2) {
		return false
	}

	for i := range x1 {
		if math.Abs(x1[i]-x2[i]) > tolerance {
			return false
		}
	}

	return true
}

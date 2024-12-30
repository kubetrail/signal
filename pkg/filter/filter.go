package filter

// Filter filters a signal
type Filter interface {
	Filter(s []float64) ([]float64, error)
}

// NewDf2sosFromJSON creates a new direct form II SOS filter object from the input
// json serialized version.
func NewDf2sosFromJSON(js []byte) (Filter, error) {
	f := new(df2sos)
	if err := f.unmarshal(js); err != nil {
		return nil, err
	}
	return f, nil
}

// NewDffirFromJSON creates a new direct form FIR filter object from the input
// json serialized version.
func NewDffirFromJSON(js []byte) (Filter, error) {
	f := new(dffir)
	if err := f.unmarshal(js); err != nil {
		return nil, err
	}

	return f, nil
}

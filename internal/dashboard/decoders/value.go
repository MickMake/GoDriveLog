package decoders

import "fmt"

const (
	ValueTypeNumber     = "number"
	ValueTypeText       = "text"
	ValueTypeBoolean    = "boolean"
	ValueTypeDigits     = "digits"
	ValueTypeFrameIndex = "frame_index"
)

type Value struct {
	Type       string
	Number     float64
	Text       string
	Bool       bool
	Digits     []string
	FrameIndex int
	Unit       string
	Status     string
	Error      string
	Min        float64
	Max        float64
	HasRange   bool
}

func (v Value) NumberValue() (float64, error) {
	switch v.Type {
	case ValueTypeNumber, ValueTypeFrameIndex:
		return v.Number, nil
	case ValueTypeBoolean:
		if v.Bool {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("value type %q is not numeric", v.Type)
	}
}

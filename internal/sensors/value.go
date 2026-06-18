package sensors

import (
	"fmt"
	"strings"
)

const (
	ValueKindNumeric = "numeric"
	ValueKindBool    = "bool"
	ValueKindString  = "string"
	ValueKindMissing = "missing"
	ValueKindError   = "error"
)

// Value is the explicit typed payload carried by v3 sensor state/events.
// Kind is mandatory. Consumers must not infer type from populated fields.
type Value struct {
	Kind    string   `json:"kind"`
	Number  *float64 `json:"number,omitempty"`
	Bool    *bool    `json:"bool,omitempty"`
	String  string   `json:"string,omitempty"`
	Unit    string   `json:"unit,omitempty"`
	Message string   `json:"message,omitempty"`
}

func NewNumericValue(number float64, unit string) Value {
	return Value{Kind: ValueKindNumeric, Number: &number, Unit: unit}
}

func NewBoolValue(value bool) Value {
	return Value{Kind: ValueKindBool, Bool: &value}
}

func NewStringValue(value string) Value {
	return Value{Kind: ValueKindString, String: value}
}

func NewMissingValue(message string) Value {
	return Value{Kind: ValueKindMissing, Message: message}
}

func NewErrorValue(message string) Value {
	return Value{Kind: ValueKindError, Message: message}
}

func (v Value) Validate() error {
	switch v.Kind {
	case ValueKindNumeric:
		if v.Number == nil {
			return fmt.Errorf("sensor value kind %q requires number", v.Kind)
		}
	case ValueKindBool:
		if v.Bool == nil {
			return fmt.Errorf("sensor value kind %q requires bool", v.Kind)
		}
	case ValueKindString:
		// Empty string is a valid string value.
	case ValueKindMissing, ValueKindError:
		// Message is optional for now; v3.1.6 will tighten status semantics.
	default:
		if strings.TrimSpace(v.Kind) == "" {
			return fmt.Errorf("sensor value kind is required")
		}
		return fmt.Errorf("sensor value kind %q is not supported", v.Kind)
	}
	return nil
}

func (v Value) IsNumeric() bool {
	return v.Kind == ValueKindNumeric && v.Number != nil
}

func (v Value) Numeric() (float64, bool) {
	if !v.IsNumeric() {
		return 0, false
	}
	return *v.Number, true
}

func (v Value) Equal(other Value) bool {
	if v.Kind != other.Kind || v.String != other.String || v.Unit != other.Unit || v.Message != other.Message {
		return false
	}
	if (v.Number == nil) != (other.Number == nil) {
		return false
	}
	if v.Number != nil && *v.Number != *other.Number {
		return false
	}
	if (v.Bool == nil) != (other.Bool == nil) {
		return false
	}
	if v.Bool != nil && *v.Bool != *other.Bool {
		return false
	}
	return true
}

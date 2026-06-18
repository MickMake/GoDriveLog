package sensors

import (
	"errors"
	"testing"
)

func TestStatusForErrorMapsExplicitUnavailableErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "missing", err: ErrSensorMissing, want: StatusMissing},
		{name: "unsupported", err: ErrSensorUnsupported, want: StatusUnsupported},
		{name: "timeout", err: ErrSensorTimeout, want: StatusTimeout},
		{name: "parse", err: ErrSensorParse, want: StatusParseError},
		{name: "generic", err: errors.New("adapter failed"), want: StatusError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := StatusForError(test.err); got != test.want {
				t.Fatalf("StatusForError() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestValueForStatusUsesMissingValueForUnavailableStates(t *testing.T) {
	for _, status := range []string{StatusMissing, StatusUnsupported} {
		value := ValueForStatus(status, "not available")
		if value.Kind != ValueKindMissing || value.Message != "not available" {
			t.Fatalf("ValueForStatus(%q) = %#v, want missing value", status, value)
		}
	}

	value := ValueForStatus(StatusTimeout, "request timed out")
	if value.Kind != ValueKindError || value.Message != "request timed out" {
		t.Fatalf("ValueForStatus(timeout) = %#v, want error value", value)
	}
}

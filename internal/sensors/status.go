package sensors

import "errors"

var (
	ErrSensorUnsupported = errors.New("sensor unsupported")
	ErrSensorMissing     = errors.New("sensor missing")
	ErrSensorTimeout     = errors.New("sensor timeout")
	ErrSensorParse       = errors.New("sensor parse error")
)

func StatusForError(err error) string {
	switch {
	case errors.Is(err, ErrSensorUnsupported):
		return StatusUnsupported
	case errors.Is(err, ErrSensorMissing):
		return StatusMissing
	case errors.Is(err, ErrSensorTimeout):
		return StatusTimeout
	case errors.Is(err, ErrSensorParse):
		return StatusParseError
	default:
		return StatusError
	}
}

func ValueForStatus(status, message string) Value {
	switch status {
	case StatusMissing, StatusUnsupported:
		return NewMissingValue(message)
	default:
		return NewErrorValue(message)
	}
}

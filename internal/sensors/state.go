package sensors

import "time"

const (
	StatusUnknown            = "unknown"
	StatusOK                 = "ok"
	StatusMissing            = "missing"
	StatusUnsupported        = "unsupported"
	StatusTimeout            = "timeout"
	StatusParseError         = "parse_error"
	StatusError              = "error"
	StatusStale              = "stale"
	StatusMissingUnsupported = StatusMissing
)

type SensorState struct {
	ID         string
	Value      float64
	TypedValue Value
	Unit       string
	Min        float64
	Max        float64
	Status     string
	Error      string
	UpdatedAt  time.Time
	StaleAfter time.Duration
}

type SensorDefinition struct {
	ID         string
	Unit       string
	Min        float64
	Max        float64
	StaleAfter time.Duration
	ValueKind  string
}

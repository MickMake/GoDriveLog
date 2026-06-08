package sensors

import "time"

const (
	StatusUnknown = "unknown"
	StatusOK      = "ok"
	StatusError   = "error"
	StatusStale   = "stale"
)

type SensorState struct {
	ID         string
	Value      float64
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
}

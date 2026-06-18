package sensors

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
)

const (
	EventKindFirstRead    = "first_read"
	EventKindValueChange  = "value_change"
	EventKindStatusChange = "status_change"
	EventKindStale        = "stale"
	EventKindError        = "error"
	EventKindRecovery     = "recovery"
)

// SensorEvent is the v3 event emitted by the central sensor polling runtime.
// It carries the latest state snapshot plus the read/check timestamp that caused
// the event, so later subscribers do not need to poll the endpoint themselves.
type SensorEvent struct {
	Kind           string
	SensorID       string
	State          SensorState
	PreviousStatus string
	Timestamp      time.Time
	ReadAt         time.Time
	Error          string
}

type RuntimeOption func(*PollingRuntime)

// WithNow is intended for deterministic tests.
func WithNow(now func() time.Time) RuntimeOption {
	return func(r *PollingRuntime) {
		if now != nil {
			r.now = now
		}
	}
}

// WithStaleCheckInterval overrides the background stale check cadence.
func WithStaleCheckInterval(interval time.Duration) RuntimeOption {
	return func(r *PollingRuntime) {
		if interval > 0 {
			r.staleCheckInterval = interval
		}
	}
}

type pollingSensor struct {
	ID         string
	PID        string
	ValueKind  string
	Unit       string
	Poll       time.Duration
	StaleAfter time.Duration
	Min        float64
	Max        float64
}

// PollingRuntime owns the one endpoint polling path for a resolved v3 runtime.
// Logs and dashboards should subscribe to events instead of reading sensors.
type PollingRuntime struct {
	reader             Reader
	sensors            []pollingSensor
	store              *StateStore
	now                func() time.Time
	staleCheckInterval time.Duration

	mu          sync.RWMutex
	subscribers []chan SensorEvent
}

func NewPollingRuntime(reader Reader, sensorConfigs map[string]v3config.SensorConfig, opts ...RuntimeOption) (*PollingRuntime, error) {
	if reader == nil {
		return nil, fmt.Errorf("sensor polling runtime requires a reader")
	}

	sensors := make([]pollingSensor, 0, len(sensorConfigs))
	definitions := make([]SensorDefinition, 0, len(sensorConfigs))
	for sensorID, cfg := range sensorConfigs {
		if cfg.Poll <= 0 {
			return nil, fmt.Errorf("sensor %q poll must be greater than zero", sensorID)
		}
		if cfg.PID == "" {
			return nil, fmt.Errorf("sensor %q pid must be set", sensorID)
		}

		valueKind, err := effectiveSensorValueKind(sensorID, cfg)
		if err != nil {
			return nil, err
		}

		poll := time.Duration(cfg.Poll) * time.Millisecond
		staleAfter := StaleAfterForPoll(poll)
		minValue := 0.0
		maxValue := 0.0
		if cfg.Min != nil {
			minValue = *cfg.Min
		}
		if cfg.Max != nil {
			maxValue = *cfg.Max
		}

		sensors = append(sensors, pollingSensor{
			ID:         sensorID,
			PID:        cfg.PID,
			ValueKind:  valueKind,
			Unit:       cfg.Unit,
			Poll:       poll,
			StaleAfter: staleAfter,
			Min:        minValue,
			Max:        maxValue,
		})
		definitions = append(definitions, SensorDefinition{
			ID:         sensorID,
			Unit:       cfg.Unit,
			Min:        minValue,
			Max:        maxValue,
			StaleAfter: staleAfter,
			ValueKind:  valueKind,
		})
	}

	runtime := &PollingRuntime{
		reader: reader,
		sensors: sensors,
		store:   NewStateStore(definitions),
		now:     time.Now,
	}
	for _, opt := range opts {
		opt(runtime)
	}
	return runtime, nil
}

func effectiveSensorValueKind(sensorID string, cfg v3config.SensorConfig) (string, error) {
	parserKind := v3config.SensorOutputValueKind(cfg)
	if parserKind == "" {
		return "", fmt.Errorf("sensor %q type %q has no declared parser output value kind", sensorID, cfg.Type)
	}
	configuredKind := v3config.SensorDeclaredValueKind(cfg)
	if configuredKind == "" {
		return parserKind, nil
	}
	if !isAllowedRuntimeValueKind(configuredKind) {
		return "", fmt.Errorf("sensor %q has invalid value_kind %q; allowed: numeric, bool, string", sensorID, configuredKind)
	}
	if configuredKind != parserKind {
		return "", fmt.Errorf("sensor %q value_kind %q is incompatible with parser output kind %q", sensorID, configuredKind, parserKind)
	}
	return configuredKind, nil
}

func isAllowedRuntimeValueKind(kind string) bool {
	switch kind {
	case ValueKindNumeric, ValueKindBool, ValueKindString:
		return true
	default:
		return false
	}
}

func StaleAfterForPoll(poll time.Duration) time.Duration {
	return time.Duration(math.Max(float64(poll*3), float64(time.Second)))
}

func (r *PollingRuntime) StateStore() *StateStore {
	return r.store
}

// Subscribe registers a subscriber for every sensor event emitted by this runtime.
// Delivery is reliable and blocking: the runtime will not silently drop events,
// and a slow subscriber may apply backpressure to sensor polling. Subscribers must
// continuously drain the returned channel until the runtime closes it.
func (r *PollingRuntime) Subscribe(buffer int) <-chan SensorEvent {
	if buffer < 0 {
		buffer = 0
	}
	ch := make(chan SensorEvent, buffer)
	r.mu.Lock()
	r.subscribers = append(r.subscribers, ch)
	r.mu.Unlock()
	return ch
}

func (r *PollingRuntime) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, sensor := range r.sensors {
		sensor := sensor
		wg.Add(2)
		go func() {
			defer wg.Done()
			r.pollLoop(ctx, sensor)
		}()
		go func() {
			defer wg.Done()
			r.staleLoop(ctx, sensor)
		}()
	}
	<-ctx.Done()
	wg.Wait()
	r.closeSubscribers()
}

func (r *PollingRuntime) pollLoop(ctx context.Context, sensor pollingSensor) {
	r.pollOnce(ctx, sensor)
	ticker := time.NewTicker(sensor.Poll)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.pollOnce(ctx, sensor)
		}
	}
}

func (r *PollingRuntime) pollOnce(ctx context.Context, sensor pollingSensor) {
	if err := ctx.Err(); err != nil {
		return
	}
	readAt := r.now()
	value, unit, err := r.reader.Read(ctx, sensor.PID)
	if unit == "" {
		unit = sensor.Unit
	}
	if err != nil {
		r.applyError(sensor.ID, err, readAt)
		return
	}
	if sensor.ValueKind != ValueKindNumeric {
		r.applyError(sensor.ID, fmt.Errorf("expected %s value, got numeric", sensor.ValueKind), readAt)
		return
	}
	r.applyTypedValue(sensor.ID, NewNumericValue(value, unit), readAt)
}

func (r *PollingRuntime) staleLoop(ctx context.Context, sensor pollingSensor) {
	interval := r.staleCheckInterval
	if interval <= 0 {
		interval = sensor.Poll
		if interval > sensor.StaleAfter/2 {
			interval = sensor.StaleAfter / 2
		}
		if interval <= 0 {
			interval = time.Second
		}
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			state, changed := r.store.MarkStale(sensor.ID, r.now())
			if changed {
				r.emit(SensorEvent{
					Kind:           EventKindStale,
					SensorID:       sensor.ID,
					State:          state,
					PreviousStatus: StatusOK,
					Timestamp:      r.now(),
					ReadAt:         state.UpdatedAt,
				})
			}
		}
	}
}

func (r *PollingRuntime) applyValue(sensorID string, value float64, unit string, readAt time.Time) {
	r.applyTypedValue(sensorID, NewNumericValue(value, unit), readAt)
}

func (r *PollingRuntime) applyTypedValue(sensorID string, value Value, readAt time.Time) {
	previous, hadPrevious := r.store.Get(sensorID)
	if err := value.Validate(); err != nil {
		r.applyError(sensorID, err, readAt)
		return
	}
	state := r.store.SetTypedValue(sensorID, value, readAt)

	kind := ""
	previousStatus := ""
	if hadPrevious {
		previousStatus = previous.Status
	}

	switch {
	case !hadPrevious || previous.UpdatedAt.IsZero():
		kind = EventKindFirstRead
	case previous.Status != StatusOK && previous.Status != StatusUnknown:
		kind = EventKindRecovery
	case previous.Status != state.Status:
		kind = EventKindStatusChange
	case !previous.TypedValue.Equal(state.TypedValue):
		kind = EventKindValueChange
	case previous.Unit != state.Unit:
		kind = EventKindValueChange
	}
	if kind == "" {
		return
	}

	r.emit(SensorEvent{
		Kind:           kind,
		SensorID:       sensorID,
		State:          state,
		PreviousStatus: previousStatus,
		Timestamp:      readAt,
		ReadAt:         readAt,
	})
}

func (r *PollingRuntime) applyError(sensorID string, readErr error, readAt time.Time) {
	previous, hadPrevious := r.store.Get(sensorID)
	state := r.store.SetError(sensorID, readErr, readAt)
	previousStatus := ""
	if hadPrevious {
		previousStatus = previous.Status
	}

	if hadPrevious && previous.Status == state.Status && previous.Error == state.Error {
		return
	}

	r.emit(SensorEvent{
		Kind:           EventKindError,
		SensorID:       sensorID,
		State:          state,
		PreviousStatus: previousStatus,
		Timestamp:      readAt,
		ReadAt:         readAt,
		Error:          state.Error,
	})
}

func (r *PollingRuntime) emit(event SensorEvent) {
	r.mu.RLock()
	subscribers := append([]chan SensorEvent(nil), r.subscribers...)
	r.mu.RUnlock()

	for _, subscriber := range subscribers {
		subscriber <- event
	}
}

func (r *PollingRuntime) closeSubscribers() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, subscriber := range r.subscribers {
		close(subscriber)
	}
	r.subscribers = nil
}

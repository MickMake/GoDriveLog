package sensors

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
)

type scriptedReader struct {
	mu       sync.Mutex
	values   []float64
	units    []string
	errors   []error
	readCount int
}

func (r *scriptedReader) Read(ctx context.Context, pid string) (float64, string, error) {
	select {
	case <-ctx.Done():
		return 0, "", ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	idx := r.readCount
	r.readCount++
	if idx < len(r.errors) && r.errors[idx] != nil {
		return 0, "", r.errors[idx]
	}
	value := 0.0
	if idx < len(r.values) {
		value = r.values[idx]
	}
	unit := ""
	if idx < len(r.units) {
		unit = r.units[idx]
	}
	return value, unit, nil
}

func (r *scriptedReader) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.readCount
}

func TestStaleAfterForPollUsesDocumentedRule(t *testing.T) {
	if got := StaleAfterForPoll(250 * time.Millisecond); got != time.Second {
		t.Fatalf("StaleAfterForPoll(250ms) = %s, want 1s", got)
	}
	if got := StaleAfterForPoll(500 * time.Millisecond); got != 1500*time.Millisecond {
		t.Fatalf("StaleAfterForPoll(500ms) = %s, want 1500ms", got)
	}
}

func TestPollingRuntimeEmitsFirstReadAndUpdatesStore(t *testing.T) {
	reader := &scriptedReader{values: []float64{88}, units: []string{"km/h"}}
	runtime, err := NewPollingRuntime(reader, map[string]v3config.SensorConfig{
		"speed": {Type: "obd", PID: "010D", Unit: "km/h", Poll: 60000},
	})
	if err != nil {
		t.Fatalf("NewPollingRuntime: %v", err)
	}
	events := runtime.Subscribe(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go runtime.Run(ctx)

	event := receiveEvent(t, events)
	cancel()
	if event.Kind != EventKindFirstRead {
		t.Fatalf("event kind = %q, want %q", event.Kind, EventKindFirstRead)
	}
	if event.SensorID != "speed" || event.State.Value != 88 || event.State.Status != StatusOK {
		t.Fatalf("unexpected event: %#v", event)
	}
	state, ok := runtime.StateStore().Get("speed")
	if !ok {
		t.Fatal("speed state missing")
	}
	if state.Value != 88 || state.Status != StatusOK {
		t.Fatalf("state = %#v, want value 88 ok", state)
	}
}

func TestPollingRuntimeDoesNotDuplicateEndpointReadsForSubscribers(t *testing.T) {
	reader := &scriptedReader{values: []float64{1234}, units: []string{"rpm"}}
	runtime, err := NewPollingRuntime(reader, map[string]v3config.SensorConfig{
		"rpm": {Type: "obd", PID: "010C", Unit: "rpm", Poll: 60000},
	})
	if err != nil {
		t.Fatalf("NewPollingRuntime: %v", err)
	}
	firstSubscriber := runtime.Subscribe(1)
	secondSubscriber := runtime.Subscribe(1)

	ctx, cancel := context.WithCancel(context.Background())
	go runtime.Run(ctx)
	defer cancel()

	first := receiveEvent(t, firstSubscriber)
	second := receiveEvent(t, secondSubscriber)
	cancel()

	if first.SensorID != "rpm" || second.SensorID != "rpm" {
		t.Fatalf("unexpected subscriber events: %#v %#v", first, second)
	}
	if count := reader.Count(); count != 1 {
		t.Fatalf("reader read count = %d, want 1 shared endpoint poll", count)
	}
}

func TestPollingRuntimeEmitsValueChangeOnlyWhenChanged(t *testing.T) {
	runtime := newRuntimeForEventTests(t)
	events := runtime.Subscribe(4)
	readAt := time.Date(2026, 6, 13, 8, 0, 0, 0, time.UTC)

	runtime.applyValue("rpm", 1000, "rpm", readAt)
	runtime.applyValue("rpm", 1000, "rpm", readAt.Add(time.Second))
	runtime.applyValue("rpm", 1100, "rpm", readAt.Add(2*time.Second))

	first := receiveEvent(t, events)
	changed := receiveEvent(t, events)
	assertNoEvent(t, events)

	if first.Kind != EventKindFirstRead {
		t.Fatalf("first kind = %q, want %q", first.Kind, EventKindFirstRead)
	}
	if changed.Kind != EventKindValueChange || changed.State.Value != 1100 {
		t.Fatalf("changed event = %#v, want value-change to 1100", changed)
	}
}

func TestPollingRuntimeEmitsErrorAndRecoveryTransitions(t *testing.T) {
	runtime := newRuntimeForEventTests(t)
	events := runtime.Subscribe(4)
	readAt := time.Date(2026, 6, 13, 8, 0, 0, 0, time.UTC)

	runtime.applyError("rpm", errors.New("adapter timeout"), readAt)
	runtime.applyValue("rpm", 1200, "rpm", readAt.Add(time.Second))

	errorEvent := receiveEvent(t, events)
	recoveryEvent := receiveEvent(t, events)

	if errorEvent.Kind != EventKindError || errorEvent.State.Status != StatusError || errorEvent.Error != "adapter timeout" {
		t.Fatalf("error event = %#v, want error status", errorEvent)
	}
	if recoveryEvent.Kind != EventKindRecovery || recoveryEvent.PreviousStatus != StatusError || recoveryEvent.State.Status != StatusOK {
		t.Fatalf("recovery event = %#v, want recovery from error to ok", recoveryEvent)
	}
}

func TestStateStoreMarkStaleTransitionsOnce(t *testing.T) {
	updatedAt := time.Date(2026, 6, 13, 8, 0, 0, 0, time.UTC)
	store := NewStateStore([]SensorDefinition{{ID: "rpm", Unit: "rpm", StaleAfter: time.Second}})
	store.SetValue("rpm", 1000, "", updatedAt)

	state, changed := store.MarkStale("rpm", updatedAt.Add(2*time.Second))
	if !changed {
		t.Fatal("expected first stale mark to change state")
	}
	if state.Status != StatusStale {
		t.Fatalf("status = %q, want %q", state.Status, StatusStale)
	}
	_, changed = store.MarkStale("rpm", updatedAt.Add(3*time.Second))
	if changed {
		t.Fatal("expected second stale mark to be unchanged")
	}
}

func newRuntimeForEventTests(t *testing.T) *PollingRuntime {
	t.Helper()
	runtime, err := NewPollingRuntime(&scriptedReader{}, map[string]v3config.SensorConfig{
		"rpm": {Type: "obd", PID: "010C", Unit: "rpm", Poll: 250},
	})
	if err != nil {
		t.Fatalf("NewPollingRuntime: %v", err)
	}
	return runtime
}

func receiveEvent(t *testing.T, events <-chan SensorEvent) SensorEvent {
	t.Helper()
	select {
	case event := <-events:
		return event
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for sensor event")
		return SensorEvent{}
	}
}

func assertNoEvent(t *testing.T, events <-chan SensorEvent) {
	t.Helper()
	select {
	case event := <-events:
		t.Fatalf("unexpected event: %#v", event)
	case <-time.After(25 * time.Millisecond):
	}
}

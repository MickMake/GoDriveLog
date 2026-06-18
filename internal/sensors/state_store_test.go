package sensors

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestStateStoreInitializesSensorDefinitions(t *testing.T) {
	store := NewStateStore([]SensorDefinition{
		{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000, StaleAfter: 500 * time.Millisecond},
		{ID: "speed", Unit: "km/h", Min: 0, Max: 220, StaleAfter: time.Second},
	})

	state, ok := store.Get("rpm")
	if !ok {
		t.Fatal("rpm state missing")
	}
	if state.ID != "rpm" {
		t.Fatalf("ID = %q, want rpm", state.ID)
	}
	if state.Unit != "rpm" {
		t.Fatalf("Unit = %q, want rpm", state.Unit)
	}
	if state.Min != 0 || state.Max != 7000 {
		t.Fatalf("range = %v..%v, want 0..7000", state.Min, state.Max)
	}
	if state.Status != StatusUnknown {
		t.Fatalf("Status = %q, want %q", state.Status, StatusUnknown)
	}
	if state.StaleAfter != 500*time.Millisecond {
		t.Fatalf("StaleAfter = %v, want 500ms", state.StaleAfter)
	}
}

func TestStateStoreSetValueUpdatesLatestState(t *testing.T) {
	updatedAt := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	store := NewStateStore([]SensorDefinition{{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000, StaleAfter: 500 * time.Millisecond}})

	store.SetValue("rpm", 1234, "", updatedAt)

	state, ok := store.Get("rpm")
	if !ok {
		t.Fatal("rpm state missing")
	}
	if state.Value != 1234 {
		t.Fatalf("Value = %v, want 1234", state.Value)
	}
	if state.Unit != "rpm" {
		t.Fatalf("Unit = %q, want rpm", state.Unit)
	}
	if state.Status != StatusOK {
		t.Fatalf("Status = %q, want %q", state.Status, StatusOK)
	}
	if state.Error != "" {
		t.Fatalf("Error = %q, want empty", state.Error)
	}
	if !state.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("UpdatedAt = %v, want %v", state.UpdatedAt, updatedAt)
	}
	if state.StaleAfter != 500*time.Millisecond {
		t.Fatalf("StaleAfter = %v, want 500ms", state.StaleAfter)
	}
}

func TestStateStoreSetTypedValueRejectsInvalidValues(t *testing.T) {
	updatedAt := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	store := NewStateStore([]SensorDefinition{{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000, StaleAfter: 500 * time.Millisecond}})

	state := store.SetTypedValue("rpm", Value{}, updatedAt)

	if state.Status != StatusError {
		t.Fatalf("Status = %q, want %q", state.Status, StatusError)
	}
	if state.TypedValue.Kind != ValueKindError {
		t.Fatalf("TypedValue = %#v, want error value", state.TypedValue)
	}
	if !strings.Contains(state.Error, "sensor value kind is required") {
		t.Fatalf("Error = %q, want missing kind error", state.Error)
	}
	if state.Value != 0 {
		t.Fatalf("Value = %v, want invalid value not to become a live numeric reading", state.Value)
	}
}

func TestStateStoreSetValueOverridesUnitWhenProvided(t *testing.T) {
	store := NewStateStore([]SensorDefinition{{ID: "coolant", Unit: "C", Min: -40, Max: 120}})

	store.SetValue("coolant", 82, "degC", time.Now())

	state, _ := store.Get("coolant")
	if state.Unit != "degC" {
		t.Fatalf("Unit = %q, want degC", state.Unit)
	}
}

func TestStateStoreSetErrorPreservesMetadata(t *testing.T) {
	updatedAt := time.Date(2026, 6, 8, 10, 5, 0, 0, time.UTC)
	store := NewStateStore([]SensorDefinition{{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000, StaleAfter: 500 * time.Millisecond}})

	store.SetError("rpm", errors.New("adapter timeout"), updatedAt)

	state, ok := store.Get("rpm")
	if !ok {
		t.Fatal("rpm state missing")
	}
	if state.Status != StatusError {
		t.Fatalf("Status = %q, want %q", state.Status, StatusError)
	}
	if state.Error != "adapter timeout" {
		t.Fatalf("Error = %q, want adapter timeout", state.Error)
	}
	if state.Unit != "rpm" || state.Min != 0 || state.Max != 7000 {
		t.Fatalf("metadata changed: unit=%q range=%v..%v", state.Unit, state.Min, state.Max)
	}
	if state.StaleAfter != 500*time.Millisecond {
		t.Fatalf("StaleAfter = %v, want 500ms", state.StaleAfter)
	}
	if !state.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("UpdatedAt = %v, want %v", state.UpdatedAt, updatedAt)
	}
}

func TestStateStoreSnapshotIsSortedAndIndependent(t *testing.T) {
	store := NewStateStore([]SensorDefinition{
		{ID: "speed", Unit: "km/h", Min: 0, Max: 220},
		{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000},
	})

	snapshot := store.Snapshot()
	if len(snapshot) != 2 {
		t.Fatalf("len(snapshot) = %d, want 2", len(snapshot))
	}
	if snapshot[0].ID != "rpm" || snapshot[1].ID != "speed" {
		t.Fatalf("snapshot order = %q, %q; want rpm, speed", snapshot[0].ID, snapshot[1].ID)
	}

	snapshot[0].Value = 9999
	state, _ := store.Get("rpm")
	if state.Value == 9999 {
		t.Fatal("snapshot mutation changed store state")
	}
}

func TestStateStoreMarksOnlyOKReadingsStale(t *testing.T) {
	updatedAt := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	now := updatedAt.Add(2 * time.Second)
	store := NewStateStore([]SensorDefinition{
		{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000, StaleAfter: time.Second},
		{ID: "speed", Unit: "km/h", Min: 0, Max: 220, StaleAfter: time.Second},
	})
	store.SetValue("rpm", 1000, "", updatedAt)
	store.SetError("speed", errors.New("read failed"), updatedAt)

	rpm, ok := store.GetWithStale("rpm", now)
	if !ok {
		t.Fatal("rpm state missing")
	}
	if rpm.Status != StatusStale {
		t.Fatalf("rpm Status = %q, want %q", rpm.Status, StatusStale)
	}

	speed, ok := store.GetWithStale("speed", now)
	if !ok {
		t.Fatal("speed state missing")
	}
	if speed.Status != StatusError {
		t.Fatalf("speed Status = %q, want %q", speed.Status, StatusError)
	}
}

func TestStateStoreUsesPerSensorStaleThresholds(t *testing.T) {
	updatedAt := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	store := NewStateStore([]SensorDefinition{
		{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000, StaleAfter: 500 * time.Millisecond},
		{ID: "fuel_level", Unit: "%", Min: 0, Max: 100, StaleAfter: 10 * time.Second},
	})
	store.SetValue("rpm", 1000, "", updatedAt)
	store.SetValue("fuel_level", 80, "", updatedAt)

	snapshot := store.SnapshotWithStale(updatedAt.Add(2 * time.Second))
	states := map[string]SensorState{}
	for _, state := range snapshot {
		states[state.ID] = state
	}

	if states["rpm"].Status != StatusStale {
		t.Fatalf("rpm Status = %q, want %q", states["rpm"].Status, StatusStale)
	}
	if states["fuel_level"].Status != StatusOK {
		t.Fatalf("fuel_level Status = %q, want %q", states["fuel_level"].Status, StatusOK)
	}
}

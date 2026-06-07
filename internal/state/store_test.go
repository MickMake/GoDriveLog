package state

import (
	"errors"
	"testing"
	"time"
)

func TestStoreInitializesSensorDefinitions(t *testing.T) {
	store := NewStore([]SensorDefinition{
		{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000},
		{ID: "speed", Unit: "km/h", Min: 0, Max: 220},
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
}

func TestStoreSetValueUpdatesLatestState(t *testing.T) {
	updatedAt := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	store := NewStore([]SensorDefinition{{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000}})

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
}

func TestStoreSetValueOverridesUnitWhenProvided(t *testing.T) {
	store := NewStore([]SensorDefinition{{ID: "coolant", Unit: "C", Min: -40, Max: 120}})

	store.SetValue("coolant", 82, "degC", time.Now())

	state, _ := store.Get("coolant")
	if state.Unit != "degC" {
		t.Fatalf("Unit = %q, want degC", state.Unit)
	}
}

func TestStoreSetErrorPreservesMetadata(t *testing.T) {
	updatedAt := time.Date(2026, 6, 8, 10, 5, 0, 0, time.UTC)
	store := NewStore([]SensorDefinition{{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000}})

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
	if !state.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("UpdatedAt = %v, want %v", state.UpdatedAt, updatedAt)
	}
}

func TestStoreSnapshotIsSortedAndIndependent(t *testing.T) {
	store := NewStore([]SensorDefinition{
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

func TestStoreMarksOnlyOKReadingsStale(t *testing.T) {
	updatedAt := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	now := updatedAt.Add(2 * time.Second)
	store := NewStore([]SensorDefinition{
		{ID: "rpm", Unit: "rpm", Min: 0, Max: 7000},
		{ID: "speed", Unit: "km/h", Min: 0, Max: 220},
	})
	store.SetValue("rpm", 1000, "", updatedAt)
	store.SetError("speed", errors.New("read failed"), updatedAt)

	rpm, ok := store.GetWithStale("rpm", time.Second, now)
	if !ok {
		t.Fatal("rpm state missing")
	}
	if rpm.Status != StatusStale {
		t.Fatalf("rpm Status = %q, want %q", rpm.Status, StatusStale)
	}

	speed, ok := store.GetWithStale("speed", time.Second, now)
	if !ok {
		t.Fatal("speed state missing")
	}
	if speed.Status != StatusError {
		t.Fatalf("speed Status = %q, want %q", speed.Status, StatusError)
	}
}

package logger

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestJSONLSubscriberWritesUnavailableStatusTypedValues(t *testing.T) {
	basePath := filepath.Join(t.TempDir(), "events.jsonl")
	loggedAt := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)
	writer, err := newJSONLEventWriter(basePath, func() time.Time { return loggedAt })
	if err != nil {
		t.Fatalf("newJSONLEventWriter: %v", err)
	}
	subscriber := NewJSONLSubscriberWithWriter("jsonl", []string{"speed"}, writer)
	defer func() {
		if err := subscriber.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	}()

	readAt := time.Date(2026, 6, 14, 9, 30, 0, 0, time.UTC)
	tests := []struct {
		status    string
		value     sensors.Value
		wantKind  string
		wantError string
	}{
		{status: sensors.StatusMissing, value: sensors.NewMissingValue("not available"), wantKind: sensors.ValueKindMissing, wantError: "not available"},
		{status: sensors.StatusUnsupported, value: sensors.NewMissingValue("sensor unsupported"), wantKind: sensors.ValueKindMissing, wantError: "sensor unsupported"},
		{status: sensors.StatusTimeout, value: sensors.NewErrorValue("sensor timeout"), wantKind: sensors.ValueKindError, wantError: "sensor timeout"},
		{status: sensors.StatusParseError, value: sensors.NewErrorValue("sensor parse error"), wantKind: sensors.ValueKindError, wantError: "sensor parse error"},
		{status: sensors.StatusError, value: sensors.NewErrorValue("adapter failed"), wantKind: sensors.ValueKindError, wantError: "adapter failed"},
	}

	for i, test := range tests {
		eventAt := readAt.Add(time.Duration(i) * time.Second)
		event := sensors.SensorEvent{
			Kind:     sensors.EventKindError,
			SensorID: "speed",
			State: sensors.SensorState{
				ID:         "speed",
				TypedValue: test.value,
				Status:     test.status,
				Error:      test.wantError,
				UpdatedAt:  eventAt,
			},
			PreviousStatus: sensors.StatusOK,
			Timestamp:      eventAt,
			ReadAt:         eventAt,
			Error:          test.wantError,
		}
		if err := subscriber.Handle(event); err != nil {
			t.Fatalf("Handle %s: %v", test.status, err)
		}
	}

	records := readRecords(t, DailyJSONLPath(basePath, loggedAt))
	if len(records) != len(tests) {
		t.Fatalf("len(records) = %d, want %d: %#v", len(records), len(tests), records)
	}
	for i, test := range tests {
		record := records[i]
		if record.Status != test.status {
			t.Fatalf("record[%d].Status = %q, want %q", i, record.Status, test.status)
		}
		if record.PreviousStatus != sensors.StatusOK {
			t.Fatalf("record[%d].PreviousStatus = %q, want ok", i, record.PreviousStatus)
		}
		if record.Error != test.wantError {
			t.Fatalf("record[%d].Error = %q, want %q", i, record.Error, test.wantError)
		}
		if record.Value.Kind != test.wantKind {
			t.Fatalf("record[%d].Value.Kind = %q, want %q from %#v", i, record.Value.Kind, test.wantKind, record.Value)
		}
		if record.Value.Message != test.wantError {
			t.Fatalf("record[%d].Value.Message = %q, want %q", i, record.Value.Message, test.wantError)
		}
	}
}

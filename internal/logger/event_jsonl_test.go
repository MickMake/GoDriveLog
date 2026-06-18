package logger

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestNewJSONLSubscribersFromPlanUsesSelectedVehicleLogs(t *testing.T) {
	selectedBasePath := filepath.Join(t.TempDir(), "selected.jsonl")
	ignoredBasePath := filepath.Join(t.TempDir(), "ignored.jsonl")
	cfg := v3config.Config{
		Vehicles: map[string]v3config.VehicleConfig{
			"van": {
				Name: "Van",
				OBD:  v3config.OBDConfig{Address: "tcp://127.0.0.1:35000", Timeout: 1000},
				Logs: []string{"selected"},
			},
		},
		Sensors: map[string]v3config.SensorConfig{
			"speed": {Type: "obd", PID: "010D", Unit: "km/h", Poll: 250},
			"rpm":   {Type: "obd", PID: "010C", Unit: "rpm", Poll: 250},
		},
		Logs: map[string]v3config.LogConfig{
			"selected": {Path: selectedBasePath, Sensors: []string{"speed"}},
			"ignored":  {Path: ignoredBasePath, Sensors: []string{"rpm"}},
		},
	}

	plan, err := v3config.Resolve(cfg, "van")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	subscribers, err := NewJSONLSubscribersFromPlan(plan)
	if err != nil {
		t.Fatalf("NewJSONLSubscribersFromPlan: %v", err)
	}
	defer closeSubscribers(t, subscribers)

	if len(subscribers) != 1 {
		t.Fatalf("len(subscribers) = %d, want 1", len(subscribers))
	}
	if subscribers[0].ID != "selected" {
		t.Fatalf("subscriber ID = %q, want selected", subscribers[0].ID)
	}
	if subscribers[0].ActivePath() != DailyJSONLPath(selectedBasePath, time.Now()) {
		t.Fatalf("subscriber path = %q, want daily path for %q", subscribers[0].ActivePath(), selectedBasePath)
	}
}

func TestJSONLSubscriberWritesSelectedSensorEvents(t *testing.T) {
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
	eventAt := readAt.Add(50 * time.Millisecond)
	if err := subscriber.Handle(sensorEvent(sensors.EventKindFirstRead, "speed", sensors.StatusOK, 42, "km/h", readAt, eventAt, "", "")); err != nil {
		t.Fatalf("Handle first read: %v", err)
	}
	if err := subscriber.Handle(sensorEvent(sensors.EventKindFirstRead, "rpm", sensors.StatusOK, 1000, "rpm", readAt, eventAt, "", "")); err != nil {
		t.Fatalf("Handle ignored rpm: %v", err)
	}
	if err := subscriber.Handle(sensorEvent(sensors.EventKindValueChange, "speed", sensors.StatusOK, 43, "km/h", readAt.Add(time.Second), eventAt.Add(time.Second), sensors.StatusOK, "")); err != nil {
		t.Fatalf("Handle value change: %v", err)
	}
	if err := subscriber.Handle(sensorEvent(sensors.EventKindError, "speed", sensors.StatusError, 43, "km/h", readAt.Add(2*time.Second), eventAt.Add(2*time.Second), sensors.StatusOK, "adapter timeout")); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	records := readRecords(t, DailyJSONLPath(basePath, loggedAt))
	if len(records) != 3 {
		t.Fatalf("len(records) = %d, want 3: %#v", len(records), records)
	}

	first := records[0]
	if first.LogID != "jsonl" || first.Kind != sensors.EventKindFirstRead || first.SensorID != "speed" {
		t.Fatalf("first record = %#v, want jsonl first_read speed", first)
	}
	if !first.ReadAt.Equal(readAt) || !first.EventAt.Equal(eventAt) {
		t.Fatalf("first timestamps = read %v event %v, want %v %v", first.ReadAt, first.EventAt, readAt, eventAt)
	}
	if first.Status != sensors.StatusOK || first.Value != 42 || first.Unit != "km/h" {
		t.Fatalf("first state = %#v, want ok 42 km/h", first)
	}

	changed := records[1]
	if changed.Kind != sensors.EventKindValueChange || changed.Value != 43 || changed.Status != sensors.StatusOK {
		t.Fatalf("changed record = %#v, want value_change 43 ok", changed)
	}

	errorRecord := records[2]
	if errorRecord.Kind != sensors.EventKindError || errorRecord.Status != sensors.StatusError || errorRecord.PreviousStatus != sensors.StatusOK || errorRecord.Error != "adapter timeout" {
		t.Fatalf("error record = %#v, want error transition from ok", errorRecord)
	}
}

func TestJSONLEventWriterRotatesDaily(t *testing.T) {
	basePath := filepath.Join(t.TempDir(), "events.jsonl")
	current := time.Date(2026, 6, 14, 23, 59, 0, 0, time.UTC)
	writer, err := newJSONLEventWriter(basePath, func() time.Time { return current })
	if err != nil {
		t.Fatalf("newJSONLEventWriter: %v", err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	}()

	if err := writer.WriteEvent(JSONLEventRecord{LogID: "jsonl", SensorID: "speed", Status: sensors.StatusOK, Value: 1}); err != nil {
		t.Fatalf("WriteEvent day one: %v", err)
	}
	dayOnePath := DailyJSONLPath(basePath, current)
	if writer.ActivePath() != dayOnePath {
		t.Fatalf("ActivePath day one = %q, want %q", writer.ActivePath(), dayOnePath)
	}

	current = current.Add(2 * time.Minute)
	if err := writer.WriteEvent(JSONLEventRecord{LogID: "jsonl", SensorID: "speed", Status: sensors.StatusOK, Value: 2}); err != nil {
		t.Fatalf("WriteEvent day two: %v", err)
	}
	dayTwoPath := DailyJSONLPath(basePath, current)
	if writer.ActivePath() != dayTwoPath {
		t.Fatalf("ActivePath day two = %q, want %q", writer.ActivePath(), dayTwoPath)
	}

	dayOne := readRecords(t, dayOnePath)
	dayTwo := readRecords(t, dayTwoPath)
	if len(dayOne) != 1 || dayOne[0].Value != 1 {
		t.Fatalf("day one records = %#v, want one value 1", dayOne)
	}
	if len(dayTwo) != 1 || dayTwo[0].Value != 2 {
		t.Fatalf("day two records = %#v, want one value 2", dayTwo)
	}
}

func TestDailyJSONLPathAddsDateBeforeExtension(t *testing.T) {
	at := time.Date(2026, 6, 14, 9, 30, 0, 0, time.UTC)
	got := DailyJSONLPath(filepath.Join("logs", "vw_caddy.jsonl"), at)
	want := filepath.Join("logs", "vw_caddy-2026-06-14.jsonl")
	if got != want {
		t.Fatalf("DailyJSONLPath = %q, want %q", got, want)
	}
}

func TestJSONLSubscriberSuppressesUnchangedDuplicateEvents(t *testing.T) {
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
	event := sensorEvent(sensors.EventKindFirstRead, "speed", sensors.StatusOK, 42, "km/h", readAt, readAt, "", "")
	if err := subscriber.Handle(event); err != nil {
		t.Fatalf("Handle first duplicate: %v", err)
	}
	if err := subscriber.Handle(event); err != nil {
		t.Fatalf("Handle second duplicate: %v", err)
	}
	if err := subscriber.Handle(sensorEvent(sensors.EventKindValueChange, "speed", sensors.StatusOK, 42, "km/h", readAt.Add(time.Second), readAt.Add(time.Second), sensors.StatusOK, "")); err != nil {
		t.Fatalf("Handle unchanged value_change: %v", err)
	}

	records := readRecords(t, DailyJSONLPath(basePath, loggedAt))
	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want duplicate suppression to leave 1", len(records))
	}
}

func TestJSONLSubscriberRunConsumesEventChannel(t *testing.T) {
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

	events := make(chan sensors.SensorEvent, 1)
	readAt := time.Date(2026, 6, 14, 9, 30, 0, 0, time.UTC)
	events <- sensorEvent(sensors.EventKindFirstRead, "speed", sensors.StatusOK, 42, "km/h", readAt, readAt, "", "")
	close(events)

	if err := subscriber.Run(context.Background(), events); err != nil {
		t.Fatalf("Run: %v", err)
	}

	records := readRecords(t, DailyJSONLPath(basePath, loggedAt))
	if len(records) != 1 || records[0].SensorID != "speed" {
		t.Fatalf("records = %#v, want one speed record", records)
	}
}

func TestJSONLSubscriberRunReturnsContextCancellation(t *testing.T) {
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

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := subscriber.Run(ctx, make(chan sensors.SensorEvent)); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run canceled err = %v, want context.Canceled", err)
	}
}

func sensorEvent(kind, sensorID, status string, value float64, unit string, readAt, eventAt time.Time, previousStatus, eventErr string) sensors.SensorEvent {
	return sensors.SensorEvent{
		Kind:     kind,
		SensorID: sensorID,
		State: sensors.SensorState{
			ID:        sensorID,
			Value:     value,
			Unit:      unit,
			Status:    status,
			Error:     eventErr,
			UpdatedAt: readAt,
		},
		PreviousStatus: previousStatus,
		Timestamp:      eventAt,
		ReadAt:         readAt,
		Error:          eventErr,
	}
}

func readRecords(t *testing.T, path string) []JSONLEventRecord {
	t.Helper()
	file, err := openForRead(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()

	var records []JSONLEventRecord
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record JSONLEventRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			t.Fatalf("unmarshal %q: %v", scanner.Text(), err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan %s: %v", path, err)
	}
	return records
}

func openForRead(path string) (*os.File, error) {
	return os.Open(path)
}

func closeSubscribers(t *testing.T, subscribers []*JSONLSubscriber) {
	t.Helper()
	for _, subscriber := range subscribers {
		if err := subscriber.Close(); err != nil {
			t.Fatalf("Close subscriber %s: %v", subscriber.ID, err)
		}
	}
}

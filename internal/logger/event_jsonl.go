package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

// JSONLEventRecord is the v3 JSON Lines representation of a sensor event.
// It preserves the sensor read timestamp separately from the logger write time
// so logs do not replace runtime timing with file-writer timing.
type JSONLEventRecord struct {
	LogID          string    `json:"log_id"`
	Kind           string    `json:"kind"`
	SensorID       string    `json:"sensor_id"`
	EventAt        time.Time `json:"event_at"`
	ReadAt         time.Time `json:"read_at"`
	LoggedAt       time.Time `json:"logged_at"`
	Status         string    `json:"status"`
	PreviousStatus string    `json:"previous_status,omitempty"`
	Value          float64   `json:"value"`
	Unit           string    `json:"unit,omitempty"`
	Error          string    `json:"error,omitempty"`
}

type jsonlEventWriter interface {
	WriteEvent(JSONLEventRecord) error
	Close() error
	ActivePath() string
}

// JSONLEventWriter writes v3 event records to a daily JSONL file derived from
// the configured base path. v3.1.4 deliberately keeps this simple: daily
// rotation is always enabled and there are no configurable rotation modes.
type JSONLEventWriter struct {
	mu       sync.Mutex
	basePath string
	path     string
	day      string
	file     *os.File
	now      func() time.Time
}

func NewJSONLEventWriter(path string) (*JSONLEventWriter, error) {
	return newJSONLEventWriter(path, time.Now)
}

func newJSONLEventWriter(path string, now func() time.Time) (*JSONLEventWriter, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return nil, fmt.Errorf("jsonl event writer requires a path")
	}
	if now == nil {
		now = time.Now
	}

	writer := &JSONLEventWriter{basePath: trimmed, now: now}
	if err := writer.rotateLocked(now()); err != nil {
		return nil, err
	}
	return writer, nil
}

// DailyJSONLPath returns the concrete path for a configured JSONL base path on
// one day. For example, logs/vw_caddy.jsonl becomes
// logs/vw_caddy-2026-06-18.jsonl.
func DailyJSONLPath(basePath string, at time.Time) string {
	day := at.Format("2006-01-02")
	ext := filepath.Ext(basePath)
	if ext == "" {
		return basePath + "-" + day
	}
	return strings.TrimSuffix(basePath, ext) + "-" + day + ext
}

func (w *JSONLEventWriter) WriteEvent(record JSONLEventRecord) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return fmt.Errorf("jsonl event writer is not open")
	}
	if record.LoggedAt.IsZero() {
		record.LoggedAt = w.now()
	}
	if err := w.rotateLocked(record.LoggedAt); err != nil {
		return err
	}

	line, err := json.Marshal(record)
	if err != nil {
		return err
	}
	_, err = w.file.Write(append(line, '\n'))
	return err
}

func (w *JSONLEventWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}

func (w *JSONLEventWriter) ActivePath() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.path
}

func (w *JSONLEventWriter) rotateLocked(at time.Time) error {
	day := at.Format("2006-01-02")
	if w.file != nil && w.day == day {
		return nil
	}

	nextPath := DailyJSONLPath(w.basePath, at)
	if err := ensureJSONLDir(nextPath); err != nil {
		return err
	}
	nextFile, err := os.OpenFile(nextPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	oldFile := w.file
	w.file = nextFile
	w.path = nextPath
	w.day = day
	if oldFile != nil {
		if err := oldFile.Close(); err != nil {
			return err
		}
	}
	return nil
}

func ensureJSONLDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

type loggedState struct {
	Status string
	Value  float64
	Unit   string
	Error  string
}

// JSONLSubscriber consumes v3 sensor events for one selected log definition.
// It never polls sensors and it never owns cadence; sensors own cadence and this
// type only records selected events delivered by the sensor runtime.
type JSONLSubscriber struct {
	ID      string
	writer  jsonlEventWriter
	sensors map[string]struct{}

	mu       sync.Mutex
	lastSeen map[string]loggedState
}

func NewJSONLSubscribersFromPlan(plan v3config.RuntimePlan) ([]*JSONLSubscriber, error) {
	subscribers := make([]*JSONLSubscriber, 0, len(plan.Logs))
	for _, resolvedLog := range plan.Logs {
		subscriber, err := NewJSONLSubscriber(resolvedLog.ID, resolvedLog.Config)
		if err != nil {
			for _, existing := range subscribers {
				_ = existing.Close()
			}
			return nil, err
		}
		subscribers = append(subscribers, subscriber)
	}
	return subscribers, nil
}

func NewJSONLSubscriber(id string, cfg v3config.LogConfig) (*JSONLSubscriber, error) {
	writer, err := NewJSONLEventWriter(cfg.Path)
	if err != nil {
		return nil, err
	}
	return NewJSONLSubscriberWithWriter(id, cfg.Sensors, writer), nil
}

func NewJSONLSubscriberWithWriter(id string, sensorIDs []string, writer jsonlEventWriter) *JSONLSubscriber {
	selectedSensors := make(map[string]struct{}, len(sensorIDs))
	for _, sensorID := range sensorIDs {
		if sensorID != "" {
			selectedSensors[sensorID] = struct{}{}
		}
	}
	return &JSONLSubscriber{
		ID:       id,
		writer:   writer,
		sensors:  selectedSensors,
		lastSeen: map[string]loggedState{},
	}
}

func (s *JSONLSubscriber) Run(ctx context.Context, events <-chan sensors.SensorEvent) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-events:
			if !ok {
				return nil
			}
			if err := s.Handle(event); err != nil {
				return err
			}
		}
	}
}

func (s *JSONLSubscriber) Handle(event sensors.SensorEvent) error {
	if !s.shouldConsider(event) {
		return nil
	}

	record := s.recordFromEvent(event)
	state := loggedState{
		Status: record.Status,
		Value:  record.Value,
		Unit:   record.Unit,
		Error:  record.Error,
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if previous, ok := s.lastSeen[record.SensorID]; ok && previous == state {
		return nil
	}
	if err := s.writer.WriteEvent(record); err != nil {
		return err
	}
	s.lastSeen[record.SensorID] = state
	return nil
}

func (s *JSONLSubscriber) Close() error {
	return s.writer.Close()
}

func (s *JSONLSubscriber) ActivePath() string {
	return s.writer.ActivePath()
}

func (s *JSONLSubscriber) shouldConsider(event sensors.SensorEvent) bool {
	if _, ok := s.sensors[event.SensorID]; !ok {
		return false
	}
	return isLoggableEventKind(event.Kind)
}

func isLoggableEventKind(kind string) bool {
	switch kind {
	case sensors.EventKindFirstRead,
		sensors.EventKindValueChange,
		sensors.EventKindStatusChange,
		sensors.EventKindStale,
		sensors.EventKindError,
		sensors.EventKindRecovery:
		return true
	default:
		return false
	}
}

func (s *JSONLSubscriber) recordFromEvent(event sensors.SensorEvent) JSONLEventRecord {
	return JSONLEventRecord{
		LogID:          s.ID,
		Kind:           event.Kind,
		SensorID:       event.SensorID,
		EventAt:        event.Timestamp,
		ReadAt:         event.ReadAt,
		Status:         event.State.Status,
		PreviousStatus: event.PreviousStatus,
		Value:          event.State.Value,
		Unit:           event.State.Unit,
		Error:          event.Error,
	}
}

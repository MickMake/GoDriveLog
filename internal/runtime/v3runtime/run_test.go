package v3runtime

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MickMake/GoDriveLog/internal/vehicle"
)

type fakeReader struct {
	mu     sync.Mutex
	reads  int
	closed bool
}

func (r *fakeReader) Read(ctx context.Context, pid string) (float64, string, error) {
	select {
	case <-ctx.Done():
		return 0, "", ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.reads++
	return float64(900 + r.reads), "rpm", nil
}

func (r *fakeReader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
	return nil
}

func (r *fakeReader) Closed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.closed
}

func TestRunLoadsResolvedVehicleAndWritesSelectedJSONLLog(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	logPath := filepath.Join(dir, "events.jsonl")
	configPath := filepath.Join(dir, "config.v3.yaml")
	config := `vehicles:
  test_vehicle:
    name: "Test vehicle"
    obd:
      address: "serial:///dev/test-obd"
      timeout: 100
    logs:
      - jsonl
sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: "rpm"
    poll: 10
assets: {}
logs:
  jsonl:
    path: "` + filepath.ToSlash(logPath) + `"
    sensors:
      - rpm
dashboards: {}
`
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}

	reader := &fakeReader{}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	summary, err := Run(ctx, Options{
		ConfigPath: configPath,
		VehicleID:  "test_vehicle",
		Connector: vehicle.Connector{
			NewSerialReader: func(target string) (vehicle.Reader, error) {
				if target != "/dev/test-obd" {
					t.Fatalf("unexpected serial target %q", target)
				}
				return reader, nil
			},
		},
		EventBuffer: 8,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if summary.VehicleID != "test_vehicle" {
		t.Fatalf("VehicleID = %q, want test_vehicle", summary.VehicleID)
	}
	if summary.SensorCount != 1 || summary.LogCount != 1 || summary.DashboardCount != 0 {
		t.Fatalf("summary counts = sensors:%d logs:%d dashboards:%d", summary.SensorCount, summary.LogCount, summary.DashboardCount)
	}
	if !reader.Closed() {
		t.Fatal("reader was not closed")
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	line := string(data)
	for _, want := range []string{`"log_id":"jsonl"`, `"sensor_id":"rpm"`, `"kind":"first_read"`, `"status":"ok"`} {
		if !strings.Contains(line, want) {
			t.Fatalf("log file did not contain %s; got %s", want, line)
		}
	}
}

func TestRunDrainsDashboardSubscriptionAfterSinkError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := writeTinyPNG(filepath.Join(dir, "assets", "test.png")); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(dir, "config.v3.yaml")
	config := `vehicles:
  test_vehicle:
    name: "Test vehicle"
    obd:
      address: "serial:///dev/test-obd"
      timeout: 100
    dashboards:
      - primary
sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: "rpm"
    poll: 1
assets:
  image_sets:
    test_image:
      image: "assets/test.png"
logs: {}
dashboards:
  primary:
    display: "primary"
    size:
      width: 1
      height: 1
    widgets:
      - id: "background"
        type: "image"
        asset: "test_image"
        position: [0, 0]
`
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}

	reader := &fakeReader{}
	sinkErr := errors.New("dashboard sink failed")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := Run(ctx, Options{
			ConfigPath: configPath,
			VehicleID:  "test_vehicle",
			RepoRoot:   dir,
			Connector: vehicle.Connector{
				NewSerialReader: func(target string) (vehicle.Reader, error) {
					return reader, nil
				},
			},
			EventBuffer: 1,
			DashboardSink: func([]Scene) error {
				time.Sleep(25 * time.Millisecond)
				return sinkErr
			},
		})
		done <- err
	}()

	select {
	case err := <-done:
		if !errors.Is(err, sinkErr) {
			t.Fatalf("Run error = %v, want %v", err, sinkErr)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after dashboard sink error; subscription draining may be blocked")
	}
}

func writeTinyPNG(path string) error {
	const png1x1 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVR4nGJgYGAAAAAEAAGjChXjAAAAAElFTkSuQmCC"
	data, err := base64.StdEncoding.DecodeString(png1x1)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

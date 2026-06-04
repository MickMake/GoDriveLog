package sensors

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Reading struct {
	Time   time.Time `json:"time"`
	Key    string    `json:"key"`
	PID    string    `json:"pid"`
	Name   string    `json:"name,omitempty"`
	Value  float64   `json:"value"`
	Unit   string    `json:"unit,omitempty"`
	Source string    `json:"source"`
}

type Reader interface {
	Read(ctx context.Context, pid string) (float64, string, error)
}

type MockReader struct {
	mu    sync.Mutex
	start time.Time
}

func NewMockReader() *MockReader {
	return &MockReader{start: time.Now()}
}

func (m *MockReader) Read(ctx context.Context, pid string) (float64, string, error) {
	select {
	case <-ctx.Done():
		return 0, "", ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	t := time.Since(m.start).Seconds()
	switch pid {
	case "010C":
		// RPM: sleeps briefly then starts, making log rotation visible without wiring anything.
		if t < 3 {
			return 0, "rpm", nil
		}
		return 900 + 1400*math.Abs(math.Sin(t/3)) + rand.Float64()*80, "rpm", nil
	case "010D":
		if t < 3 {
			return 0, "km/h", nil
		}
		return 20 + 45*math.Abs(math.Sin(t/7)) + rand.Float64()*4, "km/h", nil
	case "0105":
		return 70 + 20*math.Abs(math.Sin(t/10)) + rand.Float64()*2, "C", nil
	default:
		return rand.Float64() * 100, "", nil
	}
}

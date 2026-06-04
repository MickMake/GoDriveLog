package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"GoDriveLog/internal/sensors"
)

type JSONL struct {
	mu     sync.Mutex
	dir    string
	file   *os.File
	active string
}

func NewJSONL(dir string) (*JSONL, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	l := &JSONL{dir: dir}
	return l, l.Rotate("startup")
}

func (l *JSONL) Rotate(reason string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		_ = l.file.Close()
	}

	stamp := time.Now().Format("20060102-150405")
	path := filepath.Join(l.dir, fmt.Sprintf("%s-%s.jsonl", stamp, reason))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	l.file = f
	l.active = path
	return nil
}

func (l *JSONL) Write(r sensors.Reading) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return fmt.Errorf("logger is not open")
	}
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = l.file.Write(append(b, '\n'))
	return err
}

func (l *JSONL) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return nil
	}
	return l.file.Close()
}

func (l *JSONL) ActivePath() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.active
}

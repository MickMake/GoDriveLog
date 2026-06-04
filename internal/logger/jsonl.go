package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

const DateFormat = "20060102"

type JSONL struct {
	mu     sync.Mutex
	dir    string
	file   *os.File
	activeDate string
	active     string
}

func NewJSONL(dir string) (*JSONL, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	l := &JSONL{dir: dir}
	return l, l.openDate(time.Now().Format(DateFormat))
}

func (l *JSONL) openDate(date string) error {
	if l.file != nil {
		_ = l.file.Close()
		l.file = nil
	}

	path := filepath.Join(l.dir, fmt.Sprintf("%s.jsonl", date))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	l.file = f
	l.activeDate = date
	l.active = path
	return nil
}

func (l *JSONL) Write(r sensors.Reading) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	date := time.Now().Format(DateFormat)
	if l.file == nil || l.activeDate != date {
		if err := l.openDate(date); err != nil {
			return err
		}
	}
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
	err := l.file.Close()
	l.file = nil
	return err
}

func (l *JSONL) ActivePath() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.active
}

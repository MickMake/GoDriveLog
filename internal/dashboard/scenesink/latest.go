package scenesink

import (
	"fmt"
	"sync"
	"time"

	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

// Sink consumes one complete set of selected-dashboard scenes.
type Sink func([]v3dashboard.Scene) error

// Stats describes scene-sink producer/render activity. It is intentionally
// simple so tests, harnesses, and diagnostics can reason about producer cadence
// separately from visible render cadence.
type Stats struct {
	Submitted           uint64
	Rendered            uint64
	Superseded          uint64
	LastRenderDuration  time.Duration
	TotalRenderDuration time.Duration
}

// LatestSink coalesces dashboard scene updates so slow rendering never builds a
// stale-frame backlog. Submit stores the latest scenes and waits until that
// frame has either rendered or been superseded by a newer pending frame.
// SubmitLatest stores the latest scenes without waiting for rendering.
type LatestSink struct {
	mu           sync.Mutex
	cond         *sync.Cond
	update       Sink
	latest       []v3dashboard.Scene
	pending      bool
	closed       bool
	err          error
	done         chan struct{}
	seq          uint64
	renderingSeq uint64
	renderedSeq  uint64

	submitted           uint64
	rendered            uint64
	superseded          uint64
	lastRenderDuration  time.Duration
	totalRenderDuration time.Duration
}

// NewLatestSink starts a coalescing dashboard scene sink.
func NewLatestSink(update Sink) (*LatestSink, error) {
	if update == nil {
		return nil, fmt.Errorf("dashboard scene sink update function is required")
	}
	sink := &LatestSink{update: update, done: make(chan struct{})}
	sink.cond = sync.NewCond(&sink.mu)
	go sink.run()
	return sink, nil
}

// Submit records the latest scenes for display. If rendering is already in
// progress, older pending scenes are replaced rather than queued. Submit returns
// once its frame has rendered, has been superseded before rendering, or a render
// error is observed.
func (s *LatestSink) Submit(scenes []v3dashboard.Scene) error {
	if s == nil {
		return fmt.Errorf("dashboard scene sink is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err != nil {
		return s.err
	}
	if s.closed {
		return fmt.Errorf("dashboard scene sink is closed")
	}

	seq := s.submitLocked(scenes)

	for {
		if s.err != nil {
			return s.err
		}
		if s.renderedSeq >= seq {
			return nil
		}
		if s.renderingSeq != seq && s.seq > seq {
			return nil
		}
		if s.closed && !s.pending && s.renderingSeq != seq {
			return fmt.Errorf("dashboard scene sink is closed")
		}
		s.cond.Wait()
	}
}

// SubmitLatest records the latest scenes for display without waiting for Fyne or
// any other renderer to finish. It preserves latest-only coalescing and returns
// any previously observed render error so producers can still fail visibly.
func (s *LatestSink) SubmitLatest(scenes []v3dashboard.Scene) error {
	if s == nil {
		return fmt.Errorf("dashboard scene sink is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err != nil {
		return s.err
	}
	if s.closed {
		return fmt.Errorf("dashboard scene sink is closed")
	}
	s.submitLocked(scenes)
	return nil
}

// Stats returns a snapshot of producer/render activity.
func (s *LatestSink) Stats() Stats {
	if s == nil {
		return Stats{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return Stats{
		Submitted:           s.submitted,
		Rendered:            s.rendered,
		Superseded:          s.superseded,
		LastRenderDuration:  s.lastRenderDuration,
		TotalRenderDuration: s.totalRenderDuration,
	}
}

func (s *LatestSink) submitLocked(scenes []v3dashboard.Scene) uint64 {
	if s.pending {
		s.superseded++
	}
	s.seq++
	s.submitted++
	seq := s.seq
	s.latest = cloneScenes(scenes)
	s.pending = true
	s.cond.Broadcast()
	return seq
}

// Close waits for the worker to finish the latest pending frame and returns the
// first rendering error, if one occurred.
func (s *LatestSink) Close() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	if !s.closed {
		s.closed = true
		s.cond.Broadcast()
	}
	s.mu.Unlock()
	<-s.done
	return s.Err()
}

// Err returns the first rendering error observed by the worker.
func (s *LatestSink) Err() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

func (s *LatestSink) run() {
	defer close(s.done)
	for {
		s.mu.Lock()
		for !s.pending && !s.closed {
			s.cond.Wait()
		}
		if !s.pending && s.closed {
			s.cond.Broadcast()
			s.mu.Unlock()
			return
		}
		scenes := s.latest
		seq := s.seq
		s.latest = nil
		s.pending = false
		s.renderingSeq = seq
		s.cond.Broadcast()
		s.mu.Unlock()

		started := time.Now()
		err := s.update(scenes)
		duration := time.Since(started)

		s.mu.Lock()
		s.lastRenderDuration = duration
		s.totalRenderDuration += duration
		s.rendered++
		if err != nil && s.err == nil {
			s.err = err
			s.closed = true
		}
		if s.renderingSeq == seq {
			s.renderingSeq = 0
		}
		if s.renderedSeq < seq {
			s.renderedSeq = seq
		}
		s.cond.Broadcast()
		if s.err != nil {
			s.mu.Unlock()
			return
		}
		s.mu.Unlock()
	}
}

func cloneScenes(scenes []v3dashboard.Scene) []v3dashboard.Scene {
	if len(scenes) == 0 {
		return nil
	}
	cloned := make([]v3dashboard.Scene, len(scenes))
	for i, scene := range scenes {
		cloned[i] = cloneScene(scene)
	}
	return cloned
}

func cloneScene(scene v3dashboard.Scene) v3dashboard.Scene {
	cloned := scene
	cloned.Widgets = make([]v3dashboard.Widget, len(scene.Widgets))
	for i, widget := range scene.Widgets {
		cloned.Widgets[i] = cloneWidget(widget)
	}
	return cloned
}

func cloneWidget(widget v3dashboard.Widget) v3dashboard.Widget {
	cloned := widget
	cloned.Position = append([]int(nil), widget.Position...)
	cloned.Parts = append([]v3dashboard.Part(nil), widget.Parts...)
	return cloned
}

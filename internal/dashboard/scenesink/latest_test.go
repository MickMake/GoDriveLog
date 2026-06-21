package scenesink

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

func TestLatestSinkDropsStalePendingFrames(t *testing.T) {
	startedFirst := make(chan struct{})
	releaseFirst := make(chan struct{})
	updates := make(chan string, 8)
	var count atomic.Int32

	sink, err := NewLatestSink(func(scenes []v3dashboard.Scene) error {
		call := count.Add(1)
		if len(scenes) != 1 {
			return fmt.Errorf("scene count = %d, want 1", len(scenes))
		}
		updates <- scenes[0].DashboardID
		if call == 1 {
			close(startedFirst)
			<-releaseFirst
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	firstDone := submitAsync(sink, scene("first"))
	<-startedFirst

	staleDone := submitAsync(sink, scene("stale"))
	waitForPendingSeq(t, sink, 2)
	latestDone := submitAsync(sink, scene("latest"))

	assertSubmitReturns(t, staleDone, "stale")
	close(releaseFirst)
	assertSubmitReturns(t, firstDone, "first")
	assertSubmitReturns(t, latestDone, "latest")

	if err := sink.Close(); err != nil {
		t.Fatal(err)
	}
	close(updates)

	got := []string{}
	for update := range updates {
		got = append(got, update)
	}
	want := []string{"first", "latest"}
	if len(got) != len(want) {
		t.Fatalf("updates = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("updates = %v, want %v", got, want)
		}
	}
}

func TestLatestSinkSubmitReturnsRenderError(t *testing.T) {
	wantErr := errors.New("render failed")
	sink, err := NewLatestSink(func(scenes []v3dashboard.Scene) error {
		if scenes[0].DashboardID == "bad" {
			return wantErr
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sink.Submit(scene("bad"))
	if !errors.Is(err, wantErr) {
		t.Fatalf("Submit error = %v, want %v", err, wantErr)
	}
	if !errors.Is(sink.Err(), wantErr) {
		t.Fatalf("sink Err = %v, want %v", sink.Err(), wantErr)
	}
}

func TestLatestSinkSubmitLatestDoesNotWaitForRender(t *testing.T) {
	startedFirst := make(chan struct{})
	releaseFirst := make(chan struct{})
	var rendered atomic.Int32

	sink, err := NewLatestSink(func(scenes []v3dashboard.Scene) error {
		rendered.Add(1)
		if scenes[0].DashboardID == "first" {
			close(startedFirst)
			<-releaseFirst
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := sink.SubmitLatest(scene("first")); err != nil {
		t.Fatalf("SubmitLatest(first) error = %v", err)
	}
	<-startedFirst

	deadline := time.Now().Add(100 * time.Millisecond)
	for i := 0; i < 20; i++ {
		if err := sink.SubmitLatest(scene(fmt.Sprintf("latest-%02d", i))); err != nil {
			t.Fatalf("SubmitLatest(%d) error = %v", i, err)
		}
	}
	if time.Now().After(deadline) {
		t.Fatal("SubmitLatest calls waited for the blocked renderer")
	}

	stats := sink.Stats()
	if stats.Submitted != 21 {
		t.Fatalf("submitted stats = %d, want 21", stats.Submitted)
	}
	if stats.Rendered != 1 {
		t.Fatalf("rendered stats while first render is blocked = %d, want 1", stats.Rendered)
	}
	if stats.Superseded == 0 {
		t.Fatal("expected at least one superseded pending frame")
	}

	close(releaseFirst)
	if err := sink.Close(); err != nil {
		t.Fatal(err)
	}
	if got := rendered.Load(); got > 3 {
		t.Fatalf("rendered %d frames, want latest-only coalescing to avoid backlog", got)
	}
}

func TestLatestSinkSubmitLatestReportsPreviousRenderError(t *testing.T) {
	wantErr := errors.New("render failed")
	sink, err := NewLatestSink(func(scenes []v3dashboard.Scene) error {
		if scenes[0].DashboardID == "bad" {
			return wantErr
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := sink.SubmitLatest(scene("bad")); err != nil {
		t.Fatalf("SubmitLatest(bad) immediate error = %v", err)
	}
	waitForSinkErr(t, sink, wantErr)
	if err := sink.SubmitLatest(scene("after-error")); !errors.Is(err, wantErr) {
		t.Fatalf("SubmitLatest after render error = %v, want %v", err, wantErr)
	}
}

func TestLatestSinkStatsRecordRenderTiming(t *testing.T) {
	sink, err := NewLatestSink(func(scenes []v3dashboard.Scene) error {
		time.Sleep(2 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := sink.Submit(scene("timed")); err != nil {
		t.Fatal(err)
	}
	if err := sink.Close(); err != nil {
		t.Fatal(err)
	}
	stats := sink.Stats()
	if stats.Submitted != 1 || stats.Rendered != 1 {
		t.Fatalf("stats submitted/rendered = %d/%d, want 1/1", stats.Submitted, stats.Rendered)
	}
	if stats.LastRenderDuration <= 0 || stats.TotalRenderDuration <= 0 {
		t.Fatalf("render durations were not recorded: %+v", stats)
	}
}

func BenchmarkLatestSinkSubmitLatestNoBackpressure(b *testing.B) {
	release := make(chan struct{})
	sink, err := NewLatestSink(func(scenes []v3dashboard.Scene) error {
		<-release
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		close(release)
		_ = sink.Close()
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := sink.SubmitLatest(scene("bench")); err != nil {
			b.Fatal(err)
		}
	}
}

func waitForPendingSeq(t *testing.T, sink *LatestSink, seq uint64) {
	t.Helper()
	deadline := time.After(200 * time.Millisecond)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-deadline:
			t.Fatalf("pending seq %d was not observed", seq)
		case <-ticker.C:
			sink.mu.Lock()
			pending := sink.pending && sink.seq == seq
			sink.mu.Unlock()
			if pending {
				return
			}
		}
	}
}

func waitForSinkErr(t *testing.T, sink *LatestSink, want error) {
	t.Helper()
	deadline := time.After(200 * time.Millisecond)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-deadline:
			t.Fatalf("sink error %v was not observed", want)
		case <-ticker.C:
			if errors.Is(sink.Err(), want) {
				return
			}
		}
	}
}

func submitAsync(sink *LatestSink, scenes []v3dashboard.Scene) <-chan error {
	done := make(chan error, 1)
	go func() {
		done <- sink.Submit(scenes)
	}()
	return done
}

func assertSubmitReturns(t *testing.T, done <-chan error, label string) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Submit(%s) error = %v", label, err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Submit(%s) did not return", label)
	}
}

func scene(id string) []v3dashboard.Scene {
	return []v3dashboard.Scene{{DashboardID: id}}
}

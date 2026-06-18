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

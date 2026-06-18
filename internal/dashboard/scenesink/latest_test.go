package scenesink

import (
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
		t.Error(err)
		return
	}

	if err := sink.Submit(scene("first")); err != nil {
		t.Error(err)
		return
	}
	<-startedFirst

	if err := sink.Submit(scene("stale")); err != nil {
		t.Error(err)
		return
	}
	if err := sink.Submit(scene("latest")); err != nil {
		t.Error(err)
		return
	}
	close(releaseFirst)

	if err := sink.Close(); err != nil {
		t.Error(err)
		return
	}
	close(updates)

	got := []string{}
	for update := range updates {
		got = append(got, update)
	}
	want := []string{"first", "latest"}
	if len(got) != len(want) {
		t.Errorf("updates = %v, want %v", got, want)
		return
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("updates = %v, want %v", got, want)
			return
		}
	}
}

func TestLatestSinkSubmitReturnsWhileRendererIsBusy(t *testing.T) {
	startedFirst := make(chan struct{})
	releaseFirst := make(chan struct{})

	sink, err := NewLatestSink(func(scenes []v3dashboard.Scene) error {
		if scenes[0].DashboardID == "first" {
			close(startedFirst)
			<-releaseFirst
		}
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := sink.Submit(scene("first")); err != nil {
		t.Error(err)
		return
	}
	<-startedFirst

	done := make(chan error, 1)
	go func() {
		done <- sink.Submit(scene("second"))
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Error(err)
			return
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Submit blocked behind busy renderer")
		return
	}

	close(releaseFirst)
	if err := sink.Close(); err != nil {
		t.Error(err)
	}
}

func scene(id string) []v3dashboard.Scene {
	return []v3dashboard.Scene{{DashboardID: id}}
}

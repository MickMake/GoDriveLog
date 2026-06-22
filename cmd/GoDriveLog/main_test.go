package main

import "testing"

/*
func TestFyneShutdownSchedulesQuitOnce(t *testing.T) {
	cancelCount := 0
	quitCount := 0
	scheduled := []func(){}

	shutdown := newFyneShutdown(
		func() { cancelCount++ },
		func(fn func()) { scheduled = append(scheduled, fn) },
		func() { quitCount++ },
	)

	shutdown.CancelAndQuit()
	shutdown.CancelAndQuit()
	shutdown.Cancel()
	shutdown.Quit()

	if cancelCount != 1 {
		t.Fatalf("cancel count = %d, want 1", cancelCount)
	}
	if len(scheduled) != 1 {
		t.Fatalf("scheduled quit count = %d, want 1", len(scheduled))
	}
	if quitCount != 0 {
		t.Fatalf("quit was called directly before scheduled function ran: %d", quitCount)
	}

	scheduled[0]()
	if quitCount != 1 {
		t.Fatalf("quit count after running scheduled function = %d, want 1", quitCount)
	}
}
*/

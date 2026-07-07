package gauges

import (
	"math"
	"testing"
	"time"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestAdvanceMinMaxPointerMarkersResetsAtLocalMidnight(t *testing.T) {
	previousLocal := time.Local
	time.Local = time.FixedZone("MarkerTest", 10*60*60)
	defer func() {
		time.Local = previousLocal
	}()

	config := &PointerMarkersConfig{Min: true, Max: true}
	beforeMidnight := time.Date(2026, time.July, 7, 23, 59, 0, 0, time.Local)
	beforePosition := 0.25
	state := AdvanceMinMaxPointerMarkers(PointerMarkerState{}, config, &beforePosition, beforeMidnight, true)

	if !state.Min.Set || !state.Max.Set {
		t.Fatalf("expected daily marker state to set min/max, got %#v", state)
	}
	if state.Min.NormalizedPosition != 0.25 || state.Max.NormalizedPosition != 0.25 {
		t.Fatalf("unexpected pre-midnight marker positions: %#v", state)
	}

	afterMidnight := beforeMidnight.Add(2 * time.Minute)
	state = AdvanceMinMaxPointerMarkers(state, config, nil, afterMidnight, false)
	if state.Min.Set || state.Max.Set {
		t.Fatalf("expected midnight reset to clear min/max, got %#v", state)
	}
	if state.LocalDayKey != "2026-07-08" {
		t.Fatalf("expected local day key to advance, got %q", state.LocalDayKey)
	}

	afterPosition := 0.75
	state = AdvanceMinMaxPointerMarkers(state, config, &afterPosition, afterMidnight.Add(time.Minute), true)
	if !state.Min.Set || !state.Max.Set {
		t.Fatalf("expected new-day sample to set min/max, got %#v", state)
	}
	if state.Min.NormalizedPosition != 0.75 || state.Max.NormalizedPosition != 0.75 {
		t.Fatalf("unexpected post-midnight marker positions: %#v", state)
	}
}

func TestAdvanceMinMaxPointerMarkersRollingWindowRecalculatesAndExpires(t *testing.T) {
	window := 30 * time.Minute
	config := &PointerMarkersConfig{Min: true, Max: true, Window: &window}
	start := time.Unix(1_000, 0).UTC()

	low := 0.20
	state := AdvanceMinMaxPointerMarkers(PointerMarkerState{}, config, &low, start, true)
	high := 0.80
	state = AdvanceMinMaxPointerMarkers(state, config, &high, start.Add(10*time.Minute), true)
	middle := 0.50
	state = AdvanceMinMaxPointerMarkers(state, config, &middle, start.Add(20*time.Minute), true)

	if len(state.Samples) != 3 {
		t.Fatalf("expected three rolling samples, got %#v", state.Samples)
	}
	if !state.Min.Set || !state.Max.Set || state.Min.NormalizedPosition != 0.20 || state.Max.NormalizedPosition != 0.80 {
		t.Fatalf("unexpected initial rolling min/max: %#v", state)
	}

	state = AdvanceMinMaxPointerMarkers(state, config, nil, start.Add(35*time.Minute), false)
	if len(state.Samples) != 2 {
		t.Fatalf("expected oldest sample to expire, got %#v", state.Samples)
	}
	if !state.Min.Set || state.Min.NormalizedPosition != 0.50 {
		t.Fatalf("expected min to recalculate after expiry, got %#v", state)
	}
	if !state.Max.Set || state.Max.NormalizedPosition != 0.80 {
		t.Fatalf("expected max to remain on retained high sample, got %#v", state)
	}

	state = AdvanceMinMaxPointerMarkers(state, config, nil, start.Add(55*time.Minute), false)
	if len(state.Samples) != 0 {
		t.Fatalf("expected all rolling samples to expire, got %#v", state.Samples)
	}
	if state.Min.Set || state.Max.Set {
		t.Fatalf("expected markers to unset with no valid samples, got %#v", state)
	}
}

func TestAdvanceMinMaxPointerMarkersCoalescesUnchangedRollingSamples(t *testing.T) {
	window := 30 * time.Minute
	config := &PointerMarkersConfig{Min: true, Max: true, Window: &window}
	start := time.Unix(2_000, 0).UTC()

	value := 0.40
	state := AdvanceMinMaxPointerMarkers(PointerMarkerState{}, config, &value, start, true)
	state = AdvanceMinMaxPointerMarkers(state, config, &value, start.Add(10*time.Minute), true)

	if len(state.Samples) != 1 {
		t.Fatalf("expected unchanged sample to coalesce, got %#v", state.Samples)
	}
	if !state.Samples[0].RecordedAt.Equal(start.Add(10 * time.Minute)) {
		t.Fatalf("expected coalesced sample timestamp to advance, got %#v", state.Samples[0])
	}

	state = AdvanceMinMaxPointerMarkers(state, config, nil, start.Add(35*time.Minute), false)
	if len(state.Samples) != 1 {
		t.Fatalf("expected coalesced sample to remain within window, got %#v", state.Samples)
	}
	if !state.Min.Set || !state.Max.Set {
		t.Fatalf("expected markers to remain set after coalesced update, got %#v", state)
	}
}

func TestAdvanceMinMaxPointerMarkersRollingWindowDoesNotRefreshWithoutRecord(t *testing.T) {
	window := 30 * time.Minute
	config := &PointerMarkersConfig{Min: true, Max: true, Window: &window}
	start := time.Unix(3_000, 0).UTC()

	value := 0.40
	state := AdvanceMinMaxPointerMarkers(PointerMarkerState{}, config, &value, start, true)
	state = AdvanceMinMaxPointerMarkers(state, config, &value, start.Add(10*time.Minute), false)

	if len(state.Samples) != 1 {
		t.Fatalf("expected unchanged render to avoid new sample, got %#v", state.Samples)
	}
	if !state.Samples[0].RecordedAt.Equal(start) {
		t.Fatalf("expected unchanged render not to refresh sample timestamp, got %#v", state.Samples[0])
	}

	state = AdvanceMinMaxPointerMarkers(state, config, &value, start.Add(31*time.Minute), false)
	if len(state.Samples) != 0 {
		t.Fatalf("expected idle rolling sample to expire without refresh, got %#v", state.Samples)
	}
	if state.Min.Set || state.Max.Set {
		t.Fatalf("expected expired idle rolling markers to unset, got %#v", state)
	}
}

func TestRenderedPointerMarkerPositionUsesFinalRenderedGeometry(t *testing.T) {
	offset := 10.0
	radialPkg := Package{
		Type:   TypeRadial,
		Sensor: "rpm",
		ValueMap: ValueMap{
			Min:        0,
			Max:        7000,
			StartAngle: -135,
			EndAngle:   135,
			Clamp:      true,
		},
		Realism: Realism{CalibrationOffset: &offset},
	}

	radialPosition, ok, err := RenderedPointerMarkerPosition(radialPkg, sensors.SensorState{ID: "rpm", Status: sensors.StatusOK, Value: 3500})
	if err != nil {
		t.Fatalf("RenderedPointerMarkerPosition returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected radial rendered position to be available")
	}
	expected := clampUnit((10 - radialPkg.ValueMap.StartAngle) / (radialPkg.ValueMap.EndAngle - radialPkg.ValueMap.StartAngle))
	if math.Abs(radialPosition-expected) > 1e-9 {
		t.Fatalf("radial rendered position = %v, want %v", radialPosition, expected)
	}

	barPkg := Package{
		Type:   TypeBar,
		Sensor: "coolant_temperature",
		ValueMap: ValueMap{
			Min:   40,
			Max:   120,
			Clamp: true,
		},
	}
	barPosition, ok, err := RenderedPointerMarkerPosition(barPkg, sensors.SensorState{ID: "coolant_temperature", Status: sensors.StatusOK, Value: 80})
	if err != nil {
		t.Fatalf("RenderedPointerMarkerPosition returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected bar rendered position to be available")
	}
	if math.Abs(barPosition-0.5) > 1e-9 {
		t.Fatalf("bar rendered position = %v, want 0.5", barPosition)
	}
}

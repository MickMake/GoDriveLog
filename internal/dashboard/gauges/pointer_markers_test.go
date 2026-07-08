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

func TestAdvanceAveragePointerMarkerStaysUnsetWhenDisabled(t *testing.T) {
	config := &PointerMarkersConfig{Average: false}
	position := 0.4

	state := AdvanceAveragePointerMarker(PointerMarkerState{}, config, &position, time.Unix(1_000, 0).UTC())
	if state.Average.Set {
		t.Fatalf("expected disabled average marker to remain unset, got %#v", state.Average)
	}
}

func TestAdvanceAveragePointerMarkerInitializesFromFirstValidSample(t *testing.T) {
	config := &PointerMarkersConfig{Average: true}
	now := time.Unix(2_000, 0).UTC()
	position := 0.25

	state := AdvanceAveragePointerMarker(PointerMarkerState{}, config, &position, now)
	if !state.Average.Set {
		t.Fatalf("expected first valid sample to initialize average marker, got %#v", state.Average)
	}
	if state.Average.NormalizedPosition != 0.25 || !state.Average.RecordedAt.Equal(now) {
		t.Fatalf("unexpected initial average marker state: %#v", state.Average)
	}
}

func TestAdvanceAveragePointerMarkerUsesFixedTenSecondTimeConstant(t *testing.T) {
	config := &PointerMarkersConfig{Average: true}
	start := time.Unix(3_000, 0).UTC()
	initial := 0.0
	target := 1.0

	state := AdvanceAveragePointerMarker(PointerMarkerState{}, config, &initial, start)
	state = AdvanceAveragePointerMarker(state, config, &target, start.Add(10*time.Second))

	expected := 1 - math.Exp(-1)
	if math.Abs(state.Average.NormalizedPosition-expected) > 1e-12 {
		t.Fatalf("average marker after 10s = %v, want %v", state.Average.NormalizedPosition, expected)
	}
	if math.Abs(state.Average.NormalizedPosition-0.5) < 0.01 {
		t.Fatalf("average marker looks like an arithmetic average, got %v", state.Average.NormalizedPosition)
	}
}

func TestAdvanceAveragePointerMarkerIsFrameRateIndependent(t *testing.T) {
	config := &PointerMarkersConfig{Average: true}
	start := time.Unix(4_000, 0).UTC()
	initial := 0.0
	target := 1.0

	singleStep := AdvanceAveragePointerMarker(PointerMarkerState{}, config, &initial, start)
	singleStep = AdvanceAveragePointerMarker(singleStep, config, &target, start.Add(10*time.Second))

	multiStep := AdvanceAveragePointerMarker(PointerMarkerState{}, config, &initial, start)
	for second := 1; second <= 10; second++ {
		multiStep = AdvanceAveragePointerMarker(multiStep, config, &target, start.Add(time.Duration(second)*time.Second))
	}

	if math.Abs(singleStep.Average.NormalizedPosition-multiStep.Average.NormalizedPosition) > 1e-12 {
		t.Fatalf("frame-rate independent average mismatch: single=%v multi=%v", singleStep.Average.NormalizedPosition, multiStep.Average.NormalizedPosition)
	}
}

func TestAdvanceAveragePointerMarkerDoesNotCalculateArithmeticMean(t *testing.T) {
	config := &PointerMarkersConfig{Average: true}
	start := time.Unix(5_000, 0).UTC()
	zero := 0.0
	one := 1.0

	state := AdvanceAveragePointerMarker(PointerMarkerState{}, config, &zero, start)
	state = AdvanceAveragePointerMarker(state, config, &one, start.Add(10*time.Second))
	state = AdvanceAveragePointerMarker(state, config, &zero, start.Add(20*time.Second))

	expected := (1 - math.Exp(-1)) * math.Exp(-1)
	if math.Abs(state.Average.NormalizedPosition-expected) > 1e-12 {
		t.Fatalf("average marker after damped follow = %v, want %v", state.Average.NormalizedPosition, expected)
	}
	if math.Abs(state.Average.NormalizedPosition-(1.0/3.0)) < 0.01 || math.Abs(state.Average.NormalizedPosition-0.5) < 0.01 {
		t.Fatalf("average marker should not match arithmetic means, got %v", state.Average.NormalizedPosition)
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

func TestRenderedPointerMarkerPositionClampsRadialWhenConfigured(t *testing.T) {
	pkg := Package{
		Type:   TypeRadial,
		Sensor: "rpm",
		ValueMap: ValueMap{
			Min:        0,
			Max:        7000,
			StartAngle: -135,
			EndAngle:   135,
			Clamp:      true,
		},
	}

	above, ok, err := RenderedPointerMarkerPosition(pkg, sensors.SensorState{ID: "rpm", Status: sensors.StatusOK, Value: 9000})
	if err != nil {
		t.Fatalf("RenderedPointerMarkerPosition returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected clamped radial rendered position to be available")
	}
	if above != 1 {
		t.Fatalf("clamped radial above-range position = %v, want 1", above)
	}

	below, ok, err := RenderedPointerMarkerPosition(pkg, sensors.SensorState{ID: "rpm", Status: sensors.StatusOK, Value: -500})
	if err != nil {
		t.Fatalf("RenderedPointerMarkerPosition returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected clamped radial below-range position to be available")
	}
	if below != 0 {
		t.Fatalf("clamped radial below-range position = %v, want 0", below)
	}
}

func TestRenderedPointerMarkerPositionPreservesUnclampedRadialRange(t *testing.T) {
	pkg := Package{
		Type:   TypeRadial,
		Sensor: "rpm",
		ValueMap: ValueMap{
			Min:        0,
			Max:        7000,
			StartAngle: -135,
			EndAngle:   135,
			Clamp:      false,
		},
	}

	above, ok, err := RenderedPointerMarkerPosition(pkg, sensors.SensorState{ID: "rpm", Status: sensors.StatusOK, Value: 9000})
	if err != nil {
		t.Fatalf("RenderedPointerMarkerPosition returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected unclamped radial above-range position to be available")
	}
	if above <= 1 {
		t.Fatalf("unclamped radial above-range position = %v, want > 1", above)
	}

	below, ok, err := RenderedPointerMarkerPosition(pkg, sensors.SensorState{ID: "rpm", Status: sensors.StatusOK, Value: -500})
	if err != nil {
		t.Fatalf("RenderedPointerMarkerPosition returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected unclamped radial below-range position to be available")
	}
	if below >= 0 {
		t.Fatalf("unclamped radial below-range position = %v, want < 0", below)
	}
}

func TestRenderedPointerMarkerPositionPreservesUnclampedRadialCalibrationOffset(t *testing.T) {
	offset := 30.0
	pkg := Package{
		Type:   TypeRadial,
		Sensor: "rpm",
		ValueMap: ValueMap{
			Min:        0,
			Max:        7000,
			StartAngle: -135,
			EndAngle:   135,
			Clamp:      false,
		},
		Realism: Realism{CalibrationOffset: &offset},
	}

	position, ok, err := RenderedPointerMarkerPosition(pkg, sensors.SensorState{ID: "rpm", Status: sensors.StatusOK, Value: 7000})
	if err != nil {
		t.Fatalf("RenderedPointerMarkerPosition returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected unclamped radial calibration-offset position to be available")
	}
	if position <= 1 {
		t.Fatalf("unclamped radial calibration-offset position = %v, want > 1", position)
	}
}

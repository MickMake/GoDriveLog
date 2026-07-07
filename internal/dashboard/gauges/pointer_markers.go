package gauges

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/MickMake/GoDriveLog/internal/sensors"
	"gopkg.in/yaml.v3"
)

const pointerMarkerPositionEpsilon = 1e-9

type PointerMarkersConfig struct {
	Max     bool           `yaml:"max,omitempty"`
	Min     bool           `yaml:"min,omitempty"`
	Average bool           `yaml:"average,omitempty"`
	Window  *time.Duration `yaml:"window,omitempty"`
}

func (c *PointerMarkersConfig) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("realism pointer_markers must be a mapping")
	}

	var decoded PointerMarkersConfig
	for index := 0; index+1 < len(node.Content); index += 2 {
		keyNode := node.Content[index]
		valueNode := node.Content[index+1]
		switch keyNode.Value {
		case "max":
			value, err := decodePointerMarkerBool("max", valueNode)
			if err != nil {
				return err
			}
			decoded.Max = value
		case "min":
			value, err := decodePointerMarkerBool("min", valueNode)
			if err != nil {
				return err
			}
			decoded.Min = value
		case "average":
			value, err := decodePointerMarkerBool("average", valueNode)
			if err != nil {
				return err
			}
			decoded.Average = value
		case "window":
			value, err := decodePointerMarkerWindow(valueNode)
			if err != nil {
				return err
			}
			decoded.Window = &value
		default:
			return fmt.Errorf("realism pointer_markers field %q is not supported", keyNode.Value)
		}
	}

	*c = decoded
	return nil
}

func (c *PointerMarkersConfig) Enabled() bool {
	return c != nil && (c.Max || c.Min || c.Average)
}

func (c *PointerMarkersConfig) MinMaxEnabled() bool {
	return c != nil && (c.Max || c.Min)
}

type PointerMarkerValueState struct {
	Set                bool
	NormalizedPosition float64
	RecordedAt         time.Time
}

type PointerMarkerSample struct {
	NormalizedPosition float64
	RecordedAt         time.Time
}

type PointerMarkerState struct {
	LocalDayKey             string
	Min                     PointerMarkerValueState
	Max                     PointerMarkerValueState
	Average                 PointerMarkerValueState
	Samples                 []PointerMarkerSample
	LastRenderedPosition    float64
	LastRenderedPositionSet bool
}

func AdvanceMinMaxPointerMarkers(state PointerMarkerState, config *PointerMarkersConfig, normalizedPosition *float64, now time.Time, recordSample bool) PointerMarkerState {
	if config == nil || !config.MinMaxEnabled() {
		state.LocalDayKey = ""
		state.Min = PointerMarkerValueState{}
		state.Max = PointerMarkerValueState{}
		state.Samples = nil
		state.LastRenderedPosition = 0
		state.LastRenderedPositionSet = false
		return state
	}

	state = PruneMinMaxPointerMarkers(state, config, now)
	if recordSample {
		state = RecordMinMaxPointerMarkerSample(state, config, normalizedPosition, now)
	}
	state = UpdatePointerMarkerRenderedPosition(state, normalizedPosition)
	return state
}

func PruneMinMaxPointerMarkers(state PointerMarkerState, config *PointerMarkersConfig, now time.Time) PointerMarkerState {
	if config == nil || !config.MinMaxEnabled() {
		return PointerMarkerState{}
	}

	if config.Window == nil {
		return pruneDailyMinMaxPointerMarkers(state, now)
	}
	return pruneRollingMinMaxPointerMarkers(state, config, now)
}

func RecordMinMaxPointerMarkerSample(state PointerMarkerState, config *PointerMarkersConfig, normalizedPosition *float64, now time.Time) PointerMarkerState {
	if config == nil || !config.MinMaxEnabled() {
		return state
	}
	if config.Window == nil {
		return recordDailyMinMaxPointerMarkerSample(state, config, normalizedPosition, now)
	}
	return recordRollingMinMaxPointerMarkerSample(state, config, normalizedPosition, now)
}

func PointerMarkerRenderedPositionChanged(state PointerMarkerState, normalizedPosition *float64) bool {
	position, ok := normalizedPointerMarkerPosition(normalizedPosition)
	if !ok {
		return false
	}
	if !state.LastRenderedPositionSet {
		return true
	}
	return math.Abs(state.LastRenderedPosition-position) > pointerMarkerPositionEpsilon
}

func UpdatePointerMarkerRenderedPosition(state PointerMarkerState, normalizedPosition *float64) PointerMarkerState {
	position, ok := normalizedPointerMarkerPosition(normalizedPosition)
	if !ok {
		state.LastRenderedPosition = 0
		state.LastRenderedPositionSet = false
		return state
	}
	state.LastRenderedPosition = position
	state.LastRenderedPositionSet = true
	return state
}

func RenderedPointerMarkerPosition(pkg Package, state sensors.SensorState) (float64, bool, error) {
	if state.Status != sensors.StatusOK {
		return 0, false, nil
	}

	switch pkg.Type {
	case TypeRadial:
		angle, err := radialAngle(pkg.ValueMap, state.Value)
		if err != nil {
			return 0, false, err
		}
		angle = radialCalibrationAngle(angle, pkg.ValueMap, pkg.Realism.CalibrationOffset)
		span := pkg.ValueMap.EndAngle - pkg.ValueMap.StartAngle
		if span == 0 {
			return 0, false, fmt.Errorf("value_map start_angle and end_angle must differ")
		}
		position := (angle - pkg.ValueMap.StartAngle) / span
		if pkg.ValueMap.Clamp {
			position = clampUnit(position)
		}
		return position, true, nil
	case TypeBar:
		percent, err := barNormalizedPercent(pkg.ValueMap, state.Value)
		if err != nil {
			return 0, false, err
		}
		return clampUnit(percent / 100), true, nil
	default:
		return 0, false, nil
	}
}

func decodePointerMarkerBool(name string, node *yaml.Node) (bool, error) {
	if node.Kind != yaml.ScalarNode {
		return false, fmt.Errorf("realism pointer_markers %s must be a boolean", name)
	}

	var enabled bool
	if err := node.Decode(&enabled); err != nil {
		return false, fmt.Errorf("realism pointer_markers %s must be a boolean", name)
	}
	return enabled, nil
}

func decodePointerMarkerWindow(node *yaml.Node) (time.Duration, error) {
	if node.Kind != yaml.ScalarNode {
		return 0, fmt.Errorf("realism pointer_markers window must be a duration")
	}

	raw := strings.TrimSpace(node.Value)
	if raw == "" {
		return 0, fmt.Errorf("realism pointer_markers window must be a positive duration")
	}

	duration, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("realism pointer_markers window %q is not a valid duration", raw)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("realism pointer_markers window must be greater than zero")
	}
	return duration, nil
}

func pruneDailyMinMaxPointerMarkers(state PointerMarkerState, now time.Time) PointerMarkerState {
	if !now.IsZero() {
		dayKey := pointerMarkerLocalDayKey(now)
		if state.LocalDayKey != dayKey {
			state.LocalDayKey = dayKey
			state.Min = PointerMarkerValueState{}
			state.Max = PointerMarkerValueState{}
			state.Samples = nil
		}
	}
	return state
}

func recordDailyMinMaxPointerMarkerSample(state PointerMarkerState, config *PointerMarkersConfig, normalizedPosition *float64, now time.Time) PointerMarkerState {
	position, ok := normalizedPointerMarkerPosition(normalizedPosition)
	if !ok {
		return state
	}

	if config.Min {
		if !state.Min.Set || position < state.Min.NormalizedPosition {
			state.Min = PointerMarkerValueState{Set: true, NormalizedPosition: position, RecordedAt: now}
		}
	} else {
		state.Min = PointerMarkerValueState{}
	}
	if config.Max {
		if !state.Max.Set || position > state.Max.NormalizedPosition {
			state.Max = PointerMarkerValueState{Set: true, NormalizedPosition: position, RecordedAt: now}
		}
	} else {
		state.Max = PointerMarkerValueState{}
	}

	return state
}

func pruneRollingMinMaxPointerMarkers(state PointerMarkerState, config *PointerMarkersConfig, now time.Time) PointerMarkerState {
	state.LocalDayKey = ""
	state.Samples = prunePointerMarkerSamples(state.Samples, *config.Window, now)
	return recalculateRollingMinMaxPointerMarkers(state, config)
}

func recordRollingMinMaxPointerMarkerSample(state PointerMarkerState, config *PointerMarkersConfig, normalizedPosition *float64, now time.Time) PointerMarkerState {
	position, ok := normalizedPointerMarkerPosition(normalizedPosition)
	if !ok {
		return state
	}
	state.Samples = appendOrCoalescePointerMarkerSample(state.Samples, PointerMarkerSample{
		NormalizedPosition: position,
		RecordedAt:         now,
	})
	return recalculateRollingMinMaxPointerMarkers(state, config)
}

func recalculateRollingMinMaxPointerMarkers(state PointerMarkerState, config *PointerMarkersConfig) PointerMarkerState {
	state.Min = PointerMarkerValueState{}
	state.Max = PointerMarkerValueState{}
	for _, sample := range state.Samples {
		if config.Min && (!state.Min.Set || sample.NormalizedPosition < state.Min.NormalizedPosition) {
			state.Min = PointerMarkerValueState{
				Set:                true,
				NormalizedPosition: sample.NormalizedPosition,
				RecordedAt:         sample.RecordedAt,
			}
		}
		if config.Max && (!state.Max.Set || sample.NormalizedPosition > state.Max.NormalizedPosition) {
			state.Max = PointerMarkerValueState{
				Set:                true,
				NormalizedPosition: sample.NormalizedPosition,
				RecordedAt:         sample.RecordedAt,
			}
		}
	}

	return state
}

func pointerMarkerLocalDayKey(now time.Time) string {
	return now.In(time.Local).Format("2006-01-02")
}

func normalizedPointerMarkerPosition(normalizedPosition *float64) (float64, bool) {
	if normalizedPosition == nil {
		return 0, false
	}
	position := *normalizedPosition
	if math.IsNaN(position) || math.IsInf(position, 0) {
		return 0, false
	}
	return position, true
}

func prunePointerMarkerSamples(samples []PointerMarkerSample, window time.Duration, now time.Time) []PointerMarkerSample {
	if len(samples) == 0 || window <= 0 || now.IsZero() {
		return samples
	}

	pruned := samples[:0]
	for _, sample := range samples {
		if now.Sub(sample.RecordedAt) <= window {
			pruned = append(pruned, sample)
		}
	}
	return pruned
}

func appendOrCoalescePointerMarkerSample(samples []PointerMarkerSample, sample PointerMarkerSample) []PointerMarkerSample {
	if len(samples) == 0 {
		return append(samples, sample)
	}

	last := &samples[len(samples)-1]
	if math.Abs(last.NormalizedPosition-sample.NormalizedPosition) <= pointerMarkerPositionEpsilon {
		last.RecordedAt = sample.RecordedAt
		return samples
	}
	return append(samples, sample)
}

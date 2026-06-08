package decoders

import (
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestExecuteDecoderTypes(t *testing.T) {
	inputs := Inputs{Sensors: map[string]sensors.SensorState{
		"rpm":      {ID: "rpm", Value: 350, Unit: "rpm", Min: 0, Max: 7000, Status: sensors.StatusOK},
		"throttle": {ID: "throttle", Value: 0.50, Unit: "%", Min: 0, Max: 1, Status: sensors.StatusOK},
		"warning":  {ID: "warning", Value: 1, Min: 0, Max: 1, Status: sensors.StatusOK},
	}}

	configs := []config.DashboardDecoderConfig{
		{ID: "rpm_norm", Type: config.DashboardDecoderNormalize, Sensor: "rpm"},
		{ID: "rpm_zone", Type: config.DashboardDecoderThreshold, Sensor: "rpm", Thresholds: []config.ThresholdConfig{{At: 0, Value: "low"}, {At: 300, Value: "mid"}, {At: 6000, Value: "high"}}},
		{ID: "throttle_frame", Type: config.DashboardDecoderFrameIndex, Sensor: "throttle", FrameCount: 11},
		{ID: "rpm_text", Type: config.DashboardDecoderFormatNumber, Sensor: "rpm", Format: "0000"},
		{ID: "rpm_digits", Type: config.DashboardDecoderDigits, Input: "rpm_text"},
		{ID: "warning_bool", Type: config.DashboardDecoderBoolean, Sensor: "warning"},
	}

	values, err := Execute(configs, inputs)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if got := values["rpm_norm"].Number; got != 0.05 {
		t.Fatalf("rpm_norm = %v, want 0.05", got)
	}
	if got := values["rpm_zone"].Text; got != "mid" {
		t.Fatalf("rpm_zone = %q, want mid", got)
	}
	if got := values["throttle_frame"].FrameIndex; got != 5 {
		t.Fatalf("throttle_frame = %d, want 5", got)
	}
	if got := values["rpm_text"].Text; got != "0350" {
		t.Fatalf("rpm_text = %q, want 0350", got)
	}
	if got := values["rpm_digits"].Digits; len(got) != 4 || got[0] != "0" || got[1] != "3" || got[3] != "0" {
		t.Fatalf("rpm_digits = %#v, want [0 3 5 0]", got)
	}
	if got := values["warning_bool"].Bool; !got {
		t.Fatalf("warning_bool = false, want true")
	}
}

func TestNormalizeDecoder(t *testing.T) {
	tests := []struct {
		name  string
		state sensors.SensorState
		want  float64
	}{
		{name: "middle", state: sensorState("rpm", 3500, 0, 7000), want: 0.5},
		{name: "clamps low", state: sensorState("rpm", -10, 0, 7000), want: 0},
		{name: "clamps high", state: sensorState("rpm", 8000, 0, 7000), want: 1},
	}

	registry := NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registry.Decode(config.DashboardDecoderConfig{ID: "norm", Type: config.DashboardDecoderNormalize, Sensor: "rpm"}, sensorInputs(tt.state))
			if err != nil {
				t.Fatalf("Decode returned error: %v", err)
			}
			if got.Number != tt.want {
				t.Fatalf("normalize = %v, want %v", got.Number, tt.want)
			}
		})
	}
}

func TestThresholdDecoder(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{name: "below first keeps first value", value: -1, want: "low"},
		{name: "at boundary selects boundary", value: 3000, want: "mid"},
		{name: "between boundaries keeps previous", value: 4500, want: "mid"},
		{name: "at high boundary selects high", value: 6000, want: "high"},
	}

	decoder := config.DashboardDecoderConfig{
		ID:     "zone",
		Type:   config.DashboardDecoderThreshold,
		Sensor: "rpm",
		Thresholds: []config.ThresholdConfig{
			{At: 0, Value: "low"},
			{At: 3000, Value: "mid"},
			{At: 6000, Value: "high"},
		},
	}
	registry := NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registry.Decode(decoder, sensorInputs(sensorState("rpm", tt.value, 0, 7000)))
			if err != nil {
				t.Fatalf("Decode returned error: %v", err)
			}
			if got.Text != tt.want {
				t.Fatalf("threshold = %q, want %q", got.Text, tt.want)
			}
		})
	}
}

func TestFrameIndexDecoder(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  int
	}{
		{name: "clamps low", value: -10, want: 0},
		{name: "middle", value: 50, want: 5},
		{name: "clamps high", value: 120, want: 10},
	}

	registry := NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registry.Decode(config.DashboardDecoderConfig{ID: "frame", Type: config.DashboardDecoderFrameIndex, Sensor: "throttle", FrameCount: 11}, sensorInputs(sensorState("throttle", tt.value, 0, 100)))
			if err != nil {
				t.Fatalf("Decode returned error: %v", err)
			}
			if got.FrameIndex != tt.want {
				t.Fatalf("frame index = %d, want %d", got.FrameIndex, tt.want)
			}
		})
	}
}

func TestFormatNumberDecoder(t *testing.T) {
	tests := []struct {
		name   string
		value  float64
		format string
		want   string
	}{
		{name: "single zero", value: 7, format: "0", want: "7"},
		{name: "two zero mask", value: 7, format: "00", want: "07"},
		{name: "documented zero mask", value: 350, format: "0000", want: "0350"},
		{name: "mask does not truncate equal width", value: 7000, format: "0000", want: "7000"},
		{name: "mask does not truncate larger value", value: 12345, format: "0000", want: "12345"},
		{name: "go format still supported", value: 12.34, format: "%.1f", want: "12.3"},
	}

	registry := NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registry.Decode(config.DashboardDecoderConfig{ID: "fmt", Type: config.DashboardDecoderFormatNumber, Sensor: "rpm", Format: tt.format}, sensorInputs(sensorState("rpm", tt.value, 0, 20000)))
			if err != nil {
				t.Fatalf("Decode returned error: %v", err)
			}
			if got.Text != tt.want {
				t.Fatalf("format_number = %q, want %q", got.Text, tt.want)
			}
		})
	}
}

func TestDigitsDecoderConsumesPaddedFormatOutput(t *testing.T) {
	values, err := Execute([]config.DashboardDecoderConfig{
		{ID: "rpm_text", Type: config.DashboardDecoderFormatNumber, Sensor: "rpm", Format: "0000"},
		{ID: "rpm_digits", Type: config.DashboardDecoderDigits, Input: "rpm_text"},
	}, sensorInputs(sensorState("rpm", 350, 0, 7000)))
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got := values["rpm_digits"].Digits
	want := []string{"0", "3", "5", "0"}
	if len(got) != len(want) {
		t.Fatalf("digits length = %d, want %d (%#v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digits = %#v, want %#v", got, want)
		}
	}
}

func TestBooleanDecoder(t *testing.T) {
	tests := []struct {
		name   string
		input  Inputs
		config config.DashboardDecoderConfig
		want   bool
	}{
		{name: "numeric false", input: sensorInputs(sensorState("moving", 0, 0, 1)), config: config.DashboardDecoderConfig{ID: "bool", Type: config.DashboardDecoderBoolean, Sensor: "moving"}, want: false},
		{name: "numeric true", input: sensorInputs(sensorState("moving", 1, 0, 1)), config: config.DashboardDecoderConfig{ID: "bool", Type: config.DashboardDecoderBoolean, Sensor: "moving"}, want: true},
		{name: "bool input", input: Inputs{Values: map[string]Value{"prior": {Type: ValueTypeBoolean, Bool: true}}}, config: config.DashboardDecoderConfig{ID: "bool", Type: config.DashboardDecoderBoolean, Input: "prior"}, want: true},
		{name: "recognised text true", input: Inputs{Values: map[string]Value{"prior": {Type: ValueTypeText, Text: "yes"}}}, config: config.DashboardDecoderConfig{ID: "bool", Type: config.DashboardDecoderBoolean, Input: "prior"}, want: true},
		{name: "recognised text false", input: Inputs{Values: map[string]Value{"prior": {Type: ValueTypeText, Text: "off"}}}, config: config.DashboardDecoderConfig{ID: "bool", Type: config.DashboardDecoderBoolean, Input: "prior"}, want: false},
	}

	registry := NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registry.Decode(tt.config, tt.input)
			if err != nil {
				t.Fatalf("Decode returned error: %v", err)
			}
			if got.Bool != tt.want {
				t.Fatalf("boolean = %v, want %v", got.Bool, tt.want)
			}
		})
	}
}

func TestDecoderErrors(t *testing.T) {
	tests := []struct {
		name    string
		decoder config.DashboardDecoderConfig
		inputs  Inputs
	}{
		{name: "unknown sensor", decoder: config.DashboardDecoderConfig{ID: "missing", Type: config.DashboardDecoderNormalize, Sensor: "missing"}, inputs: Inputs{Sensors: map[string]sensors.SensorState{}}},
		{name: "sensor error", decoder: config.DashboardDecoderConfig{ID: "bad", Type: config.DashboardDecoderNormalize, Sensor: "rpm"}, inputs: sensorInputs(sensors.SensorState{ID: "rpm", Status: sensors.StatusError, Error: "read failed"})},
		{name: "invalid normalize range", decoder: config.DashboardDecoderConfig{ID: "flat", Type: config.DashboardDecoderNormalize, Sensor: "rpm"}, inputs: sensorInputs(sensorState("rpm", 1, 1, 1))},
		{name: "bad frame count", decoder: config.DashboardDecoderConfig{ID: "frame", Type: config.DashboardDecoderFrameIndex, Sensor: "rpm"}, inputs: sensorInputs(sensorState("rpm", 1, 0, 1))},
		{name: "invalid threshold config", decoder: config.DashboardDecoderConfig{ID: "zone", Type: config.DashboardDecoderThreshold, Sensor: "rpm"}, inputs: sensorInputs(sensorState("rpm", 12, 0, 100))},
		{name: "invalid format mask", decoder: config.DashboardDecoderConfig{ID: "fmt", Type: config.DashboardDecoderFormatNumber, Sensor: "rpm", Format: "00x0"}, inputs: sensorInputs(sensorState("rpm", 12, 0, 100))},
		{name: "invalid go format", decoder: config.DashboardDecoderConfig{ID: "fmt", Type: config.DashboardDecoderFormatNumber, Sensor: "rpm", Format: "%s"}, inputs: sensorInputs(sensorState("rpm", 12, 0, 100))},
		{name: "digits rejects formatted decimal", decoder: config.DashboardDecoderConfig{ID: "digits", Type: config.DashboardDecoderDigits, Sensor: "rpm", Format: "%.1f"}, inputs: sensorInputs(sensorState("rpm", 12.3, 0, 100))},
		{name: "unknown prior input", decoder: config.DashboardDecoderConfig{ID: "derived", Type: config.DashboardDecoderBoolean, Input: "missing"}, inputs: Inputs{Values: map[string]Value{}}},
		{name: "unrecognised boolean text", decoder: config.DashboardDecoderConfig{ID: "bool", Type: config.DashboardDecoderBoolean, Input: "prior"}, inputs: Inputs{Values: map[string]Value{"prior": {Type: ValueTypeText, Text: "maybe"}}}},
	}

	registry := NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := registry.Decode(tt.decoder, tt.inputs); err == nil {
				t.Fatal("Decode returned nil error, want error")
			}
		})
	}
}

func sensorState(id string, value float64, min float64, max float64) sensors.SensorState {
	return sensors.SensorState{ID: id, Value: value, Min: min, Max: max, Status: sensors.StatusOK}
}

func sensorInputs(states ...sensors.SensorState) Inputs {
	inputs := Inputs{Sensors: map[string]sensors.SensorState{}}
	for _, state := range states {
		inputs.Sensors[state.ID] = state
	}
	return inputs
}

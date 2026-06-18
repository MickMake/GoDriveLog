package sensors

import (
	"strings"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
)

func TestPollingRuntimeRejectsInvalidConfiguredValueKind(t *testing.T) {
	_, err := NewPollingRuntime(&scriptedReader{}, map[string]v3config.SensorConfig{
		"rpm": {Type: "obd", PID: "010C", ValueKind: "banana", Unit: "rpm", Poll: 250},
	})
	if err == nil {
		t.Fatal("expected invalid value_kind to fail runtime setup")
	}
	if !strings.Contains(err.Error(), `sensor "rpm" has invalid value_kind "banana"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPollingRuntimeRejectsConfiguredKindMismatchWithParserOutput(t *testing.T) {
	_, err := NewPollingRuntime(&scriptedReader{}, map[string]v3config.SensorConfig{
		"rpm": {Type: "obd", PID: "010C", ValueKind: ValueKindBool, Unit: "rpm", Poll: 250},
	})
	if err == nil {
		t.Fatal("expected value_kind/parser mismatch to fail runtime setup")
	}
	if !strings.Contains(err.Error(), `sensor "rpm" value_kind "bool" is incompatible with parser output kind "numeric"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPollingRuntimeDerivesNumericKindFromOBDParserContract(t *testing.T) {
	runtime, err := NewPollingRuntime(&scriptedReader{}, map[string]v3config.SensorConfig{
		"rpm": {Type: "obd", PID: "010C", Unit: "rpm", Poll: 250},
	})
	if err != nil {
		t.Fatalf("NewPollingRuntime: %v", err)
	}
	if len(runtime.sensors) != 1 || runtime.sensors[0].ValueKind != ValueKindNumeric {
		t.Fatalf("runtime sensors = %#v, want derived numeric kind", runtime.sensors)
	}
}

package scene

import (
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/decoders"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestEvaluateConfiguredConditionFromDecoderValue(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()
	dashboard.Blocks = append(dashboard.Blocks, config.DashboardBlockConfig{
		ID:       "redline_glow",
		Type:     config.DashboardBlockImage,
		Asset:    "background",
		Geometry: config.RectConfig{Width: 800, Height: 480},
		Condition: config.DashboardConditionConfig{
			Decoder: "rpm_warning",
			Equals:  "warning",
		},
	})
	dashboard.Layers = append(dashboard.Layers, config.DashboardLayerConfig{ID: "alerts", Z: 20, Blocks: []string{"redline_glow"}})

	values := baseDecoderValues()
	values["rpm_warning"] = decoders.Value{Type: decoders.ValueTypeText, Text: "normal"}
	scene, err := Evaluate(dashboard, registry, values, nil, Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if findElement(t, scene.Elements, "redline_glow").Visible {
		t.Fatal("redline_glow visible for normal warning state, want hidden")
	}

	values["rpm_warning"] = decoders.Value{Type: decoders.ValueTypeText, Text: "warning"}
	scene, err = Evaluate(dashboard, registry, values, nil, Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if !findElement(t, scene.Elements, "redline_glow").Visible {
		t.Fatal("redline_glow hidden for warning state, want visible")
	}
}

func TestEvaluateConfiguredConditionFromSensorStatus(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()
	dashboard.Blocks = append(dashboard.Blocks, config.DashboardBlockConfig{
		ID:       "stale_indicator",
		Type:     config.DashboardBlockImage,
		Asset:    "background",
		Geometry: config.RectConfig{Width: 100, Height: 40},
		Condition: config.DashboardConditionConfig{
			Sensor: "rpm",
			Status: sensors.StatusStale,
		},
	})
	dashboard.Layers = append(dashboard.Layers, config.DashboardLayerConfig{ID: "status", Z: 30, Blocks: []string{"stale_indicator"}})

	scene, err := Evaluate(dashboard, registry, baseDecoderValues(), map[string]sensors.SensorState{
		"rpm": {ID: "rpm", Value: 900, Status: sensors.StatusOK},
	}, Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if findElement(t, scene.Elements, "stale_indicator").Visible {
		t.Fatal("stale_indicator visible for ok state, want hidden")
	}

	scene, err = Evaluate(dashboard, registry, baseDecoderValues(), map[string]sensors.SensorState{
		"rpm": {ID: "rpm", Value: 900, Status: sensors.StatusStale},
	}, Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if !findElement(t, scene.Elements, "stale_indicator").Visible {
		t.Fatal("stale_indicator hidden for stale state, want visible")
	}
}

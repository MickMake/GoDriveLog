package v3config

import "testing"

func TestResolveSingleVehicleDefaultsSelectedVehicle(t *testing.T) {
	cfg := loadResolveConfig(t, validMinimalYAML())
	plan, err := Resolve(cfg, "")
	if err != nil {
		t.Fatalf("expected single vehicle to resolve by default: %v", err)
	}
	if plan.VehicleID != "vw_caddy" {
		t.Fatalf("expected vw_caddy, got %q", plan.VehicleID)
	}
	if plan.Endpoint.Address != "serial:///dev/ttyUSB0" {
		t.Fatalf("expected selected endpoint, got %q", plan.Endpoint.Address)
	}
	assertPlanLogs(t, plan, []string{"jsonl"})
	assertPlanDashboards(t, plan, []string{"simple_primary"})
}

func TestResolveMultipleVehiclesRequiresExplicitSelection(t *testing.T) {
	cfg := loadResolveConfig(t, runtimePlanMultiVehicleYAML())
	_, err := Resolve(cfg, "")
	if err == nil {
		t.Fatalf("expected multiple vehicles to require explicit selection")
	}
	assertErrorContains(t, err, "selected vehicle must be explicit")
}

func TestResolveSelectedVehicleMustExist(t *testing.T) {
	cfg := loadResolveConfig(t, validMinimalYAML())
	_, err := Resolve(cfg, "missing_vehicle")
	if err == nil {
		t.Fatalf("expected missing selected vehicle to fail")
	}
	assertErrorContains(t, err, "missing_vehicle")
}

func TestResolveReturnsOnlySelectedLogsAndDashboards(t *testing.T) {
	cfg := loadResolveConfig(t, runtimePlanMultiVehicleYAML())
	plan, err := Resolve(cfg, "bench_z31")
	if err != nil {
		t.Fatalf("expected bench_z31 to resolve: %v", err)
	}
	if plan.Endpoint.Address != "tcp://127.0.0.1:35000" {
		t.Fatalf("expected bench endpoint, got %q", plan.Endpoint.Address)
	}
	assertPlanLogs(t, plan, []string{"bench_jsonl"})
	assertPlanDashboards(t, plan, []string{"bench_primary"})
}

func TestResolveAllowsUnselectedDashboardsToShareDisplay(t *testing.T) {
	cfg := loadResolveConfig(t, runtimePlanMultiVehicleYAML())
	plan, err := Resolve(cfg, "vw_caddy")
	if err != nil {
		t.Fatalf("expected selected dashboard to resolve despite unselected alternative sharing display: %v", err)
	}
	assertPlanDashboards(t, plan, []string{"simple_primary"})
}

func TestResolveRejectsSelectedDashboardDisplayCollision(t *testing.T) {
	_, err := LoadBytes([]byte(runtimePlanDisplayCollisionYAML()))
	if err == nil {
		t.Fatalf("expected selected dashboard display collision to fail")
	}
	assertErrorContains(t, err, "both target display")
}

func TestResolveSingleLogAndDashboardDefaults(t *testing.T) {
	cfg := loadResolveConfig(t, runtimePlanImplicitSelectionsYAML())
	plan, err := Resolve(cfg, "vw_caddy")
	if err != nil {
		t.Fatalf("expected single log/dashboard defaults to resolve: %v", err)
	}
	assertPlanLogs(t, plan, []string{"jsonl"})
	assertPlanDashboards(t, plan, []string{"simple_primary"})
}

func TestResolveExposesSensorCatalogueForSelectedLogs(t *testing.T) {
	cfg := loadResolveConfig(t, runtimePlanMultiVehicleYAML())
	plan, err := Resolve(cfg, "bench_z31")
	if err != nil {
		t.Fatalf("expected bench_z31 to resolve: %v", err)
	}

	for _, log := range plan.Logs {
		for _, sensorID := range log.Config.Sensors {
			sensor, ok := plan.Sensors[sensorID]
			if !ok {
				t.Fatalf("expected plan sensors to expose log sensor %q", sensorID)
			}
			if sensor.Type == "" || sensor.Poll <= 0 {
				t.Fatalf("expected usable sensor config for %q, got %+v", sensorID, sensor)
			}
		}
	}
}

func TestResolveExposesCataloguesForSelectedDashboardWidgets(t *testing.T) {
	cfg := loadResolveConfig(t, runtimePlanDigitDashboardYAML())
	plan, err := Resolve(cfg, "vw_caddy")
	if err != nil {
		t.Fatalf("expected vw_caddy to resolve: %v", err)
	}
	if len(plan.Dashboards) != 1 {
		t.Fatalf("expected one dashboard, got %d", len(plan.Dashboards))
	}

	for _, widget := range plan.Dashboards[0].Config.Widgets {
		if widget.Sensor != "" {
			if _, ok := plan.Sensors[widget.Sensor]; !ok {
				t.Fatalf("expected plan sensors to expose widget sensor %q", widget.Sensor)
			}
		}
		switch widget.Type {
		case WidgetTypeImage:
			if _, ok := plan.Assets.ImageSets[widget.Asset]; !ok {
				t.Fatalf("expected plan image assets to expose %q", widget.Asset)
			}
		case WidgetTypeDigitDisplay:
			if _, ok := plan.Assets.DigitSets[widget.Asset]; !ok {
				t.Fatalf("expected plan digit assets to expose %q", widget.Asset)
			}
		default:
			t.Fatalf("unexpected widget type %q in test", widget.Type)
		}
	}
}

func loadResolveConfig(t *testing.T, text string) Config {
	t.Helper()
	cfg, err := LoadBytes([]byte(text))
	if err != nil {
		t.Fatalf("expected config to load: %v", err)
	}
	return cfg
}

func assertPlanLogs(t *testing.T, plan RuntimePlan, want []string) {
	t.Helper()
	if len(plan.Logs) != len(want) {
		t.Fatalf("expected %d logs, got %d", len(want), len(plan.Logs))
	}
	for i, id := range want {
		if plan.Logs[i].ID != id {
			t.Fatalf("expected log %d to be %q, got %q", i, id, plan.Logs[i].ID)
		}
	}
}

func assertPlanDashboards(t *testing.T, plan RuntimePlan, want []string) {
	t.Helper()
	if len(plan.Dashboards) != len(want) {
		t.Fatalf("expected %d dashboards, got %d", len(want), len(plan.Dashboards))
	}
	for i, id := range want {
		if plan.Dashboards[i].ID != id {
			t.Fatalf("expected dashboard %d to be %q, got %q", i, id, plan.Dashboards[i].ID)
		}
	}
}

func runtimePlanImplicitSelectionsYAML() string {
	return `vehicles:
  vw_caddy:
    name: VW Caddy
    obd:
      address: serial:///dev/ttyUSB0
      timeout: 1000
sensors:
  speed:
    type: obd
    pid: "010D"
    unit: km/h
    poll: 250
assets:
  image_sets:
    panel:
      image: assets/dashboard/simple/panel/background.png
logs:
  jsonl:
    path: logs/godrivelog.jsonl
    sensors:
      - speed
dashboards:
  simple_primary:
    display: HDMI-1
    size:
      width: 800
      height: 480
    widgets:
      - id: panel_backplate
        type: image
        asset: panel
        position: [0, 0]`
}

func runtimePlanMultiVehicleYAML() string {
	return `vehicles:
  vw_caddy:
    name: VW Caddy
    obd:
      address: serial:///dev/ttyUSB0
      timeout: 1000
    logs:
      - jsonl
    dashboards:
      - simple_primary
  bench_z31:
    name: Bench Z31
    obd:
      address: tcp://127.0.0.1:35000
      timeout: 1000
    logs:
      - bench_jsonl
    dashboards:
      - bench_primary
sensors:
  speed:
    type: obd
    pid: "010D"
    unit: km/h
    poll: 250
assets:
  image_sets:
    panel:
      image: assets/dashboard/simple/panel/background.png
logs:
  jsonl:
    path: logs/godrivelog.jsonl
    sensors:
      - speed
  bench_jsonl:
    path: logs/bench.jsonl
    sensors:
      - speed
dashboards:
  simple_primary:
    display: HDMI-1
    size:
      width: 800
      height: 480
    widgets:
      - id: panel_backplate
        type: image
        asset: panel
        position: [0, 0]
  bench_primary:
    display: HDMI-1
    size:
      width: 800
      height: 480
    widgets:
      - id: bench_panel
        type: image
        asset: panel
        position: [0, 0]`
}

func runtimePlanDisplayCollisionYAML() string {
	return `vehicles:
  vw_caddy:
    name: VW Caddy
    obd:
      address: serial:///dev/ttyUSB0
      timeout: 1000
    logs:
      - jsonl
    dashboards:
      - simple_primary
      - bench_primary
sensors:
  speed:
    type: obd
    pid: "010D"
    unit: km/h
    poll: 250
assets:
  image_sets:
    panel:
      image: assets/dashboard/simple/panel/background.png
logs:
  jsonl:
    path: logs/godrivelog.jsonl
    sensors:
      - speed
dashboards:
  simple_primary:
    display: HDMI-1
    size:
      width: 800
      height: 480
    widgets:
      - id: panel_backplate
        type: image
        asset: panel
        position: [0, 0]
  bench_primary:
    display: HDMI-1
    size:
      width: 800
      height: 480
    widgets:
      - id: bench_panel
        type: image
        asset: panel
        position: [0, 0]`
}

func runtimePlanDigitDashboardYAML() string {
	return `vehicles:
  vw_caddy:
    name: VW Caddy
    obd:
      address: serial:///dev/ttyUSB0
      timeout: 1000
    logs:
      - jsonl
    dashboards:
      - simple_primary
sensors:
  speed:
    type: obd
    pid: "010D"
    unit: km/h
    poll: 250
assets:
  digit_sets:
    digits:
      characters:
        "0": assets/dashboard/simple/digits/0.png
        "1": assets/dashboard/simple/digits/1.png
        "2": assets/dashboard/simple/digits/2.png
        "3": assets/dashboard/simple/digits/3.png
        "4": assets/dashboard/simple/digits/4.png
        "5": assets/dashboard/simple/digits/5.png
        "6": assets/dashboard/simple/digits/6.png
        "7": assets/dashboard/simple/digits/7.png
        "8": assets/dashboard/simple/digits/8.png
        "9": assets/dashboard/simple/digits/9.png
  image_sets:
    panel:
      image: assets/dashboard/simple/panel/background.png
logs:
  jsonl:
    path: logs/godrivelog.jsonl
    sensors:
      - speed
dashboards:
  simple_primary:
    display: HDMI-1
    size:
      width: 800
      height: 480
    widgets:
      - id: panel_backplate
        type: image
        asset: panel
        position: [0, 0]
      - id: speed_digits
        type: digit_display
        sensor: speed
        asset: digits
        position: [40, 40]
        digits: 3
        format: "%03.0f"`
}

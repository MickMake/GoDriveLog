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

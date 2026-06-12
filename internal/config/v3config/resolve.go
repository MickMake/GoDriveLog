package v3config

import "fmt"

// RuntimePlan is the resolved v3 boundary handed to later runtime packages.
// It keeps selected vehicle ownership explicit and exposes the validated global
// catalogues so callers do not need to walk raw Config maps.
type RuntimePlan struct {
	VehicleID  string
	Vehicle    VehicleConfig
	Endpoint   OBDConfig
	Sensors    map[string]SensorConfig
	Assets     AssetConfig
	Logs       []ResolvedLog
	Dashboards []ResolvedDashboard
}

type ResolvedLog struct {
	ID     string
	Config LogConfig
}

type ResolvedDashboard struct {
	ID     string
	Config DashboardConfig
}

// Resolve turns a validated v3 config and selected vehicle ID into the explicit
// runtime plan used by later migration slices.
func Resolve(cfg Config, selectedVehicleID string) (RuntimePlan, error) {
	if err := Validate(cfg); err != nil {
		return RuntimePlan{}, err
	}

	vehicleID, err := resolveVehicleID(cfg, selectedVehicleID)
	if err != nil {
		return RuntimePlan{}, err
	}

	vehicle := cfg.Vehicles[vehicleID]
	logs, err := resolveSelectedLogs(cfg, vehicle)
	if err != nil {
		return RuntimePlan{}, err
	}
	dashboards, err := resolveSelectedDashboards(cfg, vehicle)
	if err != nil {
		return RuntimePlan{}, err
	}
	if err := validateResolvedDashboardDisplays(dashboards); err != nil {
		return RuntimePlan{}, err
	}

	return RuntimePlan{
		VehicleID:  vehicleID,
		Vehicle:    vehicle,
		Endpoint:   vehicle.OBD,
		Sensors:    cfg.Sensors,
		Assets:     cfg.Assets,
		Logs:       logs,
		Dashboards: dashboards,
	}, nil
}

func resolveVehicleID(cfg Config, selectedVehicleID string) (string, error) {
	if selectedVehicleID != "" {
		if _, ok := cfg.Vehicles[selectedVehicleID]; !ok {
			return "", ValidationError{Errors: []string{fmt.Sprintf("selected vehicle %q must exist", selectedVehicleID)}}
		}
		return selectedVehicleID, nil
	}

	if len(cfg.Vehicles) == 1 {
		for vehicleID := range cfg.Vehicles {
			return vehicleID, nil
		}
	}

	return "", ValidationError{Errors: []string{"selected vehicle must be explicit when multiple vehicles exist"}}
}

func resolveSelectedLogs(cfg Config, vehicle VehicleConfig) ([]ResolvedLog, error) {
	logIDs := vehicle.Logs
	if len(logIDs) == 0 && len(cfg.Logs) == 1 {
		for logID := range cfg.Logs {
			logIDs = []string{logID}
		}
	}

	logs := make([]ResolvedLog, 0, len(logIDs))
	for _, logID := range logIDs {
		logCfg, ok := cfg.Logs[logID]
		if !ok {
			return nil, ValidationError{Errors: []string{fmt.Sprintf("selected log %q must reference top-level logs", logID)}}
		}
		logs = append(logs, ResolvedLog{ID: logID, Config: logCfg})
	}
	return logs, nil
}

func resolveSelectedDashboards(cfg Config, vehicle VehicleConfig) ([]ResolvedDashboard, error) {
	dashboardIDs := vehicle.Dashboards
	if len(dashboardIDs) == 0 && len(cfg.Dashboards) == 1 {
		for dashboardID := range cfg.Dashboards {
			dashboardIDs = []string{dashboardID}
		}
	}

	dashboards := make([]ResolvedDashboard, 0, len(dashboardIDs))
	for _, dashboardID := range dashboardIDs {
		dashboardCfg, ok := cfg.Dashboards[dashboardID]
		if !ok {
			return nil, ValidationError{Errors: []string{fmt.Sprintf("selected dashboard %q must reference top-level dashboards", dashboardID)}}
		}
		dashboards = append(dashboards, ResolvedDashboard{ID: dashboardID, Config: dashboardCfg})
	}
	return dashboards, nil
}

func validateResolvedDashboardDisplays(dashboards []ResolvedDashboard) error {
	seenDisplays := map[string]string{}
	var errs []string
	for _, dashboard := range dashboards {
		display := dashboard.Config.Display
		if display == "" {
			continue
		}
		if existing, ok := seenDisplays[display]; ok {
			errs = append(errs, fmt.Sprintf("selected dashboards %q and %q both target display %q", existing, dashboard.ID, display))
			continue
		}
		seenDisplays[display] = dashboard.ID
	}
	if len(errs) > 0 {
		return ValidationError{Errors: errs}
	}
	return nil
}

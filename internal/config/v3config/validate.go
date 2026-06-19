package v3config

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type ValidationError struct {
	Errors []string
}

func (e ValidationError) Error() string {
	return strings.Join(e.Errors, "; ")
}

func Validate(cfg Config) error {
	v := validator{}
	v.validateConfig(cfg)
	if len(v.errs) > 0 {
		return ValidationError{Errors: v.errs}
	}
	return nil
}

func ValidateSelectedVehicle(cfg Config, vehicleID string) error {
	v := validator{}
	v.validateConfig(cfg)
	if vehicleID == "" {
		if len(cfg.Vehicles) == 1 {
			for id := range cfg.Vehicles {
				vehicleID = id
			}
		} else {
			v.add("selected vehicle must be explicit when multiple vehicles exist")
		}
	}
	if vehicleID != "" {
		if _, ok := cfg.Vehicles[vehicleID]; !ok {
			v.add("selected vehicle %q must exist", vehicleID)
		}
	}
	if len(v.errs) > 0 {
		return ValidationError{Errors: v.errs}
	}
	return nil
}

type validator struct {
	errs []string
}

func (v *validator) add(format string, args ...any) {
	v.errs = append(v.errs, fmt.Sprintf(format, args...))
}

func (v *validator) validateConfig(cfg Config) {
	if len(cfg.Vehicles) == 0 {
		v.add("vehicles must contain at least one vehicle")
	}
	v.validateIDs("vehicles", keysOf(cfg.Vehicles))
	v.validateIDs("sensors", keysOf(cfg.Sensors))
	v.validateIDs("assets.digit_sets", keysOf(cfg.Assets.DigitSets))
	v.validateIDs("assets.bar_sets", keysOf(cfg.Assets.BarSets))
	v.validateIDs("assets.frame_sets", keysOf(cfg.Assets.FrameSets))
	v.validateIDs("assets.indicator_sets", keysOf(cfg.Assets.IndicatorSets))
	v.validateIDs("assets.image_sets", keysOf(cfg.Assets.ImageSets))
	v.validateIDs("logs", keysOf(cfg.Logs))
	v.validateIDs("dashboards", keysOf(cfg.Dashboards))

	for id, vehicle := range sortedMap(cfg.Vehicles) {
		v.validateVehicle("vehicles."+id, vehicle, cfg)
	}
	for id, sensor := range sortedMap(cfg.Sensors) {
		v.validateSensor("sensors."+id, sensor)
	}
	v.validateAssets(cfg.Assets)
	for id, log := range sortedMap(cfg.Logs) {
		v.validateLog("logs."+id, log, cfg.Sensors)
	}
	for id, dashboard := range sortedMap(cfg.Dashboards) {
		v.validateDashboard("dashboards."+id, dashboard, cfg)
	}
}

func (v *validator) validateIDs(path string, ids []string) {
	for _, id := range ids {
		if !idPattern.MatchString(id) {
			v.add("%s id %q must match ^[a-z][a-z0-9_]*$", path, id)
		}
	}
}

func (v *validator) validateVehicle(path string, vehicle VehicleConfig, cfg Config) {
	if strings.TrimSpace(vehicle.Name) == "" {
		v.add("%s.name must not be empty", path)
	}
	if err := validateEndpoint(vehicle.OBD.Address); err != nil {
		v.add("%s.obd.address %s", path, err)
	}
	if vehicle.OBD.Timeout < 100 || vehicle.OBD.Timeout > 30000 {
		v.add("%s.obd.timeout must be within 100..30000 milliseconds", path)
	}
	if len(cfg.Logs) > 1 && len(vehicle.Logs) == 0 {
		v.add("%s.logs must list selected logs when multiple logs are defined", path)
	}
	for _, logID := range vehicle.Logs {
		if _, ok := cfg.Logs[logID]; !ok {
			v.add("%s.logs %q must reference top-level logs", path, logID)
		}
	}
	if len(cfg.Dashboards) > 1 && len(vehicle.Dashboards) == 0 {
		v.add("%s.dashboards must list selected dashboards when multiple dashboards are defined", path)
	}
	seenDisplays := map[string]string{}
	for _, dashboardID := range vehicle.Dashboards {
		dashboard, ok := cfg.Dashboards[dashboardID]
		if !ok {
			v.add("%s.dashboards %q must reference top-level dashboards", path, dashboardID)
			continue
		}
		if dashboard.Display != "" {
			if existing, exists := seenDisplays[dashboard.Display]; exists {
				v.add("%s selected dashboards %q and %q both target display %q", path, existing, dashboardID, dashboard.Display)
			} else {
				seenDisplays[dashboard.Display] = dashboardID
			}
		}
	}
}

func validateEndpoint(address string) error {
	if strings.TrimSpace(address) == "" {
		return fmt.Errorf("must not be empty")
	}
	u, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("must parse as URL: %w", err)
	}
	switch u.Scheme {
	case "serial":
		if strings.TrimSpace(u.Path) == "" && strings.TrimSpace(u.Opaque) == "" {
			return fmt.Errorf("must include a non-empty serial path")
		}
	case "tcp":
		if strings.TrimSpace(u.Hostname()) == "" || strings.TrimSpace(u.Port()) == "" {
			return fmt.Errorf("must include host and port")
		}
	default:
		return fmt.Errorf("must use serial:// or tcp://")
	}
	return nil
}

func (v *validator) validateSensor(path string, sensor SensorConfig) {
	if sensor.Type != "obd" {
		v.add("%s.type must be obd", path)
	}
	if sensor.Type == "obd" && strings.TrimSpace(sensor.PID) == "" {
		v.add("%s.pid must not be empty for obd sensors", path)
	}
	if strings.TrimSpace(sensor.Unit) == "" {
		v.add("%s.unit must not be empty", path)
	}
	if sensor.Poll <= 0 {
		v.add("%s.poll must be greater than zero", path)
	}
	if sensor.Min != nil && sensor.Max != nil && *sensor.Min >= *sensor.Max {
		v.add("%s.min must be less than max", path)
	}
}

func (v *validator) validateAssets(assets AssetConfig) {
	for id, set := range sortedMap(assets.DigitSets) {
		path := "assets.digit_sets." + id
		for _, ch := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"} {
			if set.Characters[ch] == "" {
				v.add("%s.characters must include %q", path, ch)
			}
		}
		v.validateOptionalAssetPath(path+".background", set.Background)
		v.validateOptionalAssetPath(path+".decimal_point", set.DecimalPoint)
		v.validateOptionalAssetPath(path+".foreground", set.Foreground)
		for ch, assetPath := range set.Characters {
			v.validateRequiredAssetPath(path+".characters."+ch, assetPath)
		}
	}
	for id, set := range sortedMap(assets.BarSets) {
		path := "assets.bar_sets." + id
		if set.Cells["off"] == "" {
			v.add("%s.cells.off must exist", path)
		}
		v.validateOptionalAssetPath(path+".background", set.Background)
		v.validateOptionalAssetPath(path+".foreground", set.Foreground)
		for cell, assetPath := range set.Cells {
			v.validateRequiredAssetPath(path+".cells."+cell, assetPath)
		}
	}
	for id, set := range sortedMap(assets.FrameSets) {
		path := "assets.frame_sets." + id
		v.validateOptionalAssetPath(path+".background", set.Background)
		v.validateRequiredAssetPath(path+".frames.path", set.Frames.Path)
		v.validateOptionalAssetPath(path+".foreground", set.Foreground)
		if set.Frames.First > set.Frames.Last {
			v.add("%s.frames.first must be less than or equal to last", path)
		}
	}
	for id, set := range sortedMap(assets.IndicatorSets) {
		path := "assets.indicator_sets." + id
		for _, state := range []string{"off", "on", "unknown"} {
			if set.States[state] == "" {
				v.add("%s.states.%s must exist", path, state)
			}
		}
		v.validateOptionalAssetPath(path+".background", set.Background)
		v.validateOptionalAssetPath(path+".foreground", set.Foreground)
		for state, assetPath := range set.States {
			v.validateRequiredAssetPath(path+".states."+state, assetPath)
		}
	}
	for id, set := range sortedMap(assets.ImageSets) {
		path := "assets.image_sets." + id
		if set.Image == "" && set.Background == "" && set.Foreground == "" {
			v.add("%s must define image, background, or foreground", path)
		}
		v.validateOptionalAssetPath(path+".image", set.Image)
		v.validateOptionalAssetPath(path+".background", set.Background)
		v.validateOptionalAssetPath(path+".foreground", set.Foreground)
	}
}

func (v *validator) validateRequiredAssetPath(path, assetPath string) {
	if strings.TrimSpace(assetPath) == "" {
		v.add("%s must not be empty", path)
		return
	}
	v.validateAssetPath(path, assetPath)
}

func (v *validator) validateOptionalAssetPath(path, assetPath string) {
	if strings.TrimSpace(assetPath) == "" {
		return
	}
	v.validateAssetPath(path, assetPath)
}

func (v *validator) validateAssetPath(pathName, assetPath string) {
	trimmed := strings.TrimSpace(assetPath)
	slashPath := filepath.ToSlash(trimmed)
	cleaned := path.Clean(slashPath)

	if strings.Contains(trimmed, "://") {
		v.add("%s must be repository-root relative, not remote or URL-like", pathName)
	}
	if filepath.IsAbs(trimmed) || path.IsAbs(slashPath) {
		v.add("%s must be repository-root relative", pathName)
	}
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || hasUpwardEscapeSegment(slashPath) {
		v.add("%s must be repository-root relative", pathName)
	}
}

func (v *validator) validateGaugePath(pathName, gaugePath string) {
	v.validateAssetPath(pathName, gaugePath)
	cleaned := path.Clean(filepath.ToSlash(strings.TrimSpace(gaugePath)))
	if cleaned == "assets/gauges" || !strings.HasPrefix(cleaned, "assets/gauges/") {
		v.add("%s must reference a package directory under assets/gauges", pathName)
	}
	if strings.HasSuffix(cleaned, "/gauge.yaml") {
		v.add("%s must reference a gauge package directory, not gauge.yaml", pathName)
	}
}

func hasUpwardEscapeSegment(slashPath string) bool {
	for _, segment := range strings.Split(slashPath, "/") {
		if segment == ".." {
			return true
		}
	}
	return false
}

func (v *validator) validateLog(path string, log LogConfig, sensors map[string]SensorConfig) {
	if strings.TrimSpace(log.Path) == "" {
		v.add("%s.path must not be empty", path)
	}
	if len(log.Sensors) == 0 {
		v.add("%s.sensors must not be empty", path)
	}
	for _, sensorID := range log.Sensors {
		if _, ok := sensors[sensorID]; !ok {
			v.add("%s.sensors %q must reference top-level sensors", path, sensorID)
		}
	}
}

func (v *validator) validateDashboard(path string, dashboard DashboardConfig, cfg Config) {
	if strings.TrimSpace(dashboard.Display) == "" {
		v.add("%s.display must not be empty", path)
	}
	if dashboard.Size.Width <= 0 {
		v.add("%s.size.width must be positive", path)
	}
	if dashboard.Size.Height <= 0 {
		v.add("%s.size.height must be positive", path)
	}
	seenWidgets := map[string]bool{}
	for i, widget := range dashboard.Widgets {
		widgetPath := fmt.Sprintf("%s.widgets[%d]", path, i)
		if !idPattern.MatchString(widget.ID) {
			v.add("%s.id %q must match ^[a-z][a-z0-9_]*$", widgetPath, widget.ID)
		}
		if seenWidgets[widget.ID] {
			v.add("%s id %q must be unique within dashboard", widgetPath, widget.ID)
		}
		seenWidgets[widget.ID] = true
		v.validateWidget(widgetPath, widget, cfg)
	}
}

func (v *validator) validateWidget(path string, widget WidgetConfig, cfg Config) {
	if !isKnownWidgetType(widget.Type) {
		v.add("%s.type must be image, digit_display, bar_display, frame_gauge, indicator, or gauge", path)
	}
	if len(widget.Position) != 2 {
		v.add("%s.position must contain exactly two integers", path)
	}
	if widget.Min != nil && widget.Max != nil && *widget.Min >= *widget.Max {
		v.add("%s.min must be less than max", path)
	}

	if widget.Type == WidgetTypeGauge {
		v.validateGaugeWidget(path, widget)
		return
	}
	if strings.TrimSpace(widget.Gauge) != "" {
		v.add("%s.gauge must be empty for non-gauge widgets", path)
	}

	if widget.Type != WidgetTypeImage {
		if strings.TrimSpace(widget.Sensor) == "" {
			v.add("%s.sensor must not be empty for non-image widgets", path)
		} else if _, ok := cfg.Sensors[widget.Sensor]; !ok {
			v.add("%s.sensor %q must reference top-level sensors", path, widget.Sensor)
		}
	}
	if strings.TrimSpace(widget.Asset) == "" {
		v.add("%s.asset must not be empty", path)
		return
	}
	switch widget.Type {
	case WidgetTypeImage:
		if _, ok := cfg.Assets.ImageSets[widget.Asset]; !ok {
			v.add("%s.asset %q must reference assets.image_sets", path, widget.Asset)
		}
	case WidgetTypeDigitDisplay:
		set, ok := cfg.Assets.DigitSets[widget.Asset]
		if !ok {
			v.add("%s.asset %q must reference assets.digit_sets", path, widget.Asset)
		}
		if widget.Digits <= 0 {
			v.add("%s.digits must be greater than zero", path)
		}
		if ok && formatUsesDecimalPoint(widget.Format) && set.DecimalPoint == "" {
			v.add("%s.format requires digit set %q to define decimal_point", path, widget.Asset)
		}
	case WidgetTypeBarDisplay:
		set, ok := cfg.Assets.BarSets[widget.Asset]
		if !ok {
			v.add("%s.asset %q must reference assets.bar_sets", path, widget.Asset)
		}
		if widget.Cells <= 0 {
			v.add("%s.cells must be greater than zero", path)
		}
		if ok {
			if len(widget.Zones) == 0 && set.Cells["on"] == "" {
				v.add("%s requires bar set %q to define cells.on when zones are omitted", path, widget.Asset)
			}
			v.validateZones(path, widget.Zones, set)
		}
	case WidgetTypeFrameGauge:
		if _, ok := cfg.Assets.FrameSets[widget.Asset]; !ok {
			v.add("%s.asset %q must reference assets.frame_sets", path, widget.Asset)
		}
	case WidgetTypeIndicator:
		if _, ok := cfg.Assets.IndicatorSets[widget.Asset]; !ok {
			v.add("%s.asset %q must reference assets.indicator_sets", path, widget.Asset)
		}
	}
}

func (v *validator) validateGaugeWidget(path string, widget WidgetConfig) {
	if strings.TrimSpace(widget.Sensor) != "" {
		v.add("%s.sensor must be empty for gauge widgets", path)
	}
	if strings.TrimSpace(widget.Asset) != "" {
		v.add("%s.asset must be empty for gauge widgets", path)
	}
	if strings.TrimSpace(widget.Gauge) == "" {
		v.add("%s.gauge must not be empty", path)
	} else {
		v.validateGaugePath(path+".gauge", widget.Gauge)
	}
	if widget.Scale <= 0 {
		v.add("%s.scale must be greater than zero for gauge widgets", path)
	}
}

func (v *validator) validateZones(path string, zones []ZoneConfig, set BarSetConfig) {
	var previous float64
	for i, zone := range zones {
		zonePath := fmt.Sprintf("%s.zones[%d]", path, i)
		if i > 0 && zone.UpTo < previous {
			v.add("%s.up_to must be sorted ascending", zonePath)
		}
		previous = zone.UpTo
		if strings.TrimSpace(zone.Cell) == "" {
			v.add("%s.cell must not be empty", zonePath)
		} else if set.Cells[zone.Cell] == "" {
			v.add("%s.cell %q must exist in bar set cells", zonePath, zone.Cell)
		}
	}
}

func isKnownWidgetType(widgetType string) bool {
	switch widgetType {
	case WidgetTypeImage, WidgetTypeDigitDisplay, WidgetTypeBarDisplay, WidgetTypeFrameGauge, WidgetTypeIndicator, WidgetTypeGauge:
		return true
	default:
		return false
	}
}

func formatUsesDecimalPoint(format string) bool {
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		i++
		if i >= len(format) {
			return false
		}
		if format[i] == '%' {
			continue
		}

		for i < len(format) && strings.ContainsRune("#+- 0", rune(format[i])) {
			i++
		}
		for i < len(format) && isASCIIDigit(format[i]) {
			i++
		}

		precisionSet := false
		precision := 0
		if i < len(format) && format[i] == '.' {
			precisionSet = true
			i++
			for i < len(format) && isASCIIDigit(format[i]) {
				precision = precision*10 + int(format[i]-'0')
				i++
			}
		}

		if i < len(format) && (format[i] == 'f' || format[i] == 'F') {
			return !precisionSet || precision > 0
		}
	}
	return false
}

func isASCIIDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func keysOf[T any](items map[string]T) []string {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedMap[T any](items map[string]T) map[string]T {
	return items
}

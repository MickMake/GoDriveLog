package gauges

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	dashboardConfigEnvVar    = "GODRIVELOG_CONFIG_PATH"
	TypeNumeric              = "numeric"
	TypeRadial               = "radial"
	TypeOdometer             = "odometer"
	TypeIndicator            = "indicator"
	TypeBar                  = "bar"
	TypeSegmented            = "segmented"
	MovementInstant          = "instant"
	MovementLinear           = "linear"
	MovementEaseOut          = "ease_out"
	MovementBell             = "bell"
	MovementSmooth           = "smooth"
	MovementClick            = "click"
	MovementPolicyImmediate  = "immediate"
	MovementPolicyLinear     = "linear"
	MovementPolicyEaseOut    = "ease_out"
	OvershootSettleSmooth    = "smooth"
	OvershootSettleOscillate = "oscillate"
	WheelRoleDigit           = "digit"
	WheelRoleSubUnit         = "sub_unit"
	defaultOvershootRatio    = 0.12
	defaultNeedleShadowAlpha = 0.35
	maxOvershootRatio        = 0.25
)

type Package struct {
	ID        string            `yaml:"id"`
	Type      string            `yaml:"type"`
	Sensor    string            `yaml:"sensor"`
	Format    string            `yaml:"format,omitempty"`
	Realism   Realism           `yaml:"realism,omitempty"`
	Size      Size              `yaml:"size"`
	Layers    map[string]string `yaml:"layers,omitempty"`
	DigitSet  DigitSet          `yaml:"digit_set,omitempty"`
	Digits    Digits            `yaml:"digits,omitempty"`
	Pivot     Pivot             `yaml:"pivot,omitempty"`
	ValueMap  ValueMap          `yaml:"value_map,omitempty"`
	Odometer  Odometer          `yaml:"odometer,omitempty"`
	Bar       BarConfig         `yaml:"bar,omitempty"`
	Segmented Segmented         `yaml:"segmented,omitempty"`
	Path      string            `yaml:"-"`
	YAMLPath  string            `yaml:"-"`
	AssetRoot string            `yaml:"-"`
}

type Size struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type DigitSet struct {
	Background   string            `yaml:"background,omitempty"`
	Characters   map[string]string `yaml:"characters,omitempty"`
	DecimalPoint string            `yaml:"decimal_point,omitempty"`
	Foreground   string            `yaml:"foreground,omitempty"`
	Spacing      int               `yaml:"spacing,omitempty"`
}

type Digits struct {
	Count     int     `yaml:"count"`
	Positions [][]int `yaml:"positions,omitempty"`
}

type Pivot struct {
	Face   Point `yaml:"face,omitempty"`
	Needle Point `yaml:"needle,omitempty"`
}

type Point struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
}

type ValueMap struct {
	Min        float64 `yaml:"min"`
	Max        float64 `yaml:"max"`
	StartAngle float64 `yaml:"start_angle"`
	EndAngle   float64 `yaml:"end_angle"`
	Clamp      bool    `yaml:"clamp"`
}

type OvershootConfig struct {
	Ratio          *float64 `yaml:"ratio,omitempty"`
	MinChangeRatio *float64 `yaml:"min_change_ratio,omitempty"`
	MaxSpanRatio   *float64 `yaml:"max_span_ratio,omitempty"`
	SettleMode     string   `yaml:"settle_mode,omitempty"`
	SettleCycles   *float64 `yaml:"settle_cycles,omitempty"`
	SettleDamping  *float64 `yaml:"settle_damping,omitempty"`
	AllowExtremes  bool     `yaml:"allow_extremes,omitempty"`
}

func (o *OvershootConfig) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("realism overshoot must be a mapping")
	}
	allowedKeys := map[string]bool{
		"ratio":            true,
		"min_change_ratio": true,
		"max_span_ratio":   true,
		"settle_mode":      true,
		"settle_cycles":    true,
		"settle_damping":   true,
		"allow_extremes":   true,
	}
	for index := 0; index+1 < len(node.Content); index += 2 {
		key := node.Content[index].Value
		if !allowedKeys[key] {
			return fmt.Errorf("realism overshoot field %q is not supported", key)
		}
	}
	type rawOvershootConfig OvershootConfig
	var decoded rawOvershootConfig
	if err := node.Decode(&decoded); err != nil {
		return err
	}
	*o = OvershootConfig(decoded)
	return nil
}

type DampingConfig struct {
	Enabled   bool `yaml:"enabled,omitempty"`
	RiseMS    int  `yaml:"rise_ms,omitempty"`
	FallMS    int  `yaml:"fall_ms,omitempty"`
	RiseMSSet bool `yaml:"-"`
	FallMSSet bool `yaml:"-"`
}

func (d *DampingConfig) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var enabled bool
		if err := node.Decode(&enabled); err != nil {
			return fmt.Errorf("realism damping must be a boolean or mapping")
		}
		*d = DampingConfig{Enabled: enabled}
		return nil
	case yaml.MappingNode:
		allowedKeys := map[string]bool{
			"enabled": true,
			"rise_ms": true,
			"fall_ms": true,
		}
		riseSet := false
		fallSet := false
		enabledSet := false
		for index := 0; index+1 < len(node.Content); index += 2 {
			key := node.Content[index].Value
			if !allowedKeys[key] {
				return fmt.Errorf("realism damping field %q is not supported", key)
			}
			switch key {
			case "enabled":
				enabledSet = true
			case "rise_ms":
				riseSet = true
			case "fall_ms":
				fallSet = true
			}
		}
		type rawDampingConfig DampingConfig
		var decoded rawDampingConfig
		if err := node.Decode(&decoded); err != nil {
			return err
		}
		*d = DampingConfig(decoded)
		d.RiseMSSet = riseSet
		d.FallMSSet = fallSet
		if !enabledSet {
			d.Enabled = true
		}
		return nil
	default:
		return fmt.Errorf("realism damping must be a boolean or mapping")
	}
}

type ThermalFadeConfig struct {
	RiseMS int `yaml:"rise_ms,omitempty"`
	FallMS int `yaml:"fall_ms,omitempty"`
}

func (t *ThermalFadeConfig) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("realism thermal_fade must be a mapping")
	}
	allowedKeys := map[string]bool{
		"rise_ms": true,
		"fall_ms": true,
	}
	for index := 0; index+1 < len(node.Content); index += 2 {
		key := node.Content[index].Value
		if !allowedKeys[key] {
			return fmt.Errorf("realism thermal_fade field %q is not supported", key)
		}
	}
	type rawThermalFadeConfig ThermalFadeConfig
	var decoded rawThermalFadeConfig
	if err := node.Decode(&decoded); err != nil {
		return err
	}
	*t = ThermalFadeConfig(decoded)
	return nil
}

type NeedleShadowConfig struct {
	Offset []int    `yaml:"offset,omitempty"`
	Alpha  *float64 `yaml:"alpha,omitempty"`
}

func (n *NeedleShadowConfig) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("realism needle_shadow must be a mapping")
	}
	allowedKeys := map[string]bool{
		"offset": true,
		"alpha":  true,
	}
	for index := 0; index+1 < len(node.Content); index += 2 {
		key := node.Content[index].Value
		if !allowedKeys[key] {
			return fmt.Errorf("realism needle_shadow field %q is not supported", key)
		}
	}
	type rawNeedleShadowConfig NeedleShadowConfig
	var decoded rawNeedleShadowConfig
	if err := node.Decode(&decoded); err != nil {
		return err
	}
	*n = NeedleShadowConfig(decoded)
	return nil
}

type Realism struct {
	Wraparound        *bool               `yaml:"wraparound,omitempty"`
	CarryDrag         *bool               `yaml:"carry_drag,omitempty"`
	SnapSettle        *bool               `yaml:"snap_settle,omitempty"`
	Hysteresis        *bool               `yaml:"hysteresis,omitempty"`
	Damping           *DampingConfig      `yaml:"damping,omitempty"`
	Stiction          *float64            `yaml:"stiction,omitempty"`
	Overshoot         *OvershootConfig    `yaml:"overshoot,omitempty"`
	PegBounce         *bool               `yaml:"peg_bounce,omitempty"`
	ThermalFade       *ThermalFadeConfig  `yaml:"thermal_fade,omitempty"`
	NeedleShadow      *NeedleShadowConfig `yaml:"needle_shadow,omitempty"`
	CalibrationOffset *float64            `yaml:"calibration_offset,omitempty"`
	MovementPolicy    string              `yaml:"movement_policy,omitempty"`
	DrumSlop          []int               `yaml:"drum_slop,omitempty"`
	DrumSlopSet       bool                `yaml:"-"`
}

func (r *Realism) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("realism must be a mapping")
	}
	allowedKeys := map[string]bool{
		"wraparound":         true,
		"carry_drag":         true,
		"snap_settle":        true,
		"hysteresis":         true,
		"damping":            true,
		"stiction":           true,
		"overshoot":          true,
		"peg_bounce":         true,
		"thermal_fade":       true,
		"needle_shadow":      true,
		"calibration_offset": true,
		"movement_policy":    true,
		"drum_slop":          true,
	}
	for index := 0; index+1 < len(node.Content); index += 2 {
		key := node.Content[index].Value
		if !allowedKeys[key] {
			return fmt.Errorf("realism field %q is not supported", key)
		}
		if key == "drum_slop" {
			r.DrumSlopSet = true
		}
	}
	type rawRealism Realism
	var decoded rawRealism
	if err := node.Decode(&decoded); err != nil {
		return err
	}
	drumSlopSet := r.DrumSlopSet
	*r = Realism(decoded)
	r.DrumSlopSet = drumSlopSet
	return nil
}

type Odometer struct {
	Movement string          `yaml:"movement,omitempty"`
	Wheels   []OdometerWheel `yaml:"wheels,omitempty"`
}

type BarConfig struct {
	Mode   string `yaml:"mode,omitempty"`
	Axis   string `yaml:"axis,omitempty"`
	Origin string `yaml:"origin,omitempty"`
	Bounds []int  `yaml:"bounds,omitempty"`
}

type Segmented struct {
	Hysteresis *float64         `yaml:"hysteresis,omitempty"`
	Images     []SegmentedImage `yaml:"-"`
}

type OdometerWheel struct {
	Strip    string `yaml:"strip"`
	Position []int  `yaml:"position"`
	Window   Size   `yaml:"window"`
	Offset   []int  `yaml:"offset,omitempty"`
	Role     string `yaml:"role,omitempty"`
}

type SegmentedImage struct {
	Threshold int
	Path      string
}

type SegmentedSelection struct {
	Threshold int
	Path      string
}

func LoadPackage(packageDir string) (Package, error) {
	return loadPackage(packageDir, false)
}

func LoadPackageForPreview(packageDir string) (Package, error) {
	return loadPackage(packageDir, true)
}

func loadPackage(packageDir string, allowAnyAssetRoot bool) (Package, error) {
	if packageDir == "" {
		return Package{}, fmt.Errorf("gauge package path must not be empty")
	}
	if isRemotePath(packageDir) {
		return Package{}, fmt.Errorf("gauge package path %q must be local", packageDir)
	}

	resolvedPackageDir, err := filepath.Abs(filepath.Clean(packageDir))
	if err != nil {
		return Package{}, fmt.Errorf("gauge package path %q could not be resolved: %w", packageDir, err)
	}
	if _, err := os.Stat(filepath.Join(resolvedPackageDir, "gauge.yaml")); os.IsNotExist(err) && !filepath.IsAbs(packageDir) {
		if fallback, ok := firstGaugePackageInSearchPaths(defaultGaugeSearchPaths(), packageDir); ok {
			resolvedPackageDir = fallback
		}
	} else if err != nil {
		return Package{}, fmt.Errorf("gauge package %q could not check gauge.yaml: %w", resolvedPackageDir, err)
	}
	assetRoot, err := findAssetRoot(resolvedPackageDir)
	if err != nil {
		if !allowAnyAssetRoot {
			return Package{}, err
		}
		assetRoot = resolvedPackageDir
	}

	yamlPath := filepath.Join(resolvedPackageDir, "gauge.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return Package{}, fmt.Errorf("gauge package %q could not load gauge.yaml: %w", resolvedPackageDir, err)
	}

	pkg, err := parsePackage(data)
	if err != nil {
		return Package{}, fmt.Errorf("gauge package %q could not parse gauge.yaml: %w", resolvedPackageDir, err)
	}
	normalizePackage(&pkg)
	pkg.Path = resolvedPackageDir
	pkg.YAMLPath = yamlPath
	pkg.AssetRoot = assetRoot

	if err := validatePackage(pkg); err != nil {
		return Package{}, fmt.Errorf("gauge package %q is invalid: %w", resolvedPackageDir, err)
	}
	if err := resolvePackagePaths(&pkg, filepath.Dir(yamlPath)); err != nil {
		return Package{}, fmt.Errorf("gauge package %q is invalid: %w", resolvedPackageDir, err)
	}
	if pkg.Type == TypeSegmented {
		images, err := discoverSegmentedImages(pkg.ID, pkg.Layers["segments"])
		if err != nil {
			return Package{}, fmt.Errorf("gauge package %q is invalid: %w", resolvedPackageDir, err)
		}
		pkg.Segmented.Images = images
	}

	return pkg, nil
}

func LoadPackageWithSearchPaths(searchPaths []string, packageDir string) (Package, error) {
	if packageDir == "" {
		return Package{}, fmt.Errorf("gauge package path must not be empty")
	}
	if isRemotePath(packageDir) {
		return Package{}, fmt.Errorf("gauge package path %q must be local", packageDir)
	}
	if filepath.IsAbs(packageDir) || len(searchPaths) == 0 {
		return LoadPackage(packageDir)
	}

	if candidate, ok := firstGaugePackageInSearchPaths(searchPaths, packageDir); ok {
		return LoadPackage(candidate)
	}

	tried := make([]string, 0, len(searchPaths))
	for _, root := range searchPaths {
		if strings.TrimSpace(root) != "" {
			tried = append(tried, filepath.Join(root, filepath.Clean(packageDir)))
		}
	}
	return Package{}, fmt.Errorf("gauge package path %q could not be found in asset search paths: %s", packageDir, strings.Join(tried, ", "))
}

func firstGaugePackageInSearchPaths(searchPaths []string, packageDir string) (string, bool) {
	cleanedPackageDir := filepath.Clean(packageDir)
	for _, root := range searchPaths {
		if strings.TrimSpace(root) == "" {
			continue
		}
		candidate := filepath.Join(root, cleanedPackageDir)
		if _, err := os.Stat(filepath.Join(candidate, "gauge.yaml")); err == nil {
			return candidate, true
		}
	}
	return "", false
}

func defaultGaugeSearchPaths() []string {
	paths := []string{}
	if pwd, err := os.Getwd(); err == nil {
		paths = append(paths, pwd)
	}
	if configPath := commandLineConfigPath(os.Args[1:]); configPath != "" {
		if configDir, err := filepath.Abs(filepath.Dir(configPath)); err == nil {
			paths = append(paths, configDir)
		}
	}
	if configPath := strings.TrimSpace(os.Getenv(dashboardConfigEnvVar)); configPath != "" {
		if configDir, err := filepath.Abs(filepath.Dir(configPath)); err == nil {
			paths = append(paths, configDir)
		}
	}
	return dedupePaths(paths)
}

func commandLineConfigPath(args []string) string {
	for index, arg := range args {
		if arg == "--config" || arg == "-config" {
			if index+1 < len(args) {
				return args[index+1]
			}
			return ""
		}
		if strings.HasPrefix(arg, "--config=") {
			return strings.TrimPrefix(arg, "--config=")
		}
		if strings.HasPrefix(arg, "-config=") {
			return strings.TrimPrefix(arg, "-config=")
		}
	}
	return ""
}

func dedupePaths(paths []string) []string {
	seen := map[string]bool{}
	cleaned := make([]string, 0, len(paths))
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		abs = filepath.Clean(abs)
		if seen[abs] {
			continue
		}
		seen[abs] = true
		cleaned = append(cleaned, abs)
	}
	return cleaned
}

func parsePackage(data []byte) (Package, error) {
	var pkg Package
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&pkg); err != nil {
		return Package{}, err
	}
	return pkg, nil
}

func validatePackage(pkg Package) error {
	if pkg.ID == "" {
		return fmt.Errorf("id must not be empty")
	}
	if pkg.Type == "" {
		return fmt.Errorf("type must not be empty")
	}
	switch pkg.Type {
	case TypeNumeric, TypeRadial, TypeOdometer, TypeIndicator, TypeBar, TypeSegmented:
	default:
		return fmt.Errorf("type %q is not supported", pkg.Type)
	}
	if pkg.Sensor == "" {
		return fmt.Errorf("sensor must not be empty")
	}
	if pkg.Size.Width <= 0 || pkg.Size.Height <= 0 {
		return fmt.Errorf("size width and height must be positive")
	}
	if err := validateRealism(pkg); err != nil {
		return err
	}
	if pkg.Type == TypeOdometer {
		if err := validateOdometer(pkg.Odometer); err != nil {
			return err
		}
	}
	if pkg.Type == TypeIndicator {
		if err := validateIndicatorLayers(pkg.Layers); err != nil {
			return err
		}
	}
	if pkg.Type == TypeBar {
		if err := validateBar(pkg.Bar, pkg.Layers); err != nil {
			return err
		}
		if err := validateBarValueMap(pkg.ValueMap); err != nil {
			return err
		}
	}
	if pkg.Type == TypeSegmented {
		if err := validateSegmented(pkg.Segmented, pkg.Layers); err != nil {
			return err
		}
	}
	return nil
}

func normalizePackage(pkg *Package) {
	if strings.TrimSpace(pkg.Realism.MovementPolicy) == "" {
		pkg.Realism.MovementPolicy = MovementPolicyImmediate
	}
	if pkg.Realism.NeedleShadow != nil && pkg.Realism.NeedleShadow.Alpha == nil {
		alpha := defaultNeedleShadowAlpha
		pkg.Realism.NeedleShadow.Alpha = &alpha
	}
	if pkg.Type == TypeOdometer && strings.TrimSpace(pkg.Odometer.Movement) == "" {
		pkg.Odometer.Movement = MovementInstant
	}
	if pkg.Type == TypeOdometer {
		switch pkg.Odometer.Movement {
		case MovementSmooth, MovementClick:
			log.Printf("gauge package %q odometer movement %q is recognised but not implemented; falling back to %q", pkg.ID, pkg.Odometer.Movement, MovementInstant)
			pkg.Odometer.Movement = MovementInstant
		}
	}
	if pkg.Type == TypeSegmented && pkg.Segmented.Hysteresis == nil {
		defaultHysteresis := 25.0
		pkg.Segmented.Hysteresis = &defaultHysteresis
	}
}

func validateOdometer(odometer Odometer) error {
	switch odometer.Movement {
	case MovementInstant, MovementLinear, MovementEaseOut, MovementBell, MovementSmooth, MovementClick:
	default:
		return fmt.Errorf("odometer movement %q is not supported", odometer.Movement)
	}
	if len(odometer.Wheels) == 0 {
		return fmt.Errorf("odometer wheels must not be empty")
	}
	for index, wheel := range odometer.Wheels {
		if strings.TrimSpace(wheel.Strip) == "" {
			return fmt.Errorf("odometer wheel %d strip must not be empty", index)
		}
		if len(wheel.Position) < 2 {
			return fmt.Errorf("odometer wheel %d position must contain x and y", index)
		}
		if wheel.Window.Width <= 0 || wheel.Window.Height <= 0 {
			return fmt.Errorf("odometer wheel %d window width and height must be positive", index)
		}
		if len(wheel.Offset) != 0 && len(wheel.Offset) < 2 {
			return fmt.Errorf("odometer wheel %d offset must contain x and y", index)
		}
		switch wheel.Role {
		case "", WheelRoleDigit, WheelRoleSubUnit:
		default:
			return fmt.Errorf("odometer wheel %d role %q is not supported", index, wheel.Role)
		}
	}
	return nil
}

func validateRealism(pkg Package) error {
	switch pkg.Realism.MovementPolicy {
	case MovementPolicyImmediate, MovementPolicyLinear, MovementPolicyEaseOut:
	default:
		return fmt.Errorf("realism movement_policy %q is not supported", pkg.Realism.MovementPolicy)
	}
	if pkg.Realism.Wraparound != nil && pkg.Type != TypeOdometer {
		return fmt.Errorf("realism wraparound is only supported for odometer gauges")
	}
	if pkg.Realism.CarryDrag != nil && pkg.Type != TypeOdometer {
		return fmt.Errorf("realism carry_drag is only supported for odometer gauges")
	}
	if pkg.Realism.SnapSettle != nil && pkg.Type != TypeOdometer {
		return fmt.Errorf("realism snap_settle is only supported for odometer gauges")
	}
	if pkg.Realism.Hysteresis != nil && pkg.Type != TypeRadial && pkg.Type != TypeBar {
		return fmt.Errorf("realism hysteresis is only supported for radial and bar gauges")
	}
	if pkg.Realism.Damping != nil {
		if pkg.Type != TypeRadial && pkg.Type != TypeBar {
			return fmt.Errorf("realism damping is only supported for radial and bar gauges")
		}
		if (pkg.Realism.Damping.RiseMSSet || pkg.Realism.Damping.FallMSSet) && pkg.Type != TypeBar {
			return fmt.Errorf("realism damping rise_ms and fall_ms are only supported for bar gauges")
		}
		if !pkg.Realism.Damping.Enabled && (pkg.Realism.Damping.RiseMSSet || pkg.Realism.Damping.FallMSSet) {
			return fmt.Errorf("realism damping timing requires damping to be enabled")
		}
		if pkg.Realism.Damping.RiseMSSet && pkg.Realism.Damping.RiseMS <= 0 {
			return fmt.Errorf("realism damping rise_ms must be greater than zero")
		}
		if pkg.Realism.Damping.FallMSSet && pkg.Realism.Damping.FallMS <= 0 {
			return fmt.Errorf("realism damping fall_ms must be greater than zero")
		}
	}
	if pkg.Realism.Stiction != nil {
		if pkg.Type != TypeRadial && pkg.Type != TypeBar {
			return fmt.Errorf("realism stiction is only supported for radial and bar gauges")
		}
		if math.IsNaN(*pkg.Realism.Stiction) || math.IsInf(*pkg.Realism.Stiction, 0) {
			return fmt.Errorf("realism stiction must be a finite threshold")
		}
		if *pkg.Realism.Stiction <= 0 {
			return fmt.Errorf("realism stiction must be greater than zero")
		}
		span := pkg.ValueMap.Max - pkg.ValueMap.Min
		if span <= 0 {
			return fmt.Errorf("realism stiction requires a valid value_map range")
		}
		if *pkg.Realism.Stiction > span {
			return fmt.Errorf("realism stiction %v exceeds value_map span %v", *pkg.Realism.Stiction, span)
		}
	}
	if pkg.Realism.Overshoot != nil {
		if pkg.Type != TypeRadial && pkg.Type != TypeBar {
			return fmt.Errorf("realism overshoot is only supported for radial and bar gauges")
		}
		if pkg.ValueMap.Max <= pkg.ValueMap.Min {
			return fmt.Errorf("realism overshoot requires a valid value_map range")
		}
		if pkg.Realism.Overshoot.Ratio != nil {
			ratio := *pkg.Realism.Overshoot.Ratio
			if math.IsNaN(ratio) || math.IsInf(ratio, 0) {
				return fmt.Errorf("realism overshoot ratio must be finite")
			}
			if ratio <= 0 {
				return fmt.Errorf("realism overshoot ratio must be greater than zero")
			}
			if ratio > maxOvershootRatio {
				return fmt.Errorf("realism overshoot ratio %v exceeds maximum %v", ratio, maxOvershootRatio)
			}
		}
		if pkg.Realism.Overshoot.MinChangeRatio != nil {
			minChangeRatio := *pkg.Realism.Overshoot.MinChangeRatio
			if math.IsNaN(minChangeRatio) || math.IsInf(minChangeRatio, 0) {
				return fmt.Errorf("realism overshoot min_change_ratio must be finite")
			}
			if minChangeRatio < 0 {
				return fmt.Errorf("realism overshoot min_change_ratio must be greater than or equal to zero")
			}
		}
		if pkg.Realism.Overshoot.MaxSpanRatio != nil {
			maxSpanRatio := *pkg.Realism.Overshoot.MaxSpanRatio
			if math.IsNaN(maxSpanRatio) || math.IsInf(maxSpanRatio, 0) {
				return fmt.Errorf("realism overshoot max_span_ratio must be finite")
			}
			if maxSpanRatio <= 0 {
				return fmt.Errorf("realism overshoot max_span_ratio must be greater than zero")
			}
		}
		switch pkg.Realism.Overshoot.SettleMode {
		case "", OvershootSettleSmooth, OvershootSettleOscillate:
		default:
			return fmt.Errorf("realism overshoot settle_mode %q is not supported", pkg.Realism.Overshoot.SettleMode)
		}
		if pkg.Type == TypeBar && pkg.Realism.Overshoot.SettleMode == OvershootSettleOscillate {
			return fmt.Errorf("realism overshoot settle_mode %q is not supported for bar gauges", pkg.Realism.Overshoot.SettleMode)
		}
		if pkg.Realism.Overshoot.SettleCycles != nil {
			if pkg.Type == TypeBar {
				return fmt.Errorf("realism overshoot settle_cycles is not supported for bar gauges")
			}
			settleCycles := *pkg.Realism.Overshoot.SettleCycles
			if math.IsNaN(settleCycles) || math.IsInf(settleCycles, 0) {
				return fmt.Errorf("realism overshoot settle_cycles must be finite")
			}
			if settleCycles <= 0 {
				return fmt.Errorf("realism overshoot settle_cycles must be greater than zero")
			}
		}
		if pkg.Realism.Overshoot.SettleDamping != nil {
			if pkg.Type == TypeBar {
				return fmt.Errorf("realism overshoot settle_damping is not supported for bar gauges")
			}
			settleDamping := *pkg.Realism.Overshoot.SettleDamping
			if math.IsNaN(settleDamping) || math.IsInf(settleDamping, 0) {
				return fmt.Errorf("realism overshoot settle_damping must be finite")
			}
			if settleDamping <= 0 {
				return fmt.Errorf("realism overshoot settle_damping must be greater than zero")
			}
		}
	}
	if pkg.Realism.PegBounce != nil {
		if pkg.Type != TypeRadial {
			return fmt.Errorf("realism peg_bounce is only supported for radial gauges")
		}
		if *pkg.Realism.PegBounce && (!pkg.ValueMap.Clamp || pkg.ValueMap.Max <= pkg.ValueMap.Min) {
			return fmt.Errorf("realism peg_bounce requires a clamped radial value_map range")
		}
	}
	if pkg.Realism.NeedleShadow != nil {
		if pkg.Type != TypeRadial {
			return fmt.Errorf("realism needle_shadow is only supported for radial gauges")
		}
		if len(pkg.Realism.NeedleShadow.Offset) != 2 {
			return fmt.Errorf("realism needle_shadow offset must contain x and y")
		}
		if pkg.Realism.NeedleShadow.Alpha != nil {
			alpha := *pkg.Realism.NeedleShadow.Alpha
			if math.IsNaN(alpha) || math.IsInf(alpha, 0) {
				return fmt.Errorf("realism needle_shadow alpha must be finite")
			}
			if alpha < 0 || alpha > 1 {
				return fmt.Errorf("realism needle_shadow alpha must be between 0 and 1")
			}
		}
	}
	if pkg.Realism.ThermalFade != nil {
		if pkg.Type != TypeIndicator {
			return fmt.Errorf("realism thermal_fade is only supported for indicator gauges")
		}
		if pkg.Realism.ThermalFade.RiseMS <= 0 {
			return fmt.Errorf("realism thermal_fade rise_ms must be greater than zero")
		}
		if pkg.Realism.ThermalFade.FallMS <= 0 {
			return fmt.Errorf("realism thermal_fade fall_ms must be greater than zero")
		}
	}
	if pkg.Realism.CalibrationOffset != nil {
		if pkg.Type != TypeRadial {
			return fmt.Errorf("realism calibration_offset is only supported for radial gauges")
		}
		offset := *pkg.Realism.CalibrationOffset
		if math.IsNaN(offset) || math.IsInf(offset, 0) {
			return fmt.Errorf("realism calibration_offset must be finite")
		}
	}
	if pkg.Realism.DrumSlopSet {
		if pkg.Type != TypeOdometer {
			return fmt.Errorf("realism drum_slop is only supported for odometer gauges")
		}
		if len(pkg.Realism.DrumSlop) != len(pkg.Odometer.Wheels) {
			return fmt.Errorf("realism drum_slop must define exactly one offset per odometer wheel")
		}
		for index, slop := range pkg.Realism.DrumSlop {
			maxOffset := pkg.Odometer.Wheels[index].Window.Height / 4
			if maxOffset < 1 {
				maxOffset = 1
			}
			if slop < -maxOffset || slop > maxOffset {
				return fmt.Errorf("realism drum_slop wheel %d offset %d exceeds +/- %d", index, slop, maxOffset)
			}
		}
	}
	return nil
}

func validateIndicatorLayers(layers map[string]string) error {
	if strings.TrimSpace(layers["on"]) == "" {
		return fmt.Errorf("indicator layer on must not be empty")
	}
	return nil
}

func validateBar(bar BarConfig, layers map[string]string) error {
	switch bar.Mode {
	case "level":
	default:
		return fmt.Errorf("bar mode %q is not supported", bar.Mode)
	}
	switch bar.Axis {
	case "vertical":
	default:
		return fmt.Errorf("bar axis %q is not supported", bar.Axis)
	}
	switch bar.Origin {
	case "bottom":
	default:
		return fmt.Errorf("bar origin %q is not supported", bar.Origin)
	}
	if len(bar.Bounds) != 4 {
		return fmt.Errorf("bar bounds must contain x, y, width, and height")
	}
	if bar.Bounds[0] < 0 || bar.Bounds[1] < 0 || bar.Bounds[2] <= 0 || bar.Bounds[3] <= 0 {
		return fmt.Errorf("bar bounds x and y must be non-negative and width and height must be positive")
	}
	if strings.TrimSpace(layers["level"]) == "" {
		return fmt.Errorf("bar layer level must not be empty")
	}
	return nil
}

func validateBarValueMap(valueMap ValueMap) error {
	if valueMap.Max <= valueMap.Min {
		return fmt.Errorf("bar value_map max must be greater than min")
	}
	return nil
}

func validateSegmented(segmented Segmented, layers map[string]string) error {
	if strings.TrimSpace(layers["segments"]) == "" {
		return fmt.Errorf("segmented layer segments must not be empty")
	}
	if segmented.Hysteresis != nil {
		if *segmented.Hysteresis < 0 || *segmented.Hysteresis > 100 {
			return fmt.Errorf("segmented hysteresis must be within 0..100")
		}
	}
	if _, err := parseSegmentedPattern(strings.TrimSpace(layers["segments"])); err != nil {
		return err
	}
	return nil
}

func resolvePackagePaths(pkg *Package, yamlDir string) error {
	var err error
	pkg.Layers, err = resolvePathMap(pkg.AssetRoot, yamlDir, pkg.Layers, "layer")
	if err != nil {
		return err
	}

	pkg.DigitSet.Background, err = resolveOptionalPath(pkg.AssetRoot, yamlDir, pkg.DigitSet.Background, "digit_set background")
	if err != nil {
		return err
	}
	pkg.DigitSet.DecimalPoint, err = resolveOptionalPath(pkg.AssetRoot, yamlDir, pkg.DigitSet.DecimalPoint, "digit_set decimal_point")
	if err != nil {
		return err
	}
	pkg.DigitSet.Foreground, err = resolveOptionalPath(pkg.AssetRoot, yamlDir, pkg.DigitSet.Foreground, "digit_set foreground")
	if err != nil {
		return err
	}
	pkg.DigitSet.Characters, err = resolvePathMap(pkg.AssetRoot, yamlDir, pkg.DigitSet.Characters, "digit_set character")
	if err != nil {
		return err
	}
	for index := range pkg.Odometer.Wheels {
		pkg.Odometer.Wheels[index].Strip, err = resolveOptionalPath(pkg.AssetRoot, yamlDir, pkg.Odometer.Wheels[index].Strip, fmt.Sprintf("odometer wheel %d strip", index))
		if err != nil {
			return err
		}
	}
	return nil
}

func resolvePathMap(assetRoot string, yamlDir string, paths map[string]string, label string) (map[string]string, error) {
	if len(paths) == 0 {
		return paths, nil
	}
	resolved := make(map[string]string, len(paths))
	for key, path := range paths {
		if key == "" {
			return nil, fmt.Errorf("%s key must not be empty", label)
		}
		resolvedPath, err := resolveRequiredPath(assetRoot, yamlDir, path, fmt.Sprintf("%s %q", label, key))
		if err != nil {
			return nil, err
		}
		resolved[key] = resolvedPath
	}
	return resolved, nil
}

func resolveOptionalPath(assetRoot string, yamlDir string, path string, label string) (string, error) {
	if path == "" {
		return "", nil
	}
	return resolveRequiredPath(assetRoot, yamlDir, path, label)
}

func resolveRequiredPath(assetRoot string, yamlDir string, path string, label string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("%s path must not be empty", label)
	}
	if isRemotePath(path) {
		return "", fmt.Errorf("%s path %q must be local", label, path)
	}
	resolved := path
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(yamlDir, path)
	}
	resolved = filepath.Clean(resolved)
	if !isInside(assetRoot, resolved) {
		return "", fmt.Errorf("%s path %q escapes asset tree %q", label, path, assetRoot)
	}
	return resolved, nil
}

func findAssetRoot(packageDir string) (string, error) {
	current := packageDir
	for {
		if filepath.Base(current) == "assets" {
			return current, nil
			/*
				** OLD requirement that "gauges" need to be under "assets". Stoopid.
				rel, err := filepath.Rel(current, packageDir)
				if err == nil {
					rel = filepath.ToSlash(rel)
					if strings.HasPrefix(rel, "gauges/") {
						return current, nil
					}
				}
			*/
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("gauge package path %q must be under assets directory", packageDir)
}

func isInside(root string, path string) bool {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && !filepath.IsAbs(rel))
}

func isRemotePath(path string) bool {
	return strings.Contains(path, ":"+"//")
}

type segmentedPattern struct {
	prefix string
	suffix string
	width  int
}

var segmentedPatternRe = regexp.MustCompile(`\{percent(?::(\d+))?\}`)

func parseSegmentedPattern(path string) (segmentedPattern, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return segmentedPattern{}, fmt.Errorf("segmented layer segments must not be empty")
	}
	loc := segmentedPatternRe.FindStringSubmatchIndex(trimmed)
	if loc == nil {
		return segmentedPattern{}, fmt.Errorf("segmented layer segments path %q must include {percent}", path)
	}
	match := segmentedPatternRe.FindStringSubmatch(trimmed)
	width := 0
	if len(match) > 1 && match[1] != "" {
		parsedWidth, err := strconv.Atoi(match[1])
		if err != nil || parsedWidth <= 0 {
			return segmentedPattern{}, fmt.Errorf("segmented layer segments path %q has invalid percent width %q", path, match[1])
		}
		width = parsedWidth
	}
	return segmentedPattern{
		prefix: trimmed[:loc[0]],
		suffix: trimmed[loc[1]:],
		width:  width,
	}, nil
}

func discoverSegmentedImages(packageID string, segmentsPath string) ([]SegmentedImage, error) {
	pattern, err := parseSegmentedPattern(filepath.Base(segmentsPath))
	if err != nil {
		return nil, err
	}

	segmentDir := filepath.Dir(segmentsPath)
	entries, err := os.ReadDir(segmentDir)
	if err != nil {
		return nil, fmt.Errorf("segmented layer %q could not scan directory: %w", segmentsPath, err)
	}

	regexPattern := "^" + regexp.QuoteMeta(pattern.prefix) + `(\d+)` + regexp.QuoteMeta(pattern.suffix) + "$"
	if pattern.width > 0 {
		regexPattern = "^" + regexp.QuoteMeta(pattern.prefix) + fmt.Sprintf(`(\d{%d})`, pattern.width) + regexp.QuoteMeta(pattern.suffix) + "$"
	}
	matchPattern := regexp.MustCompile(regexPattern)

	seen := map[int]bool{}
	images := make([]SegmentedImage, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := matchPattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		threshold, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("segmented layer %q file %q has invalid percent threshold: %w", segmentsPath, entry.Name(), err)
		}
		if threshold > 100 {
			log.Printf("gauge package %q segmented threshold %d in %q ignored: above 100", packageID, threshold, entry.Name())
			continue
		}
		if seen[threshold] {
			return nil, fmt.Errorf("segmented layer %q has duplicate percent threshold %d", segmentsPath, threshold)
		}
		seen[threshold] = true
		images = append(images, SegmentedImage{
			Threshold: threshold,
			Path:      filepath.Join(segmentDir, entry.Name()),
		})
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("segmented layer %q did not find any matching percent images", segmentsPath)
	}

	sort.Slice(images, func(i, j int) bool {
		if images[i].Threshold == images[j].Threshold {
			return images[i].Path < images[j].Path
		}
		return images[i].Threshold < images[j].Threshold
	})
	return images, nil
}

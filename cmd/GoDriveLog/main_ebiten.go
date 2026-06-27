//go:build !fyne_legacy

package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/gauges"
	v3harness "github.com/MickMake/GoDriveLog/internal/dashboard/harness"
)

const dashboardConfigEnvVar = "GODRIVELOG_CONFIG_PATH"

const (
	frameworkSmokeTheme = "framework-smoke"
	ornateTimberTheme   = "ornate-timber"
	neonGridTheme       = "neon-grid"
	steamScrapTheme     = "steam-scrap"
)

var (
	dashboardRunCommand     = runV3EbitenCommand
	dashboardHarnessCommand = runV3EbitenHarnessCommand
	dashboardPreviewCommand = runV3EbitenPreviewCommand
)

type dashboardPreviewOptions struct {
	ConfigPath   string
	GaugeID      string
	InitialValue *float64
	Step         *float64
	FineStep     *float64
	CoarseStep   *float64
}

func main() {
	if err := runCLI(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		log.Fatal(err)
	}
}

func runCLI(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		writeRootHelp(stdout)
		return nil
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeRootHelp(stdout)
		return nil
	case "dashboard":
		return runDashboardCLI(args[1:], stdout, stderr)
	default:
		writeRootHelp(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runDashboardCLI(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return runBareDashboardFlags(nil, stdout, stderr)
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeDashboardHelp(stdout)
		return nil
	case "run":
		return runDashboardRunCommand(args[1:], stdout, stderr)
	case "harness":
		return runDashboardHarnessSubcommand(args[1:], stdout, stderr)
	case "preview":
		return runDashboardPreviewCommand(args[1:], stdout, stderr)
	case "examples":
		return runDashboardExamplesCommand(args[1:], stdout, stderr)
	case "validate":
		return runDashboardValidateCommand(args[1:], stdout, stderr)
	default:
		if strings.HasPrefix(args[0], "-") {
			return runBareDashboardFlags(args, stdout, stderr)
		}
		writeDashboardHelp(stderr)
		return fmt.Errorf("unknown dashboard command %q", args[0])
	}
}

func runBareDashboardFlags(args []string, stdout, stderr io.Writer) error {
	fs := newDashboardFlagSet("dashboard")
	configPath := fs.String("config", "", "path to a dashboard config; when omitted, use deterministic discovery")
	fs.Usage = func() {
		writeDashboardHelp(stdout)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			writeDashboardHelp(stdout)
			return nil
		}
		return err
	}
	if len(fs.Args()) != 0 {
		writeDashboardHelp(stderr)
		return fmt.Errorf("unexpected arguments for dashboard: %s", strings.Join(fs.Args(), " "))
	}

	overview, err := buildDashboardOverview(*configPath)
	if err != nil {
		return err
	}
	_, _ = io.WriteString(stdout, overview)
	return nil
}

func buildDashboardOverview(configPath string) (string, error) {
	resolvedConfigPath, cfg, err := loadDashboardOverviewConfig(configPath)
	if err != nil {
		return "", err
	}

	absoluteConfigPath, err := filepath.Abs(filepath.Clean(resolvedConfigPath))
	if err != nil {
		return "", fmt.Errorf("resolve dashboard config path %q: %w", resolvedConfigPath, err)
	}

	var b strings.Builder
	fmt.Fprintln(&b, "GoDriveLog dashboard overview")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Resolved config: %s\n", absoluteConfigPath)
	fmt.Fprintln(&b, "Vehicles:")

	for _, vehicleID := range sortedVehicleIDs(cfg.Vehicles) {
		searchPaths, err := v3assets.DefaultSearchPaths(absoluteConfigPath, vehicleID)
		if err != nil {
			return "", err
		}
		vehicle := cfg.Vehicles[vehicleID]
		fmt.Fprintf(&b, "- %s", vehicleID)
		if name := strings.TrimSpace(vehicle.Name); name != "" {
			fmt.Fprintf(&b, " (%s)", name)
		}
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "  obd source: %s\n", vehicle.OBD.Address)

		plan, err := v3config.Resolve(cfg, vehicleID)
		if err != nil {
			fmt.Fprintf(&b, "  warning: %s\n", compactOverviewWarning(err))
			continue
		}
		if len(plan.Dashboards) == 0 {
			fmt.Fprintln(&b, "  dashboards: none")
			continue
		}

		fmt.Fprintln(&b, "  dashboards:")
		for _, dashboard := range plan.Dashboards {
			fmt.Fprintf(&b, "    - %s\n", dashboard.ID)
			if len(dashboard.Config.Widgets) == 0 {
				fmt.Fprintln(&b, "      widgets: none")
				continue
			}

			fmt.Fprintln(&b, "      widgets:")
			for _, widget := range dashboard.Config.Widgets {
				entry := describeOverviewWidget(cfg, searchPaths, widget)
				fmt.Fprintf(&b, "        - %s: type=%s source=%s", entry.Label, entry.Type, entry.Source)
				if entry.PID != "" {
					fmt.Fprintf(&b, " pid=%s", entry.PID)
				}
				fmt.Fprintln(&b)
				if entry.Warning != "" {
					fmt.Fprintf(&b, "          warning: %s\n", entry.Warning)
				}
			}
		}
	}

	return b.String(), nil
}

func loadDashboardOverviewConfig(configPath string) (string, v3config.Config, error) {
	selectedConfigPath := strings.TrimSpace(configPath)
	if selectedConfigPath == "" {
		selection, err := resolveDashboardSelection("", "")
		if err != nil {
			return "", v3config.Config{}, err
		}
		selectedConfigPath = selection.ConfigPath
	}

	cfg, err := v3config.LoadFile(selectedConfigPath)
	if err != nil {
		return "", v3config.Config{}, fmt.Errorf("load dashboard config %q: %w", selectedConfigPath, err)
	}
	return selectedConfigPath, cfg, nil
}

type overviewWidgetEntry struct {
	Label   string
	Type    string
	Source  string
	PID     string
	Warning string
}

func describeOverviewWidget(cfg v3config.Config, searchPaths []string, widget v3config.WidgetConfig) overviewWidgetEntry {
	entry := overviewWidgetEntry{
		Label:  overviewWidgetLabel(widget),
		Type:   widget.Type,
		Source: "none",
	}

	sensorID := strings.TrimSpace(widget.Sensor)
	if widget.Type == v3config.WidgetTypeGauge {
		pkg, err := gauges.LoadPackageWithSearchPaths(searchPaths, widget.Gauge)
		if err != nil {
			entry.Type = "unknown"
			entry.Source = "unknown"
			entry.Warning = compactOverviewWarning(err)
			return entry
		}
		entry.Type = pkg.Type
		sensorID = pkg.Sensor
	}

	if sensorID == "" {
		return entry
	}

	entry.Source = sensorID
	sensor, ok := cfg.Sensors[sensorID]
	if !ok {
		entry.Warning = fmt.Sprintf("sensor %q is not defined in top-level sensors", sensorID)
		return entry
	}
	if sensor.Type == v3config.SensorTypeOBD {
		entry.PID = sensor.PID
	}
	return entry
}

func overviewWidgetLabel(widget v3config.WidgetConfig) string {
	if id := strings.TrimSpace(widget.ID); id != "" {
		return id
	}
	if widget.Type == v3config.WidgetTypeGauge {
		if gaugePath := strings.TrimSpace(widget.Gauge); gaugePath != "" {
			return gaugePath
		}
	}
	if asset := strings.TrimSpace(widget.Asset); asset != "" {
		return asset
	}
	if sensor := strings.TrimSpace(widget.Sensor); sensor != "" {
		return sensor
	}
	return widget.Type
}

func compactOverviewWarning(err error) string {
	if err == nil {
		return ""
	}
	return strings.Join(strings.Fields(err.Error()), " ")
}

func runDashboardRunCommand(args []string, stdout, stderr io.Writer) error {
	filteredArgs, positionalArgs, err := splitArgsAllowingSinglePositional(args, map[string]bool{
		"--config":   true,
		"--renderer": true,
		"--duration": true,
	})
	if err != nil {
		return err
	}

	fs := newDashboardFlagSet("dashboard run")
	configPath := fs.String("config", "", "path to a dashboard config; when omitted, use deterministic discovery")
	renderer := fs.String("renderer", v3RendererEbiten, "dashboard renderer backend")
	duration := fs.Duration("duration", 0, "optional runtime duration such as 60s; zero runs until interrupted")
	fs.Usage = func() {
		writeDashboardRunHelp(stdout)
	}

	if err := fs.Parse(filteredArgs); err != nil {
		if err == flag.ErrHelp {
			writeDashboardRunHelp(stdout)
			return nil
		}
		return err
	}
	if len(fs.Args()) != 0 {
		return fmt.Errorf("dashboard run accepts at most one vehicle id")
	}
	if len(positionalArgs) > 1 {
		return fmt.Errorf("dashboard run accepts at most one vehicle id")
	}

	normalizedRenderer, err := normalizeV3Renderer(*renderer)
	if err != nil {
		return err
	}
	if normalizedRenderer != v3RendererEbiten {
		return fmt.Errorf("unsupported dashboard renderer %q", normalizedRenderer)
	}

	vehicleID := ""
	if len(positionalArgs) == 1 {
		vehicleID = positionalArgs[0]
	}
	selection, err := resolveDashboardSelection(*configPath, vehicleID)
	if err != nil {
		return err
	}

	return withDashboardConfigPathEnv(selection.ConfigPath, func() error {
		return dashboardRunCommand(selection.ConfigPath, selection.VehicleID, *duration)
	})
}

func runDashboardHarnessSubcommand(args []string, stdout, stderr io.Writer) error {
	filteredArgs, positionalArgs, err := splitArgsAllowingSinglePositional(args, map[string]bool{
		"--config":   true,
		"--pattern":  true,
		"--interval": true,
		"--duration": true,
		"--renderer": true,
	})
	if err != nil {
		return err
	}

	fs := newDashboardFlagSet("dashboard harness")
	configPath := fs.String("config", "", "path to a dashboard config; when omitted, use deterministic discovery")
	pattern := fs.String("pattern", v3harness.PatternSweep, "harness pattern: sweep, heartbeat, or fixed")
	interval := fs.Duration("interval", 100*time.Millisecond, "harness update interval such as 50ms or 100ms")
	duration := fs.Duration("duration", 0, "optional harness duration such as 60s; zero runs until interrupted")
	renderer := fs.String("renderer", v3RendererEbiten, "dashboard renderer backend")
	fs.Usage = func() {
		writeDashboardHarnessHelp(stdout)
	}

	if err := fs.Parse(filteredArgs); err != nil {
		if err == flag.ErrHelp {
			writeDashboardHarnessHelp(stdout)
			return nil
		}
		return err
	}
	if len(fs.Args()) != 0 {
		return fmt.Errorf("dashboard harness accepts at most one vehicle id")
	}
	if len(positionalArgs) > 1 {
		return fmt.Errorf("dashboard harness accepts at most one vehicle id")
	}

	normalizedRenderer, err := normalizeV3Renderer(*renderer)
	if err != nil {
		return err
	}
	if normalizedRenderer != v3RendererEbiten {
		return fmt.Errorf("unsupported dashboard renderer %q", normalizedRenderer)
	}

	vehicleID := ""
	if len(positionalArgs) == 1 {
		vehicleID = positionalArgs[0]
	}
	selection, err := resolveDashboardSelection(*configPath, vehicleID)
	if err != nil {
		return err
	}

	return withDashboardConfigPathEnv(selection.ConfigPath, func() error {
		return dashboardHarnessCommand(selection.ConfigPath, selection.VehicleID, *pattern, *interval, *duration)
	})
}

func runDashboardExamplesCommand(args []string, stdout, stderr io.Writer) error {
	fs := newDashboardFlagSet("dashboard examples")
	configPath := fs.String("config", "", "optional source dashboard config; when omitted, use the built-in themed example")
	vehicleID := fs.String("vehicle", "", "optional vehicle id to verify in the source config")
	theme := fs.String("theme", frameworkSmokeTheme, "built-in example theme to export when --config is not supplied")
	outputDir := fs.String("output", "", "output directory for the self-contained dashboard example")
	force := fs.Bool("force", false, "overwrite an existing non-empty output directory")
	fs.Usage = func() {
		writeDashboardExamplesHelp(stdout)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			writeDashboardExamplesHelp(stdout)
			return nil
		}
		return err
	}
	if len(fs.Args()) != 0 {
		return fmt.Errorf("dashboard examples does not accept positional arguments")
	}
	if strings.TrimSpace(*outputDir) == "" {
		return fmt.Errorf("dashboard examples requires --output")
	}

	if err := exportDashboardExample(*configPath, *vehicleID, *theme, *outputDir, *force); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "exported dashboard example to %s\n", filepath.Clean(*outputDir))
	return nil
}

func runDashboardPreviewCommand(args []string, stdout, stderr io.Writer) error {
	filteredArgs, positionalArgs, err := splitArgsAllowingSinglePositional(args, map[string]bool{
		"--gauge":       true,
		"--value":       true,
		"--step":        true,
		"--fine-step":   true,
		"--coarse-step": true,
	})
	if err != nil {
		return err
	}

	fs := newDashboardFlagSet("dashboard preview")
	gaugeID := fs.String("gauge", "", "optional gauge widget id; use dashboard/widget when the file contains duplicates")
	initialValue := fs.Float64("value", 0, "optional starting value override; defaults to the midpoint of the inferred range")
	step := fs.Float64("step", 0, "optional step size for Up/Down")
	fineStep := fs.Float64("fine-step", 0, "optional fine step size for Ctrl/Cmd+Up/Down")
	coarseStep := fs.Float64("coarse-step", 0, "optional coarse step size for Shift+Up/Down")
	fs.Usage = func() {
		writeDashboardPreviewHelp(stdout)
	}

	if err := fs.Parse(filteredArgs); err != nil {
		if err == flag.ErrHelp {
			writeDashboardPreviewHelp(stdout)
			return nil
		}
		return err
	}
	if len(fs.Args()) != 0 {
		return fmt.Errorf("dashboard preview accepts exactly one positional preview file")
	}
	if len(positionalArgs) != 1 {
		writeDashboardPreviewHelp(stderr)
		return fmt.Errorf("dashboard preview requires exactly one preview file")
	}

	options := dashboardPreviewOptions{
		ConfigPath: strings.TrimSpace(positionalArgs[0]),
		GaugeID:    strings.TrimSpace(*gaugeID),
	}
	if options.ConfigPath == "" {
		return fmt.Errorf("dashboard preview requires exactly one preview file")
	}

	if value, ok, err := optionalPositiveOrNegativeFloat64Flag(fs, "value", *initialValue); err != nil {
		return err
	} else if ok {
		options.InitialValue = value
	}
	if value, ok, err := optionalPositiveFloat64Flag(fs, "step", *step); err != nil {
		return err
	} else if ok {
		options.Step = value
	}
	if value, ok, err := optionalPositiveFloat64Flag(fs, "fine-step", *fineStep); err != nil {
		return err
	} else if ok {
		options.FineStep = value
	}
	if value, ok, err := optionalPositiveFloat64Flag(fs, "coarse-step", *coarseStep); err != nil {
		return err
	} else if ok {
		options.CoarseStep = value
	}

	return withDashboardConfigPathEnv(options.ConfigPath, func() error {
		return dashboardPreviewCommand(options)
	})
}

func runDashboardValidateCommand(args []string, stdout, stderr io.Writer) error {
	fs := newDashboardFlagSet("dashboard validate")
	configPath := fs.String("config", "", "path to a dashboard config; when omitted, use deterministic discovery")
	fs.Usage = func() {
		writeDashboardValidateHelp(stdout)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			writeDashboardValidateHelp(stdout)
			return nil
		}
		return err
	}
	if len(fs.Args()) > 1 {
		return fmt.Errorf("dashboard validate accepts at most one positional config file")
	}
	if len(fs.Args()) == 1 && strings.TrimSpace(*configPath) != "" {
		return fmt.Errorf("dashboard validate accepts either a positional config file or --config, not both")
	}

	selectedConfig := strings.TrimSpace(*configPath)
	if len(fs.Args()) == 1 {
		selectedConfig = fs.Args()[0]
	}

	if selectedConfig == "" {
		discoveredConfig, err := discoverDashboardConfigForValidation()
		if err != nil {
			return err
		}
		selectedConfig = discoveredConfig
	}

	if _, err := v3config.LoadFile(selectedConfig); err != nil {
		return fmt.Errorf("validate dashboard config %q: %w", selectedConfig, err)
	}
	fmt.Fprintf(stdout, "validated dashboard config %s\n", filepath.Clean(selectedConfig))
	return nil
}

type dashboardSelection struct {
	ConfigPath string
	VehicleID  string
}

type searchedConfigRecord struct {
	Path     string
	Vehicles []string
}

func resolveDashboardSelection(configPath, vehicleID string) (dashboardSelection, error) {
	if strings.TrimSpace(configPath) != "" {
		return resolveExplicitDashboardSelection(configPath, vehicleID)
	}
	return discoverDashboardSelection(vehicleID)
}

func resolveExplicitDashboardSelection(configPath, vehicleID string) (dashboardSelection, error) {
	cfg, err := v3config.LoadFile(configPath)
	if err != nil {
		return dashboardSelection{}, fmt.Errorf("load dashboard config %q: %w", configPath, err)
	}
	vehicleIDs := sortedVehicleIDs(cfg.Vehicles)
	if vehicleID == "" {
		if len(vehicleIDs) == 1 {
			return dashboardSelection{ConfigPath: configPath, VehicleID: vehicleIDs[0]}, nil
		}
		return dashboardSelection{}, fmt.Errorf("dashboard config %q defines multiple vehicles (%s); supply a vehicle id", configPath, strings.Join(vehicleIDs, ", "))
	}
	if _, ok := cfg.Vehicles[vehicleID]; !ok {
		return dashboardSelection{}, fmt.Errorf("dashboard config %q does not define vehicle %q; available vehicles: %s", configPath, vehicleID, strings.Join(vehicleIDs, ", "))
	}
	return dashboardSelection{ConfigPath: configPath, VehicleID: vehicleID}, nil
}

func discoverDashboardSelection(vehicleID string) (dashboardSelection, error) {
	paths, err := discoverDashboardConfigPaths()
	if err != nil {
		return dashboardSelection{}, err
	}
	if len(paths) == 0 {
		return dashboardSelection{}, fmt.Errorf("no dashboard config files found in the current directory or /etc/godrivelog")
	}

	records := make([]searchedConfigRecord, 0, len(paths))
	for _, path := range paths {
		cfg, err := v3config.LoadFile(path)
		if err != nil {
			continue
		}
		vehicleIDs := sortedVehicleIDs(cfg.Vehicles)
		record := searchedConfigRecord{Path: path, Vehicles: vehicleIDs}
		records = append(records, record)

		if vehicleID == "" {
			if len(vehicleIDs) == 1 {
				return dashboardSelection{ConfigPath: path, VehicleID: vehicleIDs[0]}, nil
			}
			return dashboardSelection{}, fmt.Errorf("dashboard config %q defines multiple vehicles (%s); supply a vehicle id", path, strings.Join(vehicleIDs, ", "))
		}

		if containsString(vehicleIDs, vehicleID) {
			return dashboardSelection{ConfigPath: path, VehicleID: vehicleID}, nil
		}
	}

	if vehicleID != "" {
		if len(records) == 0 {
			return dashboardSelection{}, fmt.Errorf("vehicle %q was not found because no valid dashboard config files were discovered", vehicleID)
		}
		return dashboardSelection{}, fmt.Errorf("vehicle %q was not found in searched config files: %s", vehicleID, formatSearchedConfigRecords(records))
	}
	return dashboardSelection{}, fmt.Errorf("no valid single-vehicle dashboard config files were discovered in the current directory or /etc/godrivelog")
}

func discoverDashboardConfigForValidation() (string, error) {
	paths, err := discoverDashboardConfigPaths()
	if err != nil {
		return "", err
	}
	if len(paths) == 0 {
		return "", fmt.Errorf("no dashboard config files found in the current directory or /etc/godrivelog")
	}

	for _, path := range paths {
		if _, err := v3config.LoadFile(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no valid dashboard config files were discovered in the current directory or /etc/godrivelog")
}

func discoverDashboardConfigPaths() ([]string, error) {
	paths := []string{}
	if cwd, err := os.Getwd(); err == nil {
		if err := appendDashboardConfigPaths(cwd, false, &paths); err != nil {
			return nil, err
		}
	}
	if err := appendDashboardConfigPaths("/etc/godrivelog", true, &paths); err != nil {
		return nil, err
	}
	return paths, nil
}

func appendDashboardConfigPaths(root string, recursive bool, paths *[]string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read dashboard config directory %q: %w", root, err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		entryPath := filepath.Join(root, entry.Name())
		if entry.IsDir() {
			if recursive {
				if err := appendDashboardConfigPaths(entryPath, true, paths); err != nil {
					return err
				}
			}
			continue
		}
		if isDashboardConfigCandidate(entry.Name()) {
			*paths = append(*paths, entryPath)
		}
	}
	return nil
}

func isDashboardConfigCandidate(name string) bool {
	switch name {
	case "godrivelog.yaml", "godrivelog.yml", "dashboard.yaml", "dashboard.yml", "config.yaml", "config.yml":
		return true
	default:
		return false
	}
}

func sortedVehicleIDs(vehicles map[string]v3config.VehicleConfig) []string {
	ids := make([]string, 0, len(vehicles))
	for id := range vehicles {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func splitArgsAllowingSinglePositional(args []string, valueFlags map[string]bool) ([]string, []string, error) {
	filtered := make([]string, 0, len(args))
	positionals := []string{}

	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--" {
			positionals = append(positionals, args[index+1:]...)
			break
		}
		if flagName, isFlag := commandFlagName(arg); isFlag {
			filtered = append(filtered, arg)
			if valueFlags[flagName] && !strings.Contains(arg, "=") {
				if index+1 >= len(args) {
					return nil, nil, fmt.Errorf("flag needs an argument: %s", arg)
				}
				index++
				filtered = append(filtered, args[index])
			}
			continue
		}
		positionals = append(positionals, arg)
	}

	return filtered, positionals, nil
}

func commandFlagName(arg string) (string, bool) {
	if strings.HasPrefix(arg, "--") {
		if cut := strings.Index(arg, "="); cut >= 0 {
			return arg[:cut], true
		}
		return arg, true
	}
	return "", false
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func optionalPositiveFloat64Flag(fs *flag.FlagSet, name string, value float64) (*float64, bool, error) {
	if !flagWasProvided(fs, name) {
		return nil, false, nil
	}
	if value <= 0 {
		return nil, false, fmt.Errorf("--%s must be greater than zero", name)
	}
	return &value, true, nil
}

func optionalPositiveOrNegativeFloat64Flag(fs *flag.FlagSet, name string, value float64) (*float64, bool, error) {
	if !flagWasProvided(fs, name) {
		return nil, false, nil
	}
	return &value, true, nil
}

func flagWasProvided(fs *flag.FlagSet, name string) bool {
	provided := false
	fs.Visit(func(flag *flag.Flag) {
		if flag.Name == name {
			provided = true
		}
	})
	return provided
}

func formatSearchedConfigRecords(records []searchedConfigRecord) string {
	parts := make([]string, 0, len(records))
	for _, record := range records {
		parts = append(parts, fmt.Sprintf("%s [%s]", record.Path, strings.Join(record.Vehicles, ", ")))
	}
	return strings.Join(parts, "; ")
}

func exportDashboardExample(configPath, vehicleID, theme, outputDir string, force bool) error {
	sourceConfigPath := strings.TrimSpace(configPath)
	if sourceConfigPath == "" {
		sourceConfigPath = defaultExampleConfigPath(theme)
	}

	cfg, err := v3config.LoadFile(sourceConfigPath)
	if err != nil {
		return fmt.Errorf("load dashboard example config %q: %w", sourceConfigPath, err)
	}
	if vehicleID != "" {
		vehicleIDs := sortedVehicleIDs(cfg.Vehicles)
		if _, ok := cfg.Vehicles[vehicleID]; !ok {
			return fmt.Errorf("dashboard example config %q does not define vehicle %q; available vehicles: %s", sourceConfigPath, vehicleID, strings.Join(vehicleIDs, ", "))
		}
	}

	sourceRoot := filepath.Dir(sourceConfigPath)
	sourceAssets := filepath.Join(sourceRoot, "assets")
	if info, err := os.Stat(sourceAssets); err != nil {
		return fmt.Errorf("dashboard example assets %q: %w", sourceAssets, err)
	} else if !info.IsDir() {
		return fmt.Errorf("dashboard example assets %q must be a directory", sourceAssets)
	}

	if err := prepareExampleOutputDir(outputDir, force); err != nil {
		return err
	}
	if err := copyFile(sourceConfigPath, filepath.Join(outputDir, "dashboard.yaml")); err != nil {
		return err
	}
	if err := copyDirectory(sourceAssets, filepath.Join(outputDir, "assets")); err != nil {
		return err
	}
	return nil
}

func defaultExampleConfigPath(theme string) string {
	selectedTheme := strings.TrimSpace(theme)
	if selectedTheme == "" {
		selectedTheme = frameworkSmokeTheme
	}

	direct := filepath.Join("examples", selectedTheme, "dashboard.yaml")
	if _, err := os.Stat(direct); err == nil {
		return direct
	}

	cwd, err := os.Getwd()
	if err != nil {
		return direct
	}
	for current := cwd; ; current = filepath.Dir(current) {
		candidate := filepath.Join(current, "examples", selectedTheme, "dashboard.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return direct
}

func prepareExampleOutputDir(outputDir string, force bool) error {
	info, err := os.Stat(outputDir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("dashboard examples output %q must be a directory", outputDir)
		}
		empty, err := directoryIsEmpty(outputDir)
		if err != nil {
			return err
		}
		if empty {
			return nil
		}
		if !force {
			return fmt.Errorf("dashboard examples output %q is not empty; rerun with --force to replace it", outputDir)
		}
		return clearDirectory(outputDir)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("stat dashboard examples output %q: %w", outputDir, err)
	}
	return os.MkdirAll(outputDir, 0o755)
}

func directoryIsEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("read dashboard examples output %q: %w", path, err)
	}
	return len(entries) == 0, nil
}

func clearDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read dashboard examples output %q: %w", path, err)
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(path, entry.Name())); err != nil {
			return fmt.Errorf("clear dashboard examples output %q: %w", path, err)
		}
	}
	return nil
}

func copyDirectory(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("create dashboard example directory %q: %w", dst, err)
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == src {
			return nil
		}

		relative, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, relative)
		if d.IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("create dashboard example directory %q: %w", target, err)
			}
			return nil
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %q: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %q: %w", dst, err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}
	return nil
}

func withDashboardConfigPathEnv(configPath string, fn func() error) error {
	absolutePath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("resolve dashboard config path %q: %w", configPath, err)
	}
	previousValue, hadPrevious := os.LookupEnv(dashboardConfigEnvVar)
	if err := os.Setenv(dashboardConfigEnvVar, absolutePath); err != nil {
		return fmt.Errorf("set %s: %w", dashboardConfigEnvVar, err)
	}
	defer func() {
		if hadPrevious {
			_ = os.Setenv(dashboardConfigEnvVar, previousValue)
			return
		}
		_ = os.Unsetenv(dashboardConfigEnvVar)
	}()
	return fn()
}

func newDashboardFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

func writeRootHelp(w io.Writer) {
	fmt.Fprintln(w, "GoDriveLog command line")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  dashboard   Run dashboard-scoped commands for the active v3 Ebiten tooling")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  GoDriveLog dashboard harness vw_caddy --config ./examples/baseline-dashboard.yaml --pattern sweep --interval 50ms --duration 60s")
}

func writeDashboardHelp(w io.Writer) {
	fmt.Fprintln(w, "GoDriveLog dashboard [--config <config-file>]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Purpose:")
	fmt.Fprintln(w, "  Print a compact dashboard config overview or run dashboard-scoped commands through the active v3 command tree.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --config string")
	fmt.Fprintln(w, "        path to a dashboard config; when omitted, search the current directory and /etc/godrivelog")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Subcommands:")
	fmt.Fprintln(w, "  run       Start the active dashboard runtime for a selected vehicle")
	fmt.Fprintln(w, "  harness   Run the dashboard harness for a selected vehicle")
	fmt.Fprintln(w, "  preview   Open a single-gauge preview viewer for one preview YAML file")
	fmt.Fprintln(w, "  examples  Export a self-contained generated example dashboard")
	fmt.Fprintln(w, "  validate  Validate a dashboard config file")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  GoDriveLog dashboard --config ./examples/baseline-dashboard.yaml")
}

func writeDashboardRunHelp(w io.Writer) {
	fmt.Fprintln(w, "GoDriveLog dashboard run [vehicle-id]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Purpose:")
	fmt.Fprintln(w, "  Start the active dashboard runtime for the selected vehicle.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Positional arguments:")
	fmt.Fprintln(w, "  vehicle-id   optional vehicle id; required when the resolved config defines multiple vehicles")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --config string")
	fmt.Fprintln(w, "        path to a dashboard config; when omitted, search the current directory and /etc/godrivelog")
	fmt.Fprintln(w, "  --renderer string")
	fmt.Fprintf(w, "        dashboard renderer backend (default %q)\n", v3RendererEbiten)
	fmt.Fprintln(w, "  --duration duration")
	fmt.Fprintln(w, "        optional runtime duration such as 60s; zero runs until interrupted (default 0s)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  GoDriveLog dashboard run vw_caddy --config ./examples/baseline-dashboard.yaml --renderer ebiten")
}

func writeDashboardHarnessHelp(w io.Writer) {
	fmt.Fprintln(w, "GoDriveLog dashboard harness [vehicle-id]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Purpose:")
	fmt.Fprintln(w, "  Run the dashboard harness through the active command tree without OBD hardware.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Positional arguments:")
	fmt.Fprintln(w, "  vehicle-id   optional vehicle id; required when the resolved config defines multiple vehicles")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --config string")
	fmt.Fprintln(w, "        path to a dashboard config; when omitted, search the current directory and /etc/godrivelog")
	fmt.Fprintln(w, "  --pattern string")
	fmt.Fprintf(w, "        harness pattern: %s, %s, or %s (default %q)\n", v3harness.PatternSweep, v3harness.PatternHeartbeat, v3harness.PatternFixed, v3harness.PatternSweep)
	fmt.Fprintln(w, "  --interval duration")
	fmt.Fprintln(w, "        harness update interval such as 50ms or 100ms (default 100ms)")
	fmt.Fprintln(w, "  --duration duration")
	fmt.Fprintln(w, "        optional harness duration such as 60s; zero runs until interrupted (default 0s)")
	fmt.Fprintln(w, "  --renderer string")
	fmt.Fprintf(w, "        dashboard renderer backend (default %q)\n", v3RendererEbiten)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  GoDriveLog dashboard harness vw_caddy --config ./examples/baseline-dashboard.yaml --pattern sweep --interval 50ms --duration 60s --renderer ebiten")
}

func writeDashboardExamplesHelp(w io.Writer) {
	fmt.Fprintln(w, "GoDriveLog dashboard examples --output <directory>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Purpose:")
	fmt.Fprintln(w, "  Export a self-contained generated example dashboard with dashboard.yaml and assets/ at the output root.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --output string")
	fmt.Fprintln(w, "        required output directory for the exported dashboard")
	fmt.Fprintln(w, "  --theme string")
	fmt.Fprintf(w, "        built-in example theme to export when --config is not supplied (default %q)\n", frameworkSmokeTheme)
	fmt.Fprintln(w, "  --config string")
	fmt.Fprintln(w, "        optional source dashboard config; when supplied, it bypasses theme lookup")
	fmt.Fprintln(w, "  --vehicle string")
	fmt.Fprintln(w, "        optional vehicle id to verify in the source dashboard config")
	fmt.Fprintln(w, "  --force")
	fmt.Fprintln(w, "        replace an existing non-empty output directory")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  GoDriveLog dashboard examples --theme framework-smoke --output ./tmp/framework-smoke")
}

func writeDashboardPreviewHelp(w io.Writer) {
	fmt.Fprintln(w, "GoDriveLog dashboard preview <file>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Purpose:")
	fmt.Fprintln(w, "  Open one normal preview YAML file and interactively inspect a single gauge.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Positional arguments:")
	fmt.Fprintln(w, "  file   required preview YAML file; it must resolve to exactly one selected vehicle")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --gauge string")
	fmt.Fprintln(w, "        optional gauge widget id; use dashboard/widget when the file contains duplicate widget ids")
	fmt.Fprintln(w, "  --value float")
	fmt.Fprintln(w, "        optional starting value override; default is the midpoint of the inferred range")
	fmt.Fprintln(w, "  --step float")
	fmt.Fprintln(w, "        optional step size for Up/Down")
	fmt.Fprintln(w, "  --fine-step float")
	fmt.Fprintln(w, "        optional fine step size for Ctrl/Cmd+Up/Down")
	fmt.Fprintln(w, "  --coarse-step float")
	fmt.Fprintln(w, "        optional coarse step size for Shift+Up/Down")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Controls:")
	fmt.Fprintln(w, "  Left=min  Right=max  Up/Down=step  Shift+Up/Down=coarse  Ctrl/Cmd+Up/Down=fine")
	fmt.Fprintln(w, "  R=midpoint  Space=replay last transition  Esc/Q=quit")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  GoDriveLog dashboard preview ./examples/gauge-realism/radial/00-baseline.yaml")
}

func writeDashboardValidateHelp(w io.Writer) {
	fmt.Fprintln(w, "GoDriveLog dashboard validate [config-file]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Purpose:")
	fmt.Fprintln(w, "  Validate a dashboard config using the active config parser and validation helpers.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Positional arguments:")
	fmt.Fprintln(w, "  config-file  optional config path; when omitted, use deterministic discovery")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --config string")
	fmt.Fprintln(w, "        path to a dashboard config; must not be combined with a positional config file")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Example:")
	fmt.Fprintln(w, "  GoDriveLog dashboard validate ./examples/baseline-dashboard.yaml")
}

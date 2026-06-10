package ui

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

type PNGDigitConfig struct {
	AssetRoot string
	Glyphs    map[string]string
}

type ConfiguredInstrumentDashboard struct {
	base        *InstrumentDashboard
	store       *sensors.StateStore
	rpm         *PNGSevenSegmentDisplay
	speed       *PNGSevenSegmentDisplay
	gear        *PNGSevenSegmentDisplay
	oilTemp     *PNGSevenSegmentDisplay
	oilPressure *PNGSevenSegmentDisplay
	coolant     *PNGSevenSegmentDisplay
	battery     *PNGSevenSegmentDisplay
}

func NewInstrumentDashboard1920x480WithPNGDigitConfig(store *sensors.StateStore, options InstrumentDashboardOptions, digitConfig PNGDigitConfig) (*ConfiguredInstrumentDashboard, error) {
	base, err := NewInstrumentDashboard1920x480WithOptions(store, options)
	if err != nil {
		return nil, err
	}

	dashboard := &ConfiguredInstrumentDashboard{base: base, store: store}
	resources := loadConfiguredPNGDigitResources(digitConfig)
	if resources == nil {
		return dashboard, nil
	}
	dashboard.installPNGDigitDisplays(resources)
	dashboard.RefreshPNGDigits()
	return dashboard, nil
}

func (d *ConfiguredInstrumentDashboard) CanvasObject() fyne.CanvasObject {
	return d.base.CanvasObject()
}

func (d *ConfiguredInstrumentDashboard) Start(ctx context.Context, interval time.Duration) {
	d.base.Start(ctx, interval)
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fyne.Do(d.RefreshPNGDigits)
			}
		}
	}()
}

func (d *ConfiguredInstrumentDashboard) RefreshPNGDigits() {
	if d.rpm == nil {
		return
	}
	now := time.Now()
	states := d.base.stateMap(d.store.SnapshotWithStale(now))

	rpm := sensorValue(states, "rpm")
	speed := sensorValue(states, "speed")
	oilTemp := sensorValue(states, "oil_temperature", "oil_temp")
	oilPressure := sensorValue(states, "oil_pressure")
	coolant := sensorValue(states, "coolant_temp", "coolant_temperature")
	battery := sensorValue(states, "battery_voltage", "battery")
	gear := sensorValue(states, "gear")

	d.rpm.SetValue(fmt.Sprintf("%04.0f", rpm), 4)
	d.speed.SetValue(fmt.Sprintf("%03.0f", speed), 3)
	d.gear.SetValue(gearText(gear), 1)
	d.oilTemp.SetValue(fmt.Sprintf("%.1f", oilTemp), 4)
	d.oilPressure.SetValue(fmt.Sprintf("%.1f", oilPressure), 4)
	d.coolant.SetValue(fmt.Sprintf("%.1f", coolant), 4)
	d.battery.SetValue(fmt.Sprintf("%.1f", battery), 4)
}

func (d *ConfiguredInstrumentDashboard) installPNGDigitDisplays(resources *pngDigitResources) {
	d.base.rpmText.Hide()
	d.base.speedText.Hide()
	d.base.gearText.Hide()
	d.base.oilTempText.Hide()
	d.base.oilPressureText.Hide()
	d.base.coolantText.Hide()
	d.base.batteryText.Hide()

	d.rpm = newConfiguredPNGSevenSegmentDisplay("0000", 4, 58, 78, 560, 142, resources)
	d.speed = newConfiguredPNGSevenSegmentDisplay("000", 3, 704, 78, 390, 154, resources)
	d.gear = newConfiguredPNGSevenSegmentDisplay("0", 1, 700, 276, 110, 94, resources)
	d.oilTemp = newConfiguredPNGSevenSegmentDisplay("000.0", 4, 1240, 72, 260, 56, resources)
	d.oilPressure = newConfiguredPNGSevenSegmentDisplay("000.0", 4, 1575, 72, 260, 56, resources)
	d.coolant = newConfiguredPNGSevenSegmentDisplay("000.0", 4, 1240, 166, 260, 56, resources)
	d.battery = newConfiguredPNGSevenSegmentDisplay("00.0", 4, 1575, 166, 220, 56, resources)

	d.base.root.Add(d.rpm.CanvasObject())
	d.base.root.Add(d.speed.CanvasObject())
	d.base.root.Add(d.gear.CanvasObject())
	d.base.root.Add(d.oilTemp.CanvasObject())
	d.base.root.Add(d.oilPressure.CanvasObject())
	d.base.root.Add(d.coolant.CanvasObject())
	d.base.root.Add(d.battery.CanvasObject())
}

func newConfiguredPNGSevenSegmentDisplay(value string, digits int, x, y, width, height float32, resources *pngDigitResources) *PNGSevenSegmentDisplay {
	if resources == nil {
		resources = loadPNGDigitResources(defaultSevenSegmentAssetDir)
	}
	display := &PNGSevenSegmentDisplay{
		resources: resources,
		digitSize: fyne.NewSize(width/float32(maxInt(digits, 1)), height),
	}
	display.root = container.NewWithoutLayout()
	display.root.Move(fyne.NewPos(x, y))
	display.root.Resize(fyne.NewSize(width, height))
	display.SetValue(value, digits)
	return display
}

func loadConfiguredPNGDigitResources(config PNGDigitConfig) *pngDigitResources {
	if len(config.Glyphs) == 0 {
		return nil
	}
	resources := &pngDigitResources{dir: config.AssetRoot, bySymbol: make(map[string]fyne.Resource, 12)}
	for _, symbol := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "dash", "blank"} {
		resources.bySymbol[symbol] = loadPNGResource(resolveConfiguredAssetPath(config.AssetRoot, config.Glyphs[symbol]))
	}
	resources.dp = loadPNGResource(resolveConfiguredAssetPath(config.AssetRoot, config.Glyphs["dp"]))
	return resources
}

func resolveConfiguredAssetPath(assetRoot, path string) string {
	if path == "" || filepath.IsAbs(path) || assetRoot == "" {
		return path
	}
	return filepath.Join(assetRoot, path)
}

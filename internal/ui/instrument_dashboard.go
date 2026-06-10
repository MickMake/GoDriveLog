package ui

import (
	"context"
	"fmt"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

const (
	instrumentWidth  float32 = 1920
	instrumentHeight float32 = 480
)

type InstrumentDashboard struct {
	root  *fyne.Container
	store *sensors.StateStore

	rpmText           *canvas.Text
	speedText         *canvas.Text
	throttleText      *canvas.Text
	throttleFill      *canvas.Rectangle
	oilTempText       *canvas.Text
	oilPressureText   *canvas.Text
	gearText          *canvas.Text
	warningText       *canvas.Text
	engineFailedText  *canvas.Text
	requiresResetText *canvas.Text
	statusText        *canvas.Text
	failureOverlay    *canvas.Rectangle
}

func NewInstrumentDashboard1920x480(store *sensors.StateStore) (*InstrumentDashboard, error) {
	if store == nil {
		return nil, fmt.Errorf("state store must not be nil")
	}

	dashboard := &InstrumentDashboard{store: store}
	background := canvas.NewRectangle(color.NRGBA{R: 6, G: 8, B: 12, A: 255})
	background.Resize(fyne.NewSize(instrumentWidth, instrumentHeight))

	rpmLabel := labelText("RPM", 80, 32, 30)
	dashboard.rpmText = valueText("0000", 80, 74, 96)

	speedLabel := labelText("SPEED km/h", 720, 32, 30)
	dashboard.speedText = valueText("000", 720, 74, 128)

	gearLabel := labelText("GEAR", 1245, 32, 30)
	dashboard.gearText = valueText("N", 1245, 70, 132)

	throttleLabel := labelText("THROTTLE", 80, 275, 26)
	dashboard.throttleText = valueText("0%", 80, 330, 44)
	throttleTrack := canvas.NewRectangle(color.NRGBA{R: 28, G: 32, B: 42, A: 255})
	throttleTrack.Move(fyne.NewPos(80, 315))
	throttleTrack.Resize(fyne.NewSize(520, 36))
	dashboard.throttleFill = canvas.NewRectangle(color.NRGBA{R: 40, G: 220, B: 95, A: 255})
	dashboard.throttleFill.Move(fyne.NewPos(80, 315))
	dashboard.throttleFill.Resize(fyne.NewSize(0, 36))

	oilTempLabel := labelText("OIL TEMP", 1440, 35, 26)
	dashboard.oilTempText = valueText("-- C", 1440, 75, 48)
	oilPressureLabel := labelText("OIL PRESSURE", 1440, 155, 26)
	dashboard.oilPressureText = valueText("-- kPa", 1440, 195, 48)
	warningLabel := labelText("WARNING", 1440, 275, 26)
	dashboard.warningText = valueText("NONE", 1440, 315, 42)

	dashboard.engineFailedText = labelText("ENGINE FAILED: NO", 80, 410, 28)
	dashboard.requiresResetText = labelText("REQUIRES RESET: NO", 440, 410, 28)
	dashboard.statusText = labelText("fast instrument renderer: StateStore direct updates", 930, 410, 24)

	dashboard.failureOverlay = canvas.NewRectangle(color.NRGBA{R: 120, G: 0, B: 0, A: 90})
	dashboard.failureOverlay.Move(fyne.NewPos(0, 0))
	dashboard.failureOverlay.Resize(fyne.NewSize(instrumentWidth, instrumentHeight))
	dashboard.failureOverlay.Hide()

	root := container.NewWithoutLayout(
		background,
		rpmLabel,
		dashboard.rpmText,
		speedLabel,
		dashboard.speedText,
		gearLabel,
		dashboard.gearText,
		throttleLabel,
		throttleTrack,
		dashboard.throttleFill,
		dashboard.throttleText,
		oilTempLabel,
		dashboard.oilTempText,
		oilPressureLabel,
		dashboard.oilPressureText,
		warningLabel,
		dashboard.warningText,
		dashboard.engineFailedText,
		dashboard.requiresResetText,
		dashboard.statusText,
		dashboard.failureOverlay,
	)
	root.Resize(fyne.NewSize(instrumentWidth, instrumentHeight))
	dashboard.root = root
	dashboard.Refresh()

	return dashboard, nil
}

func (d *InstrumentDashboard) CanvasObject() fyne.CanvasObject {
	return d.root
}

func (d *InstrumentDashboard) Start(ctx context.Context, interval time.Duration) {
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
				fyne.Do(d.Refresh)
			}
		}
	}()
}

func (d *InstrumentDashboard) Refresh() {
	states := stateMap(d.store.SnapshotWithStale(time.Now()))

	rpm := sensorValue(states, "rpm")
	speed := sensorValue(states, "speed")
	throttle := sensorValue(states, "throttle_position", "throttle")
	oilTemp := sensorValue(states, "oil_temperature", "oil_temp")
	oilPressure := sensorValue(states, "oil_pressure")
	gear := sensorValue(states, "gear")
	warning := sensorValue(states, "warning_level")
	engineFailed := sensorValue(states, "engine_failed") >= 0.5
	requiresReset := sensorValue(states, "requires_reset") >= 0.5

	setText(d.rpmText, fmt.Sprintf("%04.0f", rpm))
	setText(d.speedText, fmt.Sprintf("%03.0f", speed))
	setText(d.throttleText, fmt.Sprintf("%.0f%%", throttle))
	setText(d.oilTempText, fmt.Sprintf("%.1f C", oilTemp))
	setText(d.oilPressureText, fmt.Sprintf("%.1f kPa", oilPressure))
	setText(d.gearText, gearText(gear))
	setText(d.warningText, warningText(warning))
	setText(d.engineFailedText, "ENGINE FAILED: "+boolText(engineFailed))
	setText(d.requiresResetText, "REQUIRES RESET: "+boolText(requiresReset))

	clampedThrottle := float32(math.Max(0, math.Min(100, throttle)))
	d.throttleFill.Resize(fyne.NewSize(520*(clampedThrottle/100), 36))
	d.throttleFill.Refresh()

	applyWarningColor(d.warningText, warning)
	applyFailureState(d, engineFailed, requiresReset)
}

func stateMap(states []sensors.SensorState) map[string]sensors.SensorState {
	mapped := make(map[string]sensors.SensorState, len(states))
	for _, state := range states {
		mapped[state.ID] = state
	}
	return mapped
}

func sensorValue(states map[string]sensors.SensorState, ids ...string) float64 {
	for _, id := range ids {
		if state, ok := states[id]; ok && state.Status == sensors.StatusOK {
			return state.Value
		}
	}
	return 0
}

func labelText(text string, x, y, size float32) *canvas.Text {
	label := canvas.NewText(text, color.NRGBA{R: 160, G: 180, B: 200, A: 255})
	label.TextSize = size
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Move(fyne.NewPos(x, y))
	return label
}

func valueText(text string, x, y, size float32) *canvas.Text {
	value := canvas.NewText(text, color.NRGBA{R: 235, G: 245, B: 255, A: 255})
	value.TextSize = size
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	value.Move(fyne.NewPos(x, y))
	return value
}

func setText(text *canvas.Text, value string) {
	if text.Text == value {
		return
	}
	text.Text = value
	text.Refresh()
}

func gearText(value float64) string {
	gear := int(math.Round(value))
	if gear <= 0 {
		return "N"
	}
	return fmt.Sprintf("%d", gear)
}

func warningText(value float64) string {
	switch int(math.Round(value)) {
	case 1:
		return "WARNING"
	case 2:
		return "CRITICAL"
	default:
		return "NONE"
	}
}

func boolText(value bool) string {
	if value {
		return "YES"
	}
	return "NO"
}

func applyWarningColor(text *canvas.Text, warning float64) {
	switch int(math.Round(warning)) {
	case 1:
		text.Color = color.NRGBA{R: 255, G: 210, B: 60, A: 255}
	case 2:
		text.Color = color.NRGBA{R: 255, G: 80, B: 80, A: 255}
	default:
		text.Color = color.NRGBA{R: 90, G: 230, B: 130, A: 255}
	}
	text.Refresh()
}

func applyFailureState(d *InstrumentDashboard, engineFailed, requiresReset bool) {
	if engineFailed || requiresReset {
		d.failureOverlay.Show()
		d.engineFailedText.Color = color.NRGBA{R: 255, G: 80, B: 80, A: 255}
		d.requiresResetText.Color = color.NRGBA{R: 255, G: 80, B: 80, A: 255}
		setText(d.statusText, "CRITICAL FAILURE LATCHED - RESET REQUIRED")
	} else {
		d.failureOverlay.Hide()
		d.engineFailedText.Color = color.NRGBA{R: 160, G: 180, B: 200, A: 255}
		d.requiresResetText.Color = color.NRGBA{R: 160, G: 180, B: 200, A: 255}
		setText(d.statusText, "fast instrument renderer: StateStore direct updates")
	}
	d.engineFailedText.Refresh()
	d.requiresResetText.Refresh()
}

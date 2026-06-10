package ui

import (
	"context"
	"fmt"
	"image/color"
	"math"
	"strings"
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

var (
	colourBackground         = color.NRGBA{R: 3, G: 5, B: 9, A: 255}
	colourPanel              = color.NRGBA{R: 10, G: 15, B: 23, A: 255}
	colourPanelHot           = color.NRGBA{R: 42, G: 6, B: 6, A: 255}
	colourPanelWarn          = color.NRGBA{R: 42, G: 31, B: 6, A: 255}
	colourTextDim            = color.NRGBA{R: 116, G: 135, B: 156, A: 255}
	colourTextNormal         = color.NRGBA{R: 222, G: 238, B: 255, A: 255}
	colourGreen              = color.NRGBA{R: 60, G: 235, B: 125, A: 255}
	colourAmber              = color.NRGBA{R: 255, G: 202, B: 55, A: 255}
	colourRed                = color.NRGBA{R: 255, G: 68, B: 68, A: 255}
	colourThrottle           = color.NRGBA{R: 62, G: 230, B: 110, A: 255}
	colourThrottleWarn       = color.NRGBA{R: 255, G: 190, B: 45, A: 255}
	colourTrack              = color.NRGBA{R: 30, G: 37, B: 50, A: 255}
	colourOverlay            = color.NRGBA{R: 140, G: 0, B: 0, A: 96}
	colourAlertBackground    = color.NRGBA{R: 8, G: 11, B: 17, A: 255}
)

type InstrumentDashboard struct {
	root  *fyne.Container
	store *sensors.StateStore

	states       map[string]sensors.SensorState
	statusIssues []string

	rpmPanel         *canvas.Rectangle
	speedPanel       *canvas.Rectangle
	rightPanel       *canvas.Rectangle
	bottomPanel      *canvas.Rectangle
	alertBackground *canvas.Rectangle
	failureOverlay  *canvas.Rectangle

	rpmText           *canvas.Text
	speedText         *canvas.Text
	throttleText      *canvas.Text
	throttleFill      *canvas.Rectangle
	engineLoadText    *canvas.Text
	engineLoadFill    *canvas.Rectangle
	oilTempText       *canvas.Text
	oilPressureText   *canvas.Text
	coolantText       *canvas.Text
	batteryText       *canvas.Text
	gearText          *canvas.Text
	warningText       *canvas.Text
	engineFailedText  *canvas.Text
	requiresResetText *canvas.Text
	alertText         *canvas.Text
	statusText        *canvas.Text
}

func NewInstrumentDashboard1920x480(store *sensors.StateStore) (*InstrumentDashboard, error) {
	if store == nil {
		return nil, fmt.Errorf("state store must not be nil")
	}

	dashboard := &InstrumentDashboard{
		store:        store,
		states:       make(map[string]sensors.SensorState, 16),
		statusIssues: make([]string, 0, 12),
	}
	background := rect(0, 0, instrumentWidth, instrumentHeight, colourBackground)

	dashboard.rpmPanel = rect(28, 24, 620, 350, colourPanel)
	dashboard.speedPanel = rect(670, 24, 520, 350, colourPanel)
	dashboard.rightPanel = rect(1212, 24, 680, 350, colourPanel)
	dashboard.bottomPanel = rect(28, 392, 1864, 64, colourAlertBackground)

	rpmLabel := labelText("RPM", 58, 42, 34)
	dashboard.rpmText = valueText("0000", 58, 78, 142, colourAmber)
	throttleLabel := labelText("THROTTLE", 58, 238, 24)
	dashboard.throttleText = valueText("0%", 58, 274, 42, colourTextNormal)
	throttleTrack := rect(180, 282, 420, 32, colourTrack)
	dashboard.throttleFill = rect(180, 282, 0, 32, colourThrottle)
	loadLabel := labelText("ENGINE LOAD", 58, 320, 20)
	dashboard.engineLoadText = valueText("0%", 58, 344, 28, colourTextNormal)
	loadTrack := rect(250, 350, 350, 18, colourTrack)
	dashboard.engineLoadFill = rect(250, 350, 0, 18, colourGreen)

	speedLabel := labelText("SPEED km/h", 704, 42, 34)
	dashboard.speedText = valueText("000", 704, 78, 154, colourGreen)
	gearLabel := labelText("GEAR", 705, 246, 28)
	dashboard.gearText = valueText("N", 700, 276, 94, colourTextNormal)

	oilTempLabel := labelText("OIL TEMP", 1240, 44, 23)
	dashboard.oilTempText = valueText("--.- C", 1240, 72, 40, colourTextNormal)
	oilPressureLabel := labelText("OIL PRESSURE", 1575, 44, 23)
	dashboard.oilPressureText = valueText("--.- kPa", 1575, 72, 40, colourTextNormal)
	coolantLabel := labelText("COOLANT", 1240, 138, 23)
	dashboard.coolantText = valueText("--.- C", 1240, 166, 40, colourTextNormal)
	batteryLabel := labelText("BATTERY", 1575, 138, 23)
	dashboard.batteryText = valueText("--.- V", 1575, 166, 40, colourTextNormal)
	warningLabel := labelText("WARNING", 1240, 232, 23)
	dashboard.warningText = valueText("NONE", 1240, 262, 40, colourGreen)
	dashboard.engineFailedText = labelText("ENGINE FAILED: NO", 1240, 324, 24)
	dashboard.requiresResetText = labelText("RESET REQUIRED: NO", 1575, 324, 24)

	dashboard.statusText = labelText("Race demo status messages are derived from sensor values; source text is not a sensor yet.", 52, 408, 20)
	dashboard.alertBackground = rect(28, 392, 1864, 64, colourAlertBackground)
	dashboard.alertText = valueText("SYSTEM NORMAL", 52, 424, 26, colourTextNormal)

	dashboard.failureOverlay = rect(0, 0, instrumentWidth, instrumentHeight, colourOverlay)
	dashboard.failureOverlay.Hide()

	root := container.NewWithoutLayout(
		background,
		dashboard.rpmPanel,
		dashboard.speedPanel,
		dashboard.rightPanel,
		dashboard.bottomPanel,
		rpmLabel,
		dashboard.rpmText,
		throttleLabel,
		throttleTrack,
		dashboard.throttleFill,
		dashboard.throttleText,
		loadLabel,
		loadTrack,
		dashboard.engineLoadFill,
		dashboard.engineLoadText,
		speedLabel,
		dashboard.speedText,
		gearLabel,
		dashboard.gearText,
		oilTempLabel,
		dashboard.oilTempText,
		oilPressureLabel,
		dashboard.oilPressureText,
		coolantLabel,
		dashboard.coolantText,
		batteryLabel,
		dashboard.batteryText,
		warningLabel,
		dashboard.warningText,
		dashboard.engineFailedText,
		dashboard.requiresResetText,
		dashboard.failureOverlay,
		dashboard.alertBackground,
		dashboard.statusText,
		dashboard.alertText,
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
	states := d.stateMap(d.store.SnapshotWithStale(time.Now()))

	rpm := sensorValue(states, "rpm")
	speed := sensorValue(states, "speed")
	throttle := sensorValue(states, "throttle_position", "throttle")
	engineLoad := sensorValue(states, "engine_load")
	oilTemp := sensorValue(states, "oil_temperature", "oil_temp")
	oilPressure := sensorValue(states, "oil_pressure")
	coolant := sensorValue(states, "coolant_temp", "coolant_temperature")
	battery := sensorValue(states, "battery_voltage", "battery")
	gear := sensorValue(states, "gear")
	warning := sensorValue(states, "warning_level")
	engineFailed := sensorValue(states, "engine_failed") >= 0.5
	requiresReset := sensorValue(states, "requires_reset") >= 0.5

	setText(d.rpmText, fmt.Sprintf("%04.0f", rpm))
	setText(d.speedText, fmt.Sprintf("%03.0f", speed))
	setText(d.throttleText, fmt.Sprintf("%.0f%%", throttle))
	setText(d.engineLoadText, fmt.Sprintf("%.0f%%", engineLoad))
	setText(d.oilTempText, fmt.Sprintf("%.1f C", oilTemp))
	setText(d.oilPressureText, fmt.Sprintf("%.1f kPa", oilPressure))
	setText(d.coolantText, fmt.Sprintf("%.1f C", coolant))
	setText(d.batteryText, fmt.Sprintf("%.1f V", battery))
	setText(d.gearText, gearText(gear))
	setText(d.warningText, warningText(warning))
	setText(d.engineFailedText, "ENGINE FAILED: "+boolText(engineFailed))
	setText(d.requiresResetText, "RESET REQUIRED: "+boolText(requiresReset))

	setBar(d.throttleFill, 420, 32, throttle)
	setBar(d.engineLoadFill, 350, 18, engineLoad)
	applyInstrumentColors(d, rpm, speed, throttle, engineLoad, oilTemp, oilPressure, coolant, battery, warning, engineFailed, requiresReset)
	setText(d.statusText, statusLine(d.sensorStatusText(states)))
	setText(d.alertText, alertLine(rpm, speed, throttle, oilTemp, oilPressure, warning, engineFailed, requiresReset))
}

func (d *InstrumentDashboard) stateMap(states []sensors.SensorState) map[string]sensors.SensorState {
	for key := range d.states {
		delete(d.states, key)
	}
	for _, state := range states {
		d.states[state.ID] = state
	}
	return d.states
}

func sensorValue(states map[string]sensors.SensorState, ids ...string) float64 {
	for _, id := range ids {
		if state, ok := states[id]; ok {
			return state.Value
		}
	}
	return 0
}

func (d *InstrumentDashboard) sensorStatusText(states map[string]sensors.SensorState) string {
	d.statusIssues = d.statusIssues[:0]
	d.appendSensorIssue(states, "rpm", "rpm")
	d.appendSensorIssue(states, "speed", "speed")
	d.appendSensorIssue(states, "throttle", "throttle_position", "throttle")
	d.appendSensorIssue(states, "load", "engine_load")
	d.appendSensorIssue(states, "oil_temp", "oil_temperature", "oil_temp")
	d.appendSensorIssue(states, "oil_pressure", "oil_pressure")
	d.appendSensorIssue(states, "coolant", "coolant_temp", "coolant_temperature")
	d.appendSensorIssue(states, "battery", "battery_voltage", "battery")
	d.appendSensorIssue(states, "gear", "gear")
	d.appendSensorIssue(states, "warning", "warning_level")
	d.appendSensorIssue(states, "engine_failed", "engine_failed")
	d.appendSensorIssue(states, "requires_reset", "requires_reset")

	if len(d.statusIssues) == 0 {
		return ""
	}
	return "SENSOR STATUS: " + strings.Join(d.statusIssues, ", ")
}

func (d *InstrumentDashboard) appendSensorIssue(states map[string]sensors.SensorState, label string, ids ...string) {
	for _, id := range ids {
		state, ok := states[id]
		if !ok {
			continue
		}
		switch state.Status {
		case sensors.StatusStale, sensors.StatusError:
			d.statusIssues = append(d.statusIssues, label+" "+state.Status)
		}
		return
	}
}

func statusLine(sensorStatus string) string {
	status := "fast instrument renderer: direct StateStore updates | status text is derived; race-demo source message is not exposed as a sensor yet"
	if sensorStatus != "" {
		return status + " | " + sensorStatus
	}
	return status
}

func alertLine(rpm, speed, throttle, oilTemp, oilPressure, warning float64, engineFailed, requiresReset bool) string {
	switch {
	case engineFailed || requiresReset:
		return "ENGINE FAILURE - THROWN ROD | RPM 0, oil pressure 0, reset required"
	case warning >= 2:
		return "CRITICAL OIL TEMPERATURE - EMERGENCY DOWNSHIFTS"
	case warning >= 1:
		return "OIL TEMP WARNING - DRIVER CONTINUES"
	case speed == 0 && rpm >= 2500:
		return "STATIONARY REVVING"
	case speed > 0 && speed < 15 && rpm >= 4500 && throttle >= 90:
		return "BURNOUT LAUNCH"
	case speed >= 145 && speed <= 155:
		return "150 km/h CRUISE"
	case throttle >= 90 && speed > 15:
		return "HARD ACCELERATION"
	case rpm > 0 && speed == 0:
		return "IDLE"
	case oilTemp > 0 || oilPressure > 0:
		return "SYSTEM NORMAL"
	default:
		return "WAITING FOR SENSOR DATA"
	}
}

func labelText(text string, x, y, size float32) *canvas.Text {
	label := canvas.NewText(text, colourTextDim)
	label.TextSize = size
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Move(fyne.NewPos(x, y))
	return label
}

func valueText(text string, x, y, size float32, textColor color.Color) *canvas.Text {
	value := canvas.NewText(text, textColor)
	value.TextSize = size
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	value.Move(fyne.NewPos(x, y))
	return value
}

func rect(x, y, width, height float32, fill color.Color) *canvas.Rectangle {
	rectangle := canvas.NewRectangle(fill)
	rectangle.Move(fyne.NewPos(x, y))
	rectangle.Resize(fyne.NewSize(width, height))
	return rectangle
}

func setText(text *canvas.Text, value string) {
	if text.Text == value {
		return
	}
	text.Text = value
	text.Refresh()
}

func setTextColor(text *canvas.Text, textColor color.Color) {
	if sameColor(text.Color, textColor) {
		return
	}
	text.Color = textColor
	text.Refresh()
}

func setRectColor(rectangle *canvas.Rectangle, fill color.Color) {
	if sameColor(rectangle.FillColor, fill) {
		return
	}
	rectangle.FillColor = fill
	rectangle.Refresh()
}

func setBar(fill *canvas.Rectangle, maxWidth float32, height float32, value float64) {
	clamped := float32(math.Max(0, math.Min(100, value)))
	nextSize := fyne.NewSize(maxWidth*(clamped/100), height)
	if fill.Size() == nextSize {
		return
	}
	fill.Resize(nextSize)
	fill.Refresh()
}

func setVisible(object fyne.CanvasObject, visible bool) {
	if object.Visible() == visible {
		return
	}
	if visible {
		object.Show()
		return
	}
	object.Hide()
}

func sameColor(a, b color.Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
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

func applyInstrumentColors(d *InstrumentDashboard, rpm, speed, throttle, engineLoad, oilTemp, oilPressure, coolant, battery, warning float64, engineFailed, requiresReset bool) {
	severity := int(math.Round(warning))
	if engineFailed || requiresReset {
		severity = 3
	}

	switch severity {
	case 3:
		setRectColor(d.rpmPanel, colourPanelHot)
		setRectColor(d.speedPanel, colourPanelHot)
		setRectColor(d.rightPanel, colourPanelHot)
		setRectColor(d.alertBackground, colourRed)
		setVisible(d.failureOverlay, true)
	case 2:
		setRectColor(d.rpmPanel, colourPanel)
		setRectColor(d.speedPanel, colourPanel)
		setRectColor(d.rightPanel, colourPanelHot)
		setRectColor(d.alertBackground, colourPanelHot)
		setVisible(d.failureOverlay, false)
	case 1:
		setRectColor(d.rpmPanel, colourPanel)
		setRectColor(d.speedPanel, colourPanel)
		setRectColor(d.rightPanel, colourPanelWarn)
		setRectColor(d.alertBackground, colourPanelWarn)
		setVisible(d.failureOverlay, false)
	default:
		setRectColor(d.rpmPanel, colourPanel)
		setRectColor(d.speedPanel, colourPanel)
		setRectColor(d.rightPanel, colourPanel)
		setRectColor(d.alertBackground, colourAlertBackground)
		setVisible(d.failureOverlay, false)
	}

	setTextColor(d.rpmText, rpmColor(rpm, engineFailed))
	setTextColor(d.speedText, speedColor(speed, engineFailed))
	setTextColor(d.throttleText, loadColor(throttle))
	setRectColor(d.throttleFill, loadColor(throttle))
	setTextColor(d.engineLoadText, loadColor(engineLoad))
	setRectColor(d.engineLoadFill, loadColor(engineLoad))
	setTextColor(d.oilTempText, temperatureColor(oilTemp, 118, 132))
	setTextColor(d.oilPressureText, oilPressureColor(oilPressure, engineFailed))
	setTextColor(d.coolantText, temperatureColor(coolant, 104, 110))
	setTextColor(d.batteryText, batteryColor(battery))
	setTextColor(d.warningText, warningColor(warning, engineFailed || requiresReset))
	setTextColor(d.engineFailedText, boolColor(engineFailed))
	setTextColor(d.requiresResetText, boolColor(requiresReset))
	setTextColor(d.alertText, alertTextColor(severity))
}

func rpmColor(rpm float64, engineFailed bool) color.Color {
	if engineFailed || rpm == 0 {
		return colourRed
	}
	if rpm >= 4800 {
		return colourAmber
	}
	return colourAmber
}

func speedColor(speed float64, engineFailed bool) color.Color {
	if engineFailed && speed > 0 {
		return colourAmber
	}
	return colourGreen
}

func loadColor(value float64) color.Color {
	if value >= 90 {
		return colourThrottleWarn
	}
	return colourThrottle
}

func temperatureColor(value, warnAt, criticalAt float64) color.Color {
	switch {
	case value >= criticalAt:
		return colourRed
	case value >= warnAt:
		return colourAmber
	default:
		return colourTextNormal
	}
}

func oilPressureColor(value float64, engineFailed bool) color.Color {
	if engineFailed || value <= 0 {
		return colourRed
	}
	if value < 140 {
		return colourAmber
	}
	return colourTextNormal
}

func batteryColor(value float64) color.Color {
	if value > 0 && value < 12.2 {
		return colourAmber
	}
	return colourTextNormal
}

func warningColor(warning float64, failed bool) color.Color {
	if failed || warning >= 2 {
		return colourRed
	}
	if warning >= 1 {
		return colourAmber
	}
	return colourGreen
}

func boolColor(value bool) color.Color {
	if value {
		return colourRed
	}
	return colourTextDim
}

func alertTextColor(severity int) color.Color {
	if severity >= 2 {
		return color.White
	}
	if severity == 1 {
		return colourAmber
	}
	return colourTextNormal
}

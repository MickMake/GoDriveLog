package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"GoDriveLog/internal/config"
	"GoDriveLog/internal/sensors"
)

type Dashboard struct {
	root   *fyne.Container
	panels map[string]*panel
}

type panel struct {
	cfg        config.SensorConfig
	label      *widget.Label
	value      *widget.Label
	bar        *widget.ProgressBar
	errorLabel *widget.Label
	history    []float64
	lastUpdate time.Time
}

func NewDashboard(sensors []config.SensorConfig) *Dashboard {
	root := container.NewWithoutLayout()
	d := &Dashboard{root: root, panels: map[string]*panel{}}

	for _, sc := range sensors {
		p := newPanel(sc)
		box := container.NewVBox(p.label, p.value, p.bar, p.errorLabel)
		box.Move(fyne.NewPos(sc.Display.X, sc.Display.Y))
		box.Resize(fyne.NewSize(sc.Display.Width, sc.Display.Height))
		root.Add(box)
		d.panels[sc.PID] = p
	}

	return d
}

func (d *Dashboard) CanvasObject() fyne.CanvasObject { return d.root }

func (d *Dashboard) Update(r sensors.Reading) {
	p := d.panels[r.PID]
	if p == nil {
		return
	}

	p.lastUpdate = time.Now()
	p.errorLabel.SetText("ok")
	p.value.SetText(fmt.Sprintf("%.1f %s", r.Value, r.Unit))
	norm := (r.Value - p.cfg.Min) / (p.cfg.Max - p.cfg.Min)
	if norm < 0 {
		norm = 0
	}
	if norm > 1 {
		norm = 1
	}
	p.bar.SetValue(norm)

	if strings.EqualFold(p.cfg.Style, "graph") {
		p.history = append(p.history, r.Value)
		if len(p.history) > 24 {
			p.history = p.history[len(p.history)-24:]
		}
		p.value.SetText(fmt.Sprintf("%.1f %s  %s", r.Value, r.Unit, spark(p.history, p.cfg.Min, p.cfg.Max)))
	}
}

func (d *Dashboard) SetError(pid string, err error) {
	p := d.panels[pid]
	if p == nil || err == nil {
		return
	}

	p.errorLabel.SetText("error: " + err.Error())
	if p.lastUpdate.IsZero() {
		p.value.SetText("--")
		p.bar.SetValue(0)
		return
	}

	age := time.Since(p.lastUpdate).Round(time.Second)
	p.errorLabel.SetText(fmt.Sprintf("stale %s: %v", age, err))
}

func newPanel(sc config.SensorConfig) *panel {
	label := widget.NewLabel(sc.Name + " [" + sc.PID + "]")
	value := widget.NewLabel("--")
	bar := widget.NewProgressBar()
	errorLabel := widget.NewLabel("waiting")
	bar.Min = 0
	bar.Max = 1
	return &panel{cfg: sc, label: label, value: value, bar: bar, errorLabel: errorLabel}
}

func spark(values []float64, min float64, max float64) string {
	if len(values) == 0 || max <= min {
		return ""
	}
	blocks := []rune("▁▂▃▄▅▆▇█")
	var b strings.Builder
	for _, v := range values {
		n := (v - min) / (max - min)
		if n < 0 {
			n = 0
		}
		if n > 1 {
			n = 1
		}
		idx := int(n * float64(len(blocks)-1))
		b.WriteRune(blocks[idx])
	}
	return b.String()
}

package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"pid-fyne-logger/internal/config"
	"pid-fyne-logger/internal/sensors"
)

type Dashboard struct {
	root   *fyne.Container
	panels map[string]*panel
}

type panel struct {
	cfg     config.SensorConfig
	label   *widget.Label
	value   *widget.Label
	bar     *widget.ProgressBar
	history []float64
}

func NewDashboard(sensors []config.SensorConfig) *Dashboard {
	root := container.NewWithoutLayout()
	d := &Dashboard{root: root, panels: map[string]*panel{}}

	for _, sc := range sensors {
		p := newPanel(sc)
		box := container.NewVBox(p.label, p.value, p.bar)
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

func newPanel(sc config.SensorConfig) *panel {
	label := widget.NewLabel(sc.Name + " [" + sc.PID + "]")
	value := widget.NewLabel("--")
	bar := widget.NewProgressBar()
	bar.Min = 0
	bar.Max = 1
	return &panel{cfg: sc, label: label, value: value, bar: bar}
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

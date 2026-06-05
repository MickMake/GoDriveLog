package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

type Dashboard struct {
	root   *fyne.Container
	panels map[string]*panel
}

type panel struct {
	key        string
	cfg        config.PIDConfig
	label      *widget.Label
	value      *widget.Label
	bar        *widget.ProgressBar
	errorLabel *widget.Label
	history    []float64
	lastUpdate time.Time
}

type panelBox struct {
	key string
	pid config.PIDConfig
	p   *panel
	box *fyne.Container
}

func NewDashboard(pids map[string]config.PIDConfig) *Dashboard {
	root := container.NewWithoutLayout()
	d := &Dashboard{root: root, panels: map[string]*panel{}}

	boxes := make([]panelBox, 0, len(pids))
	for key, pid := range pids {
		if !pid.Display.Enabled {
			continue
		}
		p := newPanel(key, pid)
		box := container.NewVBox(p.label, p.value, p.bar, p.errorLabel)
		box.Move(fyne.NewPos(pid.Display.Position.X, pid.Display.Position.Y))
		box.Resize(fyne.NewSize(pid.Display.Position.Width, pid.Display.Position.Height))
		boxes = append(boxes, panelBox{key: key, pid: pid, p: p, box: box})
		d.panels[key] = p
	}

	// Deterministic z-layering for overlays.
	sort.SliceStable(boxes, func(i, j int) bool {
		zi := boxes[i].pid.Display.Position.Z
		zj := boxes[j].pid.Display.Position.Z
		if zi == zj {
			return boxes[i].key < boxes[j].key
		}
		return zi < zj
	})
	for _, b := range boxes {
		root.Add(b.box)
	}

	return d
}

func (d *Dashboard) CanvasObject() fyne.CanvasObject { return d.root }

func (d *Dashboard) Update(r sensors.Reading) {
	p := d.panelForReading(r)
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

	if strings.EqualFold(p.cfg.Display.Style, "graph") {
		p.history = append(p.history, r.Value)
		if len(p.history) > 24 {
			p.history = p.history[len(p.history)-24:]
		}
		p.value.SetText(fmt.Sprintf("%.1f %s  %s", r.Value, r.Unit, spark(p.history, p.cfg.Min, p.cfg.Max)))
	}
}

func (d *Dashboard) SetError(key string, err error) {
	p := d.panels[key]
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

func (d *Dashboard) panelForReading(r sensors.Reading) *panel {
	if r.SensorKey != "" {
		return d.panels[r.SensorKey]
	}
	return d.panels[r.Name]
}

func newPanel(key string, pid config.PIDConfig) *panel {
	label := widget.NewLabel(key + " [" + pid.PID + "]")
	value := widget.NewLabel("--")
	bar := widget.NewProgressBar()
	errorLabel := widget.NewLabel("waiting")
	bar.Min = 0
	bar.Max = 1
	return &panel{key: key, cfg: pid, label: label, value: value, bar: bar, errorLabel: errorLabel}
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

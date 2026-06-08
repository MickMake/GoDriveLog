package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

type Dashboard struct {
	root  *fyne.Container
	store *sensors.StateStore
}

func NewDashboard(cfg config.DashboardConfig, store *sensors.StateStore) *Dashboard {
	label := widget.NewLabel(fmt.Sprintf("dashboard v2 placeholder (%dx%d)", cfg.Canvas.Width, cfg.Canvas.Height))
	root := container.NewCenter(label)
	return &Dashboard{root: root, store: store}
}

func (d *Dashboard) CanvasObject() fyne.CanvasObject { return d.root }

func (d *Dashboard) StateSnapshot() []sensors.SensorState {
	if d.store == nil {
		return nil
	}
	return d.store.Snapshot()
}

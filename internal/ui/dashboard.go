package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/internal/config"
)

type Dashboard struct {
	root *fyne.Container
}

func NewDashboard(cfg config.DashboardConfig) *Dashboard {
	label := widget.NewLabel(fmt.Sprintf("dashboard v2 placeholder (%dx%d)", cfg.Canvas.Width, cfg.Canvas.Height))
	root := container.NewCenter(label)
	return &Dashboard{root: root}
}

func (d *Dashboard) CanvasObject() fyne.CanvasObject { return d.root }

package ui

import (
	"context"
	"sync"
	"time"

	"fyne.io/fyne/v2"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/assets"
	"github.com/MickMake/GoDriveLog/internal/dashboard/decoders"
	fynerenderer "github.com/MickMake/GoDriveLog/internal/dashboard/renderer/fyne"
	"github.com/MickMake/GoDriveLog/internal/dashboard/scene"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

type Dashboard struct {
	renderer *fynerenderer.Renderer
	store    *sensors.StateStore
	cfg      config.DashboardConfig
	assets   *assets.Registry

	mu      sync.RWMutex
	lastErr error
}

func NewDashboard(cfg config.DashboardConfig, store *sensors.StateStore) *Dashboard {
	dashboard, err := NewDashboardWithConfigPath(cfg, "", store)
	if err != nil {
		dashboard = &Dashboard{
			renderer: fynerenderer.New(nil),
			store:    store,
			cfg:      cfg,
		}
		dashboard.setLastError(err)
	}
	dashboard.Start(context.Background(), 100*time.Millisecond)
	return dashboard
}

func NewDashboardWithConfigPath(cfg config.DashboardConfig, configPath string, store *sensors.StateStore) (*Dashboard, error) {
	assetRegistry, err := assets.Load(cfg, configPath)
	if err != nil {
		return nil, err
	}

	dashboard := &Dashboard{
		renderer: fynerenderer.New(assetRegistry),
		store:    store,
		cfg:      cfg,
		assets:   assetRegistry,
	}
	if err := dashboard.Refresh(); err != nil {
		dashboard.setLastError(err)
	}
	return dashboard, nil
}

func (d *Dashboard) CanvasObject() fyne.CanvasObject {
	return d.renderer.CanvasObject()
}

func (d *Dashboard) Start(ctx context.Context, interval time.Duration) {
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
				fyne.Do(func() {
					if err := d.Refresh(); err != nil {
						d.setLastError(err)
					}
				})
			}
		}
	}()
}

func (d *Dashboard) Refresh() error {
	sensorStates := d.sensorStateMap()
	decoderValues, err := decoders.Execute(d.cfg.Decoders, decoders.Inputs{Sensors: sensorStates})
	if err != nil {
		return err
	}

	sceneState, err := scene.Evaluate(d.cfg, d.assets, decoderValues, sensorStates, scene.Options{})
	if err != nil {
		return err
	}
	if err := d.renderer.Update(sceneState); err != nil {
		return err
	}

	d.setLastError(nil)
	return nil
}

func (d *Dashboard) StateSnapshot() []sensors.SensorState {
	if d.store == nil {
		return nil
	}
	return d.store.Snapshot()
}

func (d *Dashboard) LastError() error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.lastErr
}

func (d *Dashboard) setLastError(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lastErr = err
}

func (d *Dashboard) sensorStateMap() map[string]sensors.SensorState {
	states := map[string]sensors.SensorState{}
	for _, state := range d.StateSnapshot() {
		states[state.ID] = state
	}
	return states
}

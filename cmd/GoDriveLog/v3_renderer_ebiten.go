//go:build !fyne_legacy

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	v3ebitenadapter "github.com/MickMake/GoDriveLog/internal/dashboard/adapter/ebiten"
	v3harness "github.com/MickMake/GoDriveLog/internal/dashboard/harness"
	"github.com/MickMake/GoDriveLog/internal/dashboard/scenesink"
	v3runtime "github.com/MickMake/GoDriveLog/internal/runtime/v3runtime"
)

const v3SceneGap = 12

type ebitenWindowSize struct {
	Width  int
	Height int
}

func runV3EbitenCommand(configPath, vehicleID string, duration time.Duration) error {
	ctx, stop := newV3Context(duration)
	defer stop()

	initialSize, err := initialV3EbitenWindowSize(configPath, vehicleID)
	if err != nil {
		return err
	}
	adapter, err := v3ebitenadapter.New(".", initialSize.Width, initialSize.Height)
	if err != nil {
		return err
	}
	displaySink, err := newDirectSceneSink(adapter.UpdateScenes, "v3 dashboard Ebiten adapter update")
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		summary, err := v3runtime.Run(ctx, v3runtime.Options{
			ConfigPath:    configPath,
			VehicleID:     vehicleID,
			Logger:        log.Default(),
			DashboardSink: displaySink.SubmitLatest,
		})
		if closeErr := displaySink.Close(); err == nil {
			err = closeErr
		}
		stats := displaySink.Stats()
		if err == nil || isContextStop(err) {
			log.Printf("v3 runtime summary: vehicle=%s endpoint=%s sensors=%d logs=%d dashboards=%d renderer=%s display_submitted=%d display_rendered=%d display_superseded=%d display_last_render=%s", summary.VehicleID, summary.Endpoint, summary.SensorCount, summary.LogCount, summary.DashboardCount, v3RendererEbiten, stats.Submitted, stats.Rendered, stats.Superseded, stats.LastRenderDuration)
		} else {
			log.Printf("v3 runtime stopped with error: %v", err)
		}
		errCh <- err
		stop()
	}()

	runErr := adapter.Run(ctx, "GoDriveLog v3")
	stop()
	runtimeErr := <-errCh
	if err := ignoreContextStop(runErr); err != nil {
		return err
	}
	return ignoreContextStop(runtimeErr)
}

func runV3EbitenHarnessCommand(configPath, vehicleID, pattern string, interval time.Duration, duration time.Duration) error {
	ctx, stop := newV3Context(duration)
	defer stop()

	initialSize, err := initialV3EbitenWindowSize(configPath, vehicleID)
	if err != nil {
		return err
	}
	adapter, err := v3ebitenadapter.New(".", initialSize.Width, initialSize.Height)
	if err != nil {
		return err
	}
	displaySink, err := newDirectSceneSink(adapter.UpdateScenes, "v3 dashboard harness Ebiten adapter update")
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		summary, err := v3harness.Run(ctx, v3harness.Options{
			ConfigPath: configPath,
			VehicleID:  vehicleID,
			Pattern:    pattern,
			Interval:   interval,
			Logger:     log.Default(),
			Sink:       displaySink.SubmitLatest,
		})
		if closeErr := displaySink.Close(); err == nil {
			err = closeErr
		}
		stats := displaySink.Stats()
		if err == nil || isContextStop(err) {
			log.Printf("v3 dashboard harness summary: vehicle=%s sensors=%d dashboards=%d pattern=%s interval=%s renderer=%s events=%d display_submitted=%d display_rendered=%d display_superseded=%d display_last_render=%s", summary.VehicleID, summary.SensorCount, summary.DashboardCount, summary.Pattern, summary.Interval, v3RendererEbiten, summary.Events, stats.Submitted, stats.Rendered, stats.Superseded, stats.LastRenderDuration)
		} else {
			log.Printf("v3 dashboard harness stopped with error: %v", err)
		}
		errCh <- err
		stop()
	}()

	runErr := adapter.Run(ctx, "GoDriveLog v3 harness")
	stop()
	runtimeErr := <-errCh
	if err := ignoreContextStop(runErr); err != nil {
		return err
	}
	return ignoreContextStop(runtimeErr)
}

func newDirectSceneSink(update scenesink.Sink, label string) (*scenesink.LatestSink, error) {
	return scenesink.NewLatestSink(func(scenes []v3runtime.Scene) error {
		if err := update(scenes); err != nil {
			return err
		}
		log.Printf("%s: scenes=%d", label, len(scenes))
		return nil
	})
}

func initialV3EbitenWindowSize(configPath, vehicleID string) (ebitenWindowSize, error) {
	cfg, err := v3config.LoadFile(configPath)
	if err != nil {
		return ebitenWindowSize{}, fmt.Errorf("load v3 config for initial window size: %w", err)
	}
	plan, err := v3config.Resolve(cfg, vehicleID)
	if err != nil {
		return ebitenWindowSize{}, fmt.Errorf("resolve v3 runtime plan for initial window size: %w", err)
	}
	return selectedEbitenDashboardsSize(plan.Dashboards), nil
}

func selectedEbitenDashboardsSize(dashboards []v3config.ResolvedDashboard) ebitenWindowSize {
	var width int
	var height int
	for index, dashboard := range dashboards {
		if dashboard.Config.Size.Width > width {
			width = dashboard.Config.Size.Width
		}
		height += dashboard.Config.Size.Height
		if index < len(dashboards)-1 {
			height += v3SceneGap
		}
	}
	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 480
	}
	return ebitenWindowSize{Width: width, Height: height}
}

func cancelOnContextDone(ctx context.Context, stop func()) {
	<-ctx.Done()
	stop()
}

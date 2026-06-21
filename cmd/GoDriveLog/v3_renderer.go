package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"
	"time"

	v3ebitenadapter "github.com/MickMake/GoDriveLog/internal/dashboard/adapter/ebiten"
	v3harness "github.com/MickMake/GoDriveLog/internal/dashboard/harness"
	"github.com/MickMake/GoDriveLog/internal/dashboard/scenesink"
	v3runtime "github.com/MickMake/GoDriveLog/internal/runtime/v3runtime"
)

const (
	v3RendererFyne   = "fyne"
	v3RendererEbiten = "ebiten"
)

var (
	selectedV3Renderer = v3RendererFyne
	selectedV3Duration time.Duration
)

func normalizeV3Renderer(renderer string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(renderer)) {
	case "", v3RendererFyne:
		return v3RendererFyne, nil
	case v3RendererEbiten:
		return v3RendererEbiten, nil
	default:
		return "", fmt.Errorf("unsupported --renderer %q; expected fyne or ebiten", renderer)
	}
}

func newV3Context(duration time.Duration) (context.Context, func()) {
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	if duration <= 0 {
		return signalCtx, stopSignals
	}
	ctx, cancel := context.WithTimeout(signalCtx, duration)
	return ctx, func() {
		cancel()
		stopSignals()
	}
}

func runV3EbitenCommand(configPath, vehicleID string, duration time.Duration) error {
	ctx, stop := newV3Context(duration)
	defer stop()

	initialSize, err := initialV3WindowSize(configPath, vehicleID)
	if err != nil {
		return err
	}
	adapter, err := v3ebitenadapter.New(".", int(initialSize.Width), int(initialSize.Height))
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

	runErr := adapter.Run(ctx, "GoDriveLog v3 (Ebiten experimental)")
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

	initialSize, err := initialV3WindowSize(configPath, vehicleID)
	if err != nil {
		return err
	}
	adapter, err := v3ebitenadapter.New(".", int(initialSize.Width), int(initialSize.Height))
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

	runErr := adapter.Run(ctx, "GoDriveLog v3 harness (Ebiten experimental)")
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

func ignoreContextStop(err error) error {
	if isContextStop(err) {
		return nil
	}
	return err
}

func isContextStop(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

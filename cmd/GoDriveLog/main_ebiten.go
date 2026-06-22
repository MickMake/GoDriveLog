//go:build !fyne_legacy

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	v3harness "github.com/MickMake/GoDriveLog/internal/dashboard/harness"
)

func main() {
	configPath := flag.String("config", "config.example.yaml", "path to YAML config")
	useV3 := flag.Bool("v3", true, "run the v3 selected-vehicle runtime path")
	useHarness := flag.Bool("harness", false, "run the v3 dashboard harness without OBD; requires --v3")
	vehicleID := flag.String("vehicle", "", "v3 vehicle id; required when the v3 config contains multiple vehicles")
	harnessPattern := flag.String("pattern", v3harness.PatternSweep, "v3 dashboard harness pattern: sweep, heartbeat, or fixed")
	harnessInterval := flag.Duration("interval", 100*time.Millisecond, "v3 dashboard harness update interval, such as 50ms or 100ms")
	duration := flag.Duration("duration", 0, "optional v3 runtime or harness duration, such as 60s; zero runs until interrupted")
	renderer := flag.String("renderer", v3RendererEbiten, "v3 renderer backend: ebiten")
	flag.Parse()

	normalizedRenderer, err := normalizeV3Renderer(*renderer)
	if err != nil {
		log.Fatal(err)
	}
	if normalizedRenderer != v3RendererEbiten {
		log.Fatal("Fyne is no longer supported in the v3.3 active dashboard path; use a v3.2.x tag or branch for the last supported Fyne dashboard")
	}
	if !*useV3 {
		log.Fatal("the Ebiten command supports only the v3 dashboard path")
	}
	if *useHarness && !*useV3 {
		log.Fatal("--harness requires --v3")
	}

	if *useHarness {
		if err := runV3EbitenHarnessCommand(*configPath, *vehicleID, *harnessPattern, *harnessInterval, *duration); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := runV3EbitenCommand(*configPath, *vehicleID, *duration); err != nil {
		log.Fatal(fmt.Errorf("run v3 ebiten: %w", err))
	}
}

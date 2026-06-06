package widgets

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

// Show runs the isolated widget demo harness for commands like:
//
//	GoDriveLog widget radial1
//
// If stdin is piped, numeric values supplied on stdin update the selected widget state one line at a time.
// If stdin is a terminal, a simple demo waveform drives the widget until the window is closed.
func Show(args []string, stdout io.Writer, stdin io.Reader) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		printShowUsage(stdout)
		return nil
	}

	widgetName := strings.TrimSpace(args[0])
	w, err := New(widgetName, demoConfig(widgetName))
	if err != nil {
		printShowUsage(stdout)
		return err
	}

	obj, ok := w.(fyne.CanvasObject)
	if !ok {
		return fmt.Errorf("widget %q is not renderable yet", widgetName)
	}

	fmt.Fprintf(stdout, "GoDriveLog widget demo: %s\n", w.Style())
	fmt.Fprintln(stdout, "Tip: pipe numbers into stdin to drive updates, e.g. `seq 0 100 | ./GoDriveLog widget radial1`.")

	myApp := app.New()
	win := myApp.NewWindow("GoDriveLog widget: " + widgetName)
	win.Resize(fyne.NewSize(1920, 480))
	win.SetContent(container.NewCenter(obj))

	// Drive updates.
	if isPiped(stdin) {
		updated := make(chan bool, 1)
		go func() {
			ok, err := showFromStdin(w, stdout, stdin)
			if err != nil {
				fmt.Fprintln(stdout, err)
				myApp.Quit()
				return
			}
			updated <- ok
		}()
		go func() {
			select {
			case ok := <-updated:
				if !ok {
					// No piped input arrived - run a demo waveform so the window isn't dead-on-arrival.
					demoWaveform(myApp, w)
				}
			case <-time.After(200 * time.Millisecond):
				// If stdin is a pipe but slow to start, wait for the reader goroutine.
			}
		}()
	} else {
		demoWaveform(myApp, w)
	}

	win.SetCloseIntercept(func() {
		myApp.Quit()
		win.Close()
	})

	win.ShowAndRun()
	return nil
}

func printShowUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  GoDriveLog widget <widget-name>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Available widget names:")
	for _, style := range Styles() {
		fmt.Fprintf(w, "  %s\n", style)
	}
}

func isPiped(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice == 0
}

func showFromStdin(w Widget, stdout io.Writer, stdin io.Reader) (bool, error) {
	if stdin == nil {
		return false, nil
	}

	updated := false
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return updated, fmt.Errorf("invalid widget value %q", line)
		}
		fyne.Do(func() { w.SetValue(value) })
		printSnapshot(stdout, w.Snapshot())
		updated = true
	}
	return updated, scanner.Err()
}

func demoWaveform(app fyne.App, w Widget) {
	cfg := w.Config()
	rangeSpan := cfg.Max - cfg.Min
	if rangeSpan == 0 {
		rangeSpan = 100
	}

	start := time.Now()
	ticker := time.NewTicker(40 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			t := time.Since(start).Seconds()
			// 0..1 sine wave
			pct := (math.Sin(t*1.2) + 1) / 2
			value := cfg.Min + pct*rangeSpan
			fyne.Do(func() { w.SetValue(value) })
		}
	}()
}

func printSnapshot(w io.Writer, snap Snapshot) {
	state := "normal"
	if snap.Warning {
		state = "warning"
	}
	if snap.Danger {
		state = "danger"
	}
	fmt.Fprintf(w, "%s %-10s %9.2f %-6s %3.0f%% %s\n", snap.Style, snap.Label, snap.Value, snap.Unit, snap.Normalised*100, state)
}

func demoConfig(widgetName string) GaugeConfig {
	cfg := DefaultGaugeConfig()

	switch strings.ToLower(widgetName) {
	case "radial1", "radial2", "radial3":
		cfg.Label = "RPM"
		cfg.Unit = "rpm"
		cfg.Min = 0
		cfg.Max = 5000
		cfg.WarningRange = &Range{Min: 4000, Max: 4499.999}
		cfg.DangerRange = &Range{Min: 4500, Max: 5000}
		cfg.SmoothingWindow = 1

	case "sweep1", "sweep2", "sweep3":
		cfg.Label = "RPM"
		cfg.Unit = "rpm"
		cfg.Min = 0
		cfg.Max = 7000
		cfg.WarningRange = &Range{Min: 5500, Max: 6299.999}
		cfg.DangerRange = &Range{Min: 6300, Max: 7000}
		cfg.SmoothingWindow = 1

	case "speedhud1", "speedhud2", "speedhud3":
		cfg.Label = "SPEED"
		cfg.Unit = "km/h"
		cfg.Min = 0
		cfg.Max = 160
		cfg.WarningRange = &Range{Min: 110, Max: 129.999}
		cfg.DangerRange = &Range{Min: 130, Max: 160}
		cfg.SmoothingWindow = 1

	case "bar2", "bar3":
		cfg.Label = "LOAD"
		cfg.Unit = "%"
		cfg.Min = 0
		cfg.Max = 100
		cfg.WarningRange = &Range{Min: 75, Max: 89.999}
		cfg.DangerRange = &Range{Min: 90, Max: 100}
		cfg.SmoothingWindow = 1

	default:
		cfg.Label = strings.ToUpper(widgetName)
		cfg.Unit = "%"
		cfg.WarningRange = &Range{Min: 75, Max: 89.999}
		cfg.DangerRange = &Range{Min: 90, Max: 100}
	}

	return cfg
}

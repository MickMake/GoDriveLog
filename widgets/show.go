package widgets

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Show runs the isolated widget demo harness for commands like:
//
//	godrivelog widgets radial1
//
// Numeric values supplied on stdin update the selected widget state one line at a time.
func Show(args []string, stdout io.Writer, stdin io.Reader) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		printShowUsage(stdout)
		return nil
	}

	widgetName := strings.TrimSpace(args[0])
	widget, err := New(widgetName, demoConfig(widgetName))
	if err != nil {
		printShowUsage(stdout)
		return err
	}

	fmt.Fprintf(stdout, "GoDriveLog widget demo: %s\n", widget.Style())
	fmt.Fprintln(stdout, "Renderer: stub/state harness. Fyne drawing will attach behind this interface.")
	fmt.Fprintln(stdout, "Enter numeric values on stdin to update the widget; empty stdin prints demo samples.")

	updated, err := showFromStdin(widget, stdout, stdin)
	if err != nil {
		return err
	}
	if updated {
		return nil
	}

	for _, value := range []float64{0, 25, 50, 75, 90, 100} {
		widget.SetValue(value)
		printSnapshot(stdout, widget.Snapshot())
	}
	return nil
}

func printShowUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  godrivelog widgets <widget-name>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Available widget names:")
	for _, style := range Styles() {
		fmt.Fprintf(w, "  %s\n", style)
	}
}

func showFromStdin(widget Widget, stdout io.Writer, stdin io.Reader) (bool, error) {
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
		widget.SetValue(value)
		printSnapshot(stdout, widget.Snapshot())
		updated = true
	}
	return updated, scanner.Err()
}

func printSnapshot(w io.Writer, snap Snapshot) {
	state := "normal"
	if snap.Warning {
		state = "warning"
	}
	if snap.Danger {
		state = "danger"
	}
	fmt.Fprintf(w, "%s %-8s %7.2f %-6s %3.0f%% %s\n", snap.Style, snap.Label, snap.Value, snap.Unit, snap.Normalised*100, state)
}

func demoConfig(widgetName string) GaugeConfig {
	cfg := DefaultGaugeConfig()
	cfg.Label = strings.ToUpper(widgetName)
	cfg.Unit = "%"
	cfg.WarningRange = &Range{Min: 75, Max: 89.999}
	cfg.DangerRange = &Range{Min: 90, Max: 100}
	return cfg
}

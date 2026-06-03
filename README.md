# PID Fyne Logger v0.1

A deliberately small Go/Fyne PID dashboard for Raspberry Pi 4.

It starts, reads a JSON config, polls configured PIDs at their own refresh intervals, writes JSONL logs, rotates the log on engine start, and displays values in a Fyne window.

## What is included

- Go app using Fyne v2.
- JSON startup config.
- Per-sensor PID, name, refresh rate, display style, position, and size.
- Mock PID reader so the UI/logging can be tested without OBD hardware.
- JSON Lines logging.
- Log rotation when the configured engine-start PID crosses the configured threshold.

## Pi 4 install notes

Fyne uses Go modules and the official quick start uses:

```bash
go get fyne.io/fyne/v2@latest
go install fyne.io/tools/cmd/fyne@latest
```

On Raspberry Pi OS you will also need normal desktop/OpenGL build dependencies for Fyne. If the build complains about missing GL/X11 headers, install the Raspberry Pi OS equivalents for gcc, pkg-config, libgl, x11, xcursor, xrandr, xinerama, xi, and xxf86vm development packages.

## Build

```bash
cd pid-fyne-logger_v0.1
go mod tidy
go build ./cmd/pid-fyne-logger
```

## Run in mock mode

```bash
./pid-fyne-logger -config config.example.json
```

The mock engine sleeps for about three seconds, then RPM rises. That should trigger an `engine-start` log rotation.

## Log format

Logs are JSON Lines, one reading per line:

```json
{"time":"2026-06-03T10:15:30Z","pid":"010C","name":"RPM","value":1234.5,"unit":"rpm","source":"mock"}
```

## Config shape

```json
{
  "log_dir": "./logs",
  "engine_start_pid": "010C",
  "engine_start_threshold": 50,
  "mock_mode": true,
  "sensors": [
    {
      "pid": "010C",
      "name": "RPM",
      "refresh_ms": 250,
      "style": "gauge",
      "min": 0,
      "max": 7000,
      "display": { "x": 20, "y": 20, "width": 360, "height": 90 }
    }
  ]
}
```

## Real OBD transport

`internal/sensors/reader.go` includes an `ELM327Reader` placeholder. Keep the rest of the app unchanged and implement that `Read(ctx, pid)` method using the preferred transport: USB serial, Bluetooth serial, or TCP adapter.

## Notes

This is intentionally simple. No database, no plugin framework, no ceremony, no dashboard architecture astronautics. The app has one job and a small lunchbox.

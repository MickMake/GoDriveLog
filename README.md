# PID Fyne Logger v0.1

A deliberately small Go/Fyne PID dashboard for Raspberry Pi 4.

It starts, reads a JSON config, polls configured PIDs at their own refresh intervals, writes JSONL logs, rotates the log on engine start, and displays values in a Fyne window.

## What is included

- Go app using Fyne v2.
- JSON startup config.
- Per-sensor PID, name, refresh rate, display style, position, and size.
- App-level mock PID reader so the UI/logging can be tested without OBD hardware.
- Real OBD reader adapter using `github.com/rzetterberg/elmobd`.
- JSON Lines logging.
- Log rotation when the configured engine-start PID crosses the configured threshold.

## Pi 4 install notes

Fyne uses Go modules and the official quick start uses:

```bash
go get fyne.io/fyne/v2@latest
go install fyne.io/tools/cmd/fyne@latest
```

On Raspberry Pi OS you will also need normal desktop/OpenGL build dependencies for Fyne. If the build complains about missing GL/X11 headers, install the Raspberry Pi OS equivalents for gcc, pkg-config, libgl, x11, xcursor, xrandr, xinerama, xi, and xxf86vm development packages.

For a USB ELM327 adapter, the default address is:

```text
serial:///dev/ttyUSB0
```

`elmobd` also supports address schemes such as `tcp://host:port` and `test:///dev/ttyUSB0`.

## Build

From the repository root:

```bash
go mod tidy
go build ./cmd/pid-fyne-logger
```

The binary will be written to the current directory as `pid-fyne-logger` unless you pass `-o`.

## Run in mock mode

```bash
./pid-fyne-logger -config config.example.json
```

The mock engine sleeps for about three seconds, then RPM rises. That should trigger an `engine-start` log rotation.

## Run with a real ELM327 adapter

Set `mock_mode` to `false` and point `obd_address` at the adapter:

```json
{
  "mock_mode": false,
  "obd_address": "serial:///dev/ttyUSB0",
  "obd_debug": false
}
```

The current real OBD adapter supports these configured PIDs:

- `0105` coolant temperature, Celsius
- `010C` engine RPM
- `010D` vehicle speed, km/h

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
  "obd_address": "serial:///dev/ttyUSB0",
  "obd_debug": false,
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

`internal/sensors/elmobd_reader.go` adapts `github.com/rzetterberg/elmobd` to the app's small `Reader` interface. Add new supported PIDs there as needed.

## Notes

This is intentionally simple. No database, no plugin framework, no ceremony, no dashboard architecture astronautics. The app has one job and a small lunchbox.

# GoDriveLog v1.x

A deliberately small Go/Fyne PID dashboard for Raspberry Pi 4.

It starts, reads a YAML config, polls configured PIDs at their own refresh intervals, writes JSONL logs, rotates the log daily, and displays values in a Fyne window.

## What is included

- Go app using Fyne v2.
- YAML startup config.
- Per-sensor PID, name, refresh rate, display widget + style, position, and size.
- App-level mock PID reader so the UI/logging can be tested without OBD hardware.
- Real OBD reader adapter using `github.com/rzetterberg/elmobd`.
- JSON Lines logging.
- Log rotation will happen on a daily basis.
- Per-sensor error/stale display so failed reads are visible on screen.

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
go build ./cmd/GoDriveLog
```

The binary will be written to the current directory as `GoDriveLog` unless you pass `-o`.

## Run in mock mode

```bash
./GoDriveLog -config config.example.yaml
```

The mock engine sleeps for about three seconds, then RPM rises.

## Test the elmobd backend without hardware

Use the elmobd test address to exercise the real reader adapter without opening a serial device:

```yaml
mock_mode: true
obd_address: serial:///dev/ttyUSB0
obd_debug: false
```

You can temporarily set those fields in a copy of `config.example.yaml`.

## Run with a real ELM327 adapter

Use the real OBD example config:

```bash
./GoDriveLog -config config.example.yaml
```

Or set `mock_mode` to `false` and point `obd_address` at the adapter:

```yaml
mock_mode: false
obd_address: serial:///dev/ttyUSB0
obd_debug: false
```

The current real OBD adapter supports these configured PIDs:

| Key | PID | Unit | Meaning |
|---|---:|---|---|
| `engine_load` | `0104` | `%` | Calculated engine load |
| `coolant_temp` | `0105` | `C` | Engine coolant temperature |
| `short_fuel_trim_bank1` | `0106` | `%` | Short term fuel trim, bank 1 |
| `long_fuel_trim_bank1` | `0107` | `%` | Long term fuel trim, bank 1 |
| `intake_manifold_pressure` | `010B` | `kPa` | Intake manifold absolute pressure |
| `rpm` | `010C` | `rpm` | Engine RPM |
| `speed` | `010D` | `km/h` | Vehicle speed |
| `intake_air_temp` | `010F` | `C` | Intake air temperature |
| `throttle_position` | `0111` | `%` | Throttle position |
| `fuel_level` | `012F` | `%` | Fuel tank level |
| `control_module_voltage` | `0142` | `V` | Control module voltage |
| `engine_oil_temp` | `015C` | `C` | Engine oil temperature |


## Log output format

Log output are JSON Lines, one reading per line:

```json
{"time":"2026-06-03T10:15:30Z","pid":"010C","name":"RPM","value":1234.5,"unit":"rpm","source":"mock"}
```

## Config shape

```yaml
mock_mode: true
obd_address: serial:///dev/ttyUSB0
obd_debug: false

log:
  rotate: daily
  directory: ./log

vehicle:
  name: "VW Caddy"

  pids:
    engine_load:
      type: obd
      pid: "0104"
      unit: "%"
      refresh: 500
      min: 0
      max: 100
      log: true
      display:
        enabled: false

    coolant_temp:
      type: obd
      pid: "0105"
      unit: C
      refresh: 1000
      min: -40
      max: 140
      log: true
      display:
        enabled: true
        widget: graph1
        position:
          x: 20
          y: 240
          width: 360
          height: 120
```

## Real OBD transport

`internal/sensors/elmobd_reader.go` adapts `github.com/rzetterberg/elmobd` to the app's small `Reader` interface. Add new supported PIDs there as needed.

## Notes

This is intentionally simple. No database, no plugin framework, no ceremony, no dashboard architecture astronautics. The app has one job and a small lunchbox.

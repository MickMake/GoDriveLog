# GoDriveLog

A deliberately small Go/Fyne sensor dashboard for Raspberry Pi 4.

It starts, reads a YAML config, polls configured sensors at their own refresh intervals, writes JSONL logs, rotates the log daily, and renders a configured dashboard scene in a Fyne window.

## What is included

- Go app using Fyne v2.
- YAML startup config.
- Dashboard v2 scene config with local assets, decoders, blocks, layers, and conditions.
- App-level mock reader so the UI/logging can be tested without OBD hardware.
- Real OBD reader adapter using `github.com/rzetterberg/elmobd`.
- JSON Lines logging.
- Daily log rotation.
- Sensor status/stale/error state available to dashboard scenes.

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

## Run the dashboard v2 example in mock mode

```bash
./GoDriveLog -config config.example.yaml
```

Or from source:

```bash
go run ./cmd/GoDriveLog -config config.example.yaml
```

`config.example.yaml` runs in mock mode and loads the local SVG fixture assets under `assets/dashboard/bttf`. The example dashboard shows:

- static background
- RPM sprite digits
- throttle sprite-frame bar
- redline glow overlay from a configured threshold condition
- status/stale/error badges from sensor state

The mock reader sleeps for about three seconds, then RPM rises. The redline overlay threshold is intentionally low in the example so the visual condition is easy to see without real hardware. Tiny demo goblin, useful boots.

## Test the reader path without hardware

Mock mode is enabled in the example config:

```yaml
mock_mode: true
obd_address: serial:///dev/ttyUSB0
obd_debug: false
```

You can temporarily set those fields in a copy of `config.example.yaml`.

## Run with a real ELM327 adapter

Set `mock_mode` to `false` and point `obd_address` at the adapter:

```yaml
mock_mode: false
obd_address: serial:///dev/ttyUSB0
obd_debug: false
```

Then run:

```bash
./GoDriveLog -config config.example.yaml
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

Log output is JSON Lines, one reading per line:

```json
{"time":"2026-06-03T10:15:30Z","pid":"010C","name":"RPM","value":1234.5,"unit":"rpm","source":"mock"}
```

## Config shape

The app separates sensor state production from dashboard scene rendering:

```yaml
sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: rpm
    refresh: 250
    min: 0
    max: 7000
    log: true

dashboard:
  canvas:
    width: 800
    height: 480
  asset_root: assets/dashboard/bttf
  assets:
    - id: background
      type: image
      path: background.svg
  decoders:
    - id: rpm_text
      type: format_number
      sensor: rpm
      format: "0000"
  blocks:
    - id: background_panel
      type: image
      asset: background
      geometry:
        x: 0
        y: 0
        width: 800
        height: 480
  layers:
    - id: base
      z: 0
      blocks:
        - background_panel
```

Dashboard block visibility can be driven by configured conditions against sensor status or decoder values.

## Real OBD transport

`internal/sensors/elmobd_reader.go` adapts `github.com/rzetterberg/elmobd` to the app's small `Reader` interface. Add new supported PIDs there as needed.

## Notes

This is intentionally simple. No database, no plugin framework, no ceremony, no dashboard architecture astronautics. The app has one job and a small lunchbox.

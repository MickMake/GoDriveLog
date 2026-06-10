# GoDriveLog Config

This document defines the current GoDriveLog YAML configuration format.

## Runtime display note

The normal app display is the fast fixed 1920x480 instrument dashboard wired from `cmd/GoDriveLog/main.go` through `internal/ui/instrument_dashboard.go`.

The old config-driven scene renderer has been removed from the normal runtime path. The `dashboard` block below remains a small schema/config placeholder so existing config loading and validation have canvas/display metadata, but it is **not** the active renderer configuration.

## Top-level structure

```yaml
obd:
  provider: mock
  mock_mode: true
  address: serial:///dev/ttyUSB0
  debug: false

log:
  rotate: daily
  directory: ./log

vehicle:
  name: "VW Caddy"

dashboard:
  refresh_ms: 250
  render_min_ms: 0
  canvas:
    width: 1920
    height: 480
  assets: []
  decoders: []
  blocks:
    - id: schema_placeholder
      type: text
      geometry:
        x: 0
        y: 0
        width: 1920
        height: 480
  layers:
    - id: base
      z: 0
      blocks:
        - schema_placeholder

sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: rpm
    refresh: 250
    min: 0
    max: 7000
    log: true
```

## Root fields

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `obd` | object | yes | Adapter/mock/provider configuration. |
| `log` | object | yes | Log output configuration. |
| `vehicle` | object | yes | Vehicle metadata. |
| `dashboard` | object | yes | Minimal dashboard schema placeholder and canvas metadata. Not the active renderer. |
| `sensors` | map | yes | Sensor catalogue keyed by stable sensor key. |

## OBD config

```yaml
obd:
  provider: mock
  mock_mode: true
  address: serial:///dev/ttyUSB0
  debug: false
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `provider` | string | no | `obd`, `mock`, or `race-demo`. Defaults from `mock_mode` when omitted. |
| `mock_mode` | bool | no | Compatibility shortcut for mock provider. Defaults to `false`. |
| `address` | string | no | ELM327/elmobd address. Defaults to `serial:///dev/ttyUSB0`. |
| `debug` | bool | no | Enable verbose adapter debugging. Defaults to `false`. |

Legacy top-level `mock_mode`, `obd_address`, and `obd_debug` fields are no longer part of the schema.

## Log config

```yaml
log:
  rotate: daily
  directory: ./log
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `rotate` | string | no | Rotation mode. Currently only `daily`. Defaults to `daily`. |
| `directory` | string | no | Directory for JSONL logs. Defaults to `./log`. |

Daily rotation means GoDriveLog writes readings to a date-based JSONL file and opens a new file when the date changes.

## Vehicle config

```yaml
vehicle:
  name: "VW Caddy"
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `name` | string | yes | Human-readable vehicle name. |

## Dashboard placeholder config

The fast instrument dashboard currently uses fixed 1920x480 Fyne objects and does not read the old scene/asset/decoder renderer configuration at runtime.

Keep the dashboard block minimal:

```yaml
dashboard:
  refresh_ms: 250
  render_min_ms: 0
  canvas:
    width: 1920
    height: 480
  assets: []
  decoders: []
  blocks:
    - id: schema_placeholder
      type: text
      geometry:
        x: 0
        y: 0
        width: 1920
        height: 480
  layers:
    - id: base
      z: 0
      blocks:
        - schema_placeholder
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `refresh_ms` | int | no | Legacy schema tick value. Not used by the fast instrument dashboard runtime. |
| `render_min_ms` | int | no | Legacy schema field. Not used by the fast instrument dashboard runtime. |
| `canvas.width` | int | yes | Dashboard canvas width in pixels. |
| `canvas.height` | int | yes | Dashboard canvas height in pixels. |
| `assets` | list | yes | Keep empty unless a future renderer explicitly reintroduces asset-backed config. |
| `decoders` | list | yes | Keep empty unless a future renderer explicitly reintroduces decoder-backed config. |
| `blocks` | list | yes | Minimal placeholder block for schema validation. |
| `layers` | list | yes | Minimal placeholder layer for schema validation. |

Do not add an old/new display preference. Rollback/reference for the old config-scene dashboard is via the `legacy-config-scene-dashboard` Git ref.

## Sensor config

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
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `type` | string | yes | `obd` for live/mock OBD polling or `virtual` for non-polled derived values. |
| `pid` | string | for `obd` | OBD PID or demo PID key. |
| `unit` | string | yes | Display/logging unit string. |
| `refresh` | int | yes | Polling interval in milliseconds. |
| `min` | number | yes | Expected minimum value. |
| `max` | number | yes | Expected maximum value. |
| `log` | bool | no | Whether to write readings to JSONL logs. |

## Race demo launch

```bash
go run ./cmd/GoDriveLog --config config.example.yaml --sensor-provider race-demo
```

The race demo appends extra display-only demo sensors in `race-demo` provider mode. Live OBD mode uses the configured sensors and is unchanged by the fast dashboard cleanup.

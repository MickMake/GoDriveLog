# GoDriveLog Config

This document defines the GoDriveLog configuration format.

The v2 schema separates sensor configuration from dashboard visual configuration.
Sensors define what can be read and logged. The dashboard defines the visual canvas and, in later stages, how scenes consume sensor state.

## Current scope

Implemented now:

- YAML config format.
- Vehicle name.
- Top-level sensor catalogue.
- Per-sensor logging selection.
- Top-level dashboard canvas size.
- Daily log rotation.
- `type: virtual` accepted as a config value only.

Future release:

- Virtual sensor calculation.
- Dashboard config validation beyond canvas.
- Dashboard assets, decoders, blocks, layers, and rendering.
- Engine start/stop detection.
- Engine-state-based log rotation.
- Complex expressions or condition parsing.

## Top-level structure

```yaml
mock_mode: true
obd_address: serial:///dev/ttyUSB0
obd_debug: false

log:
  rotate: daily
  directory: ./log

vehicle:
  name: "VW Caddy"

dashboard:
  canvas:
    width: 800
    height: 480

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
| `mock_mode` | bool | yes | Use app mock data instead of a real adapter. |
| `obd_address` | string | yes | ELM327/elmobd address. Example: `serial:///dev/ttyUSB0`. |
| `obd_debug` | bool | yes | Enable verbose adapter debugging. |
| `log` | object | yes | Log output configuration. |
| `vehicle` | object | yes | Vehicle metadata. |
| `dashboard` | object | yes | Dashboard visual configuration. In v2.0.x only `canvas` is active. |
| `sensors` | map | yes | Sensor catalogue keyed by stable sensor key. |

## Log config

```yaml
log:
  rotate: daily
  directory: ./log
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `rotate` | string | yes | Rotation mode. Currently only `daily`. |
| `directory` | string | yes | Directory for JSONL logs. |

Daily rotation means GoDriveLog writes readings to a date-based JSONL file and opens a new file when the date changes.

No engine-start or engine-stop log rotation is part of the current release.

## Vehicle config

```yaml
vehicle:
  name: "VW Caddy"
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `name` | string | yes | Human-readable vehicle name. |

Vehicle config no longer owns sensors or dashboard visuals.

## Dashboard config

```yaml
dashboard:
  canvas:
    width: 800
    height: 480
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `canvas.width` | int | yes | Dashboard canvas width in pixels. |
| `canvas.height` | int | yes | Dashboard canvas height in pixels. |

In v2.0.x this only sizes the window and proves the top-level dashboard schema exists. Rendering, assets, decoders, blocks, and layers are later stages.

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
| `type` | string | yes | `obd` or `virtual`. `virtual` is accepted but not implemented yet. |
| `pid` | string | required for `obd` | Raw PID, such as `010C`. |
| `unit` | string | yes | Log/display unit, such as `rpm`, `km/h`, `%`, `C`, or `V`. |
| `refresh` | int | yes | Poll refresh interval in milliseconds. |
| `min` | number | yes | Minimum expected value. |
| `max` | number | yes | Maximum expected value. |
| `log` | bool | no | Write readings to JSONL logs. If not defined, assume `false`. |

The keys under `sensors` are internal sensor names. Use simple stable names such as `rpm`, `speed`, `engine_load`, or `coolant_temp`.

These keys are separate from raw PID values. For example:

```yaml
sensors:
  rpm:
    pid: "010C"
```

Here `rpm` is the GoDriveLog sensor key, and `010C` is the raw PID.

## Polling rule

A sensor is polled when:

```text
type == obd AND log == true
```

A sensor with `log: false` is treated as a known sensor but is not active in v2.0.x.

A sensor with `type: virtual` is valid config but is not polled or calculated in the current release.

## Logging rule

A reading is written to JSONL only when:

```text
log == true
```

## Example config

See [`config.example.yaml`](../config.example.yaml).

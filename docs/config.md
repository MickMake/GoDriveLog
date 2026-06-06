# GoDriveLog Config

This document defines the GoDriveLog configuration format.

The Go structs should mirror this YAML structure as closely as possible. Keep the format simple: one vehicle, one PID catalogue, one log config.

## Current scope

Implemented now:

- YAML config format.
- Vehicle name.
- OBD PID catalogue.
- Per-PID logging selection.
- Per-PID display selection.
- Daily log rotation.
- `type: virtual` accepted as a config value only.

Future release:

- Virtual PID calculation.
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
  pids:
    rpm:
      type: obd
      pid: "010C"
      unit: rpm
      refresh: 250
      min: 0
      max: 7000
      log: true
      display:
        enabled: true
        widget: radial1
        style:
          smoothing_window: 0
          dial_rotation: 0
          view_rotation: 0
          scale_direction: forward
        position:
          x: 20
          y: 20
          width: 360
          height: 90
          z: 0
```

## Root fields

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `mock_mode` | bool | yes | Use app mock data instead of a real OBD adapter. |
| `obd_address` | string | yes | ELM327/elmobd address. Example: `serial:///dev/ttyUSB0`. |
| `obd_debug` | bool | yes | Enable verbose OBD adapter debugging. |
| `log` | object | yes | Log output configuration. |
| `vehicle` | object | yes | Vehicle-specific PID catalogue and metadata. |

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
  pids:
    rpm:
      ...
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `name` | string | yes | Human-readable vehicle name. |
| `pids` | map | yes | PID catalogue keyed by stable sensor key. |

The keys under `vehicle.pids` are internal sensor names. Use simple stable names such as `rpm`, `speed`, `engine_load`, or `coolant_temp`.

These keys are separate from raw OBD PID values. For example:

```yaml
vehicle:
  pids:
    rpm:
      pid: "010C"
```

Here `rpm` is the GoDriveLog sensor key, and `010C` is the raw OBD PID.

## PID config

```yaml
rpm:
  type: obd
  pid: "010C"
  unit: rpm
  refresh: 250
  min: 0
  max: 7000
  log: true
  display:
    enabled: true
    widget: radial1
    style:
      smoothing_window: 0
      dial_rotation: 0
      view_rotation: 0
      scale_direction: forward
    position:
      x: 20
      y: 20
      width: 360
      height: 90
      z: 0
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `type` | string | yes | `obd` or `virtual`. `virtual` is accepted but not implemented yet. |
| `pid` | string | required for `obd` | Raw OBD PID, such as `010C`. |
| `unit` | string | yes | Display/log unit, such as `rpm`, `km/h`, `%`, `C`, or `V`. |
| `refresh` | int | yes | Poll refresh interval in milliseconds. |
| `min` | number | yes | Minimum expected value for display scaling. |
| `max` | number | yes | Maximum expected value for display scaling. |
| `log` | bool | no | Write readings to JSONL logs. If not defined assume log: FALSE. |
| `display` | object | no | Display configuration. If not defined assume display.enabled: FALSE. |

## Display config

```yaml
display:
  enabled: true
  widget: radial1
  style:
    smoothing_window: 0
    dial_rotation: 0
    view_rotation: 0
    scale_direction: forward
  position:
    x: 20
    y: 20
    width: 360
    height: 90
    z: 0
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `enabled` | bool | yes | Show this PID on screen. |
| `widget` | string | required when enabled | Display widget. Use widget ids: `radial1`, `radial2`, `radial3`, `half_top1`, `half_bottom1`, `quarter_tl1`, `quarter_tr1`, `quarter_bl1`, `quarter_br1`, `sweep1`, `sweep2`, `sweep3`, `speedhud1`, `speedhud2`, `speedhud3`, `bar1`, `bar2`, `bar3`, `graph1`, `led1`. |
| `style` | object | no | Display style settings. |
| `style.smoothing_window` | int | no | Moving average window for display smoothing. `0` or `1` disables smoothing. |
| `style.dial_rotation` | int | no | Rotate dial geometry in degrees. Allowed values: `0`, `90`, `180`, `270`. Default `0`. |
| `style.view_rotation` | int | no | Rotate view/layout in degrees (intended for display mounting/viewer orientation). Allowed values: `0`, `90`, `180`, `270`. Default `0`. |
| `style.scale_direction` | string | no | Scale direction along the sweep. Allowed values: `forward`, `reverse`. Default `forward`. |
| `position` | object | required when enabled | Widget position and size. |

If `display.enabled` is `false`, `widget` and `position` may be omitted.

## Display position

```yaml
position:
  x: 20
  y: 20
  width: 360
  height: 90
  z: 0
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `x` | number | yes | X position in the Fyne window. |
| `y` | number | yes | Y position in the Fyne window. |
| `width` | number | yes | Widget width. |
| `height` | number | yes | Widget height. |
| `z` | number | no | Z layer for overlays. Higher values render on top. Default `0`. |

## Polling rule

A PID is polled when:

```text
type == obd AND (log == true OR display.enabled == true)
```

A PID with both `log: false` and `display.enabled: false` is treated as a known PID but is not active.

A PID with `type: virtual` is valid config but is not polled or calculated in the current release.

## Logging rule

A reading is written to JSONL only when:

```text
log == true
```

A PID may be displayed without being logged.

## Display rule

A widget is created only when:

```text
display.enabled == true
```

A PID may be logged without being displayed.

## Example config

See [`config.example.yaml`](../config.example.yaml).

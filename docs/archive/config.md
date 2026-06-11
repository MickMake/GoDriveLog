# GoDriveLog Config

This document defines the GoDriveLog YAML configuration format.

The current schema keeps connection/logging concerns separate from the dashboard visual scene. Sensors define what can be read and logged. The dashboard defines canvas size, render cadence, assets, decoders, blocks, and layers.

## Top-level structure

```yaml
obd:
  mock_mode: true
  address: serial:///dev/ttyUSB0
  debug: false

log:
  rotate: daily
  directory: ./log

vehicle:
  name: "VW Caddy"

dashboard:
  refresh_ms: 1000
  render_min_ms: 5000
  canvas:
    width: 800
    height: 480
  assets:
    - id: background
      type: image
      path: assets/dashboard/bttf/background.png
  decoders: []
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
| `obd` | object | yes | Adapter/mock configuration. |
| `log` | object | yes | Log output configuration. |
| `vehicle` | object | yes | Vehicle metadata. |
| `dashboard` | object | yes | Dashboard visual and render configuration. |
| `sensors` | map | yes | Sensor catalogue keyed by stable sensor key. |

## OBD config

```yaml
obd:
  mock_mode: true
  address: serial:///dev/ttyUSB0
  debug: false
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `mock_mode` | bool | no | Use app mock data instead of a real adapter. Defaults to `false`. |
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

## Dashboard config

```yaml
dashboard:
  refresh_ms: 1000
  render_min_ms: 5000
  canvas:
    width: 800
    height: 480
  assets:
    - id: background
      type: image
      path: assets/dashboard/bttf/background.png
  decoders: []
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

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `refresh_ms` | int | no | Dashboard state refresh tick in milliseconds. Defaults to `100`. Must be positive after defaults are applied. |
| `render_min_ms` | int | no | Minimum time between expensive full renderer rebuilds. `0` disables throttling. Must not be negative. |
| `canvas.width` | int | yes | Dashboard canvas width in pixels. Must be positive. |
| `canvas.height` | int | yes | Dashboard canvas height in pixels. Must be positive. |
| `asset_root` | string | no | Base directory for dashboard assets. |
| `assets` | list | yes | Image, frame set, and charset asset definitions. |
| `decoders` | list | yes | Runtime sensor-to-display decoder definitions. |
| `blocks` | list | yes | Visual block definitions. |
| `layers` | list | yes | Draw-order layer definitions. |

`refresh_ms` controls how often the dashboard evaluates the current sensor state. `render_min_ms` controls how often the Fyne renderer is allowed to rebuild and refresh the full canvas. For Raspberry Pi use, a useful starting point is:

```yaml
dashboard:
  refresh_ms: 1000
  render_min_ms: 5000
```

This keeps state evaluation reasonably current while avoiding repeated heavy rasterisation work.

### Dashboard assets

```yaml
assets:
  - id: background
    type: image
    path: assets/dashboard/bttf/background.png
  - id: throttle_frames
    type: frame_set
    frames:
      - assets/dashboard/bttf/throttle/frame_000.png
      - assets/dashboard/bttf/throttle/frame_001.png
  - id: yellow_digits
    type: charset
    glyphs:
      "0": assets/dashboard/bttf/digits/yellow/0.png
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `id` | string | yes | Unique asset ID. |
| `type` | string | yes | `image`, `frame_set`, or `charset`. |
| `path` | string | for `image` | Image path. |
| `frames` | list | for explicit `frame_set` | Frame paths. |
| `pattern` | string | for generated `frame_set` | Frame path pattern. |
| `frame_count` | int | with generated `frame_set` | Number of generated frames. |
| `glyphs` | map | for `charset` | Character-to-image path map. |

### Dashboard decoders

```yaml
decoders:
  - id: throttle_frame
    type: frame_index
    sensor: throttle_position
    asset: throttle_frames
    frame_count: 10
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `id` | string | yes | Unique decoder ID. |
| `type` | string | yes | `normalize`, `threshold`, `frame_index`, `format_number`, `digits`, or `boolean`. |
| `sensor` | string | no | Sensor key reference. If set, must refer to `sensors`. |
| `input` | string | no | Earlier decoder reference. |
| `asset` | string | no | Asset reference. If set, must refer to `dashboard.assets`. |
| `format` | string | no | Format string for `format_number`. |
| `frame_count` | int | for `frame_index` | Must be positive. |
| `thresholds` | list | for `threshold` | Must not be empty. |

### Dashboard blocks

```yaml
blocks:
  - id: rpm_display
    type: sprite_text
    asset: yellow_digits
    decoder: rpm_digits
    geometry:
      x: 100
      y: 60
      width: 240
      height: 80
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `id` | string | yes | Unique block ID. |
| `type` | string | yes | `image`, `sprite_frame`, `sprite_text`, `group`, `text`, or a reusable dashboard block type. |
| `asset` | string | no | Asset reference. If set, must refer to `dashboard.assets`. |
| `decoder` | string | no | Decoder reference. If set, must refer to `dashboard.decoders`. |
| `blocks` | list | for `group` | Child block IDs. Must refer to configured blocks. |
| `condition` | object | no | Sensor/decoder visibility condition. |
| `geometry.width` | number | for non-group blocks | Must be positive. |
| `geometry.height` | number | for non-group blocks | Must be positive. |

### Dashboard layers

```yaml
layers:
  - id: base
    z: 0
    blocks:
      - background_panel
      - rpm_display
```

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `id` | string | yes | Unique layer ID. |
| `z` | int | no | Layer order hint. |
| `blocks` | list | yes | Block IDs. Each ID must refer to `dashboard.blocks`. |

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

A sensor is polled when:

```text
type == obd AND log == true
```

A sensor with `log: false` is treated as a known sensor but is not active in the current release. A sensor with `type: virtual` is valid config but is not polled or calculated yet.

## Example config

See [`config.example.yaml`](../config.example.yaml).

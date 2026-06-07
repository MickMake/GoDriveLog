# GoDriveLog Config

This document defines the GoDriveLog configuration format.

The v2 schema separates sensor configuration from dashboard visual configuration.
Sensors define what can be read and logged. The dashboard defines the visual canvas and scene configuration that later runtime stages will consume.

## Current scope

Implemented now:

- YAML config format.
- Vehicle name.
- Top-level sensor catalogue.
- Per-sensor logging selection.
- Top-level dashboard canvas size.
- Dashboard assets, decoders, blocks, and layers config schema.
- Dashboard config validation only; no asset loading, decoder execution, or rendering yet.
- Daily log rotation.
- `type: virtual` accepted as a config value only.

Future release:

- Virtual sensor calculation.
- Dashboard asset loading.
- Decoder execution.
- Scene rendering.
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
  assets:
    - id: background
      type: image
      path: assets/dashboard/bttf/background.png
  decoders:
    - id: rpm_digits
      type: digits
      sensor: rpm
      asset: yellow_digits
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
| `mock_mode` | bool | yes | Use app mock data instead of a real adapter. |
| `obd_address` | string | yes | ELM327/elmobd address. Example: `serial:///dev/ttyUSB0`. |
| `obd_debug` | bool | yes | Enable verbose adapter debugging. |
| `log` | object | yes | Log output configuration. |
| `vehicle` | object | yes | Vehicle metadata. |
| `dashboard` | object | yes | Dashboard visual configuration. |
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
  assets: []
  decoders: []
  blocks: []
  layers: []
```

In v2.1.x these fields are loaded and validated only. They are not rendered, and assets are not read from disk yet.

### Dashboard canvas

| Field | Type | Required | Meaning |
|---|---:|---:|---|
| `canvas.width` | int | yes | Dashboard canvas width in pixels. Must be positive. |
| `canvas.height` | int | yes | Dashboard canvas height in pixels. Must be positive. |

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
| `path` | string | for `image` | Image path. Validated as config text only. |
| `frames` | list | for `frame_set` | Frame paths. Validated as config text only. |
| `glyphs` | map | for `charset` | Character-to-image path map. Validated as config text only. |

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
| `input` | string | no | Placeholder for future decoder input references. |
| `asset` | string | no | Asset reference. If set, must refer to `dashboard.assets`. |
| `format` | string | no | Format string for later decoder execution. |
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
| `type` | string | yes | `image`, `sprite_frame`, `sprite_text`, `group`, or `text`. |
| `asset` | string | no | Asset reference. If set, must refer to `dashboard.assets`. |
| `decoder` | string | no | Decoder reference. If set, must refer to `dashboard.decoders`. |
| `blocks` | list | for `group` | Child block IDs. Must refer to configured blocks. |
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
| `z` | int | no | Layer order hint for future rendering. |
| `blocks` | list | no | Block IDs. If set, each ID must refer to `dashboard.blocks`. |

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

A sensor with `log: false` is treated as a known sensor but is not active in v2.1.x.

A sensor with `type: virtual` is valid config but is not polled or calculated in the current release.

## Logging rule

A reading is written to JSONL only when:

```text
log == true
```

## Example config

See [`config.example.yaml`](../config.example.yaml).

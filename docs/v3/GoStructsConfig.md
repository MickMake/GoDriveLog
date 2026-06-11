# GoDriveLog config Go structs

Status: draft lock document  
Schema target: v2.10.0 config/runtime boundaries  
Reference files: `config.full.yaml`, `config.example.yaml`

## 1. Purpose

This document locks the Go structs that should back the new GoDriveLog YAML config schema.

The struct layout must mirror the YAML ownership boundaries:

```text
default_vehicle  = startup vehicle profile key
active_displays  = display IDs started by default
vehicles         = named vehicle profiles and their data source
sensors          = global sensor catalogue and cache rule
log              = file logging config and logging fetch cadence
displays         = physical/logical output bindings
dashboards       = visual dashboard scene definitions
```

The runtime must load one `Config`, resolve the active vehicle from `DefaultVehicle`, start the shared sensor runtime/cache, start logging subscriptions, and start every display listed in `ActiveDisplays`.

## 2. Top-level config

### YAML source

Populated from the root of the YAML file.

```yaml
# path: /
default_vehicle: vw_caddy
active_displays:
  - main
vehicles: {}
sensors: {}
log: {}
displays: {}
dashboards: {}
```

### Go struct

```go
type Config struct {
    DefaultVehicle string                     `yaml:"default_vehicle"`
    ActiveDisplays []string                   `yaml:"active_displays"`
    Vehicles       map[string]VehicleConfig   `yaml:"vehicles"`
    Sensors        map[string]SensorConfig    `yaml:"sensors"`
    Log            LogConfig                  `yaml:"log"`
    Displays       map[string]DisplayConfig   `yaml:"displays"`
    Dashboards     map[string]DashboardConfig `yaml:"dashboards"`
}
```

### Notes

- `DefaultVehicle` is a key into `Vehicles`.
- `ActiveDisplays` contains keys into `Displays`.
- `Sensors` is global and shared by all vehicles, loggers, dashboards, and future consumers.
- `Log`, `Displays`, and `Dashboards` are app/session configuration, not vehicle-owned.

## 3. Vehicle profile config

### YAML source

Populated from each entry under `vehicles.<vehicle_id>`.

```yaml
# path: vehicles.vw_caddy
vehicles:
  vw_caddy:
    name: "VW Caddy"
    obd:
      source: real
      address: serial:///dev/ttyUSB0
      debug: false
```

### Go struct

```go
type VehicleConfig struct {
    Name string    `yaml:"name"`
    OBD  OBDConfig `yaml:"obd"`
}
```

### Notes

- The map key, for example `vw_caddy`, is the stable vehicle ID.
- `Name` is human-readable display text only.
- `OBD` belongs here because the selected vehicle controls the active data source.
- Sensors do **not** live under `VehicleConfig`; they are global under `Config.Sensors`.

## 4. OBD source config

### YAML source

Populated from `vehicles.<vehicle_id>.obd`.

```yaml
# path: vehicles.vw_caddy.obd
obd:
  source: real
  address: serial:///dev/ttyUSB0
  debug: false
```

```yaml
# path: vehicles.race_driver.obd
obd:
  source: mock
  address: mock://race-driver
  debug: false
```

### Go struct

```go
type OBDConfig struct {
    Source  string `yaml:"source"`
    Address string `yaml:"address"`
    Debug   bool   `yaml:"debug"`
}
```

### Notes

- `Source` must be either `real` or `mock`.
- `mock_mode` is removed. Do not keep the old boolean.
- `Address` is source-specific:
  - `real`: serial/Bluetooth/etc address such as `serial:///dev/ttyUSB0`.
  - `mock`: mock profile address such as `mock://race-driver`.
- `Debug` controls OBD/source-level debug output.

## 5. Sensor catalogue config

### YAML source

Populated from each entry under `sensors.<sensor_id>`.

```yaml
# path: sensors.rpm
sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: rpm
    cache: 250
    min: 0
    max: 7000
```

### Go struct

```go
type SensorConfig struct {
    Type  string  `yaml:"type"`
    PID   string  `yaml:"pid"`
    Unit  string  `yaml:"unit"`
    Cache int     `yaml:"cache"`
    Min   float64 `yaml:"min"`
    Max   float64 `yaml:"max"`
}
```

### Notes

- The map key, for example `rpm`, is the stable sensor ID used by log and dashboard config.
- `Type` currently supports `obd`. Future values may include `virtual` or calculated sensors, but do not add them until required.
- `PID` is required for `type: obd`.
- `Unit` is the unit attached to readings and display/log output.
- `Cache` is milliseconds between physical OBD requests for this sensor.
- `Min` and `Max` describe expected value bounds for scaling, validation, and display.
- `SensorConfig` does **not** contain `refresh`.
- `SensorConfig` does **not** decide whether a sensor is logged or displayed.

## 6. Cache semantics

The sensor cache is configured by `sensors.<sensor_id>.cache`.

Runtime rule:

```text
If a sensor request arrives and the cached value age is <= sensor.cache,
return the cached value.

If the cached value age is > sensor.cache,
perform a physical source read through the selected vehicle's OBD source,
then update the shared sensor state/cache.
```

All consumers use the same cache:

```text
OBD/mock source
  -> shared sensor runtime/cache
     -> logger
     -> display main
     -> display aux
```

Do not allow logger and dashboards to independently hammer the OBD source.

## 7. Log config

### YAML source

Populated from `log`.

```yaml
# path: log
log:
  rotate: daily
  directory: ./log
  sensors:
    rpm:
      refresh: 250
```

### Go struct

```go
type LogConfig struct {
    Rotate    string                     `yaml:"rotate"`
    Directory string                     `yaml:"directory"`
    Sensors   map[string]LogSensorConfig `yaml:"sensors"`
}

type LogSensorConfig struct {
    Refresh int `yaml:"refresh"`
}
```

### Notes

- `Rotate` currently supports `daily`.
- `Directory` is where JSONL logs are written.
- `Sensors` is a map of sensor IDs to logging subscription config.
- `LogSensorConfig.Refresh` is the logger's fetch cadence in milliseconds.
- `refresh` exists here because the logger is an active polling consumer.
- A sensor can exist in `Config.Sensors` and not be logged.
- A log sensor key must exist in `Config.Sensors`.

## 8. Display binding config

### YAML source

Populated from each entry under `displays.<display_id>`.

```yaml
# path: displays.main
displays:
  main:
    dashboard: bttf_primary
    output: default
    fullscreen: true
```

### Go struct

```go
type DisplayConfig struct {
    Dashboard  string `yaml:"dashboard"`
    Output     string `yaml:"output"`
    Fullscreen bool   `yaml:"fullscreen"`
}
```

### Notes

- The map key, for example `main`, is the stable display ID.
- `Dashboard` is a key into `Config.Dashboards`.
- `Output` is the physical/logical display target. Examples: `default`, `HDMI-1`, `HDMI-2`.
- `Fullscreen` controls whether that output should start fullscreen.
- Display config owns output/window behaviour only.
- Visual scene definitions belong under `dashboards`, not under `displays`.

## 9. Active displays

### YAML source

Populated from `active_displays`.

```yaml
# path: active_displays
active_displays:
  - main
  - aux
```

### Go field

```go
ActiveDisplays []string `yaml:"active_displays"`
```

### Notes

- Every listed display ID must exist in `Config.Displays`.
- Multiple displays may run at the same time.
- Each active display consumes the same shared sensor state/cache.
- A CLI override may later replace or filter this list, for example `-display main`, `-display main,aux`, or `-display all`.

## 10. Dashboard config

### YAML source

Populated from each entry under `dashboards.<dashboard_id>`.

```yaml
# path: dashboards.bttf_primary
dashboards:
  bttf_primary:
    refresh_ms: 1000
    render_min_ms: 5000
    canvas:
      width: 800
      height: 480
    asset_root: assets/dashboard/bttf
    assets: []
    decoders: []
    blocks: []
    layers: []
```

### Go struct

```go
type DashboardConfig struct {
    RefreshMS   int             `yaml:"refresh_ms"`
    RenderMinMS int             `yaml:"render_min_ms"`
    Canvas      CanvasConfig    `yaml:"canvas"`
    AssetRoot   string          `yaml:"asset_root"`
    Assets      []AssetConfig   `yaml:"assets"`
    Decoders    []DecoderConfig `yaml:"decoders"`
    Blocks      []BlockConfig   `yaml:"blocks"`
    Layers      []LayerConfig   `yaml:"layers"`
}
```

### Notes

- The map key, for example `bttf_primary`, is the stable dashboard ID.
- `RefreshMS` is the dashboard evaluation/update cadence.
- `RenderMinMS` is a renderer throttle/safety interval.
- `Canvas` defines the dashboard design coordinate space.
- `AssetRoot` is the base directory for dashboard assets.
- `Assets`, `Decoders`, `Blocks`, and `Layers` describe the visual scene.
- Dashboard config references sensors by global sensor ID, for example `rpm`.
- Dashboard config does **not** decide whether a sensor is logged.
- Dashboard config does **not** read OBD directly.

## 11. Canvas config

### YAML source

Populated from `dashboards.<dashboard_id>.canvas`.

```yaml
canvas:
  width: 800
  height: 480
```

### Go struct

```go
type CanvasConfig struct {
    Width  int `yaml:"width"`
    Height int `yaml:"height"`
}
```

### Notes

- The canvas is the logical design size for the dashboard scene.
- Renderer/window scaling is separate from the scene's logical coordinates.

## 12. Asset config

### YAML source

Populated from each item in `dashboards.<dashboard_id>.assets`.

```yaml
assets:
  - id: throttle_frames
    type: frame_set
    pattern: throttle/frame_{index:03}.svg
    frame_count: 11

  - id: yellow_digits
    type: charset
    glyphs:
      "0": digits/yellow/0.svg
```

### Go struct

```go
type AssetConfig struct {
    ID         string            `yaml:"id"`
    Type       string            `yaml:"type"`
    Path       string            `yaml:"path"`
    Pattern    string            `yaml:"pattern"`
    FrameCount int               `yaml:"frame_count"`
    Glyphs     map[string]string `yaml:"glyphs"`
}
```

### Notes

- `ID` is referenced by blocks and decoders.
- `Type` currently supports:
  - `image`
  - `frame_set`
  - `charset`
- `Path` is used by `type: image`.
- `Pattern` and `FrameCount` are used by `type: frame_set`.
- `Glyphs` is used by `type: charset`.

## 13. Decoder config

### YAML source

Populated from each item in `dashboards.<dashboard_id>.decoders`.

```yaml
decoders:
  - id: rpm_text
    type: format_number
    sensor: rpm
    format: "0000"

  - id: rpm_digits
    type: digits
    input: rpm_text
    asset: yellow_digits

  - id: rpm_warning
    type: threshold
    sensor: rpm
    thresholds:
      - at: 0
        value: normal
      - at: 2100
        value: warning
```

### Go struct

```go
type DecoderConfig struct {
    ID         string            `yaml:"id"`
    Type       string            `yaml:"type"`
    Sensor     string            `yaml:"sensor"`
    Input      string            `yaml:"input"`
    Asset      string            `yaml:"asset"`
    Format     string            `yaml:"format"`
    FrameCount int               `yaml:"frame_count"`
    Thresholds []ThresholdConfig `yaml:"thresholds"`
}

type ThresholdConfig struct {
    At    float64 `yaml:"at"`
    Value string  `yaml:"value"`
}
```

### Notes

- `ID` is referenced by blocks or later decoders.
- `Type` determines which fields are required.
- `Sensor` references a key in `Config.Sensors`.
- `Input` references another decoder output.
- `Asset` references an asset ID.
- `Format` is used by `format_number`.
- `FrameCount` is used by frame index style decoders.
- `Thresholds` is used by `threshold`.

## 14. Block config

### YAML source

Populated from each item in `dashboards.<dashboard_id>.blocks`.

```yaml
blocks:
  - id: rpm_display
    type: seven_segment_number
    asset: yellow_digits
    decoder: rpm_digits
    geometry:
      x: 92
      y: 58
      width: 300
      height: 100

  - id: main_cluster
    type: group
    blocks:
      - rpm_display
      - throttle_bar
```

### Go struct

```go
type BlockConfig struct {
    ID        string           `yaml:"id"`
    Type      string           `yaml:"type"`
    Asset     string           `yaml:"asset"`
    Decoder   string           `yaml:"decoder"`
    Geometry  GeometryConfig   `yaml:"geometry"`
    Condition ConditionConfig  `yaml:"condition"`
    Blocks    []string         `yaml:"blocks"`
}
```

### Notes

- `ID` is referenced by layers or group blocks.
- `Type` determines the rendering primitive or reusable block alias.
- `Asset` references an asset ID.
- `Decoder` references a decoder ID.
- `Geometry` defines placement and size.
- `Condition` controls whether the block is visible.
- `Blocks` is used by `type: group` to reference child block IDs.

## 15. Geometry config

### YAML source

Populated from block `geometry`.

```yaml
geometry:
  x: 0
  y: 0
  width: 800
  height: 480
```

### Go struct

```go
type GeometryConfig struct {
    X      int `yaml:"x"`
    Y      int `yaml:"y"`
    Width  int `yaml:"width"`
    Height int `yaml:"height"`
}
```

### Notes

- Coordinates use the dashboard canvas coordinate space.
- `Width` and `Height` are required for non-group visual blocks.

## 16. Condition config

### YAML source

Populated from block `condition`.

```yaml
condition:
  sensor: rpm
  status: ok
```

```yaml
condition:
  decoder: rpm_warning
  equals: warning
```

### Go struct

```go
type ConditionConfig struct {
    Sensor    string  `yaml:"sensor"`
    Decoder   string  `yaml:"decoder"`
    Status    string  `yaml:"status"`
    Equals    string  `yaml:"equals"`
    NotEquals string  `yaml:"not_equals"`
    Min       float64 `yaml:"min"`
    Max       float64 `yaml:"max"`
}
```

### Notes

- `Sensor` references a key in `Config.Sensors`.
- `Decoder` references a decoder ID.
- `Status` matches sensor state status such as `ok`, `stale`, or `error`.
- `Equals` and `NotEquals` compare decoder/sensor values.
- `Min` and `Max` support threshold-style visibility checks.
- If zero values need to be semantically meaningful for `Min` or `Max`, this struct may need pointer fields later. Do not change unless required by validation/runtime behaviour.

## 17. Layer config

### YAML source

Populated from each item in `dashboards.<dashboard_id>.layers`.

```yaml
layers:
  - id: dashboard
    z: 10
    blocks:
      - main_cluster
```

### Go struct

```go
type LayerConfig struct {
    ID     string   `yaml:"id"`
    Z      int      `yaml:"z"`
    Blocks []string `yaml:"blocks"`
}
```

### Notes

- `ID` is the stable layer name.
- `Z` controls render order.
- `Blocks` references block IDs from the same dashboard.

## 18. Resolved runtime config

The YAML structs above are raw config. Runtime should also produce a resolved config after validation.

### Go structs

```go
type ResolvedConfig struct {
    Raw            Config
    VehicleID      string
    Vehicle        VehicleConfig
    ActiveDisplays []ResolvedDisplayConfig
    Sensors        map[string]SensorConfig
    Log            LogConfig
}

type ResolvedDisplayConfig struct {
    ID        string
    Display   DisplayConfig
    Dashboard DashboardConfig
}
```

### Notes

- `ResolvedConfig` is not directly populated from YAML.
- It is produced after validating `Config`.
- It binds:
  - `DefaultVehicle` -> `Vehicles[DefaultVehicle]`
  - `ActiveDisplays[]` -> `Displays[id]`
  - `Displays[id].Dashboard` -> `Dashboards[dashboard_id]`
- Runtime code should prefer resolved config where possible.

## 19. Validation rules

At config load time, fail loudly if:

- `default_vehicle` is empty.
- `default_vehicle` is not present in `vehicles`.
- `active_displays` is empty.
- any `active_displays` item is not present in `displays`.
- any display references a missing dashboard.
- `vehicles` is empty.
- any vehicle has an empty `name`.
- any vehicle has invalid `obd.source`; allowed values are `real`, `mock`.
- `obd.address` is empty for `source: real`.
- `sensors` is empty.
- any sensor has an empty `type`.
- any `type: obd` sensor has an empty `pid`.
- any sensor has an empty `unit`.
- any sensor has `cache <= 0`.
- any sensor has `max <= min`.
- `log.rotate` is not `daily`.
- `log.directory` is empty.
- any `log.sensors` key is not present in `sensors`.
- any `log.sensors.*.refresh <= 0`.
- any dashboard decoder references a missing sensor, asset, or decoder input.
- any dashboard block references a missing asset, decoder, child block, or invalid condition target.
- any dashboard layer references a missing block.

## 20. Ownership boundaries

| Struct | YAML path | Owns | Must not own |
|---|---|---|---|
| `Config` | `/` | top-level config wiring | runtime state |
| `VehicleConfig` | `vehicles.<id>` | vehicle identity and data source | sensors, log, dashboard visuals |
| `OBDConfig` | `vehicles.<id>.obd` | selected vehicle source details | sensor definitions |
| `SensorConfig` | `sensors.<id>` | sensor metadata and cache interval | logging/display cadence |
| `LogConfig` | `log` | log output and subscriptions | sensor metadata, dashboard visuals |
| `DisplayConfig` | `displays.<id>` | output binding/window behaviour | visual scene definition |
| `DashboardConfig` | `dashboards.<id>` | visual scene | physical output binding, OBD reads |
| `ResolvedConfig` | generated | validated runtime bindings | YAML decoding |

## 21. Recommended file placement

Proposed source files:

```text
internal/config/config.go       # Config, Load, defaults, validation entrypoint
internal/config/vehicle.go      # VehicleConfig, OBDConfig
internal/config/sensor.go       # SensorConfig
internal/config/log.go          # LogConfig, LogSensorConfig
internal/config/display.go      # DisplayConfig
internal/config/dashboard.go    # DashboardConfig and visual scene structs
internal/config/resolved.go     # ResolvedConfig and binding helpers
```

Keep dashboard renderer, sensor runtime/cache, and logger implementation outside the config package.

## 22. Locked top-level schema

```yaml
default_vehicle: vw_caddy
active_displays:
  - main

vehicles:
  vw_caddy:
    name: "VW Caddy"
    obd:
      source: real
      address: serial:///dev/ttyUSB0
      debug: false

sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: rpm
    cache: 250
    min: 0
    max: 7000

log:
  rotate: daily
  directory: ./log
  sensors:
    rpm:
      refresh: 250

displays:
  main:
    dashboard: bttf_primary
    output: default
    fullscreen: true

dashboards:
  bttf_primary:
    refresh_ms: 1000
    render_min_ms: 5000
    canvas:
      width: 800
      height: 480
    asset_root: assets/dashboard/bttf
    assets: []
    decoders: []
    blocks: []
    layers: []
```

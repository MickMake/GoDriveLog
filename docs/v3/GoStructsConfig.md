# GoDriveLog v3 config Go structs

Status: draft alignment document  
Schema target: v3 simplified runtime/config boundaries  
Reference files: `docs/v3/config.full.yaml`, `docs/v3/config.example.yaml`

## 1. Purpose

This document describes the Go structs that should back the intended v3 YAML config schema.

The struct layout mirrors the v3 ownership boundaries:

```text
vehicles   = available vehicle endpoint profiles
sensors    = global sensor catalogue and polling rules
assets     = reusable dashboard asset packs
logs       = event-log subscribers
dashboards = physical display dashboards and widgets
```

The runtime should load one `Config`, select a vehicle, connect to the selected vehicle's OBD-like endpoint, start the shared sensor polling runtime, and then feed logs and dashboards from sensor events.

There is intentionally no `default_vehicle`, no `active_displays`, no top-level `displays`, no logger refresh config, and no `mock` / `real` source switch.

## 2. Top-level config

### YAML source

```yaml
vehicles: {}
sensors: {}
assets: {}
logs: {}
dashboards: {}
```

### Go struct

```go
type Config struct {
    Vehicles   map[string]VehicleConfig   `yaml:"vehicles"`
    Sensors    map[string]SensorConfig    `yaml:"sensors"`
    Assets     AssetConfig                `yaml:"assets"`
    Logs       map[string]LogConfig       `yaml:"logs"`
    Dashboards map[string]DashboardConfig `yaml:"dashboards"`
}
```

### Notes

- If exactly one vehicle exists, the runtime may use it by default.
- If multiple vehicles exist, require an explicit runtime choice such as `--vehicle <id>`.
- Sensors are global and shared by logs, dashboards, and future consumers.
- Dashboard presence means the dashboard is active.
- Logs and dashboards subscribe to sensor events; they do not poll sensors independently.

## 3. Vehicle config

### YAML source

```yaml
vehicles:
  vw_caddy:
    name: "VW Caddy"
    obd:
      address: "serial:///dev/ttyUSB0"
      timeout: 1000
```

### Go struct

```go
type VehicleConfig struct {
    Name string    `yaml:"name"`
    OBD  OBDConfig `yaml:"obd"`
}

type OBDConfig struct {
    Address string `yaml:"address"`
    Timeout int    `yaml:"timeout"`
}
```

### Notes

- The vehicle map key, for example `vw_caddy`, is the stable vehicle ID.
- `Name` is human-readable display/log text.
- `OBD.Address` is an OBD-like endpoint, such as `serial:///dev/ttyUSB0` or `tcp://127.0.0.1:35000`.
- The runtime should not know or care whether the endpoint is real hardware or a simulator.
- Do not add `source: real` or `source: mock`.

## 4. Sensor config

### YAML source

```yaml
sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: "rpm"
    poll: 250
    min: 0
    max: 7000
```

### Go struct

```go
type SensorConfig struct {
    Type string   `yaml:"type"`
    PID  string   `yaml:"pid,omitempty"`
    Unit string   `yaml:"unit"`
    Poll int      `yaml:"poll"`
    Min  *float64 `yaml:"min,omitempty"`
    Max  *float64 `yaml:"max,omitempty"`
}
```

### Notes

- The sensor map key, for example `rpm`, is the stable sensor ID used by logs and dashboard widgets.
- `Type` starts with `obd`; future values can be added only when needed.
- `PID` is required for `type: obd`.
- `Poll` is the sensor polling cadence in milliseconds.
- Sensors own timing. Logs and dashboards do not specify refresh.
- `Min` and `Max` are optional expected bounds for validation and display scaling.
- Sensor state events should include value, unit, status, original read timestamp, and a sequence/version.

## 5. Sensor event semantics

Suggested runtime event shape:

```go
type SensorEvent struct {
    SensorID      string
    Value         any
    Unit          string
    Status        SensorStatus
    ReadTimestamp time.Time
    Sequence      uint64
}
```

Suggested statuses:

```go
type SensorStatus string

const (
    SensorStatusOK          SensorStatus = "ok"
    SensorStatusStale       SensorStatus = "stale"
    SensorStatusError       SensorStatus = "error"
    SensorStatusUnsupported SensorStatus = "missing/unsupported"
)
```

Runtime should emit events on first reading, value change, status change, recovery, and stale/error/unsupported changes.

Do not encode errors as numeric values such as `0`.

## 6. Assets config

### YAML source

```yaml
assets:
  digit_sets: {}
  bar_sets: {}
  frame_sets: {}
  indicator_sets: {}
  image_sets: {}
```

### Go struct

```go
type AssetConfig struct {
    DigitSets     map[string]DigitSetConfig     `yaml:"digit_sets"`
    BarSets       map[string]BarSetConfig       `yaml:"bar_sets"`
    FrameSets     map[string]FrameSetConfig     `yaml:"frame_sets"`
    IndicatorSets map[string]IndicatorSetConfig `yaml:"indicator_sets"`
    ImageSets     map[string]ImageSetConfig     `yaml:"image_sets"`
}
```

The common asset render sandwich is:

```text
background, if present
value/state-driven dynamic layer
foreground, if present
```

### Digit sets

```go
type DigitSetConfig struct {
    Background   string            `yaml:"background,omitempty"`
    Characters   map[string]string `yaml:"characters"`
    DecimalPoint string            `yaml:"decimal_point,omitempty"`
    Foreground   string            `yaml:"foreground,omitempty"`
    Spacing      int               `yaml:"spacing,omitempty"`
}
```

Notes:

- `Characters` should include `0` through `9` and `-` where the dashboard may show negative values.
- A blank/padded slot means background-only when `Background` exists.
- Decimal points are overlays, not normal character slots.
- Prefer `characters`, not `digits`, because `-` is not a digit and config fossils are annoying.

### Bar sets

```go
type BarSetConfig struct {
    Background string            `yaml:"background,omitempty"`
    Cells      map[string]string `yaml:"cells"`
    Foreground string            `yaml:"foreground,omitempty"`
    Spacing    int               `yaml:"spacing,omitempty"`
}
```

### Frame sets

```go
type FrameSetConfig struct {
    Background string            `yaml:"background,omitempty"`
    Frames     FrameRangeConfig  `yaml:"frames"`
    Foreground string            `yaml:"foreground,omitempty"`
}

type FrameRangeConfig struct {
    Path  string `yaml:"path"`
    First int    `yaml:"first"`
    Last  int    `yaml:"last"`
}
```

Use frame sets for complex curves, sweeps, and retro gauge tricks. Do not add a dashboard drawing language until reality turns up with a receipt.

### Indicator sets

```go
type IndicatorSetConfig struct {
    Background string            `yaml:"background,omitempty"`
    States     map[string]string `yaml:"states"`
    Foreground string            `yaml:"foreground,omitempty"`
}
```

Expected states:

```text
off
on
unknown
```

Renderer rule:

```text
if sensor status != ok: unknown
else if value == true: on
else: off
```

### Image sets

```go
type ImageSetConfig struct {
    Image      string `yaml:"image,omitempty"`
    Background string `yaml:"background,omitempty"`
    Foreground string `yaml:"foreground,omitempty"`
}
```

## 7. Logs config

### YAML source

```yaml
logs:
  jsonl:
    path: "logs/godrivelog.jsonl"
    sensors:
      - speed
      - rpm
```

### Go struct

```go
type LogConfig struct {
    Path    string   `yaml:"path"`
    Sensors []string `yaml:"sensors"`
}
```

### Notes

- Logs are subscribers to sensor events.
- A listed sensor key must exist under `Config.Sensors`.
- Logs should write first readings, value changes, and status changes.
- Logs should not spam unchanged duplicate readings.
- No log refresh field.
- No per-sensor log polling field.

## 8. Dashboard config

### YAML source

```yaml
dashboards:
  primary:
    display: "HDMI-1"
    size:
      width: 1920
      height: 480
    widgets: []
```

### Go struct

```go
type DashboardConfig struct {
    Display string         `yaml:"display"`
    Size    SizeConfig     `yaml:"size"`
    Widgets []WidgetConfig `yaml:"widgets"`
}

type SizeConfig struct {
    Width  int `yaml:"width"`
    Height int `yaml:"height"`
}
```

### Notes

- The dashboard map key is the stable dashboard ID.
- Dashboard presence means active.
- A dashboard owns its physical/logical display target.
- Multiple physical regions on the same display should be widgets, not separate display configs.
- Dashboards do not define polling cadence.
- Dashboards do not read OBD directly.

## 9. Widget config

The widget model is intentionally declarative.

```go
type WidgetConfig struct {
    ID       string        `yaml:"id"`
    Type     string        `yaml:"type"`
    Sensor   string        `yaml:"sensor,omitempty"`
    Asset    string        `yaml:"asset"`
    Position [2]int        `yaml:"position"`
    Digits   int           `yaml:"digits,omitempty"`
    Format   string        `yaml:"format,omitempty"`
    Cells    int           `yaml:"cells,omitempty"`
    Min      *float64      `yaml:"min,omitempty"`
    Max      *float64      `yaml:"max,omitempty"`
    Reverse  bool          `yaml:"reverse,omitempty"`
    Zones    []ZoneConfig  `yaml:"zones,omitempty"`
}

type ZoneConfig struct {
    UpTo float64 `yaml:"up_to"`
    Cell string  `yaml:"cell"`
}
```

Expected initial widget types:

```text
image
digit_display
bar_display
frame_gauge
indicator
```

Notes:

- Use `position`, not `rect`, for native-size image-backed widgets.
- Asset packs own native visual geometry.
- Renderers should validate that related images in an asset pack have compatible dimensions.
- Keep expressions, scripts, inheritance, and clever conditions out of v3 until genuinely needed.

## 10. Validation rules

Initial validation should check:

- At least one vehicle exists.
- If multiple vehicles exist, runtime selection is explicit.
- Sensor IDs referenced by logs exist.
- Sensor IDs referenced by widgets exist, except `type: image` widgets with no sensor.
- Asset references exist in the appropriate asset family.
- Digit sets contain required displayed characters for configured formats where practical.
- Frame ranges have `first <= last`.
- Bar widgets reference valid cell names.
- Indicator sets contain `off`, `on`, and `unknown`.
- Dashboard sizes are positive.
- Widget positions are present.

## 11. Non-goals

Do not add these without strong evidence:

- plugin systems
- source orchestration
- live config reload
- dashboard scripting
- config inheritance
- generic event buses
- enable flags everywhere
- mock/real branches leaking into core runtime

Boring boundaries first. Fancy pixels second. YAML goblins never.

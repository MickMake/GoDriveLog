# GoDriveLog v3 config Go structs

Status: draft alignment document  
Schema target: v3 simplified runtime/config boundaries  
Reference files: `docs/v3/config.full.yaml`, `docs/v3/config.example.yaml`

## 1. Purpose

This document describes the Go structs that should back the intended v3 YAML config schema.

The struct layout mirrors the v3 ownership boundaries:

```text
vehicles   = runtime profiles: OBD endpoint plus selected logs and dashboards
sensors    = global sensor catalogue and polling rules
assets     = global reusable dashboard asset packs
logs       = global event-log subscriber definitions
dashboards = global physical display dashboard definitions
```

The runtime should load one `Config`, select one vehicle, connect to that vehicle's OBD-like endpoint, start the shared sensor polling runtime, and then run only the logs and dashboards selected by the vehicle.

Sensors and assets are global catalogues. Vehicles do not directly list sensors or assets. Logs reference global sensors. Dashboard widgets reference global sensors and global assets.

The v3 root schema is intentionally small. Treat the five documented top-level sections as an allow-list, and reject unknown fields during development.

Strict v3 config loading must reject unknown fields at every documented level, not only at the root. This includes nested vehicle, OBD, sensor, asset, log, dashboard, widget, frame, and zone fields.

## 2. Naming rules

IDs should be boring and stable.

Recommended ID pattern:

```text
^[a-z][a-z0-9_]*$
```

Apply this to:

- vehicle IDs
- sensor IDs
- asset IDs inside each asset family
- log IDs
- dashboard IDs
- widget IDs

Rules:

- IDs are case-sensitive, but v3 should use lowercase snake_case.
- Widget IDs must be unique within a dashboard.
- Asset IDs only need to be unique within their asset family.
- Do not rely on display names as IDs.

## 3. Top-level config

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
- If multiple vehicles exist, require an explicit runtime vehicle choice such as `--vehicle <id>`.
- The selected vehicle defines the OBD endpoint, log definitions, and dashboard definitions to run.
- Sensors are global and shared by logs, dashboards, and future consumers.
- Assets are global and shared by dashboards.
- Logs and dashboards subscribe to sensor events; they do not poll sensors independently.
- Unknown root fields fail validation.
- Unknown nested fields fail validation.

## 4. Vehicle config

### YAML source

```yaml
vehicles:
  vw_caddy:
    name: "VW Caddy"
    obd:
      address: "serial:///dev/ttyUSB0"
      timeout: 1000
    logs:
      - jsonl
    dashboards:
      - simple_primary
```

### Go struct

```go
type VehicleConfig struct {
    Name       string    `yaml:"name"`
    OBD        OBDConfig `yaml:"obd"`
    Logs       []string  `yaml:"logs,omitempty"`
    Dashboards []string  `yaml:"dashboards,omitempty"`
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
- `Logs` selects global log definitions to run for this vehicle.
- `Dashboards` selects global dashboard definitions to render for this vehicle.
- The runtime should not know or care whether the endpoint is real hardware or a simulator.
- Prefer endpoint addresses over endpoint-type branching.
- Vehicles do not directly list sensors.
- Vehicles do not directly list assets.

### Vehicle log/dashboard selection

Rules:

- If a vehicle lists `logs`, every listed log ID must exist under top-level `logs`.
- If a vehicle lists `dashboards`, every listed dashboard ID must exist under top-level `dashboards`.
- If a vehicle omits `logs` and exactly one log is defined, the runtime may use that single log automatically.
- If a vehicle omits `dashboards` and exactly one dashboard is defined, the runtime may use that single dashboard automatically.
- If multiple logs are defined, each vehicle should list the logs it runs.
- If multiple dashboards are defined, each vehicle should list the dashboards it renders.
- Display collision validation applies to the dashboards selected by the selected vehicle.
- Within one selected vehicle's dashboard set, no two dashboards may target the same physical display.
- Multiple dashboard definitions may target the same display when they are alternatives selected by different vehicle profiles.

### Endpoint validation

Initial validation should require:

- `address` is present.
- `serial://` endpoints include a non-empty path.
- `tcp://` endpoints include host and port.
- `timeout > 0`.
- Recommended initial sanity range: `100 <= timeout <= 30000` milliseconds.

The timeout range is a guardrail, not a tuning API. Change it only if real adapters prove the range wrong.

## 5. Sensor config

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
- `Type` starts with `obd`; future values can be added only when needed and documented first.
- `PID` is required for `type: obd`.
- `Poll` is the sensor polling cadence in milliseconds.
- Sensors own timing. Logs and dashboards do not specify cadence.
- `Min` and `Max` are optional expected bounds for validation and display scaling.
- If both `Min` and `Max` are present, require `Min < Max`.
- Sensor state events should include value, unit, status, original read timestamp, and a sequence/version.
- Sensors are global definitions, not vehicle-owned definitions.

## 6. Sensor event semantics

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

### Value expectations

Initial v3 value expectations:

- Numeric OBD readings should use a numeric type after decoding, normally `float64` unless an integer type is useful internally.
- Boolean/status OBD readings should use `bool`.
- Unsupported, stale, or error states should be represented by `Status`, not by fake values.
- `Unit` should remain the configured unit label for usable values.

## 7. Stale timing

Initial runtime stale rule:

```text
stale_after = max(sensor.poll * 3, 1000ms)
```

A sensor becomes `stale` when no successful read has occurred within that window.

Rules:

- Stale is derived by the runtime.
- Stale timing is not a log setting.
- Stale timing is not a dashboard setting.
- Do not add YAML fields for stale timing unless a later reviewed design proves they are needed.
- Stale transitions must emit sensor events.
- Recovery from stale to ok must emit sensor events.

## 8. Assets config

Asset paths in v3 config are repository-root relative unless a later document explicitly says otherwise.

Example:

```yaml
assets/dashboard/bttf/amber_digits/amber0.png
```

Do not teach multiple path dialects in active v3 examples.

Assets are global definitions. Vehicles do not directly list assets; dashboards and widgets reference them.

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
- Prefer `characters`, not `digits`, because formatted display output is not limited to numeric digits.

### Digit display semantics

For widgets that use a digit set:

- `digits` is the number of character slots, excluding decimal point overlays.
- Formatted decimal separators do not consume a character slot.
- If the configured `format` can emit a decimal separator, the digit set must provide `decimal_point`.
- Formatted output must fit the configured slot count after decimal separators are removed.
- Every non-decimal formatted character must have a `characters` entry, unless it is a blank/padded slot rendered as background-only.

### Bar sets

```go
type BarSetConfig struct {
    Background string            `yaml:"background,omitempty"`
    Cells      map[string]string `yaml:"cells"`
    Foreground string            `yaml:"foreground,omitempty"`
    Spacing    int               `yaml:"spacing,omitempty"`
}
```

Bar set rules:

- `off` is required for unfilled cells.
- If a bar widget omits `zones`, `on` is required for filled cells.
- Zone cell names must exist in the bar set's `cells` map.

### Frame sets

```go
type FrameSetConfig struct {
    Background string           `yaml:"background,omitempty"`
    Frames     FrameRangeConfig `yaml:"frames"`
    Foreground string           `yaml:"foreground,omitempty"`
}

type FrameRangeConfig struct {
    Path  string `yaml:"path"`
    First int    `yaml:"first"`
    Last  int    `yaml:"last"`
}
```

Use frame sets for complex curves, sweeps, and retro gauge tricks. Do not add a dashboard drawing language until reality turns up with a receipt.

Frame range validation requires `first <= last`.

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

## 9. Logs config

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

- Logs are global subscriber definitions.
- Vehicles select which log definitions run.
- A listed sensor key must exist under `Config.Sensors`.
- Logs should write first readings, value changes, and status changes.
- Logs should not spam unchanged duplicate readings.
- Logs do not define polling/cadence fields.

## 10. Dashboard config

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

- Dashboards are global display definitions.
- Vehicles select which dashboard definitions render.
- The dashboard map key is the stable dashboard ID.
- A dashboard owns its physical/logical display target.
- Multiple physical regions on the same display should be widgets.
- Multiple dashboard definitions may share a display if they are alternatives selected by different vehicles.
- Within one selected vehicle's dashboard set, no two dashboards may target the same display.
- Dashboards do not define polling cadence.
- Dashboards do not read OBD directly.

## 11. Widget config

The widget model is intentionally declarative.

```go
type WidgetConfig struct {
    ID       string       `yaml:"id"`
    Type     string       `yaml:"type"`
    Sensor   string       `yaml:"sensor,omitempty"`
    Asset    string       `yaml:"asset"`
    Position [2]int       `yaml:"position"`
    Digits   int          `yaml:"digits,omitempty"`
    Format   string       `yaml:"format,omitempty"`
    Cells    int          `yaml:"cells,omitempty"`
    Min      *float64     `yaml:"min,omitempty"`
    Max      *float64     `yaml:"max,omitempty"`
    Reverse  bool         `yaml:"reverse,omitempty"`
    Zones    []ZoneConfig `yaml:"zones,omitempty"`
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

Widget type to asset family mapping:

| Widget type | Required asset family |
|---|---|
| `image` | `assets.image_sets` |
| `digit_display` | `assets.digit_sets` |
| `bar_display` | `assets.bar_sets` |
| `frame_gauge` | `assets.frame_sets` |
| `indicator` | `assets.indicator_sets` |

Notes:

- Use `position`, not `rect`, for native-size image-backed widgets.
- Asset packs own native visual geometry.
- Renderers should validate that related images in an asset pack have compatible dimensions.
- Keep expressions, scripts, inheritance, and clever conditions out of v3 until genuinely needed.
- If both widget `Min` and `Max` are present, require `Min < Max`.

## 12. Bar widget semantics

Bar widgets map one sensor value onto repeated cells.

Rules:

- `cells` is the number of visible cells and must be greater than zero.
- `min` and `max` define the value mapping range.
- If both `min` and `max` are present, require `min < max`.
- Values below widget `min` render zero filled cells.
- Values above widget `max` render all cells filled.
- `off` is required for unfilled cells.
- If `zones` is omitted, `on` is required for filled cells.
- Bar zones must be sorted ascending by `up_to`.
- A filled cell uses the first zone where `value <= up_to`.
- Values above the final zone use the final zone.
- `reverse` reverses fill direction only, not zone interpretation.

## 13. Validation rules

Initial validation should check:

- Unknown fields fail at all documented levels.
- IDs match `^[a-z][a-z0-9_]*$`.
- At least one vehicle exists.
- If multiple vehicles exist, runtime vehicle selection is explicit.
- Vehicle endpoint address is present.
- `serial://` endpoints include a non-empty path.
- `tcp://` endpoints include host and port.
- OBD timeout is greater than zero, preferably within `100..30000` milliseconds.
- Vehicle `logs` references exist under top-level `logs`.
- Vehicle `dashboards` references exist under top-level `dashboards`.
- If multiple logs are defined, each vehicle lists the logs it runs.
- If multiple dashboards are defined, each vehicle lists the dashboards it renders.
- For each selected vehicle, no two selected dashboards target the same display.
- Sensor `poll > 0`.
- Sensor `min < max` when both are present.
- Sensor IDs referenced by logs exist.
- Sensor IDs referenced by widgets exist, except `type: image` widgets with no sensor.
- Asset references exist in the correct asset family for the widget type.
- Asset paths are repository-root relative.
- Digit sets contain required displayed characters for configured formats where practical.
- Digit formats that can emit decimal separators have a `decimal_point` asset.
- Digit formatted output fits configured slot count after decimal separators are removed.
- Frame ranges have `first <= last`.
- Bar sets contain `off`.
- Bar widgets without `zones` use a bar set containing `on`.
- Bar zones are sorted ascending.
- Bar zones reference valid cell names.
- Indicator sets contain `off`, `on`, and `unknown`.
- Dashboard sizes are positive.
- Widget positions are present.
- Widget IDs are unique within each dashboard.
- Widget `min < max` when both are present.

## 14. Non-goals

Do not add these without strong evidence:

- plugin systems
- source orchestration
- live config reload
- dashboard scripting
- config inheritance
- generic event buses
- enable flags everywhere
- endpoint-type branches leaking into core runtime
- dashboard refresh fields
- widget refresh fields
- stale timeout YAML fields
- asset root config fields

Boring boundaries first. Fancy pixels second. YAML goblins never.

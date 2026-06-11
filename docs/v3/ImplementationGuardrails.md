# GoDriveLog v3 implementation guardrails

Status: implementation guidance  
Applies to: v3 config/runtime/dashboard work  
References: `config.example.yaml`, `config.full.yaml`, `GoStructsConfig.md`, `DirectoryStructure.md`

## 1. Purpose

These guardrails exist to keep v3 implementation smooth, boring, and aligned with the docs.

When implementation and docs disagree, stop and resolve the disagreement before adding compatibility glue. Compatibility glue is how small tools turn into haunted furniture.

## 2. Core rule

Implement this pipeline first:

```text
vehicle endpoint
-> sensor polling runtime
-> sensor events
-> logs and dashboards as subscribers
```

Do not build sideways features until that path works end-to-end.

## 3. Config boundaries

The only intended v3 top-level config sections are:

```yaml
vehicles:
sensors:
assets:
logs:
dashboards:
```

Do not reintroduce:

```text
default_vehicle
active_displays
displays
log
cache
refresh
refresh_ms
render_min_ms
source: mock
source: real
mock_mode
```

If one of those seems necessary, the design needs review before code changes.

## 4. Vehicle endpoint guardrails

Vehicles own endpoint configuration only:

```yaml
vehicles:
  vw_caddy:
    name: "VW Caddy"
    obd:
      address: "serial:///dev/ttyUSB0"
      timeout: 1000
```

Rules:

- Treat real hardware and bench simulators as OBD-like endpoints.
- Use `serial://...` for serial adapters.
- Use `tcp://...` for simulator/bench endpoints.
- Do not branch the core runtime on mock versus real.
- Do not leak simulator concepts into sensors, logs, or dashboards.

Implementation shape:

```text
address string -> endpoint connector -> reader/runtime
```

Bad shape:

```text
if config.Source == "mock" { ... } else { ... }
```

## 5. Sensor runtime guardrails

Sensors own polling.

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

Rules:

- `poll` is milliseconds.
- Reject `poll <= 0`.
- Keep sensor definitions global.
- Do not put sensor definitions under vehicles.
- Do not let logs or dashboards request direct OBD reads.
- Do not create one polling loop per consumer.
- Do not use `0` as an error value.

The polling runtime should maintain latest state per sensor and emit events on:

```text
first reading
value change
status change
recovery
stale/error/unsupported transition
```

A sensor event should preserve the original read timestamp.

## 6. Sensor status guardrails

Use explicit status. Do not smuggle status into values.

Recommended statuses:

```text
ok
stale
error
missing/unsupported
```

Rules:

- `ok` means the value is usable.
- `stale` means last value exists but is too old.
- `error` means a read failed.
- `missing/unsupported` means the sensor cannot currently be provided.
- Dashboards must not silently display stale/error/missing values as live values.
- Indicators should prefer `unknown` display state when sensor status is not `ok`.

## 7. Log guardrails

Logs are subscribers.

```yaml
logs:
  jsonl:
    path: "logs/godrivelog.jsonl"
    sensors:
      - speed
      - rpm
```

Rules:

- Logs subscribe to sensor events.
- Logs do not poll.
- Logs do not own refresh cadence.
- Logs should write first reading, value changes, and status changes.
- Logs should not spam unchanged duplicate readings.
- Logs should include the sensor read timestamp.
- Log writer timestamp may also be recorded, but it is not a substitute for sensor read timestamp.

## 8. Dashboard guardrails

Dashboards are subscribers and renderers.

```yaml
dashboards:
  primary:
    display: "HDMI-1"
    size:
      width: 1920
      height: 480
    widgets: []
```

Rules:

- Dashboard presence means active.
- A dashboard owns its display target.
- Do not add top-level display bindings.
- Dashboards do not poll sensors.
- Dashboards do not own refresh cadence in config.
- Dashboards consume current sensor state produced by the sensor runtime.
- Keep dashboard config declarative.
- Avoid conditions, scripts, formulas, templates, inheritance, and expression languages.

If a visual behaviour needs code, put it in a widget implementation, not a YAML mini-language.

## 9. Asset guardrails

The asset model is descriptive, not procedural.

Allowed asset families:

```yaml
assets:
  digit_sets:
  bar_sets:
  frame_sets:
  indicator_sets:
  image_sets:
```

Common render order:

```text
background, if present
value/state-driven dynamic layer
foreground, if present
```

Rules:

- Assets describe images.
- Widgets decide how to map sensor state to rendered content.
- Do not put rules, conditions, formulas, or scripts inside assets.
- Validate asset references at startup.
- Validate related image dimensions where the renderer requires alignment.
- Missing required assets should be a clear config/asset error, not a nil panic with jazz hands.

## 10. Digit display guardrails

Digit displays render formatted strings as characters.

```yaml
characters:
  "0": "..."
  "1": "..."
  "-": "..."
decimal_point: "..."
```

Rules:

- Use `characters`, not `digits`, in v3 config.
- Support `-` as a first-class character.
- Treat blank/padded slots as background-only when background exists.
- Decimal point is an overlay.
- Renderer should report a useful error when a formatted character has no asset.

Example error:

```text
digit set bttf_amber_digits cannot render character "E" for widget speed_digits
```

## 11. Bar display guardrails

Bar widgets map one sensor value onto repeated cells.

Rules:

- `cells` is the number of visible cells.
- `min` and `max` define the value mapping range.
- `zones`, if present, select cell image names by range.
- `reverse`, if present, reverses fill direction only.
- Do not create curved bar geometry in YAML.
- Use `frame_gauge` for fancy curved/sweeping visuals.

## 12. Frame gauge guardrails

Frame gauges map one sensor value onto a frame sequence.

Rules:

- Use frame sets for complex visuals.
- Keep frame selection deterministic.
- Clamp values outside min/max unless a later design explicitly says otherwise.
- Do not build a vector drawing language into config.

## 13. Indicator guardrails

Indicators map boolean/status data onto image states.

Required states:

```text
off
on
unknown
```

Runtime mapping:

```text
if sensor status != ok:
  state = unknown
else if value == true:
  state = on
else:
  state = off
```

Rules:

- Do not show `off` for missing/error/stale data.
- Missing `unknown` asset should be a validation error unless a deliberate fallback is documented.
- Sensor values should remain boolean, not UI-state strings.

## 14. Validation guardrails

Validation should fail early and loudly.

Minimum checks:

- At least one vehicle exists.
- Multiple vehicles require explicit runtime selection.
- Vehicle endpoint address is present.
- Sensor IDs are unique by map key.
- `poll > 0` for every sensor.
- Logs reference existing sensors.
- Dashboards have positive size.
- Widgets have IDs, types, assets, and positions.
- Non-image widgets reference existing sensors.
- Widget asset references exist.
- Indicator assets contain `off`, `on`, and `unknown`.
- Frame set ranges have `first <= last`.
- Bar widget `cells > 0`.
- Bar zones reference valid cell names.
- Digit displays can render expected configured characters where practical.

Unknown fields should fail validation during development. Silent config typos are tiny assassins.

## 15. Implementation order

Recommended order:

1. Config structs and strict loading.
2. Config validation.
3. Vehicle endpoint abstraction.
4. Sensor event type and state store.
5. Sensor polling runtime.
6. JSONL log subscriber.
7. Asset registry and asset validation.
8. `image` widget.
9. `digit_display` widget.
10. `indicator` widget.
11. `bar_display` widget.
12. `frame_gauge` widget.
13. Dashboard renderer integration.
14. Simulator endpoint.

Do not start with the fancy renderer. That is dessert. Eat the vegetables first.

## 16. Testing guardrails

Prefer small tests that prove boundaries.

Config tests:

- valid minimal config
- valid full config
- missing vehicles
- multiple vehicles without explicit selection
- bad OBD address
- poll zero
- log references unknown sensor
- dashboard widget references unknown sensor
- dashboard widget references unknown asset
- digit display missing formatted character
- indicator set missing unknown state

Runtime tests:

- sensor first reading emits event
- unchanged value does not spam logs
- value change emits event
- status change emits event
- stale transition emits event
- logger preserves read timestamp
- dashboard receives state without polling endpoint

Asset tests:

- digit set dimension mismatch reports useful error
- frame range invalid reports useful error
- bar zone unknown cell reports useful error

## 17. Refusal rules for future complexity

Say no, or at least not yet, to:

- plugin systems
- generic event buses
- config inheritance
- dashboard scripting
- YAML formulas
- live config reload
- source orchestration
- enable flags everywhere
- per-log or per-dashboard polling knobs
- mock/real branches in core runtime

These may become real requirements later. They are not starting requirements.

## 18. Definition of done for first v3 implementation slice

A first useful v3 slice is done when:

- a minimal v3 config loads strictly
- a selected vehicle endpoint connects
- configured sensors poll on their own cadence
- sensor events update state
- JSONL logs receive selected events
- one dashboard displays at least image + digit_display + indicator
- missing/stale/error states are visible instead of silently lying
- old v2 config keys are rejected in v3 mode

That is enough. Anything beyond that should earn its keep.

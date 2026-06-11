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

Guardrails:

- Treat those five sections as the complete v3 root schema.
- Unknown fields should fail validation at every documented level during v3 development.
- Do not add compatibility aliases for undocumented fields.
- Do not add timing knobs outside `sensors.<id>.poll` unless the design is reviewed first.
- Do not add endpoint-type switches when an endpoint address can express the same thing.

The goal is an allow-list, not a blacklist. The schema should say what is valid, not make imaginary alternatives sound official.

## 4. Naming guardrails

IDs should match:

```text
^[a-z][a-z0-9_]*$
```

Apply this to vehicle IDs, sensor IDs, asset IDs, log IDs, dashboard IDs, and widget IDs.

Rules:

- Use lowercase snake_case.
- Widget IDs must be unique within a dashboard.
- Asset IDs only need to be unique within their own asset family.
- Do not use human display names as IDs.

## 5. Vehicle endpoint guardrails

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
- `serial://` endpoints must include a non-empty path.
- `tcp://` endpoints must include host and port.
- `timeout` is milliseconds and must be greater than zero.
- Initial timeout sanity range is `100..30000` milliseconds.
- Do not branch the core runtime on endpoint type unless a real implementation constraint proves it is needed.
- Do not leak simulator concepts into sensors, logs, or dashboards.

Implementation shape:

```text
address string -> endpoint connector -> reader/runtime
```

Bad shape:

```text
switch config.ProviderKind { ... }
```

## 6. Sensor runtime guardrails

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
- If `min` and `max` are both present, reject `min >= max`.
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

## 7. Sensor status and stale guardrails

Use explicit status. Do not smuggle status into values.

Recommended statuses:

```text
ok
stale
error
missing/unsupported
```

Initial stale rule:

```text
stale_after = max(sensor.poll * 3, 1000ms)
```

Rules:

- `ok` means the value is usable.
- `stale` means last value exists but is too old.
- `error` means a read failed.
- `missing/unsupported` means the sensor cannot currently be provided.
- Stale timing is runtime-derived.
- Do not add YAML stale timing fields unless reviewed later.
- Dashboards must not silently display stale/error/missing values as live values.
- Indicators should prefer `unknown` display state when sensor status is not `ok`.
- Stale transitions and recovery transitions must emit sensor events.

## 8. Log guardrails

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
- Logs do not own cadence.
- Logs should write first reading, value changes, and status changes.
- Logs should not spam unchanged duplicate readings.
- Logs should include the sensor read timestamp.
- Log writer timestamp may also be recorded, but it is not a substitute for sensor read timestamp.

## 9. Dashboard guardrails

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
- In the initial v3 implementation, two active dashboards must not target the same display.
- Multiple physical regions on one display should be widgets inside one dashboard.
- Dashboards do not poll sensors.
- Dashboards do not own cadence in config.
- Dashboards consume current sensor state produced by the sensor runtime.
- Keep dashboard config declarative.
- Avoid conditions, scripts, formulas, templates, inheritance, and expression languages.

If a visual behaviour needs code, put it in a widget implementation, not a YAML mini-language.

## 10. Asset guardrails

The asset model is descriptive, not procedural.

Asset paths in v3 config are repository-root relative.

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

Widget type to asset family mapping:

| Widget type | Required asset family |
|---|---|
| `image` | `assets.image_sets` |
| `digit_display` | `assets.digit_sets` |
| `bar_display` | `assets.bar_sets` |
| `frame_gauge` | `assets.frame_sets` |
| `indicator` | `assets.indicator_sets` |

## 11. Digit display guardrails

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
- `digits` is the number of character slots, excluding decimal point overlays.
- Formatted decimal separators do not consume a character slot.
- If the configured format can emit a decimal separator, `decimal_point` is required.
- Formatted output must fit the configured slot count after decimal separators are removed.
- Renderer should report a useful error when a formatted non-decimal character has no asset.

Example error:

```text
digit set bttf_amber_digits cannot render character "E" for widget speed_digits
```

## 12. Bar display guardrails

Bar widgets map one sensor value onto repeated cells.

Rules:

- `cells` is the number of visible cells and must be greater than zero.
- `min` and `max` define the value mapping range.
- If `min` and `max` are both present, reject `min >= max`.
- Values below widget `min` render zero filled cells.
- Values above widget `max` render all cells filled.
- `off` is required for unfilled cells.
- If `zones` is omitted, `on` is required for filled cells.
- `zones`, if present, select cell image names by range.
- Bar zones must be sorted ascending by `up_to`.
- A filled cell uses the first zone where `value <= up_to`.
- Values above the final zone use the final zone.
- `reverse`, if present, reverses fill direction only.
- `reverse` does not change zone interpretation.
- Do not create curved bar geometry in YAML.
- Use `frame_gauge` for fancy curved/sweeping visuals.

## 13. Frame gauge guardrails

Frame gauges map one sensor value onto a frame sequence.

Rules:

- Use frame sets for complex visuals.
- Keep frame selection deterministic.
- Frame ranges require `first <= last`.
- Clamp values outside min/max unless a later design explicitly says otherwise.
- Do not build a vector drawing language into config.

## 14. Indicator guardrails

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

## 15. Validation guardrails

Validation should fail early and loudly.

Minimum checks:

- Root config contains only documented v3 sections.
- Unknown fields fail at every documented level.
- IDs match `^[a-z][a-z0-9_]*$`.
- At least one vehicle exists.
- Multiple vehicles require explicit runtime selection.
- Vehicle endpoint address is present.
- `serial://` endpoints include a non-empty path.
- `tcp://` endpoints include host and port.
- `timeout > 0`, preferably within `100..30000` milliseconds.
- Sensor IDs are unique by map key.
- `poll > 0` for every sensor.
- Sensor `min < max` when both are present.
- Logs reference existing sensors.
- Dashboards have positive size.
- No two active dashboards target the same display.
- Widgets have IDs, types, assets, and positions.
- Widget IDs are unique within each dashboard.
- Non-image widgets reference existing sensors.
- Widget asset references exist in the correct asset family for the widget type.
- Widget `min < max` when both are present.
- Indicator assets contain `off`, `on`, and `unknown`.
- Frame set ranges have `first <= last`.
- Bar widget `cells > 0`.
- Bar sets contain `off`.
- Bar widgets without zones use a bar set containing `on`.
- Bar zones are sorted ascending.
- Bar zones reference valid cell names.
- Digit displays can render expected configured characters where practical.
- Decimal-capable digit formats have `decimal_point`.
- Digit formatted output fits configured slot count after decimal separators are removed.

Silent config typos are tiny assassins.

## 16. Implementation order

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

## 17. Testing guardrails

Prefer small tests that prove boundaries.

Config tests:

- valid minimal config
- valid full config
- all standalone v3 examples validate
- missing vehicles
- multiple vehicles without explicit runtime selection
- bad OBD address
- timeout zero
- poll zero
- sensor min >= max
- widget min >= max
- unknown root field
- nested unknown field
- log references unknown sensor
- dashboard widget references unknown sensor
- dashboard widget references unknown asset
- duplicate widget IDs fail
- two dashboards on same display fail
- digit display missing formatted character
- decimal format without `decimal_point` fails
- indicator set missing unknown state
- bar set missing `off` fails
- bar widget without zones requires `on`
- unsorted bar zones fail

Runtime tests:

- sensor first reading emits event
- unchanged value does not spam logs
- value change emits event
- status change emits event
- stale transition emits event
- recovery transition emits event
- logger preserves read timestamp
- dashboard receives state without polling endpoint

Asset tests:

- digit set dimension mismatch reports useful error
- frame range invalid reports useful error
- bar zone unknown cell reports useful error
- repository-root relative asset paths resolve consistently

## 18. Refusal rules for future complexity

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
- endpoint-type branches in core runtime
- stale timeout YAML fields
- asset root config fields

These may become real requirements later. They are not starting requirements.

## 19. Definition of done for first v3 implementation slice

A first useful v3 slice is done when:

- a minimal v3 config loads strictly
- all active v3 examples validate against the same schema rules
- a selected vehicle endpoint connects
- configured sensors poll on their own cadence
- sensor events update state
- stale/error/recovery transitions are visible as status changes
- JSONL logs receive selected events
- one dashboard displays at least image + digit_display + indicator
- missing/stale/error states are visible instead of silently lying
- undocumented config fields are rejected in v3 mode

That is enough. Anything beyond that should earn its keep.

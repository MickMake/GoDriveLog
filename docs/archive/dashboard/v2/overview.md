# Dashboard Implementation Overview - v2.x.x series

## Purpose

This document describes the staged implementation plan for the GoDriveLog dashboard rewrite across the v2.x.x series.

The core goal is:

> Isolate coding from visuals.

The old model made sensor/PID configuration own its visual presentation. The new model makes sensor data available as state, and lets a configurable dashboard scene decide how to render that state using assets, decoders, reusable blocks, layers, and conditions.

No legacy compatibility is required. GoDriveLog is not in production, so the old dashboard model was removed once the new scene renderer could render a real dashboard.

Tiny note from the goblin department: do not build two dashboard engines and hope they remain friends. They will not.

---

## Current architecture result

The app now uses the dashboard v2 scene path. Sensor configuration produces runtime state; dashboard configuration describes how that state becomes visuals.

The model:

```text
sensor state -> decoders -> visual blocks -> layered scene
```

Example target shape:

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

dashboard:
  canvas:
    width: 800
    height: 480

  assets:
    - id: background
      type: image
      path: background.svg

  decoders:
    - id: rpm_text
      type: format_number
      sensor: rpm
      format: "0000"

  blocks:
    - id: background_panel
      type: image
      asset: background

  layers:
    - id: base
      z: 0
      blocks:
        - background_panel
```

---

## Version map

| Stage | Version series | Name | Main result |
|---:|---|---|---|
| 1 | v2.0.x | New config schema | Sensors separated from dashboard visuals |
| 2 | v2.1.x | Config validation only | Dashboard config loads and validates without rendering |
| 3 | v2.2.x | Sensor state boundary | Dashboard reads neutral state, not PID config |
| 4 | v2.3.x | Decoder engine | Configured values convert to states, text, frame indexes |
| 5 | v2.4.x | Asset registry | Images, frame sets, and charsets load and validate |
| 6 | v2.5.x | Scene primitives | Image, sprite frame, sprite text, groups, conditions, z-order |
| 7 | v2.6.x | Fyne scene renderer | Configured scene renders instead of old panels |
| 8 | v2.7.x | First real dashboard | Background, RPM digits, bar frame, redline glow |
| 9 | v2.8.x | Remove old widgets | Old widget package tree removed |
| 10 | v2.9.x | Reusable block library | Practical dashboard building blocks shipped |

---

## Stage 1 - v2.0.x - New config schema

### Goal

Separate sensor definition from dashboard visual definition.

### Work

- Introduce a new top-level `sensors` section.
- Introduce a new top-level `dashboard` section.
- Remove visual ownership from sensor/PID config.
- Keep OBD, mock mode, logging, vehicle metadata, PID polling, and existing sensor reading behaviour.
- Do not render the new dashboard yet.

### Acceptance

- Config loads with `sensors` and `dashboard`.
- Sensor polling/logging does not require visual config.
- Active polling/logging still works.
- Tests prove config loading and sensor extraction.

### Do not

- Do not build dashboard rendering yet.
- Do not keep legacy visual config compatibility.
- Do not add visual implementation names to the new sensor schema.

---

## Stage 2 - v2.1.x - Config validation only

### Goal

Load and validate the full dashboard schema without drawing it.

### Add schema areas

```yaml
dashboard:
  canvas:
    width: 800
    height: 480

assets: {}

decoders: {}

blocks: {}

layers: []
```

### Work

- Add structs for dashboard config.
- Add validation for IDs, references, duplicate names, empty layers, invalid types, missing canvas sizes, and impossible geometry.
- Validate that layer references point to known assets, decoders, or blocks.
- Make errors human-readable.

### Acceptance

- Bad dashboard config fails early.
- Good dashboard config passes.
- App can start with a valid dashboard config even if renderer ignores it temporarily.

### Do not

- Do not implement asset loading yet.
- Do not render yet.
- Do not quietly ignore invalid dashboard sections.

---

## Stage 3 - v2.2.x - Sensor state boundary

### Goal

Create a clean runtime boundary between data collection and display rendering.

### Model

```go
type SensorState struct {
    ID        string
    Value     float64
    Unit      string
    Min       float64
    Max       float64
    Status    string
    Error     string
    UpdatedAt time.Time
}
```

### Work

- Introduce a `StateStore` or equivalent.
- Polling writes latest readings into state.
- Logging still writes readings.
- Dashboard reads from state.
- Errors/stale status are represented in state.
- Mock mode and real OBD mode continue working.

### Acceptance

- Dashboard code no longer needs `PIDConfig`.
- State can answer current value, unit, stale/error status, and min/max.
- Tests can update fake sensor values and inspect state.

### Do not

- Do not make the dashboard poll sensors directly.
- Do not pass raw PID config into the new dashboard.

---

## Stage 4 - v2.3.x - Decoder engine

### Goal

Implement reusable configured decoders.

### First decoder types

```text
normalize
threshold
frame_index
format_number
digits
boolean
```

### Work

- Add decoder registry.
- Evaluate decoders from state.
- Support named decoder outputs.
- Support simple comparison operators.
- Support clean errors for bad input references.

### Acceptance

- Decoders can be unit-tested without Fyne.
- Fake state values produce expected decoded outputs.
- Decoder outputs can be used by future scene elements.

### Do not

- Do not implement a full programming language.
- Do not use arbitrary `eval`.
- Do not bind decoders to visual widgets.

---

## Stage 5 - v2.4.x - Asset registry

### Goal

Load, cache, and validate visual assets.

### Asset types

```text
image
frame_set
charset
```

Potential later:

```text
nine_slice
font
animation
```

### Work

- Resolve asset paths relative to config location or project asset root.
- Validate missing files.
- Validate complete frame ranges.
- Validate charsets.
- Cache loaded images.
- Surface asset load errors clearly.

### Acceptance

- Missing asset fails at startup.
- Frame set validates expected frame count.
- Charset validates requested characters exist.

### Do not

- Do not silently replace bad images with blank placeholders.
- Do not load assets repeatedly every frame.

---

## Stage 6 - v2.5.x - Scene primitives

### Goal

Implement the small set of visual primitives that can fake almost anything.

### Primitive types

```text
image
sprite_frame
sprite_text
group
condition
z-order
```

### Work

- Implement scene element model.
- Implement condition evaluation.
- Implement z-order.
- Implement coordinate placement.
- Implement basic visibility toggling.

### Acceptance

- A static background can display.
- Conditional image can appear/disappear.
- Sprite text can render digits from a value.
- Sprite frame can select an image from a decoded frame index.

### Do not

- Do not add gauges.
- Do not hardcode RPM, speed, throttle, or coolant.
- Do not add animation beyond frame selection unless required for the vertical slice.

---

## Stage 7 - v2.6.x - Fyne scene renderer

### Goal

Render configured scene elements in Fyne.

### Work

- Replace panel-specific Fyne layout with a scene renderer.
- Render layers in z-order.
- Use asset registry images.
- Update scene on state changes.
- Keep the renderer dumb: it draws resolved visual elements.

### Acceptance

- Fyne window shows a scene from config.
- Sensor updates change displayed sprite text/frame images.
- Conditions update visibility.
- No old panel dashboard is required for the new scene.

### Do not

- Do not make the Fyne renderer responsible for decoding logic.
- Do not make scene elements query OBD.
- Do not add old widget compatibility.

---

## Stage 8 - v2.7.x - First real dashboard

### Goal

Deliver one real asset-driven dashboard proving the whole architecture.

### Minimum dashboard

```text
static background
RPM seven-segment sprite text
throttle sprite-frame bar
redline glow overlay
status/error indicator
```

### Acceptance

- Start app in mock mode.
- Background appears.
- RPM digits update.
- Throttle bar frame updates.
- Redline glow appears when threshold is crossed.
- Errors/stale state are visible.

---

## Stage 9 - v2.8.x - Remove old widgets

### Goal

Delete the old standalone widget package tree.

### Acceptance

- No old widget package tree remains.
- Build/tests pass.
- App boots only through new dashboard scene config.

### Do not

- Do not keep a compatibility shim.
- Do not leave dead widget names in docs.

---

## Stage 10 - v2.9.x - Reusable block library

### Goal

Ship useful building blocks so complex dashboards can be built from config and images.

### Initial block set

```text
seven_segment_number
percent_frame_bar
state_lamp
glowing_number_box
labelled_sensor_value
warning_overlay
stale_overlay
static_panel
```

### Acceptance

- Complex dashboard config can reuse blocks.
- Blocks are sensor-agnostic.
- Same block can render RPM, speed, voltage, or temperature depending on inputs.

---

## Deferred until after v2.9.x

```text
remote dashboard config
Google Drive dashboard sync
visual editor
SVG import
full expression language
plugin system
advanced animation
theme packs
dashboard marketplace, unless the goblins get funding
```

---

## Final architecture statement

Old:

```text
PID owns visuals
```

New:

```text
Dashboard scene binds to sensor state
```

That is the cut. Keep it clean.

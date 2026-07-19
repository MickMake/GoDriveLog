# Gauge Presets

Index: 15

Status: desired

Area: dashboard config, per-gauge config expansion, gauge assets, gauge realism, config loading, validation, reusable gauge profiles

Effort: 5-9 Codex hours

Add named reusable gauge presets for visual and realism configuration.

A gauge preset is a reusable config block applied before gauge-local config. It should make a gauge easier to describe without copying the same realism, asset, lighting, power, or imperfection settings into every gauge.

This is a preset/profile system, not text substitution.

## Concept

A preset captures reusable gauge character.

Examples:

```text
clean
old_analogue
worn_cable_speedo
cheap_led
weak_battery
bad_ground
gas_discharge
race_car_vibration
```

Project-specific presets may also exist:

```text
torana_speedo
vk_commodore_cluster
old_school_vdo
micks_dash_default
```

A gauge opts into a preset, then may override any preset values locally.

## Proposed config shape

Top-level preset definitions:

```yaml
gauge_presets:
  worn_cable_speedo:
    realism:
      imperfections:
        mechanical:
          eddy_speedo_wobble:
            enabled: true
            low_speed_wobble: 14
            high_speed_wobble: 2
            asymmetry: 70/30
            bias: under_read
            damping: 0.6

          idle_needle_vibration:
            enabled: true
            amplitude: 0.8
```

Gauge usage:

```yaml
gauges:
  speed:
    type: radial
    preset: worn_cable_speedo
    source: vehicle.speed
```

Gauge-local config overrides preset config:

```yaml
gauges:
  speed:
    type: radial
    preset: worn_cable_speedo
    source: vehicle.speed
    realism:
      imperfections:
        mechanical:
          eddy_speedo_wobble:
            low_speed_wobble: 8
```

Meaning:

```text
Use the worn cable speedo preset, but tone the wobble down for this gauge.
```

## Naming

Prefer `preset` over `macro` or `alias`.

Reasons:

- `macro` sounds like text substitution;
- `alias` sounds like another name for the same object;
- `preset` implies a reusable bundle of normal config values.

Use:

```yaml
preset: worn_cable_speedo
```

Not:

```yaml
macro: worn_cable_speedo
```

## Scope

Presets should be allowed to affect visual and display-only realism behaviour.

Safe preset areas:

```text
assets
realism
display-only realism quirks
realism.power
realism.lighting
realism.imperfections
animation behaviour
render styling
```

Presets should not silently affect data semantics.

Dangerous preset areas:

```text
source
sensor identity
unit conversion
logging
exports
configured value range
calibration
business logic
```

A preset should make a gauge look and behave like a worn cable speedo. It should not secretly change what speed means.

## Movement-related config

Movement-related display behaviour may eventually be presettable, but only through the current supported config keys for the target gauge family.

Do not use presets to introduce a new generic movement model.

For now, examples in this document intentionally use safer display-only quirk examples rather than movement-shaped config.

## Precedence

Use a predictable merge order.

Suggested order:

```text
built-in defaults
built-in preset
project preset
gauge-local config
runtime state
```

Gauge-local config wins over preset config.

Runtime state, such as power on/off or lights on/off, may affect final rendering but should not rewrite the loaded preset.

## One preset first

Start with one preset per gauge:

```yaml
gauges:
  speed:
    preset: worn_cable_speedo
```

Multiple composed presets may come later if needed:

```yaml
gauges:
  speed:
    presets:
      - old_analogue
      - worn_cable_speedo
```

Do not start with multi-preset composition unless the implementation needs it. It adds ordering and conflict complexity.

## Inline vs external presets

Support inline project presets first:

```yaml
gauge_presets:
  old_analogue:
    realism:
      damping: 0.6
      stiction: 0.2
      pointer_markers:
        average: true
```

Later, support external preset files:

```text
config/gauge_presets/worn_cable_speedo.yaml
config/gauge_presets/cheap_led.yaml
config/gauge_presets/old_analogue.yaml
```

The runtime can load built-in presets first, then project presets, then gauge-local overrides.

Project presets should be able to override built-in presets of the same name only if that behaviour is explicit and documented. Otherwise, reject duplicate preset names during validation.

## Built-in presets

Possible built-in presets:

```text
clean
old_analogue
worn_cable_speedo
cheap_led
weak_battery
bad_ground
gas_discharge
race_car_vibration
```

Built-in presets should be conservative. They should be useful starting points, not extreme novelty effects.

## Project presets

Project presets allow a dashboard to define its own named gauge styles.

Examples:

```yaml
gauge_presets:
  torana_speedo:
    realism:
      imperfections:
        preset: worn_cable_speedo

  micks_dash_default:
    realism:
      lighting:
        on:
          mode: dark
          asset_suffix: _dark
```

Project presets should be versioned with the dashboard config so that a dashboard remains reproducible.

## Presets referencing presets

A later implementation may allow a preset to extend another preset.

Example direction:

```yaml
gauge_presets:
  torana_speedo:
    extends: worn_cable_speedo
    realism:
      imperfections:
        mechanical:
          eddy_speedo_wobble:
            low_speed_wobble: 16
```

If this is implemented, guard against cycles:

```text
A extends B
B extends A
```

For the first slice, avoid preset inheritance unless needed.

## Validation rules

Validate presets before applying them to gauges.

Suggested checks:

- preset name is valid and unique;
- referenced preset exists;
- preset does not set forbidden fields;
- gauge-local config can override allowed preset fields;
- unsupported preset fields produce clear errors;
- external preset files are deterministic and loaded in a documented order;
- preset expansion can be printed or inspected for debugging.

## Debugging support

Add a way to inspect expanded gauge config.

Possible command/output direction:

```text
gatorlog config expand --gauge speed
```

or dashboard preview output:

```text
Gauge speed:
  preset: worn_cable_speedo
  expanded realism.imperfections.mechanical.eddy_speedo_wobble.low_speed_wobble: 8
```

The point is to make preset behaviour visible. Hidden config inheritance will otherwise become painful.

## Do not

- Do not implement presets as raw text substitution.
- Do not let presets silently change `source`.
- Do not let presets silently change units, logging, exports, value ranges, calibration, or sensor identity.
- Do not make multi-preset ordering complicated in the first slice.
- Do not allow cyclic preset inheritance.
- Do not hide the final expanded config from users.
- Do not use presets to introduce a generic movement model.

## Possible future slices

```text
gauge preset schema and validation
single preset application per gauge
preset override precedence
gauge config expansion/debug output
external gauge preset files
built-in conservative preset library
multi-preset composition
preset inheritance with cycle detection
```

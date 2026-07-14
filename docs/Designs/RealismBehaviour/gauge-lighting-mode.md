# Gauge Lighting Mode

Index: 13

Status: desired

Area: `gauge/radial`, `gauge/bar`, indicator, numeric, odometer, dashboard runtime, lights-state events, gauge realism, alternate asset sets

Effort: 4-7 Codex hours

Add gauge-level lighting-mode realism driven by a dashboard lights-state signal, such as vehicle lights on/off.

The dashboard runtime detects the external lights state. Each gauge owns how it changes appearance when that state changes.

This keeps gauges self-contained: a realistic gauge should know how it looks in daytime and how it looks when the vehicle lights are on.

## Boundary

Separate the feature into two layers:

```text
vehicle lights input -> dashboard lights signal provider -> runtime lights events -> gauge realism.lighting behaviour
```

The dashboard/runtime layer is responsible for detecting and debouncing a lights-state signal.

The gauge layer is responsible for presentation behaviour when the runtime broadcasts lights events.

Do not wire GPIO details directly into gauge rendering.

## Electrical note

Vehicle lighting circuits should not be treated as raw Raspberry Pi GPIO inputs.

The vehicle lights signal should be electrically isolated and level-shifted before the Pi sees it. The intended software design should allow a GPIO-backed provider, but the implementation docs should still warn that vehicle wiring is noisy and must not be connected directly to a Pi GPIO.

## Primary design: alternate asset set

Prefer a simple alternate asset-set convention over complex dynamic lighting code.

When lights are off, use the normal asset names.

When lights are on, try to use matching dark/illuminated assets first. If the dark asset does not exist, fall back to the normal asset.

Preferred suffix:

```text
_dark
```

Examples:

```text
face.png          -> face_dark.png
needle.png        -> needle_dark.png
foreground.png    -> foreground_dark.png
bar.png           -> bar_dark.png
fill.png          -> fill_dark.png
digits.png        -> digits_dark.png
segment.png       -> segment_dark.png
```

This keeps the first implementation mostly as deterministic asset selection. Gauge designers can make the illuminated look in PNGs instead of forcing the renderer to simulate lighting.

## Proposed dashboard input config

```yaml
dashboard:
  lights:
    input:
      type: gpio
      pin: 27
      active: high
      debounce_ms: 100
      stable_ms: 250
```

Supported input provider types should eventually include:

```text
gpio
manual
mock
none
```

`manual` and `mock` are useful for desktop preview, tests, and CI where GPIO is not available.

## Proposed gauge config

Gauge behaviour should live under gauge realism:

```yaml
gauges:
  rpm:
    realism:
      lighting:
        off:
          mode: normal

        on:
          mode: dark
          asset_suffix: _dark
```

Global defaults may exist, but per-gauge config should win.

```yaml
dashboard:
  lights:
    defaults:
      off:
        mode: normal
      on:
        mode: dark
        asset_suffix: _dark
```

Use dashboard defaults only as a convenience layer. The actual visual response belongs to the gauge.

## Modes

Start with simple named modes:

```text
normal
dark
```

`normal` uses the existing asset set.

`dark` tries the configured suffix asset set and falls back to the normal assets when a matching file is not present.

Avoid making the first implementation a full theme engine. This is gauge illumination realism, not general dashboard skinning.

## Asset lookup rules

For each asset used by a gauge:

1. If lights are off, use the configured normal asset.
2. If lights are on, look for the same asset name with the configured suffix before the extension.
3. If the suffixed asset exists, use it.
4. If the suffixed asset does not exist, use the normal asset.

Example:

```text
needle.png + _dark -> needle_dark.png
```

This should work for all ordinary image layers without requiring every gauge to declare every dark asset explicitly.

## Gauge-family examples

### Radial gauges

Normal assets:

```text
face.png
needle.png
foreground.png
```

Lights-on assets:

```text
face_dark.png
needle_dark.png
foreground_dark.png
```

Radial lighting should use the same geometry and pivot model. Only the selected images change.

### Bar gauges

Normal assets:

```text
bar.png
fill.png
foreground.png
```

Lights-on assets:

```text
bar_dark.png
fill_dark.png
foreground_dark.png
```

Bar lighting should keep the same value-to-extent mapping. Only the selected images change.

### Indicator gauges

Normal assets:

```text
indicator_on.png
indicator_off.png
```

Lights-on assets:

```text
indicator_on_dark.png
indicator_off_dark.png
```

Indicator lighting should not be confused with the indicator's own on/off signal. Lighting mode affects presentation, not the indicator's logical state.

### Numeric or segmented gauges

Normal assets:

```text
digits.png
segments.png
```

Lights-on assets:

```text
digits_dark.png
segments_dark.png
```

Numeric lighting may use alternate digit/segment assets. Do not change the displayed numeric source value.

### Odometer gauges

Normal assets:

```text
digits.png
wheel.png
```

Lights-on assets:

```text
digits_dark.png
wheel_dark.png
```

Odometer lighting should be conservative. Do not change odometer source values, wheel state, or digit state because the vehicle lights changed.

## Optional transition

A later slice may add a simple bounded fade between normal and dark asset sets.

```yaml
realism:
  lighting:
    on:
      mode: dark
      asset_suffix: _dark
      transition: fade
      duration_ms: 300
```

Keep this optional. The core feature is asset-set selection.

## Composition rules

Lighting mode is display-only.

It must not mutate:

```text
source values
logs
exports
configured ranges
stat marker samples
odometer source values
sensor state
```

Lighting should compose with power lifecycle:

- if the gauge is off, lights changes may update the pending asset set but should not force the gauge visibly on;
- if the gauge is powering on, the final settled appearance should reflect the current lights state;
- if the gauge is live, lights changes may switch between normal and dark asset sets;
- if the gauge is powering off, lights changes should not restart the gauge unless power state changes.

## Sampling rules

Other display features should not treat lighting changes as live signal data.

In particular:

- stat markers should ignore lighting-only changes;
- logs and exports should continue to reflect real source values;
- replay should be able to simulate lights events deterministically.

## Preview and test support

Add a mock/manual lights signal provider so lights-on and lights-off behaviour can be previewed without vehicle hardware.

Useful preview controls:

```text
lights on
lights off
toggle lights
replay scripted lights events
```

Possible mock script shape:

```yaml
dashboard:
  lights:
    input:
      type: mock
      script:
        - at_ms: 0
          state: off
        - at_ms: 2000
          state: on
        - at_ms: 7000
          state: off
```

## Do not

- Do not connect vehicle lighting wiring assumptions directly to gauge code.
- Do not require GPIO for desktop preview or tests.
- Do not change source values, logs, exports, stat marker samples, configured ranges, or odometer values.
- Do not build a complex lighting simulation in the first slice.
- Do not require every normal asset to have a dark variant.
- Do not make all gauges behave the same by dashboard fiat; each gauge should own its realistic lighting response.

## Possible future slices

```text
lights signal provider and runtime lights events
gauge realism.lighting config model
suffix-based dark asset lookup
radial gauge dark asset preview
bar gauge dark asset preview
indicator/numeric dark asset preview
optional fade transition between asset sets
```

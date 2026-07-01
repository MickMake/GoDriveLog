# Gauge Stat Markers

Index: 12

Status: desired

Area: `gauge/radial`, `gauge/bar`, renderer, rolling-window statistics, marker assets

Effort: 6-10 Codex hours

Add optional `realism.stat_markers` support for rolling-window minimum, maximum, and average markers.

Stat markers display separate gauge markers at statistical values calculated over a trailing time window. They are display-only markers, not source value changes, not fading peak-hold animation, and not a second live value indicator.

This feature should support radial and bar gauges. Radial and bar renderers may share rolling-window/statistics logic, but their drawing semantics are gauge-specific.

## Proposed config shape

```yaml
realism:
  stat_markers:
    window: 1h
    min: true
    max: true
    average: true
```

## Window semantics

`window` defines the trailing time range used to calculate enabled stat markers.

Examples:

- `window: 1h` means use stable displayed values from the last hour.
- `window: 24h` means use stable displayed values from the last day.
- `window: 0` means use all stable displayed values since runtime start.

Keep marker history bounded by the configured rolling window. Do not store unbounded history except for the explicit `window: 0` runtime-start case, and keep that case intentionally simple.

## Radial rendering

Radial gauges should use separate marker needle assets:

```text
needle_min.png
needle_max.png
needle_average.png
```

Each radial marker asset should use the same pivot/rotation geometry model as the live radial needle.

Render radial stat marker needles above the live needle and below existing foreground, overlay, or bezel layers.

## Bar rendering

Bar gauges should display stat markers as value-position markers on the bar fill/reveal axis.

Possible asset names:

```text
marker_min.png
marker_max.png
marker_average.png
```

If the bar renderer already has a clearer asset convention, prefer the existing bar-gauge naming style, but keep min/max/average distinct.

Bar stat markers should:

- use the same value-to-extent mapping as the live bar;
- support horizontal and vertical bars;
- respect bar origin/direction configuration;
- remain visually distinct from the live fill/reveal extent;
- render within the bar’s sensible visual bounds;
- render above the live bar fill/reveal layer and below foreground, overlay, or frame layers.

## Min/max behaviour

- `min: true` renders a min marker at the lowest stable displayed value inside the rolling window.
- `max: true` renders a max marker at the highest stable displayed value inside the rolling window.
- When a new higher maximum enters the window, move the max marker to that value.
- When a new lower minimum enters the window, move the min marker to that value.
- When the currently displayed min/max sample leaves the rolling window, fall back to the next valid min/max sample still inside the window, or hide that marker if none exists.

## Average behaviour

- `average: true` renders an average marker at the rolling average stable displayed value inside the configured window.
- Use a simple deterministic rolling average over the stable displayed values retained for the configured window.
- If the runtime has a clear timestamped sample model, bound retained samples by timestamp.
- If sample cadence questions arise, keep the implementation simple and document the chosen behaviour in the relevant implementation state or future implementation note.

## Sampling rules

- Track stable displayed values after normal value mapping and display-only calibration/offset behaviour for the gauge type.
- Do not capture temporary overshoot excursions as stat marker values.
- Do not capture temporary peg-bounce/end-stop bounce excursions as stat marker values.
- Stat markers should represent the stable display state, not transient animation tails.

## Rules

- Support radial and bar gauges.
- Disabled by default.
- Display-only.
- Preserve current gauge rendering when `realism.stat_markers` is absent or all markers are disabled.
- Keep source values, logs, exported values, configured ranges, and input data unchanged.
- Keep retained history bounded by the configured rolling window.
- Marker assets must remain visually distinguishable from the live needle, bar fill, or live value indicator.
- Add visual inspection fixtures that make min, max, and average markers easy to judge by eye.

## Do not

- Do not implement fade or decay behaviour.
- Do not implement trip or session-lifetime windows unless a later spec explicitly defines them.
- Do not use `needle_peak.png` for this feature.
- Do not place stat markers under `movement`.
- Do not mutate source values, logs, exports, configured ranges, or input data.
- Do not render markers over foreground, overlay, bezel, frame, or cover layers.

## Possible future slices

```text
radial stat markers min/max
radial stat marker average
bar stat markers min/max
bar stat marker average
```

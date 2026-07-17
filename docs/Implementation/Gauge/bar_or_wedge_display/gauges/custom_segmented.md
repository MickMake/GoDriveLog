# Custom segmented percent-threshold gauge implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `segmented` |
| New Gauge group | `bar_or_wedge_display` |
| Paired design | `docs/Designs/Gauge/bar_or_wedge_display/gauges/custom_segmented.md` |
| Runtime code impact | None |

## Current implementation model

The current GoDriveLog `segmented` implementation is a sparse percent-threshold image-selection gauge.

It is not a segmented character display. It selects a complete pre-rendered level image according to the current normalised percentage.

## Configuration shape

The current model uses a value layer pattern containing a `{percent}` placeholder. Matching filenames provide available threshold images.

Example pattern shape:

```yaml
layers:
  segments: levels/rpm_{percent:03}.png
```

Files such as `rpm_000.png`, `rpm_010.png`, and `rpm_030.png` become sparse display thresholds.

## Rendering approach

The renderer:

- discovers matching threshold filenames;
- ignores non-matching files;
- normalises the sensor value to `0..100`;
- selects the highest threshold reached by the current value;
- uses hysteresis based on the adjacent threshold gap;
- lazy-loads selected images rather than eagerly decoding every threshold image.

## Current limitations and boundaries

This implementation record deliberately keeps the old GoDriveLog name and the new Gauge group separate:

```text
Old GoDriveLog type: segmented
New Gauge group:    bar_or_wedge_display
```

It does not document seven-segment, dot-matrix, LCD, or individual segment rendering.


## Documentation boundary

This file records current GoDriveLog implementation behaviour only.

It does not:
- record implementation status;
- describe intended future work as implemented;
- rename runtime package types;
- replace or migrate existing documentation.

Implementation status belongs only in `docs/Status.md`.

## Historical source basis

- `docs/v3.4/ReleasePlan.md`
- `docs/v3.4/ImplementationState.md`
- `docs/v3.4/prompts/v3.4.0-gauge-type-docs.md`


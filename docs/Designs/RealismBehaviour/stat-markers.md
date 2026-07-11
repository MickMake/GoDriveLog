# Gauge Stat Markers

Status: superseded by `docs/v3.6/ReleasePlan.md`

This note was the first cut of the marker idea. It is retained only as historical context.

Do not implement this file as written.

The v3.6 feature is now **pointer markers**, not statistical markers. The current spec lives in:

```text
docs/v3.6/ReleasePlan.md
docs/v3.6/ImplementationState.md
docs/v3.6/prompts/
```

## What changed

The original version of this note described rolling-window statistical min, max, and average markers.

The v3.6 release instead implements physical-style pointer markers:

- the live rendered needle/bar pushes marker state;
- markers sample the final rendered indicator position;
- rendered overshoot, bounce, damping lag, stiction, and other visible movement can affect marker state;
- `average` is a highly damped physical pointer, not an arithmetic average;
- no `window` means daily local reset for min/max;
- `window` means rolling min/max history only.

## What survived into v3.6

The following ideas from this note remain valid and are now part of the v3.6 pointer marker spec:

- support radial and bar gauges;
- keep marker behaviour display-only;
- do not mutate source values, logs, exports, configured ranges, or input data;
- share marker logic where practical;
- keep rendering gauge-family specific;
- render markers above the live needle/bar and below overlay/glass/bezel/frame layers;
- use explicit marker PNG assets.

## Asset contract preserved

Radial gauges use explicit marker needle assets where provided:

```text
needle_min.png
needle_max.png
needle_average.png
```

Bar gauges use explicit marker assets where provided:

```text
marker_min.png
marker_max.png
marker_average.png
```

These asset names and the requirement that marker assets remain visually distinct from live indicators are preserved in the v3.6 pointer marker release.

## Do not use from this historical note

The following older ideas are superseded and must not be used for v3.6 pointer markers:

- `realism.stat_markers` as the config key;
- requiring `window` for all marker behaviour;
- calculating `average` as a rolling arithmetic average;
- ignoring rendered overshoot/bounce in favour of stable-only values;
- treating the feature as a statistics overlay.

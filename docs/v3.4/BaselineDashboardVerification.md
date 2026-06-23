# GoDriveLog v3.4 baseline dashboard verification

Status: planning placeholder

## Purpose

This file records how v3.4 gauge package changes should be checked against the reusable baseline dashboard.

The baseline dashboard remains:

```text
examples/baseline-dashboard.yaml
```

The v3.4 docs should not duplicate active runnable config unless a slice specifically needs a frozen test fixture.

## Current baseline workload

| Gauge | Type target | Sensor | Notes |
|---|---|---|---|
| Temperature | `numeric` | `coolant_temperature` | Existing old `seven_segment` package should become `numeric`; still exercises minus/format handling. |
| Speed | `numeric` | `speed` | Normal changing numeric display. |
| RPM numeric | `numeric` | `rpm` | High-frequency digit changes. |
| RPM radial | `radial` | `rpm` | Existing radial transform renderer remains valid. |

## Verification goals

- Existing numeric behaviour survives the hard rename from `seven_segment` to `numeric`.
- Existing radial transform behaviour is not disturbed.
- New transform gauge families (`odometer`, `bar`) are added without breaking the dashboard scene/display-sink boundary.
- New image-selection gauge families (`indicator`, `segmented`) are added without turning renderer-private state into gauge config.
- Renderer changes do not turn asset style into a code concern.

## Future verification additions

Add explicit examples and checks as slices land:

| Slice | Verification addition |
|---|---|
| v3.4.1 numeric rename | Run baseline dashboard with all numeric packages renamed. |
| v3.4.2 odometer | Add a harness-driven odometer example covering default `smooth` movement and optional `click` movement. |
| v3.4.3 indicator | Add off/on state example. |
| v3.4.4 bar | Add continuous transform example such as clipping/revealing a level layer. |
| v3.4.5 segmented | Add sparse percent-threshold image example, including missing-`000` no-layer behaviour and threshold-gap hysteresis. |

## Notes

Do not add renderer-private fake data for verification. Keep using the real v3 dashboard path. Tiny fake dashboards have a habit of becoming folk tales with bugs in them.

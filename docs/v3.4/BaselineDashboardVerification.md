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
| RPM radial | `radial` | `rpm` | Existing radial renderer remains valid. |

## Verification goals

- Existing numeric behaviour survives the hard rename from `seven_segment` to `numeric`.
- Existing radial behaviour is not disturbed.
- New gauge families are added without breaking the dashboard scene/display-sink boundary.
- Renderer changes do not turn asset style into a code concern.

## Future verification additions

Add explicit examples and checks as slices land:

| Slice | Verification addition |
|---|---|
| v3.4.1 numeric rename | Run baseline dashboard with all numeric packages renamed. |
| v3.4.2 odometer | Add a harness-driven odometer example. |
| v3.4.3 indicator | Add off/on state example. |
| v3.4.4 bar | Add continuous fill/reveal example. |
| v3.4.5 segmented | Add sparse percent-threshold image example, such as `rpm_000.png`, `rpm_010.png`, `rpm_030.png`. |

## Notes

Do not add renderer-private fake data for verification. Keep using the real v3 dashboard path. Tiny fake dashboards have a habit of becoming folk tales with bugs in them.

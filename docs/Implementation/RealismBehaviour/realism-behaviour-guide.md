# Realism Behaviour Guide

Design reference: [`docs/Designs/RealismBehaviour/realism-behaviour-guide.md`](../../Designs/RealismBehaviour/realism-behaviour-guide.md)

## Purpose
Tracks how much of the canonical realism guide currently exists in code.

## Implementation Status
Status: **Partially implemented**.

The guide is broader than the current implementation: many behaviours are implemented, many remain backlog, and the guide intentionally avoids per-feature status.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/gauges/pointer_markers.go`](../../../internal/dashboard/gauges/pointer_markers.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`
- `JSONLEventRecord` is unrelated here; runtime realism stays in dashboard packages only.
- `PointerMarkersConfig`
- `DampingConfig`
- `OvershootConfig`
- `ThermalFadeConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveMovementState`
- `resolveIndicatorThermalFadeState`
- `radialCalibrationAngle`

## Runtime Flow
Current runtime supports implemented realism for radial, bar, odometer, and indicator families through package validation plus dashboard movement state resolution.

## Configuration
Implemented keys include `wraparound`, `carry_drag`, `snap_settle`, `drum_slop`, `damping`, `hysteresis`, `stiction`, `overshoot`, `peg_bounce`, `pointer_markers`, `thermal_fade`, `needle_shadow`, `calibration_offset`, and `movement_policy`.

## Behaviour
Several guide behaviours are live, but others such as `backlash`, `needle_trail`, `lighting_mode`, `imperfections`, and most numeric/segmented candidate features remain unimplemented.

## Rendering
Implemented features are rendered by family-specific scene composition rather than one global realism layer.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/gauges/pointer_markers_test.go`](../../../internal/dashboard/gauges/pointer_markers_test.go)
- [`internal/dashboard/scenesink/latest_test.go`](../../../internal/dashboard/scenesink/latest_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
Because the guide is intentionally status-free, readers must use `docs/Status.md` and the per-feature implementation records for truth.

## Deviations from Design
The guide remains canonical design intent, not a promise that every described behaviour already exists.

## Remaining Work
Keep the implementation records and status register aligned as more guide behaviours are added or deferred.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

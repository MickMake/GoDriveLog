# `calibration_offset`

Design reference: [`docs/Designs/RealismBehaviour/calibration-offset.md`](../../Designs/RealismBehaviour/calibration-offset.md)

## Purpose
Tracks the fixed angular display offset for radial needles.

## Implementation Status
Status: **Implemented**.

Radial packages support `realism.calibration_offset`, and runtime rendering applies the offset without changing source values.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `radialCalibrationAngle`

## Runtime Flow
The runtime preserves source values and hands the final displayed state to radial scene building, where the configured offset shifts the displayed angle.

## Configuration
Radial packages accept `realism.calibration_offset` as a display-only degree offset.

## Behaviour
The needle can render slightly high or low while the underlying sensor value, value map, and stored state remain unchanged.

## Rendering
Offset is applied at radial scene composition time, including out-of-range handling based on the package value-map clamp rules.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
This is a fixed display offset only; there is no drift, noise, or dynamic recalibration model.

## Deviations from Design
The implementation matches the design intent closely.

## Remaining Work
No known design work remains beyond routine maintenance.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

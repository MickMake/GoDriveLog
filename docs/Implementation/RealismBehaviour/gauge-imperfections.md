# Gauge Imperfections

Design reference: [`docs/Designs/RealismBehaviour/gauge-imperfections.md`](../../Designs/RealismBehaviour/gauge-imperfections.md)

## Purpose
Tracks the broader backlog for visible gauge imperfections across multiple gauge families.

## Implementation Status
Status: **Partially implemented**.

Some specific imperfection-style features exist, but there is no umbrella implementation for the backlog described here.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`
- `NeedleShadowConfig`
- `ThermalFadeConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `radialCalibrationAngle`
- `resolveIndicatorThermalFadeState`

## Runtime Flow
Current realism support covers a subset of related behaviours such as calibration offset, needle shadow, and thermal fade, each as separate keys.

## Configuration
There is no shared `imperfections` group or cross-family backlog switch. Only individual implemented keys are accepted.

## Behaviour
GoDriveLog can already add a few deterministic display artefacts, but the larger wear/noise/ageing backlog is still absent.

## Rendering
Existing imperfections are rendered through family-specific scene logic, not a unified display-layer system.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The document describes a broad catalogue that is much larger than the current code surface.

## Deviations from Design
Implementation is fragmented into specific options rather than the broader umbrella backlog named here.

## Remaining Work
Decide whether to keep the umbrella document as backlog context or split it into concrete implementation slices.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

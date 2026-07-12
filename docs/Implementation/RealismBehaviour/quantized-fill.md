# `quantized_fill`

Design reference: [`docs/Designs/RealismBehaviour/quantized-fill.md`](../../Designs/RealismBehaviour/quantized-fill.md)

## Purpose
Tracks the planned discrete-resolution fill behaviour for bar and segmented displays.

## Implementation Status
Status: **Not implemented**.

Current bar rendering remains continuous, and there is no `quantized_fill` parser or runtime support.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- None in current code.

## Functions and Methods
- `validateRealismForGaugeFamily` rejects unsupported realism keys.

## Runtime Flow
Bar movement resolves continuous displayed extents rather than snapping them to discrete thresholds.

## Configuration
There is no `realism.quantized_fill` key in package loading.

## Behaviour
Displayed fill updates can be smooth or animated, but not quantised into fixed visible increments.

## Rendering
Bar scenes render the final reveal height directly from continuous geometry.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)

## Limitations
Segmented gauges exist, but this specific realism behaviour is still absent.

## Deviations from Design
The candidate design remains backlog only.

## Remaining Work
Add parser support, quantisation rules, and family-specific rendering if the feature stays desirable.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

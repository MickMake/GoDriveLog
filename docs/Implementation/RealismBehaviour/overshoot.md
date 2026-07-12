# `overshoot`

Design reference: [`docs/Designs/RealismBehaviour/overshoot.md`](../../Designs/RealismBehaviour/overshoot.md)

## Purpose
Tracks bounded pass-and-settle movement for radial and bar gauges.

## Implementation Status
Status: **Implemented**.

Radial and bar packages support `realism.overshoot`, and runtime movement resolves a bounded overshoot before settling at target.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`
- `OvershootConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveMovementState`

## Runtime Flow
The movement resolver can extend the displayed path beyond the target and then settle back while preserving the raw source value.

## Configuration
Radial and bar packages accept shared overshoot tuning, including thresholding and bounded travel parameters.

## Behaviour
Displayed movement can briefly pass the target then return, remaining finite and display-only.

## Rendering
Rendering uses the current interpolated displayed state; there is no separate overshoot art layer.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
Only radial and bar families implement overshoot.

## Deviations from Design
The implementation matches the documented feature well.

## Remaining Work
No design-specific work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

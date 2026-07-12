# `damping`

Design reference: [`docs/Designs/RealismBehaviour/damping.md`](../../Designs/RealismBehaviour/damping.md)

## Purpose
Tracks the lag-and-catch-up behaviour for radial and bar gauges.

## Implementation Status
Status: **Implemented**.

Damping is supported for radial and bar gauges, including directional timing for bars.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`
- `DampingConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveMovementState`
- `barDampingDuration`

## Runtime Flow
The movement resolver can retain a prior displayed position and advance it over time toward the target according to the configured damping policy.

## Configuration
Radial packages accept shared damping config. Bar packages also accept directional `rise_ms` and `fall_ms` timing.

## Behaviour
Displayed motion lags the source value and settles smoothly at the final position without altering stored raw input.

## Rendering
Rendering consumes the current interpolated displayed state; there is no special damping-only art layer.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
Only radial and bar families implement damping today.

## Deviations from Design
The implementation matches the documented family scope.

## Remaining Work
No design-specific work is required unless additional families adopt damping later.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

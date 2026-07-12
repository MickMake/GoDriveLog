# `hysteresis`

Design reference: [`docs/Designs/RealismBehaviour/hysteresis.md`](../../Designs/RealismBehaviour/hysteresis.md)

## Purpose
Tracks direction-dependent displayed offsets for radial and bar gauges.

## Implementation Status
Status: **Implemented**.

Hysteresis is supported for radial and bar gauges in both package validation and runtime movement resolution.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`
- `HysteresisConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveMovementState`

## Runtime Flow
The movement resolver applies rising or falling display bias while preserving the stored raw source value.

## Configuration
Radial and bar packages accept `realism.hysteresis` with family-appropriate tuning.

## Behaviour
The same source value can render slightly differently depending on approach direction, matching the display-only contract.

## Rendering
Rendering uses the resolved displayed state after hysteresis adjustments and optional clamp handling.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
Only radial and bar families implement hysteresis.

## Deviations from Design
The implementation aligns with the documented scope.

## Remaining Work
No design-specific work remains unless more families adopt the effect.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

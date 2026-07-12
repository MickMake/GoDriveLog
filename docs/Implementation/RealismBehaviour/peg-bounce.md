# `peg_bounce`

Design reference: [`docs/Designs/RealismBehaviour/peg-bounce.md`](../../Designs/RealismBehaviour/peg-bounce.md)

## Purpose
Tracks the tap-rebound-settle behaviour when radial or bar gauges hit display limits.

## Implementation Status
Status: **Implemented**.

Radial and bar packages support `realism.peg_bounce`, and runtime applies stop-hit rebound while keeping in-range targets immediate.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveMovementState`

## Runtime Flow
When the displayed path reaches a configured min or max stop with momentum, the movement resolver schedules a short rebound before settling at the limit.

## Configuration
Radial and bar packages accept `realism.peg_bounce` as a boolean display-only effect.

## Behaviour
The effect only appears when the displayed indicator actually hits the stop. In-range changes do not bounce.

## Rendering
Scene composition renders the resolved displayed position; the bounce is entirely a runtime movement effect.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The effect is intentionally subtle and only applies at display bounds.

## Deviations from Design
The implementation matches the reviewed contract, including the later fix that kept ordinary in-range movement immediate.

## Remaining Work
No design-specific work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

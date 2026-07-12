# Bar Gauge Overshoot Follow-Up

Design reference: [`docs/Designs/RealismBehaviour/bar-overshoot-follow-up.md`](../../Designs/RealismBehaviour/bar-overshoot-follow-up.md)

## Purpose
Records the follow-up idea for bar-gauge overshoot and how it now exists on `main`.

## Implementation Status
Status: **Implemented**.

Bar gauges now support `realism.overshoot` with bounded pass-and-settle behaviour in runtime and tests.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- `Realism`
- `OvershootConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveMovementState`

## Runtime Flow
Bar movement can schedule overshoot in the shared movement resolver before settling at the final displayed extent.

## Configuration
Bar packages accept `realism.overshoot` and the shared overshoot tuning keys validated in package loading.

## Behaviour
The displayed bar extent can briefly pass the target then settle back, matching the approved bounded display-only behaviour.

## Rendering
Bar rendering consumes the final resolved reveal height; the overshoot effect appears through runtime target interpolation, not separate art.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The follow-up note did not define extra bar-only tuning beyond the shared overshoot contract.

## Deviations from Design
This document describes something intentionally deferred from an earlier slice. The current code has since implemented it.

## Remaining Work
No design-specific code work remains beyond any future refinements to overshoot tuning.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

# `snap_settle`

Design reference: [`docs/Designs/RealismBehaviour/snap-settle.md`](../../Designs/RealismBehaviour/snap-settle.md)

## Purpose
Tracks the short landing snap for odometer wheels.

## Implementation Status
Status: **Implemented**.

Odometer packages support `realism.snap_settle`, and runtime movement adds a bounded settle tail before exact rest.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`

## Functions and Methods
- `validateRealismForGaugeFamily`

## Runtime Flow
After the main odometer movement reaches the target, the effect can add a short forward settle that returns exactly to the final digit position.

## Configuration
Odometer packages accept `realism.snap_settle` as a display-only effect.

## Behaviour
The wheel lands with a small mechanical-feeling snap rather than gliding silently into place.

## Rendering
Wheel-strip rendering uses the transient settle offset, then returns to the exact target offset.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
This effect is odometer-only.

## Deviations from Design
The implementation matches the design intent.

## Remaining Work
No known design work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

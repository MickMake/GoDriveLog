# `carry_drag`

Design reference: [`docs/Designs/RealismBehaviour/odometer-carry-drag.md`](../../Designs/RealismBehaviour/odometer-carry-drag.md)

## Purpose
Tracks the early-coupling movement of higher odometer digits near rollover.

## Implementation Status
Status: **Implemented**.

Odometer packages support `realism.carry_drag`, and wheel interpolation applies the effect near rollover.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`

## Functions and Methods
- `validateRealismForGaugeFamily`

## Runtime Flow
Odometer movement resolves the displayed wheel positions, then carry-drag logic advances the higher wheel slightly as the lower wheel approaches rollover.

## Configuration
Odometer packages accept `realism.carry_drag` as a display-only realism option.

## Behaviour
Higher wheels start to creep before the lower wheel fully rolls over, but still settle on the exact target value.

## Rendering
Wheel-strip rendering uses the adjusted fractional offsets so the effect appears as physical pre-load rather than digit substitution.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
Only odometer families implement this effect.

## Deviations from Design
The implementation matches the documented effect closely.

## Remaining Work
No known design work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

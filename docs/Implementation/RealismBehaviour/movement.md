# `movement`

Design reference: [`docs/Designs/RealismBehaviour/movement.md`](../../Designs/RealismBehaviour/movement.md)

## Purpose
Tracks the single movement knob and its family-specific behaviour across gauges.

## Implementation Status
Status: **Partially implemented**.

Movement support exists, but it is family-specific and incomplete relative to the guide.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`

## Functions and Methods
- `validateOdometerMovementMode`
- `resolveMovementState`

## Runtime Flow
Odometers implement `instant`, `linear`, `ease_out`, and `bell`. Radial and bar gauges still use `movement_policy` rather than the scalar `movement` key described here. Other families remain immediate.

## Configuration
The parser accepts odometer movement values, warns and falls back for reserved `smooth` and `click`, and separately accepts shared `movement_policy` for radial/bar movement runtime.

## Behaviour
Movement is not a single uniform feature yet. Odometer movement is explicit, radial/bar movement is policy-based, and several families still remain immediate.

## Rendering
Scene composition uses the currently resolved displayed state, with odometer wheel offsets or radial/bar positions driven by whichever movement path applies.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The documented scalar contract is not implemented consistently across families.

## Deviations from Design
Current code preserves older shared movement-policy behaviour for radial/bar gauges instead of the planned scalar key.

## Remaining Work
Either unify movement semantics across families or keep the guide explicit about the split contracts.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

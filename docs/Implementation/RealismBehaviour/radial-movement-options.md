# Radial Movement Options

Design reference: [`docs/Designs/RealismBehaviour/radial-movement-options.md`](../../Designs/RealismBehaviour/radial-movement-options.md)

## Purpose
Tracks the planned scalar movement options for radial gauges.

## Implementation Status
Status: **Partially implemented**.

Radial gauges animate through `movement_policy`, not through the scalar `movement` options proposed by this design.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`

## Functions and Methods
- `resolveMovementState`

## Runtime Flow
Radial movement can already be immediate, linear, bell-shaped, or ease-out through the shared movement planner used by `movement_policy`.

## Configuration
The implemented key is `realism.movement_policy`, not the scalar `movement` contract proposed here.

## Behaviour
The behaviour exists in substance, but the config surface and compatibility story differ from the design note.

## Rendering
Rendered needle motion follows the resolved movement plan regardless of which config key selected it.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
Users cannot configure radial movement through the proposed scalar key on `main`.

## Deviations from Design
The runtime behaviour is present, but the public config contract remains the older policy form.

## Remaining Work
Either adopt the scalar `movement` options for radial gauges or update the design to bless `movement_policy` as the long-term contract.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

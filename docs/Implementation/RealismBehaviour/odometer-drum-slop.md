# `drum_slop`

Design reference: [`docs/Designs/RealismBehaviour/odometer-drum-slop.md`](../../Designs/RealismBehaviour/odometer-drum-slop.md)

## Purpose
Tracks the fixed per-wheel alignment imperfection for odometers.

## Implementation Status
Status: **Implemented**.

Odometer packages support `realism.drum_slop`, and scene composition applies stable per-wheel offsets.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- `Realism`

## Functions and Methods
- `validateRealismForGaugeFamily`

## Runtime Flow
The runtime resolves the displayed odometer position and scene composition applies configured wheel offsets without changing the value.

## Configuration
Odometer packages accept `realism.drum_slop` with per-wheel alignment configuration.

## Behaviour
Each wheel can sit slightly high or low in a stable, deterministic way.

## Rendering
Wheel slices are shifted by the configured slop during odometer scene assembly.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)

## Limitations
This is static alignment only; it does not create movement jitter or backlash.

## Deviations from Design
The implementation matches the design intent.

## Remaining Work
No known design work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

# `uneven_brightness`

Design reference: [`docs/Designs/RealismBehaviour/uneven-brightness.md`](../../Designs/RealismBehaviour/uneven-brightness.md)

## Purpose
Tracks the planned stable slot/region brightness variation for numeric and segmented displays.

## Implementation Status
Status: **Not implemented**.

There is no `uneven_brightness` config or per-slot brightness model on `main`.

## Packages and Files
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
Runtime does not calculate per-slot brightness adjustments for numeric or segmented displays.

## Configuration
Package loading does not accept `realism.uneven_brightness`.

## Behaviour
All slots render at the same brightness unless the asset art itself differs.

## Rendering
No post-processing or per-slot alpha/brightness pass exists in scene composition.

## Tests
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)

## Limitations
A future implementation would need stable deterministic brightness offsets and likely family-specific asset assumptions.

## Deviations from Design
Still a candidate only.

## Remaining Work
Define the brightness model and rendering hooks before implementation.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

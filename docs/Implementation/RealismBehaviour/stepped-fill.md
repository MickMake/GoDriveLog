# `stepped_fill`

Design reference: [`docs/Designs/RealismBehaviour/stepped-fill.md`](../../Designs/RealismBehaviour/stepped-fill.md)

## Purpose
Tracks the planned block-style fill behaviour for bar and segmented displays.

## Implementation Status
Status: **Not implemented**.

Current bar rendering remains continuous, and segmented rendering does not expose `stepped_fill` as a realism option.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
Display state is resolved as continuous bar extent or current segmented display state, not a visible block-step progression keyed by this feature.

## Configuration
There is no `realism.stepped_fill` key in package loading.

## Behaviour
Displayed fill does not deliberately advance in discrete visible steps under a dedicated realism control.

## Rendering
Bar and segmented scenes do not switch into a stepped/block fill mode.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)

## Limitations
Segmented rendering exists, but not as the configurable stepped-fill behaviour described here.

## Deviations from Design
Still a candidate only.

## Remaining Work
Define the config and family-specific render model before implementation.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

# `load_sag`

Design reference: [`docs/Designs/RealismBehaviour/numeric-load-sag.md`](../../Designs/RealismBehaviour/numeric-load-sag.md)

## Purpose
Tracks the planned brightness sag effect for high-load numeric and segmented values.

## Implementation Status
Status: **Not implemented**.

There is no brightness model tied to lit-segment count or value-dependent display load.

## Packages and Files
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
Runtime does not calculate display electrical load or vary brightness by current value.

## Configuration
There is no `realism.load_sag` key in package loading.

## Behaviour
Values such as `888` and `111` render with the same brightness unless their art already differs.

## Rendering
Numeric and segmented scenes do not have a brightness pass that varies with active segment count.

## Tests
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)

## Limitations
The design would need per-family brightness control and perhaps per-segment awareness.

## Deviations from Design
Still a candidate only.

## Remaining Work
Define a stable brightness model and family-specific rendering hooks before implementation.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

# Gauge Stat Markers

Design reference: [`docs/Designs/RealismBehaviour/stat-markers.md`](../../Designs/RealismBehaviour/stat-markers.md)

## Purpose
Preserves the historical note for an older statistical marker concept that should not be implemented as written.

## Implementation Status
Status: **Not implemented**.

The document is superseded. Current code implements pointer markers, not the original statistical-marker idea.

## Packages and Files
- [`internal/dashboard/gauges/pointer_markers.go`](../../../internal/dashboard/gauges/pointer_markers.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `PointerMarkersConfig`

## Functions and Methods
- Pointer-marker update and render helpers exist, but there is no stat-marker feature keyed to this historical note.

## Runtime Flow
Runtime tracks pointer markers for min, max, and average rendered positions. It does not implement the earlier historical concept as a separate feature.

## Configuration
Supported config is `realism.pointer_markers`, not a stat-marker schema.

## Behaviour
Users get the newer pointer-marker feature set rather than the superseded design.

## Rendering
Pointer markers render as dedicated marker assets above the gauge, using final displayed geometry.

## Tests
- [`internal/dashboard/gauges/pointer_markers_test.go`](../../../internal/dashboard/gauges/pointer_markers_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
This implementation record exists mainly to prevent the historical note from being mistaken for live scope.

## Deviations from Design
The design explicitly says not to implement the file as written. The code correctly implements the replacement feature instead.

## Remaining Work
No work should target this superseded design. Use the pointer-marker records instead.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

# `pointer_markers`

Design reference: [`docs/Designs/RealismBehaviour/witness-markers.md`](../../Designs/RealismBehaviour/witness-markers.md)

## Purpose
Tracks the pointer-marker feature that records and renders reference positions for final displayed indicator geometry.

## Implementation Status
Status: **Implemented**.

Radial and bar gauges support `realism.pointer_markers` with min, max, and average marker behaviour.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/pointer_markers.go`](../../../internal/dashboard/gauges/pointer_markers.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`
- `PointerMarkersConfig`

## Functions and Methods
- Pointer-marker advance helpers in `pointer_markers.go`
- `resolveMovementState`

## Runtime Flow
Runtime records markers from final displayed geometry, including overshoot and other movement effects, and can keep average markers active across idle ticks until settled.

## Configuration
Radial and bar packages accept `realism.pointer_markers` with per-marker enablement and supporting options.

## Behaviour
Markers can track min, max, and running average positions while remaining display-only and deterministic.

## Rendering
Scene composition renders marker assets above the live pointer/bar and before overlay layers, using family-appropriate geometry.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/pointer_markers_test.go`](../../../internal/dashboard/gauges/pointer_markers_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The feature applies only to radial and bar families.

## Deviations from Design
The implementation matches the current pointer-marker design rather than the older superseded stat-marker note.

## Remaining Work
No known design work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.

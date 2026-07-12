# Gauge Stat Markers — Implementation

## Purpose
Audits the superseded statistical-marker note against current code.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/pointer_markers.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`

## Types
- `PointerMarkersConfig`
- `PointerMarkerState`

## Functions and Methods
- `AdvanceMinMaxPointerMarkers`
- `AdvanceAveragePointerMarker`
- `updatePointerMarkerState`

## Runtime Flow
Current code implements pointer markers as a separate design. No statistical-marker feature or config named by this historical note was found.

## Configuration
Current code accepts `realism.pointer_markers`. No stat-marker config key was found.

## Behaviour
Pointer markers are implemented. The older stat-marker design is not.

## Rendering
Pointer markers render through the radial and bar scene helpers. No stat-marker-specific render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record distinguishes the superseded design from its replacement feature.

## Deviations from Design
The design explicitly says not to implement the file as written. Current code instead implements the replacement pointer-marker feature.

## Remaining Work
No work should target this superseded design unless the design itself changes.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/pointer_markers.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`

Symbols verified:
- `PointerMarkersConfig`
- `PointerMarkerState`
- `AdvanceMinMaxPointerMarkers`
- `AdvanceAveragePointerMarker`
- `updatePointerMarkerState`

Searches performed:
- `stat markers`
- `pointer_markers`
- `stat_markers`

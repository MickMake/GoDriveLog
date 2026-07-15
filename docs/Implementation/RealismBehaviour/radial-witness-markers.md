# `pointer_markers` — Implementation

## Purpose
Audits current pointer-marker support for radial and bar gauges.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/pointer_markers.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `PointerMarkersConfig`
- `PointerMarkerState`
- `PointerMarkerValueState`

## Functions and Methods
- `validateRealism`
- `AdvanceMinMaxPointerMarkers`
- `AdvanceAveragePointerMarker`
- `RenderedPointerMarkerPosition`
- `updatePointerMarkerState`
- `appendRadialPointerMarkerParts`
- `appendBarPointerMarkerParts`

## Runtime Flow
Runtime stores pointer-marker state separately from movement state and updates it through `updatePointerMarkerState` using the current rendered gauge position.

## Configuration
`PointerMarkersConfig` accepts `min`, `max`, `average`, and optional `window`. `validateRealism` restricts pointer markers to radial and bar gauges.

## Behaviour
Current code supports min, max, and average pointer markers. Daily reset or rolling-window behaviour is selected by whether `window` is configured.

## Rendering
Radial and bar scenes append dedicated pointer-marker parts after the live indicator and before overlay layers.

## Tests
- `TestLoadPackageLoadsPointerMarkersConfig`
- `TestLoadPackageLoadsDisabledPointerMarkersConfig`
- `TestAdvanceMinMaxPointerMarkersResetsAtLocalMidnight`
- `TestAdvanceAveragePointerMarkerUsesFixedTenSecondTimeConstant`
- `TestRenderedPointerMarkerPositionUsesFinalRenderedGeometry`
- `TestRuntimeInitializesPointerMarkerStateStore`
- `TestRuntimeRadialGaugeWidgetRendersPointerMarkersAboveNeedleBeforeOverlay`
- `TestRuntimeBarGaugeWidgetRendersPointerMarkersAboveBarBeforeGlass`

## Limitations
Only radial and bar gauges implement pointer markers.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/pointer_markers.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `PointerMarkersConfig`
- `PointerMarkerState`
- `PointerMarkerValueState`
- `validateRealism`
- `AdvanceMinMaxPointerMarkers`
- `AdvanceAveragePointerMarker`
- `RenderedPointerMarkerPosition`
- `updatePointerMarkerState`
- `appendRadialPointerMarkerParts`
- `appendBarPointerMarkerParts`

Configuration verified:
- `realism.pointer_markers`
- `min`
- `max`
- `average`
- `window`

Tests inspected:
- `TestLoadPackageLoadsPointerMarkersConfig`
- `TestLoadPackageLoadsDisabledPointerMarkersConfig`
- `TestAdvanceMinMaxPointerMarkersResetsAtLocalMidnight`
- `TestAdvanceAveragePointerMarkerUsesFixedTenSecondTimeConstant`
- `TestRenderedPointerMarkerPositionUsesFinalRenderedGeometry`
- `TestRuntimeInitializesPointerMarkerStateStore`
- `TestRuntimeRadialGaugeWidgetRendersPointerMarkersAboveNeedleBeforeOverlay`
- `TestRuntimeBarGaugeWidgetRendersPointerMarkersAboveBarBeforeGlass`

Searches performed:
- `pointer_markers`
- `PointerMarkersConfig`
- `RenderedPointerMarkerPosition`

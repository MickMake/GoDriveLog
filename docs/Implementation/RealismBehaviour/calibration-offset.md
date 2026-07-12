# `calibration_offset` — Implementation

## Purpose
Audits radial calibration offset support in current code.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/gauges/pointer_markers.go`

## Types
- `Realism`
- `ValueMap`

## Functions and Methods
- `validateRealism`
- `radialCalibrationAngle`
- `RenderedPointerMarkerPosition`

## Runtime Flow
The current display angle is adjusted by `radialCalibrationAngle` when radial scenes and radial pointer-marker positions are built.

## Configuration
`Realism` declares `CalibrationOffset *float64`. `validateRealism` accepts it for radial gauges only and rejects non-finite values.

## Behaviour
A configured offset changes the displayed radial angle without changing the source value.

## Rendering
The offset is applied in radial scene building and in radial pointer-marker position calculation. Clamp behaviour follows `ValueMap.Clamp`.

## Tests
- `TestLoadPackageAcceptsRadialCalibrationOffset`
- `TestLoadPackageRejectsInvalidRadialCalibrationOffset`
- `TestRadialSceneCalibrationOffsetZeroPreservesAngle`
- `TestRadialSceneCalibrationOffsetAppliesPositiveAndNegativeDegrees`
- `TestRadialSceneCalibrationOffsetClampsToDialBounds`
- `TestRuntimeRadialGaugeWidgetCalibrationOffsetChangesOnlyDisplayedAngle`
- `TestRenderedPointerMarkerPositionPreservesUnclampedRadialCalibrationOffset`

## Limitations
Only radial gauges implement this option.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/gauges/pointer_markers.go`

Symbols verified:
- `Realism`
- `ValueMap`
- `validateRealism`
- `radialCalibrationAngle`
- `RenderedPointerMarkerPosition`

Configuration verified:
- `realism.calibration_offset`

Tests inspected:
- `TestLoadPackageAcceptsRadialCalibrationOffset`
- `TestLoadPackageRejectsInvalidRadialCalibrationOffset`
- `TestRadialSceneCalibrationOffsetZeroPreservesAngle`
- `TestRadialSceneCalibrationOffsetAppliesPositiveAndNegativeDegrees`
- `TestRadialSceneCalibrationOffsetClampsToDialBounds`
- `TestRuntimeRadialGaugeWidgetCalibrationOffsetChangesOnlyDisplayedAngle`
- `TestRenderedPointerMarkerPositionPreservesUnclampedRadialCalibrationOffset`

Searches performed:
- `calibration_offset`
- `realism.calibration_offset`
- `radialCalibrationAngle`

# `drum_slop` — Implementation

## Purpose
Audits current odometer drum-slop support.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`

## Types
- `Realism`
- `Odometer`

## Functions and Methods
- `validateRealism`
- `odometerDrumSlop`
- `OdometerSceneWithWheelOffsets`

## Runtime Flow
No separate runtime state is required. Drum slop is read when odometer scenes are built from wheel offsets.

## Configuration
`Realism` declares `DrumSlop []int` and tracks whether the field was set through `DrumSlopSet`. `validateRealism` restricts the field to odometers, requires one entry per wheel, and bounds each offset by wheel window height.

## Behaviour
Each wheel can render with a fixed vertical offset.

## Rendering
`OdometerSceneWithWheelOffsets` adds the configured per-wheel slop to each wheel position before building the `wheel_strip` part.

## Tests
- `TestLoadPackageLoadsOdometerDrumSlopRealism`
- `TestLoadPackageRejectsExplicitEmptyDrumSlopOnNonOdometerGauge`
- `TestLoadPackageRejectsExplicitEmptyOdometerDrumSlop`
- `TestLoadPackageRejectsInvalidOdometerDrumSlop`
- `TestOdometerSceneAppliesConfiguredDrumSlopToWheelPositions`
- `TestOdometerSceneDefaultsToNoDrumSlop`

## Limitations
This feature is static; it does not create motion or reversal slack.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`

Symbols verified:
- `Realism`
- `Odometer`
- `validateRealism`
- `odometerDrumSlop`
- `OdometerSceneWithWheelOffsets`

Configuration verified:
- `realism.drum_slop`

Tests inspected:
- `TestLoadPackageLoadsOdometerDrumSlopRealism`
- `TestLoadPackageRejectsExplicitEmptyDrumSlopOnNonOdometerGauge`
- `TestLoadPackageRejectsExplicitEmptyOdometerDrumSlop`
- `TestLoadPackageRejectsInvalidOdometerDrumSlop`
- `TestOdometerSceneAppliesConfiguredDrumSlopToWheelPositions`
- `TestOdometerSceneDefaultsToNoDrumSlop`

Searches performed:
- `drum_slop`
- `DrumSlopSet`

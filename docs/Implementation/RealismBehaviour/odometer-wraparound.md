# `wraparound` — Implementation

## Purpose
Audits current wraparound behaviour and configuration for odometer wheels.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `Odometer`
- `ScenePart`

## Functions and Methods
- `validateRealism`
- `OdometerWheelStripOffsets`
- `OdometerTravelWheelOffsets`
- `odometerRoutedTargetOffsets`
- `odometerWheelCircular`
- `resolveOdometerMovementState`

## Runtime Flow
Odometer movement routing uses `OdometerTravelWheelOffsets` and `odometerRoutedTargetOffsets` so wheel offsets continue across strip boundaries.

## Configuration
`Realism` declares `Wraparound *bool`, and `validateRealism` accepts the field for odometer gauges only. Current scene rendering does not branch on that field. `odometerWheelCircular` returns `true` unconditionally and is documented in code as compatibility-only for `realism.wraparound`.

## Behaviour
Current odometer rendering is circular across wheel boundaries, including `9 -> 0` and `0 -> 9` style transitions.

## Rendering
`OdometerSceneWithWheelOffsets` always marks wheel-strip parts as wraparound, and the wheel-position helpers route offsets as circular strips.

## Tests
- `TestLoadPackageLoadsOdometerWraparoundRealism`
- `TestLoadPackageRejectsWraparoundOnNonOdometerGauge`
- `TestOdometerSceneAlwaysUsesCircularWheelRendering`
- `TestOdometerInterpolatedWheelOffsetsUseInfiniteForwardRoutingAcrossNineToZero`
- `TestOdometerInterpolatedWheelOffsetsUseInfiniteBackwardRoutingAcrossZeroToNine`

## Limitations
Current code does not provide a verified way to disable circular wheel rendering with `realism.wraparound: false`.

## Deviations from Design
The visual behaviour exists, but the `realism.wraparound` field is compatibility-only in scene code rather than an active on/off switch.

## Remaining Work
If the design requires a true configuration switch, current code would need to honor `false` explicitly.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `Odometer`
- `ScenePart`
- `validateRealism`
- `OdometerWheelStripOffsets`
- `OdometerTravelWheelOffsets`
- `odometerRoutedTargetOffsets`
- `odometerWheelCircular`
- `resolveOdometerMovementState`

Configuration verified:
- `realism.wraparound`

Tests inspected:
- `TestLoadPackageLoadsOdometerWraparoundRealism`
- `TestLoadPackageRejectsWraparoundOnNonOdometerGauge`
- `TestOdometerSceneAlwaysUsesCircularWheelRendering`
- `TestOdometerInterpolatedWheelOffsetsUseInfiniteForwardRoutingAcrossNineToZero`
- `TestOdometerInterpolatedWheelOffsetsUseInfiniteBackwardRoutingAcrossZeroToNine`

Searches performed:
- `wraparound`
- `odometerWheelCircular`
- `OdometerTravelWheelOffsets`

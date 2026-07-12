# Bar Realism Scope — Implementation

## Purpose
Audits the current bar-gauge realism coverage against the scope note.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/gauges/pointer_markers.go`

## Types
- `Realism`
- `DampingConfig`
- `OvershootConfig`
- `PointerMarkersConfig`

## Functions and Methods
- `validateRealism`
- `resolveBarMovementState`
- `barDampingDuration`
- `BarSceneWithPointerMarkers`
- `RenderedPointerMarkerPosition`

## Runtime Flow
`resolveBarMovementState` implements bar damping, hysteresis, stiction, overshoot, and peg bounce. Pointer marker state is maintained separately through `updatePointerMarkerState` and the helpers in `pointer_markers.go`.

## Configuration
`Realism` includes `Damping`, `Hysteresis`, `Stiction`, `Overshoot`, `PegBounce`, `PointerMarkers`, and `MovementPolicy`. It does not include `stepped_fill` or `quantized_fill`.

## Behaviour
Several planned bar realism behaviours are present. The fill still resolves as a continuous reveal height; no bar-specific stepped or quantized fill behaviour was found.

## Rendering
`BarSceneWithPointerMarkers` renders the resolved reveal height and optional pointer marker parts. It does not render stepped or quantized fill states.

## Tests
- `TestLoadPackageLoadsBarDampingWithDirectionalTiming`
- `TestLoadPackageAcceptsBarHysteresis`
- `TestLoadPackageAcceptsBarStiction`
- `TestLoadPackageAcceptsBarOvershoot`
- `TestLoadPackageAcceptsBarPegBounce`
- `TestLoadPackageLoadsPointerMarkersConfig`
- `TestRuntimeBarGaugeDampingAnimatesRisingReveal`
- `TestRuntimeBarGaugeHysteresisOffsetsRisingApproach`
- `TestRuntimeBarGaugeStictionBelowThresholdHoldsDisplay`
- `TestRuntimeBarGaugeOvershootAnimatesRisingReveal`
- `TestRuntimeBarGaugePegBounceAtMaxStopSettlesBackToLimit`

## Limitations
This scope note includes behaviours that are not in current code.

## Deviations from Design
`stepped_fill` and `quantized_fill` are named in the design note but were not found in current code.

## Remaining Work
Implement or explicitly defer `stepped_fill` and `quantized_fill` if this scope note remains active.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/gauges/pointer_markers.go`

Symbols verified:
- `Realism`
- `DampingConfig`
- `OvershootConfig`
- `PointerMarkersConfig`
- `validateRealism`
- `resolveBarMovementState`
- `barDampingDuration`
- `BarSceneWithPointerMarkers`
- `RenderedPointerMarkerPosition`

Configuration verified:
- `realism.damping`
- `realism.hysteresis`
- `realism.stiction`
- `realism.overshoot`
- `realism.peg_bounce`
- `realism.pointer_markers`
- `realism.movement_policy`

Tests inspected:
- `TestLoadPackageLoadsBarDampingWithDirectionalTiming`
- `TestLoadPackageAcceptsBarHysteresis`
- `TestLoadPackageAcceptsBarStiction`
- `TestLoadPackageAcceptsBarOvershoot`
- `TestLoadPackageAcceptsBarPegBounce`
- `TestLoadPackageLoadsPointerMarkersConfig`
- `TestRuntimeBarGaugeDampingAnimatesRisingReveal`
- `TestRuntimeBarGaugeHysteresisOffsetsRisingApproach`
- `TestRuntimeBarGaugeStictionBelowThresholdHoldsDisplay`
- `TestRuntimeBarGaugeOvershootAnimatesRisingReveal`
- `TestRuntimeBarGaugePegBounceAtMaxStopSettlesBackToLimit`

Searches performed:
- `stepped_fill`
- `quantized_fill`
- `pointer_markers`

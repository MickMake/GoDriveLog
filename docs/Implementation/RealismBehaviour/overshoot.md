# `overshoot` — Implementation

## Purpose
Audits current overshoot support for radial and bar gauges.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `OvershootConfig`
- `ValueMap`

## Functions and Methods
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `radialOvershootTarget`
- `radialOvershootTravelDuration`

## Runtime Flow
Radial overshoot is handled in `resolveMovementState`. Bar overshoot is handled in `resolveBarMovementState`. Both use `radialOvershootTarget` and the overshoot travel/settle timing helpers.

## Configuration
`OvershootConfig` accepts `ratio`, `min_change_ratio`, `max_span_ratio`, `settle_mode`, `settle_cycles`, `settle_damping`, and `allow_extremes`. `validateRealism` restricts the configuration by gauge family.

## Behaviour
The displayed value can move past the target and settle back while the stored raw source value remains unchanged.

## Rendering
The render path uses the current movement state; overshoot itself is resolved before scene generation.

## Tests
- `TestLoadPackageAcceptsRadialOvershoot`
- `TestLoadPackageAcceptsBarOvershoot`
- `TestLoadPackageRejectsUnknownRadialOvershootKey`
- `TestLoadPackageRejectsInvalidRadialOvershoot`
- `TestRuntimeRadialGaugeOvershootAnimatesWithoutDamping`
- `TestRuntimeRadialGaugeOvershootStaysBoundedAndSettlesOnTarget`
- `TestRuntimeBarGaugeOvershootAnimatesRisingReveal`
- `TestRuntimeBarGaugeOvershootStaysBoundedAndSettlesOnTarget`

## Limitations
Only radial and bar gauges implement overshoot.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `OvershootConfig`
- `ValueMap`
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `radialOvershootTarget`
- `radialOvershootTravelDuration`

Configuration verified:
- `realism.overshoot`
- `ratio`
- `min_change_ratio`
- `max_span_ratio`
- `settle_mode`
- `settle_cycles`
- `settle_damping`
- `allow_extremes`

Tests inspected:
- `TestLoadPackageAcceptsRadialOvershoot`
- `TestLoadPackageAcceptsBarOvershoot`
- `TestLoadPackageRejectsUnknownRadialOvershootKey`
- `TestLoadPackageRejectsInvalidRadialOvershoot`
- `TestRuntimeRadialGaugeOvershootAnimatesWithoutDamping`
- `TestRuntimeRadialGaugeOvershootStaysBoundedAndSettlesOnTarget`
- `TestRuntimeBarGaugeOvershootAnimatesRisingReveal`
- `TestRuntimeBarGaugeOvershootStaysBoundedAndSettlesOnTarget`

Searches performed:
- `overshoot`
- `OvershootConfig`
- `radialOvershootTarget`

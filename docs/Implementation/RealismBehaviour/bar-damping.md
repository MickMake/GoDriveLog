# `damping` — Implementation

## Purpose
Audits current damping support for radial and bar gauges.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `DampingConfig`

## Functions and Methods
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `barDampingDuration`
- `effectiveMovementPolicy`

## Runtime Flow
Radial damping is resolved in `resolveMovementState`. Bar damping is resolved in `resolveBarMovementState`, with `barDampingDuration` choosing direction-specific durations when configured.

## Configuration
`DampingConfig` accepts a boolean scalar or a mapping. Mapping fields are `enabled`, `rise_ms`, and `fall_ms`. `normalizePackage` does not create damping by default. `validateRealism` restricts `rise_ms` and `fall_ms` to bar gauges and requires positive durations.

## Behaviour
When enabled, the displayed value moves over time instead of jumping immediately. For bars, rising and falling moves can use different durations.

## Rendering
The scene layer renders the current resolved display value; damping itself is handled in runtime movement state.

## Tests
- `TestLoadPackageAcceptsRadialDamping`
- `TestLoadPackageLoadsBarDampingWithDirectionalTiming`
- `TestLoadPackageRejectsInvalidDamping`
- `TestRuntimeRadialGaugeDampingAnimatesWithDefaultLinearCurve`
- `TestRuntimeBarGaugeDampingAnimatesRisingReveal`
- `TestRuntimeBarGaugeDampingAnimatesFallingReveal`
- `TestRuntimeBarGaugeDampingSettlesAtFinalReveal`

## Limitations
Only radial and bar gauges implement damping.

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
- `DampingConfig`
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `barDampingDuration`
- `effectiveMovementPolicy`

Configuration verified:
- `realism.damping`
- `enabled`
- `rise_ms`
- `fall_ms`

Tests inspected:
- `TestLoadPackageAcceptsRadialDamping`
- `TestLoadPackageLoadsBarDampingWithDirectionalTiming`
- `TestLoadPackageRejectsInvalidDamping`
- `TestRuntimeRadialGaugeDampingAnimatesWithDefaultLinearCurve`
- `TestRuntimeBarGaugeDampingAnimatesRisingReveal`
- `TestRuntimeBarGaugeDampingAnimatesFallingReveal`
- `TestRuntimeBarGaugeDampingSettlesAtFinalReveal`

Searches performed:
- `damping`
- `rise_ms`
- `fall_ms`
- `barDampingDuration`

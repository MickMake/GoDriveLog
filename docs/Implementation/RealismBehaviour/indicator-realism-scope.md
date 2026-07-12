# Indicator Realism Scope — Implementation

## Purpose
Audits current indicator realism support against the scope note.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`

## Types
- `Realism`
- `ThermalFadeConfig`

## Functions and Methods
- `validateRealism`
- `resolveIndicatorThermalFadeState`
- `IndicatorSceneWithOnAlpha`

## Runtime Flow
Indicator-specific movement state is handled by `resolveIndicatorThermalFadeState` when `ThermalFade` is configured.

## Configuration
`Realism` declares `ThermalFade *ThermalFadeConfig`. No other indicator-specific realism field from this scope note was found.

## Behaviour
Indicator gauges can fade on and off over time. No additional indicator realism behaviour from this scope note was found.

## Rendering
`IndicatorSceneWithOnAlpha` blends the `off` and `on` layers using the resolved alpha value.

## Tests
- `TestLoadPackageLoadsIndicatorThermalFade`
- `TestLoadPackageRejectsInvalidIndicatorThermalFade`
- `TestRuntimeIndicatorThermalFadeDisabledRemainsImmediate`
- `TestRuntimeIndicatorThermalFadeOnTransition`
- `TestRuntimeIndicatorThermalFadeOffTransition`

## Limitations
The audited scope is limited to `thermal_fade`.

## Deviations from Design
The broader indicator backlog described by the scope note was not found in current code.

## Remaining Work
Add additional indicator realism only if later designs define it explicitly.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`

Symbols verified:
- `Realism`
- `ThermalFadeConfig`
- `validateRealism`
- `resolveIndicatorThermalFadeState`
- `IndicatorSceneWithOnAlpha`

Configuration verified:
- `realism.thermal_fade`
- `rise_ms`
- `fall_ms`

Tests inspected:
- `TestLoadPackageLoadsIndicatorThermalFade`
- `TestLoadPackageRejectsInvalidIndicatorThermalFade`
- `TestRuntimeIndicatorThermalFadeDisabledRemainsImmediate`
- `TestRuntimeIndicatorThermalFadeOnTransition`
- `TestRuntimeIndicatorThermalFadeOffTransition`

Searches performed:
- `thermal_fade`
- `indicator realism`

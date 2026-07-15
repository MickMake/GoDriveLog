# `thermal_fade` — Implementation

## Purpose
Audits current indicator thermal-fade support.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `ThermalFadeConfig`

## Functions and Methods
- `validateRealism`
- `resolveIndicatorThermalFadeState`
- `indicatorFadeDuration`
- `IndicatorSceneWithOnAlpha`

## Runtime Flow
Indicator fade state is handled by `resolveIndicatorThermalFadeState`, which updates display alpha over time using `indicatorFadeDuration`.

## Configuration
`ThermalFadeConfig` accepts `rise_ms` and `fall_ms`. `validateRealism` restricts it to indicator gauges and requires both values to be positive.

## Behaviour
Indicator transitions can warm up and cool down over time instead of changing state immediately.

## Rendering
`IndicatorSceneWithOnAlpha` renders the current alpha by blending the `off` and `on` layers.

## Tests
- `TestLoadPackageLoadsIndicatorThermalFade`
- `TestLoadPackageRejectsInvalidIndicatorThermalFade`
- `TestRuntimeIndicatorThermalFadeDisabledRemainsImmediate`
- `TestRuntimeIndicatorThermalFadeOnTransition`
- `TestRuntimeIndicatorThermalFadeOffTransition`

## Limitations
Only indicator gauges implement thermal fade.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `ThermalFadeConfig`
- `validateRealism`
- `resolveIndicatorThermalFadeState`
- `indicatorFadeDuration`
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
- `indicatorFadeDuration`

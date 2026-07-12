# `needle_shadow` — Implementation

## Purpose
Audits current radial needle-shadow support.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`

## Types
- `Realism`
- `NeedleShadowConfig`

## Functions and Methods
- `validateRealism`
- `RadialSceneWithPointerMarkers`
- `needleShadowEnabled`
- `needleShadowAlpha`

## Runtime Flow
The option does not create its own runtime state. Radial scene generation checks `NeedleShadow` directly when building the current scene.

## Configuration
`NeedleShadowConfig` accepts `offset` and `alpha`. `normalizePackage` fills in the default alpha when a shadow is configured without one. `validateRealism` restricts the option to radial gauges, requires a two-element offset, and requires `alpha` to be finite and between 0 and 1.

## Behaviour
When enabled, a shadow copy of the needle is drawn behind the main needle using the configured offset and alpha.

## Rendering
`RadialSceneWithPointerMarkers` appends a `needle_shadow` scene part before the live needle. Angle calculation happens first through `radialAngle` and `radialCalibrationAngle`; the shadow does not perform its own clamping logic.

## Tests
- `TestLoadPackageAcceptsRadialNeedleShadow`
- `TestLoadPackageRejectsInvalidRadialNeedleShadow`
- `TestRadialSceneAddsNeedleShadowBeforeNeedleWhenConfigured`
- `TestRuntimeRadialGaugeWidgetIncludesNeedleShadowBeforeNeedle`

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

Symbols verified:
- `Realism`
- `NeedleShadowConfig`
- `validateRealism`
- `RadialSceneWithPointerMarkers`
- `needleShadowEnabled`
- `needleShadowAlpha`

Configuration verified:
- `realism.needle_shadow`
- `offset`
- `alpha`

Tests inspected:
- `TestLoadPackageAcceptsRadialNeedleShadow`
- `TestLoadPackageRejectsInvalidRadialNeedleShadow`
- `TestRadialSceneAddsNeedleShadowBeforeNeedleWhenConfigured`
- `TestRuntimeRadialGaugeWidgetIncludesNeedleShadowBeforeNeedle`

Searches performed:
- `needle_shadow`
- `NeedleShadowConfig`
- `needleShadowAlpha`

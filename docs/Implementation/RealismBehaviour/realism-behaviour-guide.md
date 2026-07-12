# Realism Behaviour Guide — Implementation

## Purpose
Audits which parts of the realism behaviour guide have matching implementation evidence in current code.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/gauges/pointer_markers.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `Odometer`
- `PointerMarkersConfig`
- `DampingConfig`
- `OvershootConfig`
- `ThermalFadeConfig`
- `NeedleShadowConfig`

## Functions and Methods
- `validateRealism`
- `normalizePackage`
- `resolveMovementState`
- `resolveBarMovementState`
- `resolveOdometerMovementState`
- `resolveIndicatorThermalFadeState`
- `AdvanceMinMaxPointerMarkers`
- `AdvanceAveragePointerMarker`

## Runtime Flow
Current guide-backed behaviour is spread across gauge package parsing, dashboard runtime movement state, pointer marker state, and scene generation.

## Configuration
Implemented current-code keys include `wraparound`, `carry_drag`, `snap_settle`, `drum_slop`, `hysteresis`, `damping`, `stiction`, `overshoot`, `peg_bounce`, `pointer_markers`, `thermal_fade`, `needle_shadow`, `calibration_offset`, and `movement_policy`.

## Behaviour
Many guide entries have matching code, but several guide entries do not. This guide remains broader than the current implementation surface.

## Rendering
Guide-backed features that exist are rendered through the current scene builders rather than a shared realism layer.

## Tests
- `TestLoadPackageAcceptsSharedMovementPolicies`
- `TestLoadPackageLoadsPointerMarkersConfig`
- `TestLoadPackageLoadsIndicatorThermalFade`
- `TestLoadPackageAcceptsRadialNeedleShadow`
- `TestRuntimeGaugeMovementLifecycle`
- `TestRuntimeIndicatorThermalFadeOnTransition`
- `TestRuntimeRadialGaugeWidgetRendersPointerMarkersAboveNeedleBeforeOverlay`

## Limitations
This record is intentionally a summary. Per-feature records contain narrower evidence.

## Deviations from Design
The guide includes behaviours that are not implemented in current code, including `backlash`, `needle_trail`, `realism.imperfections`, `lighting-mode`, and several numeric/segmented candidates.

## Remaining Work
Keep the per-feature records and this summary aligned as implementation changes.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/gauges/pointer_markers.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `Odometer`
- `PointerMarkersConfig`
- `DampingConfig`
- `OvershootConfig`
- `ThermalFadeConfig`
- `NeedleShadowConfig`
- `validateRealism`
- `normalizePackage`
- `resolveMovementState`
- `resolveBarMovementState`
- `resolveOdometerMovementState`
- `resolveIndicatorThermalFadeState`
- `AdvanceMinMaxPointerMarkers`
- `AdvanceAveragePointerMarker`

Tests inspected:
- `TestLoadPackageAcceptsSharedMovementPolicies`
- `TestLoadPackageLoadsPointerMarkersConfig`
- `TestLoadPackageLoadsIndicatorThermalFade`
- `TestLoadPackageAcceptsRadialNeedleShadow`
- `TestRuntimeGaugeMovementLifecycle`
- `TestRuntimeIndicatorThermalFadeOnTransition`
- `TestRuntimeRadialGaugeWidgetRendersPointerMarkersAboveNeedleBeforeOverlay`

Searches performed:
- `backlash`
- `needle_trail`
- `realism.imperfections`
- `lighting_mode`
- `pointer_markers`

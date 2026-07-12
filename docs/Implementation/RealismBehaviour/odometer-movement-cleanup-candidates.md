# Odometer Movement Cleanup Candidates — Implementation

## Purpose
Audits the cleanup note for reserved odometer movement values.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Odometer`

## Functions and Methods
- `normalizePackage`
- `validateOdometer`
- `resolveOdometerMovementState`
- `applyOdometerMovementCurve`

## Runtime Flow
`resolveOdometerMovementState` implements `instant`, `linear`, `ease_out`, and `bell` through `applyOdometerMovementCurve`.

## Configuration
`validateOdometer` accepts `smooth` and `click`, but `normalizePackage` logs that both are recognised but not implemented and falls them back to `instant`.

## Behaviour
Reserved values remain loadable for compatibility, but they do not create distinct movement behaviour.

## Rendering
Rendering follows the actual movement mode after normalization.

## Tests
- `TestLoadPackageWarnsAndFallsBackForRecognizedOdometerMovementValues`
- `TestRuntimeOdometerGaugeRecognizedMovementFallbacksStayInstant`

## Limitations
The reserved values are compatibility paths, not distinct implementations.

## Deviations from Design
The design note leaves `smooth` and `click` undefined until a later focused slice. Current code follows that by normalizing them to `instant`.

## Remaining Work
Define or remove the reserved values only if a later design requires that work.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Odometer`
- `normalizePackage`
- `validateOdometer`
- `resolveOdometerMovementState`
- `applyOdometerMovementCurve`

Configuration verified:
- `smooth`
- `click`
- `instant`
- `linear`
- `ease_out`
- `bell`

Tests inspected:
- `TestLoadPackageWarnsAndFallsBackForRecognizedOdometerMovementValues`
- `TestRuntimeOdometerGaugeRecognizedMovementFallbacksStayInstant`

Searches performed:
- `smooth`
- `click`
- `normalizePackage`
